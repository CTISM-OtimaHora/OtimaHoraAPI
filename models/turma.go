package models

type Turma struct {
    Id      int
    Curso_id int
    Etapa_idx int
    Idx_in_etapa int
    Nome    string
    Dispo Disponibilidade
}

func (t Turma) GetId() int {
    return t.Id
}
func (t Turma) GetNome() string {
    return t.Nome
}
func (t Turma) GetDispo() Disponibilidade {
    return t.Dispo
}
