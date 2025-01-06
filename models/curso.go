package models

import (
	"encoding/binary"
	"github.com/google/uuid"
)

type HorasAula struct {
	Horas   int
	Formato string // exemplo 2+2, 4 1+3
}

type Curso struct {
	Id     int
	Nome   string
	Etapas []Etapa
}

type Etapa struct {
	Idx_in_Curso int
	Curso_id     int
	Curriculo    map[int]HorasAula // ID para Horas aula
	Turmas       []Turma
}

func (c Curso) GetId() int {
	return c.Id
}

func NewEtapa() Etapa {
	return Etapa{
		Curso_id:     0,
		Idx_in_Curso: 0,
		Turmas:       make([]Turma, 0),
		Curriculo:    map[int]HorasAula{},
	}
}


func NewCurso(id int, nome string) Curso {
	return Curso{
		Id:     id,
		Nome:   nome,
		Etapas: make([]Etapa, 0),
	}
}

func (c *Curso) AddEtapa(e Etapa) int {
	e.Idx_in_Curso = len(c.Etapas)
	e.Curso_id = c.Id

	c.Etapas = append(c.Etapas, e)
	return e.Idx_in_Curso
}

func (c *Curso) AddTurma(etapa int, t Turma) int {
	for etapa >= len(c.Etapas) {
		c.AddEtapa(NewEtapa())
	}

	t.Id = int(binary.BigEndian.Uint32([]byte(uuid.NewString())[:4]))
	t.Curso_id = c.Id
	t.Etapa_idx = etapa
    t.Idx_in_etapa = len(c.Etapas[etapa].Turmas)
	c.Etapas[etapa].Turmas = append(c.Etapas[etapa].Turmas, t)
	return t.Id
}
