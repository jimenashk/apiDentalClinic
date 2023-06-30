package domain

type Dentist struct {
	Id       int    `json:"id"`
	Name     string `json:"name,omitempty"`
	LastName string `json:"lastname,omitempty"`
	License  string `json:"license,omitempty"`
}
