package dentist

import "apiDentalClinic/internal/domain"

type Service interface {
	ReadAll() ([]domain.Dentist, error)
	Read(id int) (domain.Dentist, error)
	Create(dentist domain.Dentist) error
	Update(id int, dentist domain.Dentist) error
	Delete(id int) error
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) ReadAll() ([]domain.Dentist, error) {
	l, err := s.r.ReadAll()
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (s *service) Read(id int) (domain.Dentist, error) {
	d, err := s.r.Read(id)
	if err != nil {
		return domain.Dentist{}, err
	}
	return d, nil
}

func (s *service) Create(dentist domain.Dentist) error {
	err := s.r.Create(dentist)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) Update(id int, dentist domain.Dentist) error {
	err := s.r.Update(id, dentist)
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
