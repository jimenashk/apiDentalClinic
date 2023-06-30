package handler

import (
	"errors"
	"os"
	"strconv"

	"apiDentalClinic/internal/domain"
	"apiDentalClinic/internal/patient"
	"apiDentalClinic/pkg/web"

	"github.com/gin-gonic/gin"
)

type patienthandler struct {
	s patient.Service
}

func NewPatientHandler(s patient.Service) *patienthandler {
	return &patienthandler{s: s}
}

// ListPatient godoc
// @Summary Get All Patient
// @Tags Patient
// @Description Get Patient
// @Accept json
// @Produce json
// @Success 200 {object} web.response
// @Failure 400 {object} web.response
// @Router /patients [get]
func (h *patienthandler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		patients, _ := h.s.ReadAll()
		if len(patients) == 0 {
			web.Failure(c, 400, errors.New("there is no patient"))
		}
		web.Success(c, 200, patients)
	}
}

// PatientbyID godoc
// @Summary Get Patient by ID
// @Tags Patient
// @Description Get Patient
// @Accept json
// @Produce json
// @Param id path int true "Patient ID"
// @Success 200 {object} web.response
// @Failure 404 {object} web.response
// @Router /patients/{id} [get]
func (h *patienthandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, 404, errors.New("invalid id"))
			return
		}
		patient, err := h.s.Read(id)
		if err != nil {
			web.Failure(c, 404, errors.New("patient not found"))
			return
		}
		web.Success(c, 200, patient)
	}

}

// PostPatient godoc
// @Summary Post Patient
// @Tags Patient
// @Description Post a Patient
// @Accept json
// @Produce json
// @Param token header string true "token"
// @Param dentist body domain.Patient true "Patient to data"
// @Success 201 {object} web.response
// @Failure 401 {object} web.response
// @Failure 400 {object} web.response
// @Router /patients [post]
func (h *patienthandler) Post() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.GetHeader("token")
		if token == "" {
			web.Failure(c, 401, errors.New("token Not Found"))
			return
		}

		if token != os.Getenv("TOKEN") {
			web.Failure(c, 401, errors.New("invalid Token"))
			return
		}

		var patient domain.Patient
		err := c.ShouldBindJSON(&patient)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid Patient"))
			return
		}

		valid, err := validateEmptysPatient(&patient)
		if !valid {
			web.Failure(c, 400, errors.New(err.Error()))
			return
		}

		err = h.s.Create(patient)
		if err != nil {
			web.Failure(c, 400, errors.New(err.Error()))
			return
		}

		web.Success(c, 201, patient)
	}
}

// DeletePatient godoc
// @Summary Delete Patient
// @Tags Patient
// @Description Delete a Patient
// @Accept json
// @Produce json
// @Param token header string true "token"
// @Param id path int true "Patient ID"
// @Success 200 {object} web.response
// @Failure 401 {object} web.response
// @Failure 400 {object} web.response
// @Router /patients/{id} [delete]
func (h *patienthandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Token")
		if token == "" {
			web.Failure(c, 401, errors.New("token Not Found"))
			return
		}
		if token != os.Getenv("TOKEN") {
			web.Failure(c, 401, errors.New("invalid Token"))
			return
		}
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid ID"))
			return
		}
		err = h.s.Delete(id)
		if err != nil {
			web.Failure(c, 400, errors.New(err.Error()))
			return
		}

		web.Success(c, 200, "patient Deleted")
	}
}

// UpdatePatient godoc
// @Summary Put Patient
// @Tags Patient
// @Description Put a Patient
// @Accept json
// @Produce json
// @Param token header string true "token"
// @Param dentist body domain.Patient true "Patient to data"
// @Param id path int true "Patient ID"
// @Success 200 {object} web.response
// @Failure 401 {object} web.response
// @Failure 400 {object} web.response
// @Failure 409 {object} web.response
// @Router /patients/{id} [put]
func (h *patienthandler) Put() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("TOKEN")
		if token == "" {
			web.Failure(c, 401, errors.New("token Not Found"))
			return
		}

		if token != os.Getenv("TOKEN") {
			web.Failure(c, 401, errors.New("invalid Token"))
			return
		}
		var patient domain.Patient
		err := c.ShouldBindJSON(&patient)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid patient"))
			return
		}
		valid, err := validateEmptysPatient(&patient)
		if !valid {
			web.Failure(c, 400, errors.New(err.Error()))
			return
		}
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid ID"))
			return
		}
		err = h.s.Update(id, patient)
		if err != nil {
			web.Failure(c, 409, errors.New(err.Error()))
			return
		}
		web.Success(c, 200, "patient Updated")
	}
}

// UpdatePatient godoc
// @Summary Patch Patient
// @Tags Patient
// @Description Patch a Patient
// @Accept json
// @Produce json
// @Param token header string true "token"
// @Param dentist body domain.Patient true "Patient to data"
// @Param id path int true "Patient ID"
// @Success 200 {object} web.response
// @Failure 401 {object} web.response
// @Failure 400 {object} web.response
// @Failure 409 {object} web.response
// @Router /patients/{id} [patch]
func (h *patienthandler) Patch() gin.HandlerFunc {
	type Request struct {
		Name     string `json:"name,omitempty"`
		LastName string `json:"lastname,omitempty"`
		Address  string `json:"address,omitempty"`
		DNI      string `json:"dni,omitempty"`
		DateUp   string `json:"dateup,omitempty"`
	}
	return func(c *gin.Context) {
		token := c.GetHeader("TOKEN")
		if token == "" {
			web.Failure(c, 401, errors.New("token Not Found"))
			return
		}

		if token != os.Getenv("TOKEN") {
			web.Failure(c, 401, errors.New("invalid Token"))
			return
		}
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid ID"))
			return
		}
		var r Request
		if err := c.ShouldBindJSON(&r); err != nil {
			web.Failure(c, 400, errors.New("invalid Request"))
			return
		}
		update := domain.Patient{
			Name:     r.Name,
			LastName: r.LastName,
			Address:  r.Address,
			DNI:      r.DNI,
			DateUp:   r.DateUp,
		}
		err = h.s.Update(id, update)
		if err != nil {
			web.Failure(c, 409, errors.New(err.Error()))
			return
		}
		web.Success(c, 200, "Patient Updated")
	}
}

// validateEmptys valida que los campos no esten vacios
func validateEmptysPatient(patient *domain.Patient) (bool, error) {
	if patient.Name == "" || patient.LastName == "" || patient.Address == "" || patient.DNI == "" || patient.DateUp == "" {
		return false, errors.New("fields can't be empty")
	}
	return true, nil
}
