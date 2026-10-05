package main

import (
	"encoding/json"
	"context"
	"errors"
	"net/http"
	"youdotodo/ghost"
	. "youdotodo/repo"
)

func NewTodoApi(tc *TodoCollection, sm *SessionManager, uc *UserCollection) *ghost.Api[TodoDto, Todo] {
	todoApi := ghost.NewApi("todos", tc)
	todoApi.SetMiddleware(func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			sid, err := r.Cookie("sid")
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			uid, err := sm.Get(sid.Value)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			user, err := uc.Get(uid.(int64))
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), ctxUserKey, user)
			handler(w, r.WithContext(ctx))
		}
	})
	todoApi.UpdateFrom(func(id int64, new TodoDto, r *http.Request) error {
		user, ok := r.Context().Value(ctxUserKey).(User)
		if !ok {
			return errors.New("user id missing")
		}
		old, err := tc.Get(id)
		if err != nil {
			return err
		}
		if old.UserId != user.Id {
			return errors.New("permission denied")
		}
		todo := new.ToD()
		if err := todo.Validate(); err != nil {
			return err
		}
		todo.UserId = old.UserId
		return tc.Update(id, todo)
	})
	todoApi.PatchFrom(func(id int64, todo Todo, r *http.Request) (Todo, error) {
		user, ok := r.Context().Value(ctxUserKey).(User)
		if !ok {
			return Todo{}, errors.New("user id missing")
		}
		if todo.UserId != user.Id {
			return Todo{}, errors.New("permission denied")
		}

		var patch struct {
			Content *string `json:"content"`
			Idx     *int64  `json:"idx"`
			Done    *bool   `json:"done"`
		}
		err := json.NewDecoder(r.Body).Decode(&patch)
		if err != nil {
			return Todo{}, errors.New("Invalid JSON")
		}

		if patch.Content != nil {
			todo.Content = *patch.Content
		}
		if patch.Idx != nil {
			todo.Idx = *patch.Idx
		}
		if patch.Done != nil {
			todo.Done = *patch.Done
		}

		if err := todo.Validate(); err != nil {
			return Todo{}, err
		}

		return todo, nil
	})
	todoApi.DeleteFrom(func(id int64, r *http.Request) error {
		user, ok := r.Context().Value(ctxUserKey).(User)
		if !ok {
			return errors.New("user id missing")
		}
		todo, err := tc.Get(id)
		if err != nil {
			return err
		}
		if todo.UserId != user.Id {
			return errors.New("permission denied")
		}
		return tc.Delete(id)
	})
	todoApi.CreateFrom(func(t TodoDto, r *http.Request) (Todo, error) {
		user, ok := r.Context().Value(ctxUserKey).(User)
		if !ok {
			return Todo{}, errors.New("login first")
		}
		todo := t.ToD()
		todo.UserId = user.Id
		if err := todo.Validate(); err != nil {
			return Todo{}, err
		}
		return todo, nil
	})
	return todoApi
}
