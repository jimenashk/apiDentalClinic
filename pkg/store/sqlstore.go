package store

import (
	"database/sql"
	"errors"
	"log"

	"apiDentalClinic/internal/domain"
	"apiDentalClinic/internal/domain/dto"
)

type store struct {
	db *sql.DB
}

func NewSQLStore(db *sql.DB) Store {
	return &store{db: db}
}

type Store interface {
	// ---------------- DENTIST
	// ReadAllDentist Trae todos los dentistas
	ReadAllDentists() ([]domain.Dentist, error)
	// Read devuelve un dentista por su id
	ReadDentist(id int) (domain.Dentist, error)
	// CreateDentist agrega un nuevo dentista
	CreateDentist(dentist domain.Dentist) error
	// UpdateDentist actualiza un dentista
	UpdateDentist(id int, dentist domain.Dentist) error
	// DeleteDentist elimina un dentista
	DeleteDentist(id int) error

	// ---------------- PATIENT
	// ReadAll Trae todos los pacientes
	ReadAllPatient() ([]domain.Patient, error)
	// Read devuelve un paciente por su id
	ReadPatient(id int) (domain.Patient, error)
	// Create agrega un nuevo paciente
	CreatePatient(dentist domain.Patient) error
	// Update actualiza un paciente
	UpdatePatient(id int, patient domain.Patient) error
	// Delete elimina un paciente
	DeletePatient(id int) error

	// ---------------- APPOINTMENTS
	// Devuelve todos los turnos
	ReadAllAppointments() ([]domain.Appointments, error)
	// Create crea un turno
	CreateAppointment(appointment dto.AppointmentInsert) error
	// Leer un turno x ID
	ReadAppointment(id int) (domain.Appointments, error)
	// Update actualiza un turno
	UpdateAppointment(id int, appointmentDTO dto.AppointmentInsert) error
	// Delete elimina un turno
	DeleteAppointment(id int) error
	//Crea un turno x DNI del Paciente y la Matricula del Dentista
	CreateAppointmentByDniAndLicense(appointment dto.AppointmentPost) (dto.AppointmentInsert, error)
	//Devuelve un turno x DNI de paciente
	ReadAppointmentbyDNI(dni string) ([]dto.AppointmentGet, error)
}

// 		----------DENTIST--------------

func (s *store) ReadAllDentists() ([]domain.Dentist, error) {
	var list []domain.Dentist
	var dentist domain.Dentist
	rows, err := s.db.Query("SELECT * FROM dentist")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		if err := rows.Scan(&dentist.Id, &dentist.Name, &dentist.LastName, &dentist.License); err != nil {
			return nil, err
		}
		list = append(list, dentist)
	}
	rows.Close()
	return list, nil
}

func (s *store) ReadDentist(id int) (domain.Dentist, error) {
	//defer func () {s.db.Close()}()
	var dentist domain.Dentist
	row := s.db.QueryRow("SELECT * FROM dentist WHERE id=?", id)

	if err := row.Scan(&dentist.Id, &dentist.Name, &dentist.LastName, &dentist.License); err != nil {
		return domain.Dentist{}, errors.New("The dentist with this id not exist")
		//panic(err)
	}

	return dentist, nil
}

