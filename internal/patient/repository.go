package patient

import (
	"errors"

	"apiDentalClinic/internal/domain"
	"apiDentalClinic/pkg/store"
)

type Repository interface {
	ReadAll() ([]domain.Patient, error)
	Read(id int) (domain.Patient, error)
	Create(patient domain.Patient) error
	Update(id int, patient domain.Patient) error
	Delete(id int) error
}

type repository struct {
	store store.Store
}

func NewRepository(store store.Store) Repository {
	return &repository{store: store}
}

func (r *repository) ReadAll() ([]domain.Patient, error) {
	list, err := r.store.ReadAllPatient()
	if err != nil {
		return []domain.Patient{}, err
	}
	return list, nil

}
func (r *repository) Read(id int) (domain.Patient, error) {
	patient, err := r.store.ReadPatient(id)
	if err != nil {
		return domain.Patient{}, errors.New("Patient not found")
	}
	return patient, nil
}

func (r *repository) Create(patient domain.Patient) error {
	if !r.ValidateDNI(patient.DNI) {
		return errors.New("DNI already exists")
	}
	err := r.store.CreatePatient(patient)
	if err != nil {
		return errors.New("Error creating a new Patient")
	}
	return nil
}
func (r *repository) Update(id int, patient domain.Patient) error {
	if !r.ValidateDNI(patient.DNI) {
		return errors.New("DNI already exists")
	}
	original, err := r.store.ReadPatient(id)
	if err != nil {
		return errors.New("The Patient does not exists")
	}
	complete := unchangeEmptysPatient(patient, original)
	err = r.store.UpdatePatient(id, complete)
	if err != nil {
		return errors.New("Error updating a Patient")
	}
	return nil
}

func (r *repository) Delete(id int) error {
	err := r.store.DeletePatient(id)
	if err != nil {
		return errors.New("Error deleting a Patient  Cause 1:he have still appointment. Cause 2: He doest exist.")
	}
	return nil
}

// Validation Functions
func (r *repository) ValidateDNI(DNI string) bool {
	list, err := r.store.ReadAllPatient()
	if err != nil {
		return false
	}
	for _, patient := range list {
		if patient.DNI == DNI {
			return false
		}
	}
	return true
}
func unchangeEmptysPatient(patient domain.Patient, original domain.Patient) domain.Patient {

	if patient.Name == "" {
		patient.Name = original.Name
	}
	if patient.LastName == "" {
		patient.LastName = original.LastName
	}
	if patient.Address == "" {
		patient.Address = original.Address
	}
	if patient.DNI == "" {
		patient.DNI = original.DNI
	}
	if patient.DateUp == "" {
		patient.DateUp = original.DateUp
	}

	return patient
}
