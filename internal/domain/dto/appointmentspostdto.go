package dto

type AppointmentPost struct {
	DentistLicense string `json:"license"`
	PatientDni     string `json:"dni"`
	DateUp         string `json:"dateup"`
	Hour           string `json:"hour"`
	Description    string `json:"description"`
}
