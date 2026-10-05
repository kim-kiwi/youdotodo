package main

import (
	"net/http"
	"strconv"
	"youdotodo/pages"
	. "youdotodo/repo"
)

func indexHandler(sm *SessionManager, tc *TodoCollection, uc *UserCollection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sid, err := r.Cookie("sid")
		if err != nil {
			pages.Index(w, r, false, 0, nil)
			return
		}
		uid, err := sm.Get(sid.Value)
		if err != nil {
			pages.Index(w, r, false, 0, nil)
			return
		}
		_, err = uc.Get(uid.(int64))
		if err != nil {
			pages.Index(w, r, false, 0, nil)
			return
		}

		query := r.URL.Query()

		if !query.Has("year") || !query.Has("month") || !query.Has("day") {
			http.ServeFile(w, r, "pages/gotohome.html")
			return
		}

		parseInt := func(key string) int {
			if err != nil {
				return 0
			}
			var value int
			value, err = strconv.Atoi(query.Get(key))
			return value
		}
		year := int16(parseInt("year"))
		month := int8(parseInt("month"))
		day := int8(parseInt("day"))
		if err != nil {
			http.Error(w, "invalid query", http.StatusInternalServerError)
			return
		}

		todos, err := tc.ListByDate(year, month, day)
		if err != nil {
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}

		todoDtos := make([]TodoDto, len(todos))
		for i, v := range todos {
			t := v.ToT()
			user, err := uc.Get(v.UserId)
			if err == nil {
				t.UserName = user.Name
			}
			todoDtos[i] = t
		}

		todosByUser := make(map[int64][]TodoDto)
		for _, t := range todoDtos {
			todosByUser[t.UserId] = append(todosByUser[t.UserId], t)
		}

		pages.Index(w, r, true, uid.(int64), todosByUser)
	}
}
