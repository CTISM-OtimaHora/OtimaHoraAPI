import pandas
from pandas import DataFrame
import sys
import requests
# 296 ate 329
SERVER = "http://localhost:3000"

sess = requests.session()
if len(sys.argv) >= 6:
    sess.cookies.set("id", sys.argv[5])
else:
    sess.get(SERVER+"/add-session")
res = sess.post(SERVER+"/curso/add", json={'Id':0, 'Nome':sys.argv[2], 'Etapas':[]})
curso_id = int(res.text)
curso = sess.get(SERVER+'/curso/get/'+str(curso_id)).json()


new_dispo = [[1 for _ in range(5)] for _ in range(5)]

cs = ['disciplina', 'nome da turma', 'turma', 'ch', 'semestre/ano',
       'professor', 'Encargo semanal',
       'Divisão carga horária disciplina na semana',
       'Laboratório ou Recurso Didático', 'Divisões de turmas']

f = pandas.read_excel(sys.argv[1])
f: DataFrame = f.iloc[int(sys.argv[3])-2:int(sys.argv[4])-1]
f.columns = [x.strip() for x in f.columns]
rem = []
for c in f.columns:
    s = c.strip()
    if s not in cs:
        rem.append(s)
f.drop(columns=rem, inplace=True)


disciplinas={}
professores={}
etapas={}
turmas={}
recursos={}


for i in f.iterrows():
    row = dict(i[1])

    if row['disciplina'] not in disciplinas:
        res = sess.post(SERVER+"/disciplina/add", json={'Id':0, 'Nome':row['disciplina'], 'Dispo': new_dispo})
        if res.status_code != 200:
            print("ERRO adicionar disciplina")
            print(res.status_code)
            print(res.text)
            exit()
        disciplinas[row['disciplina']] = int(res.text)

    if row["professor"] not in professores:
        res = sess.post(SERVER+"/professor/add", json={'Id':0, 'Nome':row['professor'], 'Dispo': new_dispo, 'Disciplinas_ids':[disciplinas[row['disciplina']]]})
        if res.status_code != 200:
            print("ERRO adicionar professor")
            print(res.status_code)
            print(res.text)
            exit()
        professores[row['professor']] = int(res.text)
    else:
        prof = sess.get(SERVER+f"/professor/get/{professores[row['professor']]}").json()
        prof["Disciplinas_ids"].append(disciplinas[row['disciplina']])
        res = sess.put(SERVER+f"/professor/set/{professores[row['professor']]}", json=prof)
        if res.status_code != 200:
            print("ERRO adicionar disciplina ao professor")
            print(res.status_code)
            print(res.text)
            exit()

    if row["semestre/ano"] not in etapas:
        j={
            'Idx_in_Curso': len(curso['Etapas']),
            'Curso_id':curso_id,
            'Curriculo': {disciplinas[row['disciplina']]: {"Horas": int(row['Encargo semanal']), "Formato":str(row['Divisão carga horária disciplina na semana']).replace(" ", "")}},
            'Turmas':[]
        }
        curso['Etapas'].append(j)
        res = sess.put(SERVER+f"/curso/set/{curso_id}", json=curso)
        if res.status_code != 200:
            print("ERRO adicionar etapa")
            print(res.status_code)
            print(res.text)
            exit()
        etapas[row['semestre/ano']] = len(curso['Etapas'])-1
    else:
        curso['Etapas'][etapas[row['semestre/ano']]]["Curriculo"][disciplinas[row['disciplina']]] = {"Horas": int(row['Encargo semanal']), "Formato":str(row['Divisão carga horária disciplina na semana']).replace(" ", "")}
        res = sess.put(SERVER+f"/curso/set/{curso_id}", json=curso)
        if res.status_code != 200:
            print("ERRO adicionar disciplina a curriculo")
            print(res.status_code)
            print(res.text)
            exit()
    if row["nome da turma"] not in turmas:
        t = {
            "Id": 0,
            "Curso_id": curso_id,
            "Etapa_idx": etapas[row['semestre/ano']],
            "Nome": str(int(row["nome da turma"])),
            "Dispo": new_dispo
        }

        res = sess.post(SERVER+f"/turma/add", json=t)
        if res.status_code != 200:
            print("ERRO adicionar turma")
            print(res.status_code)
            print(res.text)
            exit()
        turmas[row['nome da turma']] =  int(res.text)
        curso = sess.get(SERVER+f"/curso/get/{curso_id}").json()
    if not pandas.isna(row['Laboratório ou Recurso Didático']) and row['Laboratório ou Recurso Didático'] not in recursos:
        res = sess.post(SERVER+"/recurso/add", json={'Id':0, 'Nome':row['Laboratório ou Recurso Didático'], 'Dispo': new_dispo})
        if res.status_code != 200:
            print("ERRO adicionar recurso")
            print(res.status_code)
            print(res.text)
            exit()
        recursos[row['Laboratório ou Recurso Didático']] = int(res.text)

    contrato = []
    contrato.append({"Id": disciplinas[row['disciplina']], "Nome": row['disciplina'], "Tipo": 'disciplina'})
    contrato.append({"Id": professores[row['professor']], "Nome": row['professor'], "Tipo": 'professor'})
    contrato.append({"Id": turmas[row['nome da turma']], "Nome": str(int(row['nome da turma'])), "Tipo": 'turma'})
    if not pandas.isna(row['Laboratório ou Recurso Didático']):
        contrato.append({"Id": recursos[row['Laboratório ou Recurso Didático']], "Nome": row['Laboratório ou Recurso Didático'], "Tipo": 'recurso'})
    res = sess.post(f"{SERVER}/contrato/add", json=contrato)
    if res.status_code != 200:
        print("ERRO adicionar contrato")
        print(contrato)
        print(res.status_code)
        print(res.text)
        exit()

print(sess.cookies)

