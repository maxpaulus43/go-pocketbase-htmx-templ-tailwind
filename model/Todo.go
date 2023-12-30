package model

type Todo struct {
	Id         string `db:"id"`
	Content    string `db:"content"`
	IsComplete bool   `db:"is_complete"`
}
