package handler

import (
	"errors"
	"os"
	"strconv"

	"apiDentalClinic/internal/dentist"
	"apiDentalClinic/internal/domain"
	"apiDentalClinic/pkg/web"

	"github.com/gin-gonic/gin"
)

type dentistHandler struct {
	s dentist.Service
}

func NewDentistHandler(s dentist.Service) *dentistHandler {
	return &dentistHandler{s: s}
}

func (h *dentistHandler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		dentists, _ := h.s.ReadAll()
		if len(dentists) == 0 {
			web.Failure(c, 400, errors.New("there is no dentist"))
		}
		web.Success(c, 200, dentists)
	}
}

func (h *dentistHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, 404, errors.New("invalid id"))
			return
		}
		dentist, err := h.s.Read(id)
		if err != nil {
			web.Failure(c, 404, errors.New("dentist not found"))
			return
		}
		web.Success(c, 200, dentist)
	}

}

func (h *dentistHandler) Post() gin.HandlerFunc {
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

		var dentist domain.Dentist
		err := c.ShouldBindJSON(&dentist)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid Dentist"))
			return
		}

		valid, err := validateEmptys(&dentist)
		if !valid {
			web.Failure(c, 400, errors.New(err.Error()))
			return
		}

		err = h.s.Create(dentist)
		if err != nil {
			web.Failure(c, 400, errors.New(err.Error()))
			return
		}

		web.Success(c, 201, dentist)
	}
}

func (h *dentistHandler) Delete() gin.HandlerFunc {
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

		web.Success(c, 200, "Dentist Deleted")
	}
}

func (h *dentistHandler) Put() gin.HandlerFunc {
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
		var dentist domain.Dentist
		err := c.ShouldBindJSON(&dentist)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid dentist"))
			return
		}
		valid, err := validateEmptys(&dentist)
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
		err = h.s.Update(id, dentist)
		if err != nil {
			web.Failure(c, 409, errors.New(err.Error()))
			return
		}
		web.Success(c, 200, "Dentist Updated")
	}
}

func (h *dentistHandler) Patch() gin.HandlerFunc {
	type Request struct {
		Name     string `json:"name,omitempty"`
		LastName string `json:"lastname,omitempty"`
		License  string `json:"license,omitempty" `
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
		update := domain.Dentist{
			Name:     r.Name,
			LastName: r.LastName,
			License:  r.License,
		}

		err = h.s.Update(id, update)
		if err != nil {
			web.Failure(c, 409, errors.New(err.Error()))
			return
		}
		web.Success(c, 200, "Dentist Updated")
	}
}

func validateEmptys(dentist *domain.Dentist) (bool, error) {
	if dentist.Name == "" || dentist.LastName == "" || dentist.License == "" {
		return false, errors.New("fields can't be empty")
	}
	return true, nil
}
