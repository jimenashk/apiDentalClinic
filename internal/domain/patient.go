package domain

type Patient struct {
	Id       int    `json:"id"`
	Name     string `json:"name,omitempty"`
	LastName string `json:"lastname,omitempty"`
	Address  string `json:"address,omitempty"`
	DNI      string `json:"dni,omitempty" binding:"required"`
	DateUp   string `json:"dateup,omitempty" binding:"required"`
}
