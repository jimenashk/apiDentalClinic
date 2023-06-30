package domain

type Appointments struct {
	Id          int     `json:"id"`
	Patient     Patient `json:"patient,omitempty" xml:"patient,omitempty" `
	Dentist     Dentist `json:"dentist,omitempty" xml:"dentist,omitempty" `
	DateUp      string  `json:"date" binding:"required"`
	Hour        string  `json:"hour" binding:"required"`
	Description string  `json:"description"`
}
