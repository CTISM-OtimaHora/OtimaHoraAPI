package models

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ParticipanteContrato interface {
    GetId() int
    GetNome() string
    GetDispo() Disponibilidade
}

type Participante_to_json struct {
    Id int
    Nome string
    Tipo string
}

func Pjson_to_P(p Participante_to_json, s * Session) (ParticipanteContrato, error) {
        var e ParticipanteContrato
        switch p.Tipo {
        case "disciplina":
            item, ok := s.Disciplinas[p.Id]
            e = item
            if !ok{
                return nil, errors.New("Invalid Json participant" + fmt.Sprint(p.Id))
            }
            return e, nil
        case "professor":
            item, ok := s.Professores[p.Id]
            e = item
            if !ok{
                return nil, errors.New("Invalid Json participant" + fmt.Sprint(p.Id))
            }
            return e, nil
        case "turma":
            for _, c := range s.Cursos {
                for _, et := range c.Etapas {
                    for _, t := range  et.Turmas {
                        if t.Id == p.Id && t.Nome == p.Nome {
                            e = t 
                            return e, nil
                        }
                    }
                }
            }
            return nil, errors.New("Invalid Json participant" + fmt.Sprint(p.Id))
        case "recurso":
            item, ok := s.Recursos[p.Id]
            e = item
            if !ok{
                return nil, errors.New("Invalid Json participant" + fmt.Sprint(p.Id))
            }
            return e, nil
        }
    return nil, errors.New("Unsuported participant type: " + p.Tipo)
}

func Json_contrato_to_contrato(j Json_contrato, s * Session) (Contrato, error) {
    c := Contrato{
        Id: j.Id, 
        Dispo: j.Dispo,
        Tipo_por_participante: j.Tipos_por_participante,
        Participantes: make([]ParticipanteContrato, len(j.Participantes)),
    }
    for i := range c.Participantes {
        p, err := Pjson_to_P(j.Participantes[i], s)
        if err != nil {
            return Contrato{}, err
        }
        c.Participantes[i] = p
    }
    return c, nil
}

func Contrato_to_json_contrato(c Contrato) Json_contrato {
    ps := make([]Participante_to_json, len(c.Participantes))
    for i := range ps {
        ps[i] = Participante_to_json{
            Id: c.Participantes[i].GetId(),
            Tipo: c.Tipo_por_participante[i],
            Nome: c.Participantes[i].GetNome(),
        }
    }
    

    jc := Json_contrato{
        Id: c.Id,
        Participantes: ps,
        Tipos_por_participante: c.Tipo_por_participante,
        Dispo: c.Dispo,
    }
    return jc
}

type Contrato struct {
    Id int
    Participantes []ParticipanteContrato
    Tipo_por_participante []string      // Tipo_por_participante[i] é o tipo do participante Participantes[i]
    Dispo       Disponibilidade
}


type Json_contrato struct {
    Id int
    Participantes []Participante_to_json
    Tipos_por_participante []string
    Dispo Disponibilidade
}

func (c Contrato) GetId() int{
    return c.Id
}


func (c * Contrato) UnmarshallJSON(b []byte) (error) {
    c.Id = 0
    c.Dispo = make(Disponibilidade, 0)
    c.Tipo_por_participante = []string{"call the custom ContratoUnmarshal function with the same []byte"}
    return nil
}

func (c * Contrato) CustomUnmarshallJSON(b []byte, s * Session) (error) {
    var jp Json_contrato
    err := json.Unmarshal(b, &jp)
    if err != nil {
        return err
    }

    ps := make([]ParticipanteContrato, len(jp.Participantes))
    for i := range ps {
        ps[i], err = Pjson_to_P(jp.Participantes[i], s)
        if err != nil {
            return err
        }
    }


    c.Id = jp.Id
    c.Dispo = jp.Dispo
    c.Tipo_por_participante = jp.Tipos_por_participante

    return nil
}

func (c Contrato) MarshalJSON() ([]byte, error) {
    ps := make([]Participante_to_json, len(c.Participantes))
    for i := range ps {
        ps[i] = Participante_to_json{
            Id: c.Participantes[i].GetId(),
            Tipo: c.Tipo_por_participante[i],
            Nome: c.Participantes[i].GetNome(),
        }
    }
    

    jc := Json_contrato{
        Id: c.Id,
        Participantes: ps,
        Tipos_por_participante: c.Tipo_por_participante,
        Dispo: c.Dispo,
    }
    
    return json.Marshal(jc)
}

func (c Contrato) HasParticipante(find ParticipanteContrato) (int, error) {
    for i, p := range c.Participantes {
        if p.GetId() == find.GetId() && p.GetNome() == find.GetNome() {
            return i, nil
        }
    }
    return -1, errors.New("Cant find item in Contrato")
}



