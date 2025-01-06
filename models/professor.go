package models

type Professor struct {
    Id int
    Nome string
    Dispo Disponibilidade
    Disciplinas_ids []int
}


func NewProfessor(id int, nome string) Professor {
    return Professor {Id: id, Nome: nome, Dispo: NewDisponibilidade(), Disciplinas_ids: make([]int, 0)}
}

func (p Professor) GetId() int {
    return p.Id
}
func (p  Professor) GetNome() string {
    return p.Nome
}
func (p  Professor) GetDispo() Disponibilidade {
    return p.Dispo
}
