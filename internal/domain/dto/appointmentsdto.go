package dto

type AppointmentInsert struct {
	PatientId   int    `json:"patientid" binding:"required"`
	DentistId   int    `json:"dentistid" binding:"required"`
	DateUp      string `json:"dateup" binding:"required"`
	Hour        string `json:"hour" binding:"required"`
	Description string `json:"description"`
}
