package models

import (
    "net/http"
    "strconv"
)



var sessions =  make(map[int]*Session)

var next_id int = 0

func Session_or_nil(r * http.Request) *Session {
    session_cookie, err := r.Cookie("id")
    if err != nil {
        return nil
    }

    idx, _ := strconv.Atoi(session_cookie.Value)

    s, ok := sessions[idx]

    if ok {
        return s
    }
    return nil
}

func CreateSession () *Session {
    s := NewSession(next_id)
    sessions[next_id] = &s 
    next_id += 1
    return &s
}

func DeleteSession(id int) {
    delete(sessions, id)
}

