package handler

import (
	"errors"
	"os"
	"strconv"

	"apiDentalClinic/internal/appointments"
	"apiDentalClinic/internal/domain/dto"
	"apiDentalClinic/pkg/web"

	"github.com/gin-gonic/gin"
)

type appointmenthandler struct {
	s appointments.Service
}

func NewAppointmentHandler(s appointments.Service) *appointmenthandler {
	return &appointmenthandler{s: s}
}

func (h *appointmenthandler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		appointments, _ := h.s.ReadAll()
		if len(appointments) == 0 {
			web.Failure(c, 400, errors.New("There is no Appointment"))
		}
		web.Success(c, 200, appointments)
	}
}

func (h *appointmenthandler) Post() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.GetHeader("Token")
		if token == "" {
			web.Failure(c, 401, errors.New("Token Not Found"))
			return
		}

		if token != os.Getenv("TOKEN") {
			web.Failure(c, 401, errors.New("Invalid Token"))
			return
		}

		var appointment dto.AppointmentInsert
		err := c.ShouldBindJSON(&appointment)
		if err != nil {
			web.Failure(c, 400, errors.New("Invalid Appointment"))
			return
		}

		valid, err := validateEmptysAppointment(&appointment)
		if !valid {
			web.Failure(c, 400, errors.New(err.Error()))
			return
		}

		err = h.s.Create(appointment)
		if err != nil {
			web.Failure(c, 400, errors.New(err.Error()))
			return
		}

		web.Success(c, 201, appointment)
	}
}

func (h *appointmenthandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, 404, errors.New("Invalid id"))
			return
		}
		appointment, err := h.s.Read(id)
		if err != nil {
			web.Failure(c, 404, errors.New("Appointment not found"))
			return
		}
		web.Success(c, 200, appointment)
	}
}

func (h *appointmenthandler) Put() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("TOKEN")
		if token == "" {
			web.Failure(c, 401, errors.New("Token Not Found"))
			return
		}

		if token != os.Getenv("TOKEN") {
			web.Failure(c, 401, errors.New("Invalid Token"))
			return
		}
		var appointment dto.AppointmentInsert
		err := c.ShouldBindJSON(&appointment)
		if err != nil {
			web.Failure(c, 400, errors.New("Invalid patient"))
			return
		}
		valid, err := validateEmptysAppointment(&appointment)
		if !valid {
			web.Failure(c, 400, errors.New(err.Error()))
			return
		}
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, 400, errors.New("Invalid ID"))
			return
		}
		err = h.s.Update(id, appointment)
		if err != nil {
			web.Failure(c, 409, errors.New(err.Error()))
			return
		}
		web.Success(c, 200, "Patient Updated")
	}
}

func (h *appointmenthandler) Patch() gin.HandlerFunc {
	type AppointmentRequest struct {
		PatientId   int    `json:"patientid,omitempty"`
		DentistId   int    `json:"dentistid,omitempty"`
		DateUp      string `json:"dateup,omitempty"`
		Hour        string `json:"hour,omitempty"`
		Description string `json:"description,omitempty"`
	}
	return func(c *gin.Context) {
		token := c.GetHeader("TOKEN")
		if token == "" {
			web.Failure(c, 401, errors.New("Token Not Found"))
			return
		}

		if token != os.Getenv("TOKEN") {
			web.Failure(c, 401, errors.New("Invalid Token"))
			return
		}
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, 400, errors.New("Invalid ID"))
			return
		}
		var r AppointmentRequest
		if err := c.ShouldBindJSON(&r); err != nil {
			web.Failure(c, 400, errors.New("Invalid Request"))
			return
		}
		update := dto.AppointmentInsert{
			DentistId:   r.DentistId,
			PatientId:   r.PatientId,
			DateUp:      r.DateUp,
			Hour:        r.Hour,
			Description: r.Description,
		}
		err = h.s.Update(id, update)
		if err != nil {
			web.Failure(c, 409, errors.New(err.Error()))
			return
		}
		web.Success(c, 200, "Appointment Updated")
	}
}

func (h *appointmenthandler) PostxLicenseAndDni() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.GetHeader("Token")
		if token == "" {
			web.Failure(c, 401, errors.New("Token Not Found"))
			return
		}

		if token != os.Getenv("TOKEN") {
			web.Failure(c, 401, errors.New("Invalid Token"))
			return
		}

		var appointment dto.AppointmentPost
		err := c.ShouldBindJSON(&appointment)
		if err != nil {
			web.Failure(c, 400, errors.New("Invalid appointment"))
			return
		}

		if appointment.DentistLicense == "" || appointment.PatientDni == "" {
			web.Failure(c, 400, errors.New("We need DNI and License to create an appointment"))
			return
		}
		r, err := h.s.CreateAppointmentByDniAndLicense(appointment)
		if err != nil {
			web.Failure(c, 400, errors.New(err.Error()))
			return
		}

		web.Success(c, 201, r)
	}
}

func (h *appointmenthandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Token")
		if token == "" {
			web.Failure(c, 401, errors.New("Token Not Found"))
			return
		}
		if token != os.Getenv("TOKEN") {
			web.Failure(c, 401, errors.New("Invalid Token"))
			return
		}
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, 400, errors.New("Invalid ID"))
			return
		}
		err = h.s.Delete(id)
		if err != nil {
			web.Failure(c, 400, errors.New(err.Error()))
			return
		}

		web.Success(c, 200, "Appointment Deleted")
	}
}

func (h *appointmenthandler) GetByDNI() gin.HandlerFunc {
	return func(c *gin.Context) {
		dniParam := c.Param("dni")
		dni, err := strconv.Atoi(dniParam)
		if err != nil {
			web.Failure(c, 404, errors.New("Appointment not found"))
			return
		}
		appointments, err := h.s.Read(dni)
		if err != nil {
			web.Failure(c, 404, errors.New("Appointment not found"))
			return
		}
		web.Success(c, 200, appointments)
	}

}

func validateEmptysAppointment(appointment *dto.AppointmentInsert) (bool, error) {
	if appointment.Description == "" || appointment.Hour == "" || appointment.DateUp == "" || appointment.DentistId == 0 || appointment.PatientId == 0 {
		return false, errors.New("fields can't be empty")
	}
	return true, nil
}
