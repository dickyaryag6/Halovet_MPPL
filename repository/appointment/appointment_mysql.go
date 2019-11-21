package repository

import (
	dbCon "Halovet/driver"
	models "Halovet/models"
	"time"

	"database/sql"
	. "fmt"
	"log"
	"strconv"
)

var err error
var db *sql.DB

func init() {
	// KONEK KE DATABASE
	db, err = dbCon.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func Insert(time_appointment string, doctor_name string, pet_owner_name string, pet_owner_id int, pet_type string, complaint string) (models.Appointment, bool) {
	var appointment models.Appointment

	sql_statement := "insert into appointment (appointment_time,doctor_name, pet_owner_name, pet_owner_id, pet_type, complaint_description) values (?,?,?,?,?,?)"
	row, err := db.Exec(sql_statement, time_appointment, doctor_name, pet_owner_name, pet_owner_id, pet_type, complaint)

	if err != nil {
		Println(err.Error())
		return appointment, false
	}

	id, err := row.LastInsertId()
	if err != nil {
		log.Fatal(err.Error())
	}
	appointment.AppointmentID = id
	appointment.Doctor_name = doctor_name
	appointment.Time_Appointment = time_appointment
	appointment.Pet_Owner_Name = pet_owner_name
	appointment.Pet_owner_id = pet_owner_id
	appointment.Pet_Type = pet_type
	appointment.Complaint = complaint
	appointment.CreatedAt = Sprintf(time.Now().Format("2006-01-02 15:04:05"))
	appointment.UpdatedAt = Sprintf(time.Now().Format("2006-01-02 15:04:05"))

	return appointment, true

}

func FindAllAppointment(limitstart string, limit string) ([]models.Appointment, error) {
	var Appointment models.Appointment
	var Appointments []models.Appointment

	realLimitStart, err := strconv.Atoi(limitstart)
	if err != nil {
		Println("format limit salah")
		return Appointments, err
	}
	realLimit, err := strconv.Atoi(limit)
	if err != nil {
		Println("format limit salah")
		return Appointments, err
	}

	sqlStatement := "select * from appointment order by created_at limit ?, ?"
	results, err := db.Query(sqlStatement, realLimitStart, realLimit)
	if err != nil {
		panic(err.Error())
		return Appointments, err
	}
	for results.Next() {
		err = results.Scan(&Appointment.AppointmentID,
			&Appointment.Time_Appointment,
			&Appointment.Doctor_name,
			&Appointment.Pet_Owner_Name,
			&Appointment.Pet_Type,
			&Appointment.Complaint,
			&Appointment.IsPaid,
			&Appointment.CreatedAt,
			&Appointment.UpdatedAt,
			&Appointment.Pet_owner_id)
		if err != nil {
			panic(err.Error())
		} else {
			Appointments = append(Appointments, Appointment)
		}

	}

	return Appointments, nil
}

func FindbyID(id string) (models.Appointment, bool) {
	var appointment models.Appointment
	appointmentid, err := strconv.Atoi(id)
	if err != nil {
		Println("format ID salah")
	}
	sql_statement := "select * from appointment where id = ?"
	err = db.QueryRow(sql_statement, appointmentid).
		Scan(&appointment.AppointmentID,
			&appointment.Time_Appointment,
			&appointment.Doctor_name,
			&appointment.Pet_Owner_Name,
			&appointment.Pet_Type,
			&appointment.Complaint,
			&appointment.IsPaid,
			&appointment.CreatedAt,
			&appointment.UpdatedAt,
			&appointment.Pet_owner_id,
		)

	if err != nil {
		Println(err.Error())
		return appointment, false
	}

	return appointment, true
}

// func FindAll()

func Remove(id string) bool {

	appointmentid, err := strconv.Atoi(id)

	if err != nil {
		Println("format ID salah")
	}
	sql_statement := "delete from appointment where id = ?"
	_, err = db.Exec(sql_statement, appointmentid)

	if err != nil {
		Println(err.Error())
		return false
	}
	return true
}

func Update(id string, doctor_name string, pet_type string, complaint string, appointment_time string) bool {
	// var appointment models.Appointment
	appointmentid, err := strconv.Atoi(id)
	if err != nil {
		Println("format ID salah")
	}
	timeNow := Sprintf(time.Now().Format("2006-01-02 15:04:05"))
	sql_statement := "update appointment set doctor_name = ?, pet_type = ?, complaint_description = ?, appointment_time = ?, updated_at = ? where id = ?"
	_, err = db.Exec(sql_statement, doctor_name, pet_type, complaint, appointment_time, timeNow, appointmentid)

	if err != nil {
		Println(err.Error())
		return false
	}

	// //ROW AFFECTED
	// id, err := row.RowsAffected()
	// if err != nil {
	//   log.Fatal(err.Error())
	// }

	return true

}
