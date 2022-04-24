package model

type Person struct {
	Firstname string   `json:"firstname"`
	Lastname  string   `json:"lastname"`
	Friends   []Person `json:"friends"`
}
