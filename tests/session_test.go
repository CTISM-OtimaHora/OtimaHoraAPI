package tests

import (
	"encoding/json"
	"fmt"
	"slices"
	"testing"

	m "github.com/CTISM-OtimaHora/OtimaHora/models"
)

func Test_session_custom_json(t * testing.T) {
    s := m.NewSession(0)
    p := m.NewProfessor(0, "artur")
    d := m.NewDisciplina(0, "matemática")
    c := m.NewCurso(0, "CTISM")
    r := m.NewRecurso(0, "lab-matemática")

    m.AddProfessor(&s, p)
    m.AddDisciplina(&s, d)
    m.AddCurso(&s, c)
    m.AddRecurso(&s, r)

    s.AddContrato(m.NewContrato(0, []m.Entidade{p, d, c, r}))

    bs, err := json.Marshal(s)
    
    if err != nil {
        t.Error(err.Error())
    }

    var from m.Session
    err = json.Unmarshal(bs, &from)

    if err != nil {
        t.Error(err.Error())
    }
    fmt.Println("passesd TesteTest_session_custom_json")   
}

func Test_session_contratos_update(t * testing.T) {
    s := m.NewSession(0)
    p := m.NewProfessor(0, "artur")
    d := m.NewDisciplina(0, "matemática")
    c := m.NewCurso(0, "CTISM")
    r := m.NewRecurso(0, "lab-matemática")

    p.Id =  m.AddProfessor(&s, p)
    d.Id =  m.AddDisciplina(&s, d)
    c.Id =  m.AddCurso(&s, c)
    r.Id =  m.AddRecurso(&s, r)

    c.Curriculo[d.Id] = m.HorasAula{Horas: 10, Formato: "5+5"}

    _ = s.AddContrato(m.NewContrato(0, []m.Entidade{p, d, c, r}))
    _ = s.AddContrato(m.NewContrato(1, []m.Entidade{p, d}))
    _ = s.AddContrato(m.NewContrato(2, []m.Entidade{p, d, c}))
    id4 := s.AddContrato(m.NewContrato(3, []m.Entidade{p, r}))


    delete(s.Professores, d.Id)
    s.UpdateSessionFromDelete(d)

    cons := make([]int, 0, len(s.Contratos))
    for _, c := range s.Contratos {
        cons = append(cons, c.Id)
    }

    if slices.Compare(cons, []int{id4}) != 0  {
        t.Errorf("expected all contratos with the removed element to be deleted, no true have %v", cons)
    }

    if len(c.Curriculo) != 0 {
        t.Errorf("Removing, Discilplina should clear it from Curriculos, didnt: %v", c.Curriculo)
    }

    p.Nome = "lol artur"
    s.UpdateContratos(p)

    for _, contrato := range s.Contratos {
        for _, ent := range  contrato.Participantes {
            if ent.Tipo == "professor" && ent.Nome != "lol artur" {
                t.Errorf("expect change to entidade to modify all contratos with it, not true")
            }
        }
    }

    fmt.Println("passed Test_session_contratos_update")
}
