package models

type Question struct {
	Id    string `json:"id"`
	Text  string `json:"text"`
	Votes int
}
