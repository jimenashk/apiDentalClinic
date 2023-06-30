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

func (h *patienthandler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		patients, _ := h.s.ReadAll()
		if len(patients) == 0 {
			web.Failure(c, 400, errors.New("there is no patient"))
		}
		web.Success(c, 200, patients)
	}
}

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

func validateEmptysPatient(patient *domain.Patient) (bool, error) {
	if patient.Name == "" || patient.LastName == "" || patient.Address == "" || patient.DNI == "" || patient.DateUp == "" {
		return false, errors.New("fields can't be empty")
	}
	return true, nil
}
