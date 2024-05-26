package models

type Filter struct {
}

type FilterItem struct {
	Field string `json:"field"`
	Value string `json:"value"`
}
