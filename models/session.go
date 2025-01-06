package models

import (
	"encoding/binary"
	"encoding/json"

	"github.com/google/uuid"
)

type SessionItem interface {
    GetId() int
}

type Session struct {
	Id          int
	Cursos      map[int]Curso
	Professores map[int]Professor
	Disciplinas map[int]Disciplina
	Contratos   map[int]Contrato
	Recursos    map[int]Recurso
}

func NewSession(id int) Session {
	return Session{
		Id:          id,
		Cursos:      make(map[int]Curso),
		Professores: make(map[int]Professor),
		Disciplinas: make(map[int]Disciplina),
		Contratos:   make(map[int]Contrato),
		Recursos:    make(map[int]Recurso),
	}
}

func AddCurso(s *Session, c Curso) int {
	// c.Dispo = NewDisponibilidade()

	c.Id = len(s.Cursos)
	s.Cursos[c.Id] = c
	return c.Id
}

func AddEtapa(s *Session, e Etapa) int {
	curso, ok := s.Cursos[e.Curso_id]
	if ok {
		return curso.AddEtapa(e)
	}
	return -1
}

func AddTurma(s *Session, t Turma) int {
	curso, _ := s.Cursos[t.Curso_id] // must exist
	return curso.AddTurma(t.Etapa_idx, t)

}

func AddProfessor(s *Session, p Professor) int {
	p.Dispo = NewDisponibilidade()

	// gera uuid, pega o final, converte pra binario 32, converte pra 32bit sinalizado
	p.Id = int(binary.BigEndian.Uint32([]byte(uuid.NewString())[:4]))
	s.Professores[p.Id] = p
	return p.Id
}

func AddRecurso(s *Session, r Recurso) int {
	r.Dispo = NewDisponibilidade()

	r.Id = int(binary.BigEndian.Uint32([]byte(uuid.NewString())[:4]))
	s.Recursos[r.Id] = r
	return r.Id
}

func AddDisciplina(s *Session, d Disciplina) int {
	d.Dispo = NewDisponibilidade()

	d.Id = int(binary.BigEndian.Uint32([]byte(uuid.NewString())[:4]))
	s.Disciplinas[d.Id] = d
	return d.Id
}

func (s *Session) AddContrato(c Contrato) int {
	c.Id = int(binary.BigEndian.Uint32([]byte(uuid.NewString())[:4]))
	s.Contratos[c.Id] = c
	return c.Id
}

func (s *Session) FindContratosWith(p ParticipanteContrato) map[int]*Contrato {
    r := make(map[int]*Contrato)
    for _, c := range s.Contratos {
        if idx, err := c.HasParticipante(p); err == nil {
            r[idx] = &c
        }
    }
    return r
}

func ParticipanteUpdater[T ParticipanteContrato](s *Session, old T, novo *T) {
    delete_mode := novo == nil
    for id, c := range s.Contratos {
        if idx_in, err := c.HasParticipante(old); err == nil {
            if delete_mode {
                delete(s.Contratos, id)
            } else {
                con := c
                con.Participantes[idx_in] = *novo
                con.Dispo = AndDisp(con.Participantes)
                s.Contratos[id] = con
            }
        }
    }
}

func CursoUpdater(s * Session, old Curso, novo *Curso) {
    delete_mode := novo == nil
    turmas := make([]Turma, 0)
    for _, e := range old.Etapas {
        for _, t := range e.Turmas {
            turmas = append(turmas, t)
        }
    }

    for _, t := range turmas {
        for id, c := range s.Contratos {
            if idx_in, err := c.HasParticipante(t); err == nil {
                if delete_mode {
                    delete(s.Contratos, id)
                } else {
                    con := c
                    con.Participantes[idx_in] = novo.Etapas[t.Etapa_idx].Turmas[t.Idx_in_etapa]
                    con.Dispo = AndDisp(con.Participantes)
                    s.Contratos[id] = con
                }
            }
        }
    }
}



func ProfessorGeter(s *Session) map[int]Professor {
	return s.Professores
}
func DisciplinaGeter(s *Session) map[int]Disciplina {
	return s.Disciplinas
}
func CursoGetter(s *Session) map[int]Curso {
	return s.Cursos
}
func TurmaGetter(s *Session) map[int]Turma {
    ts := make(map[int]Turma)
    for _, c := range s.Cursos {
        for _, e := range c.Etapas {
            for _, t := range e.Turmas {
                ts[t.Id] = t
            }
        }
    }
    return ts
}

func ContratoGetter(s *Session) map[int]Contrato {
	return s.Contratos
}
func RecursoGetter(s *Session) map[int]Recurso {
	return s.Recursos
}

func map_to_slice[T SessionItem, A comparable](m map[A]T) []T {
	slice := make([]T, 0, len(m))
	for _, v := range m {
		slice = append(slice, v)
	}
	return slice
}

func slice_to_map[T SessionItem](s []T) map[int]T {
	m := make(map[int]T)
	for _, v := range s {
		m[v.GetId()] = v
	}
	return m
}

type json_session struct {
	Id          int
	Cursos      []Curso
	Professores []Professor
	Disciplinas []Disciplina
	Contratos   []Json_contrato
	Recursos    []Recurso
}

func (s Session) MarshalJSON() ([]byte, error) {
    contratos := map_to_slice(s.Contratos)
    jc := make([]Json_contrato, len(contratos))
    for i := range jc {
        jc[i] = Contrato_to_json_contrato(contratos[i])
    }

	a := json_session{
		s.Id,
		map_to_slice(s.Cursos),
		map_to_slice(s.Professores),
		map_to_slice(s.Disciplinas),
		jc,
		map_to_slice(s.Recursos),
	}
	return json.Marshal(a)
}

func (s *Session) UnmarshalJSON(bs []byte) error {
	var a json_session
	err := json.Unmarshal(bs, &a)
	if err != nil {
		return err
	}
    
	*s = Session{
		a.Id,
		slice_to_map(a.Cursos),
		slice_to_map(a.Professores),
		slice_to_map(a.Disciplinas),
		make(map[int]Contrato, len(a.Contratos)),
		slice_to_map(a.Recursos),
	}
    
    for _, c := range a.Contratos {
        c2, err := Json_contrato_to_contrato(c, s)
        if err != nil {
            return err
        }
        s.Contratos[c2.Id] = c2
    }

	return nil
}
