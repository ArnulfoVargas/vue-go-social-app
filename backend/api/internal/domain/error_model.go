package domain

type Error struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
}
