package repository

import(
  models "Halovet/models"
  dbCon "Halovet/driver"
  // "time"
  . "fmt"
  "database/sql"
  "log"
  "strconv"
)

var err error
var db *sql.DB

func init(){
  // KONEK KE DATABASE
  db, err = dbCon.Connect()
  if err != nil {
    log.Fatal(err.Error())
  }
}

func Insert(time string, doctor_name string, pet_owner_name string, pet_name string,complaint string) (models.Appointment, bool) {
  var appointment models.Appointment

  sql_statement := "insert into appointment (time,doctor_name, pet_owner_name, pet_name, complaint) values (?,?,?,?,?)"
  row ,err := db.Exec(sql_statement, time,doctor_name, pet_owner_name, pet_name, complaint)

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
   appointment.Time_Appointment = time
   appointment.Pet_owner_name = pet_owner_name
   appointment.Pet_name = pet_name
   appointment.Complaint = complaint

   return appointment, true

}

func FindbyID(id string) (models.Appointment,bool) {
  var appointment models.Appointment
  appointmentid, err := strconv.Atoi(id)
  if err != nil {
	   Println("format ID salah")
	}
  sql_statement := "select * from appointment where id = ?"
  err = db.QueryRow(sql_statement, appointmentid).
            Scan(&appointment.AppointmentID, &appointment.Time_Appointment, &appointment.Doctor_name, &appointment.Pet_owner_name, &appointment.Pet_name, &appointment.Complaint)

  if err != nil {
     Println(err.Error())
     return appointment, false
   }

  //  for rows.Next() {
  //       var each = models.Appointment{}
  //       var err = rows.Scan(&each.ID, &each.time, &each.doctor_name, &each.pet_owner_name, &each.pet_name, &each.complaint)
  //
  //       if err != nil {
  //           fmt.Println(err.Error())
  //           return
  //       }
  //
  //       result = append(result, each)
  // }

   return appointment, true
}

// func FindAll()

func Remove(id string) bool {
  Println(id)
  appointmentid, err:=strconv.Atoi(id)
  Println(appointmentid)
  if err != nil {
    Println("format ID salah")
  }
  sql_statement := "delete from appointment where id = ?"
  _ ,err = db.Exec(sql_statement, appointmentid)

  if err != nil {
     Println(err.Error())
     return false
   }
  return true
}

func Update(id string, doctor_name string, pet_name string, complaint string, time string) bool {
  // var appointment models.Appointment
  appointmentid, err:=strconv.Atoi(id)
  if err != nil {
    Println("format ID salah")
  }
  sql_statement := "update appointment set doctor_name = ?, pet_name = ?, complaint = ?, time = ? where id = ?"
  _, err = db.Exec(sql_statement, doctor_name, pet_name, complaint, time, appointmentid)

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
