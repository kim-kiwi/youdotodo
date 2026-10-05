package ghost

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Endpoints uint8

const (
    EndpointCreate Endpoints = 1 << iota
    EndpointGet
    EndpointUpdate
    EndpointPatch
    EndpointDelete
    EndpointList
)

const EndpointAll = EndpointCreate | EndpointGet | EndpointUpdate | EndpointPatch | EndpointDelete | EndpointList
const EndpointCGDL = EndpointCreate | EndpointGet | EndpointDelete | EndpointList
const EndpointCGL = EndpointCreate | EndpointGet | EndpointList

type Repo[T any] interface {
	Create(T) (int64, error)
	Get(int64) (T, error)
	Update(int64, T) error
	Delete(int64) error
	List() ([]T, error)
	Find(func(T) bool) (T, error)
	FindMany(func(T) bool) ([]T, error)
}

type CreateRes struct {
	Id int64 `json:"id"`
}

type Api[T interface{ ToD() D }, D interface{ ToT() T }] struct {
	path       string
	repo       Repo[D]
	createFrom func(T, *http.Request) (D, error)
	updateFrom func(int64, T, *http.Request) error
	patchFrom  func(int64, D, *http.Request) (D, error)
	listFrom   func(*http.Request) ([]T, error)
	deleteFrom func(int64, *http.Request) error
	middleware func(handler http.HandlerFunc) http.HandlerFunc
}

func NewApi[T interface{ ToD() D }, D interface{ ToT() T }](path string, repo Repo[D]) *Api[T, D] {
	return &Api[T, D]{
		path: path,
		repo: repo,
	}
}

func (a *Api[T, D]) CreateFrom(fn func(T, *http.Request) (D, error)) {
	a.createFrom = fn
}

func (a *Api[T, D]) ListFrom(fn func(*http.Request) ([]T, error)) {
	a.listFrom = fn
}

func (a *Api[T, D]) UpdateFrom(fn func(int64, T, *http.Request) error) {
	a.updateFrom = fn
}

func (a *Api[T, D]) PatchFrom(fn func(int64, D, *http.Request) (D, error)) {
	a.patchFrom = fn
}

func (a *Api[T, D]) DeleteFrom(fn func(int64, *http.Request) error) {
	a.deleteFrom = fn
}

func (a *Api[T, D]) SetMiddleware(fn func(handler http.HandlerFunc) http.HandlerFunc) {
	a.middleware = fn
}

func AddApiToMux[T interface{ ToD() D }, D interface{ ToT() T }](mux *http.ServeMux, endpoints Endpoints, api *Api[T, D]) {
	var createHandler http.HandlerFunc
	if (endpoints & EndpointCreate) != 0 {
		createHandler = func(w http.ResponseWriter, r *http.Request) {
			var req T
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}

			var data D
			if api.createFrom != nil {
				data, err = api.createFrom(req, r)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			} else {
				data = req.ToD()
			}

			id, err := api.repo.Create(data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			res := CreateRes{id}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(res)
		}
	}
	var getHandler http.HandlerFunc
	if (endpoints & EndpointGet) != 0 {
		getHandler = func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
			if err != nil {
				http.Error(w, "Invalid ID", http.StatusBadRequest)
				return
			}
			f, err := api.repo.Get(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(f.ToT())
		}
	}
	var updateHandler http.HandlerFunc
	if (endpoints & EndpointUpdate) != 0 {
		updateHandler = func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
			if err != nil {
				http.Error(w, "Invalid ID", http.StatusBadRequest)
				return
			}

			var req T
			err = json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}

			if api.updateFrom != nil {
				err = api.updateFrom(id, req, r)
			} else {
				err = api.repo.Update(id, req.ToD())
			}
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
		}
	}
	var patchHandler http.HandlerFunc
	if (endpoints & EndpointPatch) != 0 && api.patchFrom != nil {
		patchHandler = func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
			if err != nil {
				http.Error(w, "Invalid ID", http.StatusBadRequest)
				return
			}

			old, err := api.repo.Get(id)
			if err != nil {
				http.Error(w, "Invalid ID", http.StatusBadRequest)
				return
			}

			data, err := api.patchFrom(id, old, r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = api.repo.Update(id, data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
		}
	}
	var deleteHandler http.HandlerFunc
	if (endpoints & EndpointDelete) != 0 {
		deleteHandler = func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
			if err != nil {
				http.Error(w, "Invalid ID", http.StatusBadRequest)
				return
			}
			if api.deleteFrom != nil {
				err = api.deleteFrom(id, r)
			} else {
				err = api.repo.Delete(id)
			}
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
		}
	}
	var listHandler http.HandlerFunc
	if (endpoints & EndpointList) != 0 {
		listHandler = func(w http.ResponseWriter, r *http.Request) {
			var l []T
			var err error
			if api.listFrom != nil {
				l, err = api.listFrom(r)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				dl, err := api.repo.List()
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				l = make([]T, len(dl))
				for i, v := range dl {
					l[i] = v.ToT()
				}
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(l)
		}
	}
	addMiddleware := func(handler http.HandlerFunc) http.HandlerFunc {
		if handler != nil {
			return api.middleware(handler)
		}
		return nil
	}
	if api.middleware != nil {
		createHandler = addMiddleware(createHandler)
		getHandler = addMiddleware(getHandler)
		updateHandler = addMiddleware(updateHandler)
		patchHandler = addMiddleware(patchHandler)
		deleteHandler = addMiddleware(deleteHandler)
		listHandler = addMiddleware(listHandler)
	}
	addToMux := func(path string, handler http.HandlerFunc) {
		if handler != nil {
			mux.HandleFunc(path, handler)
		}
	}
	addToMux(fmt.Sprintf("POST /%s", api.path), createHandler)
	addToMux(fmt.Sprintf("GET /%s/{id}", api.path), getHandler)
	addToMux(fmt.Sprintf("PUT /%s/{id}", api.path), updateHandler)
	addToMux(fmt.Sprintf("PATCH /%s/{id}", api.path), patchHandler)
	addToMux(fmt.Sprintf("DELETE /%s/{id}", api.path), deleteHandler)
	addToMux(fmt.Sprintf("GET /%s", api.path), listHandler)
}
