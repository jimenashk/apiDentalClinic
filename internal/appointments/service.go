package appointments

import (
	"apiDentalClinic/internal/domain"
	"apiDentalClinic/internal/domain/dto"
)

type Service interface {
	ReadAll() ([]domain.Appointments, error)
	Create(appointment dto.AppointmentInsert) error
	Read(id int) (domain.Appointments, error)
	Update(id int, appointment dto.AppointmentInsert) error
	Delete(id int) error
	CreateAppointmentByDniAndLicense(appointment dto.AppointmentPost) (dto.AppointmentInsert, error)
	ReadbyDni(dni string) ([]dto.AppointmentGet, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) ReadAll() ([]domain.Appointments, error) {
	l, err := s.r.ReadAll()
	if err != nil {
		return nil, err
	}
	return l, nil
}
func (s *service) Create(appointment dto.AppointmentInsert) error {
	err := s.r.Create(appointment)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) Read(id int) (domain.Appointments, error) {
	t, err := s.r.Read(id)
	if err != nil {
		return domain.Appointments{}, err
	}
	return t, nil
}
func (s *service) Update(id int, appointment dto.AppointmentInsert) error {
	err := s.r.Update(id, appointment)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) Delete(id int) error {
	err := s.r.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) CreateAppointmentByDniAndLicense(appointment dto.AppointmentPost) (dto.AppointmentInsert, error) {
	t, err := s.r.CreateAppointmentByDniAndLicense(appointment)
	if err != nil {
		return dto.AppointmentInsert{}, err
	}
	return t, nil
}
func (s *service) ReadbyDni(dni string) ([]dto.AppointmentGet, error) {
	l, err := s.r.ReadAppointmentbyDNI(dni)
	if err != nil {
		return []dto.AppointmentGet{}, err
	}
	return l, nil
}
