package tests

import (
	"fmt"
	"slices"
	"testing"

	m "github.com/CTISM-OtimaHora/OtimaHora/models"
)

func Test_And_Disponibilidade(t * testing.T) {
    disp1 := m.NewDisponibilidade()

    // sets first column of disp1 to NAO
    for i := range len(disp1) {
        disp1[i][0] = m.NO
    }

    disp2 := m.NewDisponibilidade()
    // sets last column of disp2 to NAO
    for i := range len(disp2) {
        disp1[i][len(disp2[i])-1] = m.NO
    }

    ent1 := m.NewProfessor(0, "a")
    ent2 := m.NewProfessor(1, "a")
    ent1.Dispo = disp1
    ent2.Dispo = disp2

    and := m.AndDisp([]m.Entidade{ent1, ent2})
    
    check := m.NewDisponibilidade()
    for i := range check {
        check[i][0] = m.NO
        check[i][len(check[i])-1] = m.NO
    } 

    if slices.CompareFunc(check, and, slices.Compare) != 0 {
        t.Error("Fail to and two disponibilidades")
    }
    fmt.Println("passed Test_And_Disponibilidade")
}
