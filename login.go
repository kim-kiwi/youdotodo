package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	. "youdotodo/repo"
)

type LoginReq struct {
	Name string `json:"name"`
}
type LoginRes struct {
	Ok  bool   `json:"ok"`
	Msg string `json:"msg"`
}

func loginHandler(sm *SessionManager, uc *UserCollection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginReq
		json.NewDecoder(r.Body).Decode(&req)

		res := LoginRes{Ok: true, Msg: "successfully logged in."}
		if uid, err := uc.FindByName(req.Name); err != nil {
			res.Ok = false
			res.Msg = "user not found"
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
		} else {
			sid := sm.New()
			cookie := http.Cookie{
				Name:     "sid",
				Value:    sid,
				Path:     "/",
				MaxAge:   24 * 3600,
				HttpOnly: true,
				Secure:   false,
			}
			sm.Set(sid, uid)
			http.SetCookie(w, &cookie)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
		}

		json.NewEncoder(w).Encode(res)
	}
}
func logoutHandler(sm *SessionManager, uc *UserCollection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if sid, err := r.Cookie("sid"); err != nil {
			http.Error(w, "you are not logged in :(", http.StatusUnauthorized)
		} else {
			if uid, err := sm.Get(sid.Value); err != nil {
				http.Error(w, "you are not logged in :(", http.StatusUnauthorized)
			} else {
				_, err := uc.Get(uid.(int64))
				if err != nil {
					http.Error(w, "you are not logged in :(", http.StatusUnauthorized)
				} else {
					cookie := http.Cookie{
						Name:     "sid",
						Value:    "",
						Path:     "/",
						HttpOnly: true,
						Secure:   false,
					}
					sm.Delete(sid.Value)
					http.SetCookie(w, &cookie)
					w.WriteHeader(http.StatusOK)
					fmt.Fprintf(w, "successfully logged out")
				}
			}
		}
	}
}
