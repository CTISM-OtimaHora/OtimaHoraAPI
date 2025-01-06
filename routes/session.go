package routes

import (
    "net/http"
    "time"
    "fmt"
    "encoding/json"

    ."github.com/CTISM-OtimaHora/OtimaHora/models"
)

func AddSession(w http.ResponseWriter, r * http.Request) {
    s := Session_or_nil(r)
    if s != nil {
        http.SetCookie(w, &http.Cookie{Name: "id", Value: fmt.Sprint(s.Id), Path: "/", Expires: time.Now().Add(2*time.Hour)})
        return
    }

    new_s := CreateSession() 
    http.SetCookie(w, &http.Cookie{Name: "id", Value: fmt.Sprint(new_s.Id), Path: "/", Expires: time.Now().Add(2*time.Hour)})
    return
}

func AddSessionFromDocument(w http.ResponseWriter, r * http.Request) {
    server_session := Session_or_nil(r)
    if server_session == nil {
        server_session = CreateSession()
    }

    file, _, err := r.FormFile("document")
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Bad Request: " + err.Error()))
        return
    }
    
    var file_session Session
    if err := json.NewDecoder(file).Decode(&file_session); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Bad request: " + err.Error()))
        println(err.Error())
        return
    }

    // just ignore the ID
    id_saver := server_session.Id
    (*server_session) = file_session
    server_session.Id =  id_saver
}

func GetSession(w http.ResponseWriter, r * http.Request) {
    s := Session_or_nil(r)
    if s == nil {
        w.WriteHeader(http.StatusUnauthorized)
        w.Write([]byte("No session or Session expired"))
        return 
    }

    json.NewEncoder(w).Encode(s)
    return
}
 
