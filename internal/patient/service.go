package patient

import "apiDentalClinic/internal/domain"

type Service interface {
	ReadAll() ([]domain.Patient, error)
	Read(id int) (domain.Patient, error)
	Create(patient domain.Patient) error
	Update(id int, patient domain.Patient) error
	Delete(id int) error
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) ReadAll() ([]domain.Patient, error) {
	l, err := s.r.ReadAll()
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (s *service) Read(id int) (domain.Patient, error) {
	d, err := s.r.Read(id)
	if err != nil {
		return domain.Patient{}, err
	}
	return d, nil
}

func (s *service) Create(patient domain.Patient) error {
	err := s.r.Create(patient)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) Update(id int, patient domain.Patient) error {
	err := s.r.Update(id, patient)
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
