package repo

import (
	"time"
	"errors"
	"sync"
	"youdotodo/ghost/repo"
)

type User struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	LastUpload time.Time `json:"lastUpload"`
}

func (u User) ToD() User {
	return u
}
func (u User) ToT() User {
	return u
}

type UserCollection struct {
	*repo.OnJson[User]
	mu sync.Mutex
}

func NewUserCollection() (*UserCollection, error) {
	onjson, err := repo.NewOnJson[User]("users.json", func(v *User, id int64) {
		v.Id = id
	})
	if err != nil {
		return nil, err
	}
	return &UserCollection{OnJson: onjson}, nil
}

func (uc *UserCollection) Create(u User) (int64, error) {
	if len(u.Name) > 16 {
		return 0, errors.New("name too long")
	}
	uc.mu.Lock()
	defer uc.mu.Unlock()
	if _, err := uc.FindByName(u.Name); err == nil {
		return 0, errors.New("name already exists")
	}
	return uc.OnJson.Create(u)
}

func (uc *UserCollection) List() ([]User, error) {
	users, err := uc.OnJson.List()
	if err != nil {
		return err
	}
	sort.Slice(users, func(i, j int) bool {
		return todos[i].Id < todos[j].Id
	})
	return users, nil
}

func (uc *UserCollection) FindByName(name string) (int64, error) {
	user, err := uc.Find(func(u User) bool {
		return u.Name == name
	})
	if err != nil {
		return 0, errors.New("User not found")
	}
	return user.Id, nil
}
