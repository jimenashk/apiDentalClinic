package dto

type AppointmentGet struct {
	DNIPatient  string `json:"dni" binding:"required"`
	PatientName string `json:"patient" binding:"required"`
	DentistName string `json:"dentist" binding:"required"`
	Date        string `json:"date" binding:"required"`
	Hour        string `json:"hour" binding:"required"`
	Description string `json:"description"`
}
