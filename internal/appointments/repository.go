package appointments

import (
	"errors"
	"strings"
	"time"

	"apiDentalClinic/internal/domain"
	"apiDentalClinic/internal/domain/dto"
	"apiDentalClinic/pkg/store"
)

type Repository interface {
	ReadAll() ([]domain.Appointments, error)
	Create(appointment dto.AppointmentInsert) error
	Read(id int) (domain.Appointments, error)
	Update(id int, appointments dto.AppointmentInsert) error
	Delete(id int) error
	CreateAppointmentByDniAndLicense(appointment dto.AppointmentPost) (dto.AppointmentInsert, error)
	ReadAppointmentbyDNI(dni string) ([]dto.AppointmentGet, error)
}

type repository struct {
	store store.Store
}

func NewRepository(store store.Store) Repository {
	return &repository{store: store}
}

func (r *repository) ReadAll() ([]domain.Appointments, error) {
	list, err := r.store.ReadAllAppointments()
	if err != nil {
		return []domain.Appointments{}, err
	}
	return list, nil
}

func (r *repository) Create(appointment dto.AppointmentInsert) error {
	if !r.ValidatePatient(appointment.PatientId) {
		return errors.New("The patient does not exist")
	}
	if !r.ValidateDentist(appointment.DentistId) {
		return errors.New("The dentist does not exist")
	}
	err := r.store.CreateAppointment(appointment)
	if err != nil {
		return errors.New("Error creating a new Appointment: " + err.Error())
	}
	return nil
}

func (r *repository) Read(id int) (domain.Appointments, error) {
	appo, err := r.store.ReadAppointment(id)
	if err != nil {
		return domain.Appointments{}, err
	}
	return appo, nil
}

func (r *repository) Update(id int, appointment dto.AppointmentInsert) error {
	original, err := r.store.ReadAppointment(id)
	if err != nil {
		return errors.New("The Appointment does not exists")
	}
	if !r.ValidateDentist(original.Dentist.Id) {
		return errors.New("The Dentist does not exists")
	}
	if !r.ValidatePatient(original.Patient.Id) {
		return errors.New("The Patient does not exists")
	}
	complete := unchangeEmptysAppointment(appointment, original)
	err = r.store.UpdateAppointment(id, complete)
	if err != nil {
		return errors.New("Error updating Appointment")
	}
	return nil
}
func (r *repository) Delete(id int) error {
	err := r.store.DeleteAppointment(id)
	if err != nil {
		return errors.New("Error deleting a Dentist")
	}
	return nil
}

func (r *repository) CreateAppointmentByDniAndLicense(appointment dto.AppointmentPost) (dto.AppointmentInsert, error) {
	updateappointment := ChargeDataAppointment(appointment)
	appointmentOk, err := r.store.CreateAppointmentByDniAndLicense(updateappointment)
	if err != nil {
		return dto.AppointmentInsert{}, err
	}
	return appointmentOk, nil
}

func (r *repository) ReadAppointmentbyDNI(dni string) ([]dto.AppointmentGet, error) {
	if !r.ValidateDNI(dni) {
		return []dto.AppointmentGet{}, errors.New("There is no patient with that DNI")
	}
	list, err := r.store.ReadAppointmentbyDNI(dni)
	if err != nil {
		return []dto.AppointmentGet{}, errors.New("Error getting an appointment by DNI")
	}
	return list, nil
}

// Validate Functions

func (r *repository) ValidatePatient(id int) bool {
	_, err := r.store.ReadPatient(id)
	if err != nil {
		return false
	}
	return true
}

func (r *repository) ValidateDentist(id int) bool {
	_, err := r.store.ReadDentist(id)
	if err != nil {
		return false
	}
	return true
}

func (r *repository) ValidateDNI(dni string) bool {
	list, err := r.store.ReadAllPatient()
	if err != nil {
		return false
	}
	for _, patient := range list {
		if patient.DNI == dni {
			return true
		}
	}
	return false
}

func unchangeEmptysAppointment(appointmentDTO dto.AppointmentInsert, original domain.Appointments) dto.AppointmentInsert {

	if appointmentDTO.DateUp == "" {
		appointmentDTO.DateUp = original.DateUp
	}
	if appointmentDTO.Hour == "" {
		appointmentDTO.Hour = original.Hour
	}
	if appointmentDTO.Description == "" {
		appointmentDTO.Description = original.Description
	}
	if appointmentDTO.DentistId == 0 {
		appointmentDTO.DentistId = original.Dentist.Id
	}
	if appointmentDTO.PatientId == 0 {
		appointmentDTO.PatientId = original.Patient.Id
	}

	return appointmentDTO
}

func ChargeDataAppointment(appointment dto.AppointmentPost) dto.AppointmentPost {
	t := time.Now()
	s := t.Format("2006-01-02 15:04:05")
	fechayhora := strings.Split(s, " ")
	if appointment.DateUp == "" {
		appointment.DateUp = fechayhora[0]
	}
	if appointment.Hour == "" {
		appointment.Hour = fechayhora[1]
	}
	if appointment.Description == "" {
		appointment.Description = "Appointment created: " + fechayhora[0] + " " + fechayhora[1]
	}
	return appointment
}
