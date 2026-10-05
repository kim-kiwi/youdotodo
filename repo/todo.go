package repo

import (
	"errors"
	"sort"
	"youdotodo/ghost/repo"
)

type TodoDto struct {
	Id       int64  `json:"id"`
	Idx      int64  `json:"idx"`
	UserId   int64  `json:"userId"`
	UserName string `json:"userName"`
	Year     int16  `json:"year"`
	Month    int8   `json:"month"`
	Day      int8   `json:"day"`
	Content  string `json:"content"`
	Done     bool   `json:"done"`
}

func (td TodoDto) ToD() Todo {
	return Todo{
		Id:      td.Id,
		Idx:     td.Idx,
		UserId:  td.UserId,
		Year:    td.Year,
		Month:   td.Month,
		Day:     td.Day,
		Content: td.Content,
		Done:    td.Done,
	}
}

func (t Todo) Validate() error {
    if t.Content == "" {
        return errors.New("content is empty")
    }
    if len(t.Content) >= 64 {
        return errors.New("content too long")
    }
	if t.Month > 12 || t.Month < 1 {
		return errors.New("invalid month")
	}
	max_days := []int8{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	if (t.Year%4 == 0 && t.Year%100 != 0) || (t.Year%400 == 0) {
		max_days[1] = 29
	}
	if t.Day < 1 || t.Day > max_days[t.Month-1] {
		return errors.New("invalid day")
	}
    return nil
}

type Todo struct {
	Id      int64  `json:"id"`
	Idx     int64  `json:"idx"`
	UserId  int64  `json:"userId"`
	Year    int16  `json:"year"`
	Month   int8   `json:"month"`
	Day     int8   `json:"day"`
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

func (t Todo) ToT() TodoDto {
	return TodoDto{
		Id:      t.Id,
		Idx:     t.Idx,
		UserId:  t.UserId,
		Year:    t.Year,
		Month:   t.Month,
		Day:     t.Day,
		Content: t.Content,
		Done:    t.Done,
	}
}

type TodoCollection struct {
	*repo.OnJson[Todo]
}

func NewTodoCollection() (*TodoCollection, error) {
	onjson, err := repo.NewOnJson[Todo]("todos.json", func(v *Todo, id int64) {
		v.Id = id
	})
	if err != nil {
		return nil, err
	}
	return &TodoCollection{OnJson: onjson}, nil
}

func (tc *TodoCollection) Create(t Todo) (int64, error) {
	if err := t.Validate(); err != nil {
		return 0, err
	}
	return tc.OnJson.Create(t)
}

func (tc *TodoCollection) ListByDate(year int16, month int8, day int8) ([]Todo, error) {
	todos, err := tc.FindMany(func(t Todo) bool {
		return t.Year == year && t.Month == month && t.Day == day
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(todos, func(i, j int) bool {
		return todos[i].Idx < todos[j].Idx
	})
	return todos, nil
}
