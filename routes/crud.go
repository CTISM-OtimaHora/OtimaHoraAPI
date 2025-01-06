package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	. "github.com/CTISM-OtimaHora/OtimaHora/models"
)

func GetBuilder[T SessionItem] (map_geter func(*Session) map[int]T) func(http.ResponseWriter, * http.Request) {
    return func (w http.ResponseWriter, r * http.Request) {
        s := Session_or_nil(r)
        if s == nil {
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte("No session or Session expired"))
            return 
        }

        id, conv_err := strconv.Atoi(r.PathValue("id"))
        if conv_err != nil {
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        if json.NewEncoder(w).Encode(map_geter(s)[id]) != nil {
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        return   
    }
}

func AddBuilder[T SessionItem] (slice_adder func(sess *Session, en T) int) func (http.ResponseWriter, *http.Request) {
    return func (w http.ResponseWriter, r * http.Request) {
        s := Session_or_nil(r)
        if s == nil {
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte("No session or Session expired"))
            return 
        }

        var e T 

        if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte("malformed body 1: " + err.Error()))
            return
        }
        fmt.Printf("Addded  %v\n", e)

        w.Write([]byte(fmt.Sprint(slice_adder(s, e))))
    }
}

func SetBuilder[T SessionItem] (map_geter func(*Session) map[int]T, contrato_updater func(*Session, T, *T)) func(http.ResponseWriter, * http.Request) {
    return func (w http.ResponseWriter, r * http.Request) {
        s := Session_or_nil(r)
        if s == nil {
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte("No session or Session expired"))
            return 
        }

        id, conv_err := strconv.Atoi(r.PathValue("id"))
        if conv_err != nil {
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        old, ok := map_geter(s)[id]
        if !ok {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte("Item not found"))
            return
        }
        
        var e T 

        if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte("malformed body 1: " + err.Error()))
            return
        }

        delete (map_geter(s), id)
        map_geter(s)[id] = e

        contrato_updater(s, old, &e)
        return   
    }
}

func DeleteBuilder[T SessionItem] (map_geter func(*Session) map[int]T, contrato_updater func(*Session, T, *T)) func(http.ResponseWriter, * http.Request) {
    return func (w http.ResponseWriter, r * http.Request) {
        s := Session_or_nil(r)
        if s == nil {
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte("No session or Session expired"))
            return 
        }

        id, conv_err := strconv.Atoi(r.PathValue("id"))
        if conv_err != nil {
            w.WriteHeader(http.StatusBadRequest)
            return
        }
        
        // tipo := strings.Split(strings.TrimSpace(r.RequestURI), "/")[1] // 1 e não 0 pq a URI é uma porra
        // ent := SearchSessionItem{Id: id, Tipo: tipo}.GetEntidadeOrNil(s)
        // if ent == nil {
        //     w.WriteHeader(http.StatusBadRequest)
        //     return
        // }
        
        old, ok := map_geter(s)[id]
        if !ok {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte("Item not found"))
            return
        }

        delete(map_geter(s), id)
        contrato_updater(s, old, nil)
    
        return   
    }
}

func SliceGetBuilder[T SessionItem] (map_geter func(*Session) map[int]T) func(http.ResponseWriter, * http.Request) {
    return func (w http.ResponseWriter, r * http.Request) {
        s := Session_or_nil(r)
        if s == nil {
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte("No session or Session expired"))
            return 
        }

        m := map_geter(s)
        slice := make([]T, 0, len(m))
        for _, v := range m {
            slice = append(slice, v)
        }

        if json.NewEncoder(w).Encode(slice) != nil {
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        return   
    }
}