func (s *store) CreateDentist(dentist domain.Dentist) error {

	st, err := s.db.Prepare("INSERT INTO dentist (name, lastname, license) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer st.Close()

	res, err := st.Exec(dentist.Name, dentist.LastName, dentist.License)
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (s *store) UpdateDentist(id int, dentist domain.Dentist) error {

	st, err := s.db.Prepare("UPDATE dentist SET name = ?, lastName = ?, license = ? WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer st.Close()

	_, err = st.Exec(dentist.Name, dentist.LastName, dentist.License, id)
	if err != nil {
		log.Fatal(err)
	}

	return nil

}

func (s *store) DeleteDentist(id int) error {
	//Preguntar si esta bien usar un metodo
	var idselect int
	row := s.db.QueryRow("SELECT id FROM dentist WHERE id = ?", id)
	if err := row.Scan(&idselect); err != nil {
		return errors.New("The dentist doest exists.")
	}
	query := "DELETE FROM dentist WHERE id = ?"
	_, err := s.db.Exec(query, idselect)
	if err != nil {
		return err
	}
	return nil

}

// 		----------PATIENT--------------

func (s *store) ReadAllPatient() ([]domain.Patient, error) {
	var list []domain.Patient
	var patient domain.Patient
	rows, err := s.db.Query("SELECT * FROM patient")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		if err := rows.Scan(&patient.Id, &patient.Name, &patient.LastName, &patient.Address, &patient.DNI, &patient.DateUp); err != nil {
			return nil, err
		}
		list = append(list, patient)
	}
	rows.Close()
	return list, nil

}

func (s *store) ReadPatient(id int) (domain.Patient, error) {
	//defer func () {s.db.Close()}()
	var patient domain.Patient
	row := s.db.QueryRow("SELECT * FROM patient WHERE id=?", id)

	if err := row.Scan(&patient.Id, &patient.Name, &patient.LastName, &patient.Address, &patient.DNI, &patient.DateUp); err != nil {
		return domain.Patient{}, err
		//panic(patient
	}
	return patient, nil
}

func (s *store) CreatePatient(patient domain.Patient) error {
	st, err := s.db.Prepare("INSERT INTO patient (name, lastname, address, dni, dateup) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer st.Close()

	res, err := st.Exec(patient.Name, patient.LastName, patient.Address, patient.DNI, patient.DateUp)
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (s *store) UpdatePatient(id int, patient domain.Patient) error {
	st, err := s.db.Prepare("UPDATE patient SET name = ?, lastName = ?, address = ?, dni = ?, dateup = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer st.Close()

	res, err := st.Exec(patient.Name, patient.LastName, patient.Address, patient.DNI, patient.DateUp, id)
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (s *store) DeletePatient(id int) error {
	//Preguntar si esta bien usar un metodo
	var idselect int
	row := s.db.QueryRow("SELECT id FROM patient WHERE id = ?", id)
	if err := row.Scan(&idselect); err != nil {
		return errors.New("The patient doest exists.")
	}
	query := "DELETE FROM patient WHERE id = ?"
	_, err := s.db.Exec(query, idselect)
	if err != nil {
		return err
	}
	return nil

}

// 		----------APPOINTMENTS--------------

func (s *store) ReadAllAppointments() ([]domain.Appointments, error) {
	var list []domain.Appointments
	var appointment domain.Appointments
	rows, err := s.db.Query("SELECT t.id, t.date, t.hour, t.description, p.id, d.id FROM appointments AS t JOIN patient AS p ON t.patientid = p.id JOIN dentist AS d ON t.dentistid = d.id")
	if err != nil {
		return nil, err
	}
	for rows.Next() {

		if err := rows.Scan(&appointment.Id, &appointment.DateUp, &appointment.Hour, &appointment.Description, &appointment.Patient.Id, &appointment.Dentist.Id); err != nil {
			return nil, err
		}
		list = append(list, appointment)
	}
	rows.Close()
	return list, nil
}

func (s *store) ReadAppointment(id int) (domain.Appointments, error) {
	var appointment domain.Appointments
	row := s.db.QueryRow("SELECT t.id, t.date, t.hour, t.description, p.id, d.id FROM appointments AS t JOIN patient AS p ON t.patientid = p.id JOIN dentist AS d ON t.dentistid = d.id WHERE t.id = ?", id)

	if err := row.Scan(&appointment.Id, &appointment.DateUp, &appointment.Hour, &appointment.Description, &appointment.Patient.Id, &appointment.Dentist.Id); err != nil {
		return domain.Appointments{}, err
	}
	return appointment, nil
}

func (s *store) CreateAppointment(appointmentDTO dto.AppointmentInsert) error {
	st, err := s.db.Prepare("INSERT INTO appointments (date, hour, description, patientid, dentistid) VALUES (?, ?, ?, ? ,?)")
	if err != nil {
		return err
	}

	res, err := st.Exec(appointmentDTO.DateUp, appointmentDTO.Hour, appointmentDTO.Description, appointmentDTO.PatientId, appointmentDTO.DentistId)
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	defer st.Close()
	return nil

}
func (s *store) UpdateAppointment(id int, appointmentDTO dto.AppointmentInsert) error {
	st, err := s.db.Prepare("UPDATE appointments SET date = ?, hour = ?, description = ?, patientid = ?, dentistid = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer st.Close()

	res, err := st.Exec(&appointmentDTO.DateUp, &appointmentDTO.Hour, &appointmentDTO.Description, &appointmentDTO.PatientId, &appointmentDTO.DentistId, id)
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (s *store) CreateAppointmentByDniAndLicense(appointment dto.AppointmentPost) (dto.AppointmentInsert, error) {
	var idpatient int
	row := s.db.QueryRow("SELECT id FROM patient WHERE dni=?", appointment.PatientDni)
	if err := row.Scan(&idpatient); err != nil {
		return dto.AppointmentInsert{}, errors.New("The patient doesnt exists.")
	}
	var iddentist int
	row2 := s.db.QueryRow("SELECT id FROM dentist WHERE license=?", appointment.DentistLicense)
	if err := row2.Scan(&iddentist); err != nil {
		return dto.AppointmentInsert{}, errors.New("The dentist doesnt exists.")
	}
	newappointmentinsert := dto.AppointmentInsert{
		DateUp:      appointment.DateUp,
		Hour:        appointment.Hour,
		Description: appointment.Description,
		PatientId:   idpatient,
		DentistId:   iddentist,
	}

	st, err := s.db.Prepare("INSERT INTO appointments (date, hour, description, patientid, dentistid) VALUES (?, ?, ?, ? ,?)")
	if err != nil {
		return dto.AppointmentInsert{}, err
	}

	res, err := st.Exec(newappointmentinsert.DateUp, newappointmentinsert.Hour, newappointmentinsert.Description, newappointmentinsert.PatientId, newappointmentinsert.DentistId)
	if err != nil {
		return dto.AppointmentInsert{}, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return dto.AppointmentInsert{}, err
	}
	return newappointmentinsert, nil

}

func (s *store) DeleteAppointment(id int) error {
	var idselect int
	row := s.db.QueryRow("SELECT id FROM appointments WHERE id = ?", id)
	if err := row.Scan(&idselect); err != nil {
		return errors.New("The patient doesnt exists.")
	}
	query := "DELETE FROM appointments WHERE id = ?"
	_, err := s.db.Exec(query, idselect)
	if err != nil {
		return err
	}
	return nil
}

func (s *store) ReadAppointmentbyDNI(dni string) ([]dto.AppointmentGet, error) {
	var appointment dto.AppointmentGet
	var list []dto.AppointmentGet
	rows, err := s.db.Query("SELECT t.date, t.hour, t.description, p.name, p.dni, d.name  FROM appointments AS t JOIN patient AS p ON t.patientid = p.id JOIN dentist AS d ON t.dentistid = d.id WHERE p.dni = ?", dni)
	if err != nil {
		return []dto.AppointmentGet{}, err
	}
	for rows.Next() {

		if err := rows.Scan(&appointment.Date, &appointment.Hour, &appointment.Description, &appointment.PatientName, &appointment.DNIPatient, &appointment.DentistName); err != nil {
			return nil, err
		}
		list = append(list, appointment)
	}
	rows.Close()
	return list, nil
}
