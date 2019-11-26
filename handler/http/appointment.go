package handler

import (
	// . "fmt"

	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	// "time"
	models "Halovet/models"
	method "Halovet/repository/appointment"
	. "fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

func GetAppointmentByUserID(w http.ResponseWriter, r *http.Request) {
	Println("Endpoint Hit: GetAppointmentByUserID")

	var response models.Response
	// var result []models.Appointment

	vars := mux.Vars(r)

	realResult, status := method.FindAppointmentbyUserID(vars["userid"])

	if status == false {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		response.Status = false
		response.Message = "Failed to Get Appointment"
		json.NewEncoder(w).Encode(response)
	} else {
		// result = append(result, realResult)
		// result = realResult

		data := map[string]interface{}{
			"Appointments": realResult,
		}

		w.Header().Set("Content-Type", "application/json")
		message := "Appointments Get Succesfully"
		w.WriteHeader(202)
		response.Status = true
		response.Message = message
		response.Data = data
		json.NewEncoder(w).Encode(response)

	}

}

//YANG DIKEMBALIKAN CUMA APPOINTMENT DENGAN USER ID TERTENTU
func GetAllAppointment(w http.ResponseWriter, r *http.Request) {
	Println("Endpoint Hit: GetAllAppointment")
	// Println("GET params were:", r.URL.Query())

	var response models.Response
	// var result []models.Appointment

	querymap := r.URL.Query()
	limitstart := querymap["limitstart"][0]
	// Printf("%T\n", limitstart)
	limit := querymap["limit"][0]
	// Println(limit)

	realResult, err := method.FindAllAppointment(limitstart, limit)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		response.Status = false
		response.Message = "Failed to Get Appointment"
		json.NewEncoder(w).Encode(response)
	} else {
		// result = append(result, realResult)

		data := map[string]interface{}{
			"Appointments": realResult,
		}

		w.Header().Set("Content-Type", "application/json")
		message := "Appointments Get Succesfully"
		w.WriteHeader(202)
		response.Status = true
		response.Message = message
		response.Data = data
		json.NewEncoder(w).Encode(response)
	}
}

// GetAppointmentByID : nyari appointment dengan id tertentu
func GetAppointmentByID(w http.ResponseWriter, r *http.Request) {
	Println("Endpoint Hit: GetAppointmentByID")
	var response models.Response
	var result []models.Appointment

	vars := mux.Vars(r)

	realResult, status := method.FindbyID(vars["id"])

	w.Header().Set("Content-Type", "application/json")
	if status == false {
		message := "Appointment Failed to Get"
		w.WriteHeader(400)
		response.Status = false
		response.Message = message
		json.NewEncoder(w).Encode(response)
	} else {
		result = append(result, realResult)

		data := map[string]interface{}{
			"Appointment": result,
		}

		message := "Appointment Get Succesfully"
		w.WriteHeader(202)
		response.Status = true
		response.Message = message
		response.Data = data
		json.NewEncoder(w).Encode(response)
	}
}

//CUMA BOLEH ADMIN
// func GetAllAppointments(w http.ResponseWriter, r *http.Request) {
//
// }

// CreateAppointment : fungsi createappointment
func CreateAppointment(w http.ResponseWriter, r *http.Request) {

	//INSERT

	Println("Endpoint Hit: createAppointment")
	// reqBody, _ := ioutil.ReadAll(r.Body)
	var appointment models.Appointment
	var result []models.Appointment
	var response models.Response
	// json.Unmarshal(reqBody, &appointment)
	appointment.Doctor_name = r.FormValue("doctor_name")
	appointment.Pet_Type = r.FormValue("pet_type")
	appointment.Complaint = r.FormValue("complaint")
	appointment.Time_Appointment = r.FormValue("time")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	user := userInfo["User"]
	userReal, _ := user.(map[string]interface{})
	appointment.Pet_Owner_Name = Sprintf("%v", userReal["Name"])
	appointment.Pet_owner_id, err = strconv.Atoi(Sprintf("%v", userReal["ID"]))
	if err != nil {
		Println("format ID salah")
		// return Topic, false
	}

	realResult, status := method.Insert(
		appointment.Time_Appointment,
		appointment.Doctor_name,
		appointment.Pet_Owner_Name,
		appointment.Pet_owner_id,
		appointment.Pet_Type,
		appointment.Complaint,
	)

	if status == false {
		message := "Appointment Failed to Create"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		response.Status = false
		response.Message = message
		// response.Data = result
		json.NewEncoder(w).Encode(response)
	} else {

		result = append(result, realResult)

		data := map[string]interface{}{
			"Appointment": result,
		}

		//KIRIM EMAIL

		message := "Appointment Booked Succesfully"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		response.Status = true
		response.Message = message
		response.Data = data
		json.NewEncoder(w).Encode(response)

		s := Sprintf("%s Succesfully Created New Appointment", appointment.Pet_Owner_Name)
		log.Println(s)
	}

}

