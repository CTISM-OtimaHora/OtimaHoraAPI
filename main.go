package main

import (
	"net/http"

	m "github.com/CTISM-OtimaHora/OtimaHora/models"
	. "github.com/CTISM-OtimaHora/OtimaHora/routes"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Origin", origin)

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// If it's an OPTIONS request, just return OK status
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	r := http.NewServeMux()

	r.HandleFunc("GET /add-session", AddSession)
	r.HandleFunc("POST /add-session-document", AddSessionFromDocument)
	r.HandleFunc("GET /session", GetSession)

	// professor routes
	r.HandleFunc("POST /professor/add", AddBuilder(m.AddProfessor))
	r.HandleFunc("DELETE /professor/delete/{id}", DeleteBuilder(m.ProfessorGeter, m.ParticipanteUpdater))
	r.HandleFunc("GET /professor/get/{id}", GetBuilder(m.ProfessorGeter))
	r.HandleFunc("PUT /professor/set/{id}", SetBuilder(m.ProfessorGeter, m.ParticipanteUpdater))
	r.HandleFunc("GET /professor/slice", SliceGetBuilder(m.ProfessorGeter))

	// disciplina routes
	r.HandleFunc("POST /disciplina/add", AddBuilder(m.AddDisciplina))
	r.HandleFunc("DELETE /disciplina/delete/{id}", DeleteBuilder(m.DisciplinaGeter, m.ParticipanteUpdater))
	r.HandleFunc("GET /disciplina/get/{id}", GetBuilder(m.DisciplinaGeter))
	r.HandleFunc("PUT /disciplina/set/{id}", SetBuilder(m.DisciplinaGeter, m.ParticipanteUpdater))
	r.HandleFunc("GET /disciplina/slice", SliceGetBuilder(m.DisciplinaGeter))

	// recurso routes
	r.HandleFunc("POST /recurso/add", AddBuilder(m.AddRecurso))
	r.HandleFunc("DELETE /recurso/delete/{id}", DeleteBuilder(m.RecursoGetter, m.ParticipanteUpdater))
	r.HandleFunc("GET /recurso/get/{id}", GetBuilder(m.RecursoGetter))
	r.HandleFunc("PUT /recurso/set/{id}", SetBuilder(m.RecursoGetter, m.ParticipanteUpdater))
	r.HandleFunc("GET /recurso/slice", SliceGetBuilder(m.RecursoGetter))

	// curso routes
	r.HandleFunc("POST /curso/add", AddBuilder(m.AddCurso))
	r.HandleFunc("DELETE /curso/delete/{id}", DeleteBuilder(m.CursoGetter, m.CursoUpdater))
	r.HandleFunc("GET /curso/get/{id}", GetBuilder(m.CursoGetter))
	r.HandleFunc("PUT /curso/set/{id}", SetBuilder(m.CursoGetter, m.CursoUpdater))
	r.HandleFunc("GET /curso/slice", SliceGetBuilder(m.CursoGetter))

	//  Turma routes
	r.HandleFunc("POST /turma/add", AddBuilder(m.AddTurma))

	// contrato routes
	r.HandleFunc("POST /contrato/add", AddContrato)
	r.HandleFunc("GET /contrato/get/{id}", GetContrato)
	// TODO
	// r.HandleFunc("PUT /contrato/set/{id}", SetContrato)

	with_cors := corsMiddleware(r)
	http.ListenAndServe("localhost:3000", with_cors)
}
