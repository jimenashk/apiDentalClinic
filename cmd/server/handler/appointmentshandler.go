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

// ListAppointments godoc
// @Summary Get All Appointments
// @Tags Appointments
// @Description Get Appointments
// @Accept json
// @Produce json
// @Success 200 {object} web.response
// @Failure 400 {object} web.response
// @Router /appointments [get]
func (h *appointmenthandler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		appointments, _ := h.s.ReadAll()
		if len(appointments) == 0 {
			web.Failure(c, 400, errors.New("There is no Appointment"))
		}
		web.Success(c, 200, appointments)
	}
}

// PostAppointments godoc
// @Summary Post Appointments
// @Tags Appointments
// @Description Post a Appointments
// @Accept json
// @Produce json
// @Param token header string true "token"
// @Param dentist body dto.AppointmentsInsert true "Appointments to data"
// @Success 201 {object} web.response
// @Failure 401 {object} web.response
// @Failure 400 {object} web.response
// @Router /appointments [post]
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

// AppointmentsbyID godoc
// @Summary Get Appointments by ID
// @Tags Appointments
// @Description Get Appointments
// @Accept json
// @Produce json
// @Param id path int true "Appointments ID"
// @Success 200 {object} web.response
// @Failure 404 {object} web.response
// @Router /appointments/{id} [get]
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

// UpdateAppointment godoc
// @Summary Put Appointments
// @Tags Appointments
// @Description Put a Appointments
// @Accept json
// @Produce json
// @Param token header string true "token"
// @Param dentist body dto.AppointmentInsert true "Appointments to data"
// @Param id path int true "Appointments ID"
// @Success 200 {object} web.response
// @Failure 401 {object} web.response
// @Failure 400 {object} web.response
// @Failure 409 {object} web.response
// @Router /appointments/{id} [put]
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

// UpdateAppointment godoc
// @Summary Patch Appointments
// @Tags Appointments
// @Description Patch a Appointments
// @Accept json
// @Produce json
// @Param token header string true "token"
// @Param dentist body dto.AppointmentInsert true "Appointments to data"
// @Param id path int true "Appointments ID"
// @Success 200 {object} web.response
// @Failure 401 {object} web.response
// @Failure 400 {object} web.response
// @Failure 409 {object} web.response
// @Router /appointments/{id} [patch]
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

// PostAppointmentsWithLicenseAndDNI godoc
// @Summary PostAppointmentsWithLicenseAndDNI Appointments
// @Tags Appointments
// @Description Post a Appointments
// @Accept json
// @Produce json
// @Param token header string true "token"
// @Param dentist body dto.AppointmentPost true "Appointments to data"
// @Success 201 {object} web.response
// @Failure 401 {object} web.response
// @Failure 400 {object} web.response
// @Router /appointments/post [post]
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

// DeleteAppointment godoc
// @Summary Delete Appointments
// @Tags Appointments
// @Description Delete a Appointments
// @Accept json
// @Produce json
// @Param token header string true "token"
// @Param id path int true "Appointments ID"
// @Success 200 {object} web.response
// @Failure 401 {object} web.response
// @Failure 400 {object} web.response
// @Router /appointments/{id} [delete]
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

// AppointmentsbyID godoc
// @Summary Get Appointments by ID
// @Tags Appointments
// @Description Get Appointments
// @Accept json
// @Produce json
// @Param dni query string false "DNI"
// @Success 200 {object} web.response
// @Failure 404 {object} web.response
// @Router /appointments/dni [get]
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