// DeleteAppointment : fungsi DeleteAppointment
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
		response.Status = false
		response.Message = message
		// response.Data = result
		json.NewEncoder(w).Encode(response)
	} else {
		message := "Appointment Deleted Succesfully"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(301)
		response.Status = true
		response.Message = message
		json.NewEncoder(w).Encode(response)
	}

}

// UpdateAppointment : fungsi UpdateAppointment
func UpdateAppointment(w http.ResponseWriter, r *http.Request) {
	Println("Endpoint Hit: UpdateAppointment")

	vars := mux.Vars(r) //ambil isi dari w,isinya parameter dari endpoint

	var appointment models.Appointment
	var response models.Response
	// json.Unmarshal(reqBody, &appointment)
	appointment.Doctor_name = r.FormValue("doctor_name")
	appointment.Pet_Type = r.FormValue("pet_type")
	appointment.Complaint = r.FormValue("complaint")
	appointment.Time_Appointment = r.FormValue("time")
	// vars := mux.Vars(r) //ambil isi dari w,isinya parameter dari endpoint
	// id := "art"+vars["id"] //ambil id dari endpoint

	// reqBody, _ := ioutil.ReadAll(r.Body)
	// var appointment models.Appointment
	// json.Unmarshal(reqBody, &appointment)

	status := method.Update(
		vars["id"],
		// appointment.Time_Appointment,
		appointment.Doctor_name,
		// appointment.Pet_owner_name,
		appointment.Pet_Type,
		appointment.Complaint,
		appointment.Time_Appointment)

	if status == false {
		message := "Appointment Failed to Update"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		response.Status = false
		response.Message = message
		// response.Data = result
		json.NewEncoder(w).Encode(response)
	} else {
		message := "Appointment Succesfully Updated"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(202)
		response.Status = true
		response.Message = message
		// response.Data = result
		json.NewEncoder(w).Encode(response)
	}

}

func UploadPayment(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(1024); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	uploadedFile, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer uploadedFile.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	dir, err := os.Getwd()
	// dir == folder Project
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//DAPETIN ID APPOINTMENT
	var appointment models.Appointment
	vars := mux.Vars(r)
	appointmentid, err := strconv.Atoi(vars["id"])
	// Print(appointmentid)

	//NGAMBIL DATA PET OWNER, APPOINTMENT TIME, DAN NAMA DOKTER
	sqlStatement := "select pet_owner_name, appointment_time, doctor_name from appointment where id = ?"
	err = db.QueryRow(sqlStatement, appointmentid).
		Scan(&appointment.Pet_Owner_Name, &appointment.Time_Appointment, &appointment.Doctor_name)
	if err != nil {
		Println(err.Error())
	}
	Println(appointment.Pet_Owner_Name, appointment.Time_Appointment[0:10], appointment.Doctor_name)
	filename := fmt.Sprintf("%s-%s-%s%s",
		appointment.Pet_Owner_Name,
		appointment.Time_Appointment[0:10],
		appointment.Doctor_name,
		filepath.Ext(handler.Filename))

	fileLocation := filepath.Join(dir, "payment", filename)
	//hasil join Halovet/payment/namafile
	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer targetFile.Close()

	if _, err := io.Copy(targetFile, uploadedFile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Berhasil Upload Bukti Pembayaran"))
	return
}

func ValidatePay(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	appointmentid, err := strconv.Atoi(vars["appointmentid"])
	if err != nil {
		Println("format ID salah")
	}

	//CEK APAKAH SUDAH PEMBAYARAN SUDAH DILAKUKAN
	status := method.CheckValidation(appointmentid)
	if status == false {
		w.Write([]byte("Validasi Pembayaran sudah dilakukan"))
		return
	}

	status = method.ValidatePayment(appointmentid)
	if status == false {
		w.Write([]byte("Gagal Validasi Pembayaran"))
		return
	} else {
		w.Write([]byte("Berhasil Validasi Pembayaran"))
		return
	}

}
