 package handler

import(
  // . "fmt"
  "net/http"
  "encoding/json"
  // "time"
  method "Halovet/repository/appointment"
  models "Halovet/models"
  jwt "github.com/dgrijalva/jwt-go"
  . "fmt"
  "github.com/gorilla/mux"
)

//YANG DIKEMBALIKAN CUMA APPOINTMENT DENGAN USER ID TERTENTU
// func GetAppointmentByUserID(w http.ResponseWriter, r *http.Request) {
//
// }


func GetAppointmentByID(w http.ResponseWriter, r *http.Request) {
  Println("Endpoint Hit: GetAppointmentByID")
  var response models.AppointmentResponse
  var result []models.Appointment

  vars := mux.Vars(r)

  real_result, status := method.FindbyID(vars["id"])

  w.Header().Set("Content-Type", "application/json")
  if status == false {
    message := "Appointment Failed to Get"
    w.WriteHeader(400)
    response.Status   = false
    response.Message  = message
    // response.Data = result
    json.NewEncoder(w).Encode(response)
  } else {
    result = append(result,real_result)

    data := map[string]interface{}{
      "Appointment" : result,
    }

    message := "Appointment Get Succesfully"
    w.WriteHeader(202)
    response.Status   = true
    response.Message  = message
    response.Data = data
    json.NewEncoder(w).Encode(response)
  }
}

//CUMA BOLEH ADMIN
// func GetAllAppointments(w http.ResponseWriter, r *http.Request) {
//
// }

func CreateAppointment(w http.ResponseWriter, r *http.Request) {

  Println("Endpoint Hit: createAppointment")
  // reqBody, _ := ioutil.ReadAll(r.Body)
  var appointment models.Appointment
  var result []models.Appointment
  var response models.AppointmentResponse
  // json.Unmarshal(reqBody, &appointment)
  appointment.Doctor_name = r.FormValue("doctor_name")
  appointment.Pet_name = r.FormValue("pet_name")
  appointment.Complaint = r.FormValue("complaint")
  appointment.Time_Appointment = r.FormValue("time")
  //ISI FORM CUMA doctor_name,pet_name,complaint, sama time_appointment
  // currentTime := time.Now()
  // appointment.Time_Appointment = currentTime.Format("2006-01-02 15:04:05")

  userInfo:=r.Context().Value("userInfo").(jwt.MapClaims)
  user := userInfo["User"]
  user_real, _ := user.(map[string]interface{})
  appointment.Pet_owner_name = Sprintf("%v", user_real["Name"])

  //INSERT OBJEK KE DB

  real_result, status := method.Insert(
    appointment.Time_Appointment,
    appointment.Doctor_name,
    appointment.Pet_owner_name,
    appointment.Pet_name,
    appointment.Complaint,
  )

  if status == false {
    message := "Appointment Failed to Update"
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(400)
    response.Status   = false
    response.Message  = message
    // response.Data = result
    json.NewEncoder(w).Encode(response)
  } else {

    result = append(result,real_result)

    data := map[string]interface{}{
      "Appointment" : result,
    }

    //KIRIM EMAIL

    message := "Appointment Booked Succesfully"
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)
    response.Status   = true
    response.Message  = message
    response.Data     = data
    json.NewEncoder(w).Encode(response)
  }

}

func DeleteAppointment(w http.ResponseWriter, r *http.Request) {
  //OBJEK RESPONSE
  Println("Endpoint Hit: DeleteAppointment")
  var response models.AppointmentResponse

  vars := mux.Vars(r)

  status := method.Remove(vars["id"])

  if status == false {
    message := "Appointment Failed to Update"
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(400)
    response.Status   = false
    response.Message  = message
    // response.Data = result
    json.NewEncoder(w).Encode(response)
  } else {
    message := "Appointment Deleted Succesfully"
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(301)
    response.Status   = true
    response.Message  = message
    json.NewEncoder(w).Encode(response)
  }

}

func UpdateAppointment(w http.ResponseWriter, r *http.Request) {
  Println("Endpoint Hit: UpdateAppointment")

  vars := mux.Vars(r) //ambil isi dari w,isinya parameter dari endpoint


  var appointment models.Appointment
  var response models.AppointmentResponse
  // json.Unmarshal(reqBody, &appointment)
  appointment.Doctor_name = r.FormValue("doctor_name")
  appointment.Pet_name = r.FormValue("pet_name")
  appointment.Complaint = r.FormValue("complaint")
  appointment.Time_Appointment = r.FormValue("time")
  // vars := mux.Vars(r) //ambil isi dari w,isinya parameter dari endpoint
  // id := "art"+vars["id"] //ambil id dari endpoint
  //
  // reqBody, _ := ioutil.ReadAll(r.Body)
  // var appointment models.Appointment
  // json.Unmarshal(reqBody, &appointment)

  status := method.Update(
      vars["id"],
      // appointment.Time_Appointment,
      appointment.Doctor_name,
      // appointment.Pet_owner_name,
      appointment.Pet_name,
      appointment.Complaint,
      appointment.Time_Appointment)

  if status == false {
    message := "Appointment Failed to Update"
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(400)
    response.Status   = false
    response.Message  = message
    // response.Data = result
    json.NewEncoder(w).Encode(response)
  } else {
    message := "Appointment Succesfully Updated"
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(202)
    response.Status   = true
    response.Message  = message
    // response.Data = result
    json.NewEncoder(w).Encode(response)
  }


}
