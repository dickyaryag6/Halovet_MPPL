package main

import (
	handler "Halovet/handler/http"
	mid "Halovet/middleware"
	"fmt"

	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func index(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: homepage")
	fmt.Fprintf(w, "Welcome to Homepage")
}

func handleRequest() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)

	// ACCOUNT
	account := router.PathPrefix("/account").Subrouter()
	account.HandleFunc("/register", handler.Register).Methods("POST") //account/register
	account.HandleFunc("/login", handler.Login).Methods("POST")       //account/login
	account.HandleFunc("/logout", handler.Logout)                     //account/logout

	// APPOINTMENT
	appointment := router.PathPrefix("/appointment").Subrouter()
	appointment.Use(mid.JWTAuthorization, mid.PetOwner)
	//kalo mau apply lebih dari satu middleware tambahin aja di dalam kurung
	appointment.HandleFunc("", handler.CreateAppointment).Methods("POST")                //appointment
	appointment.HandleFunc("/{id}/uploadPayment", handler.UploadPayment).Methods("POST") //appointment/{id}/uploadPayment
	// appointment.HandleFunc("", handler.GetAllAppointments).Methods("GET")				//appointment
	appointment.HandleFunc("/{id}", handler.GetAppointmentByID).Methods("GET")   //appointment/{id}
	appointment.HandleFunc("/{id}", handler.DeleteAppointment).Methods("DELETE") //appointment/{id}
	appointment.HandleFunc("/{id}", handler.UpdateAppointment).Methods("PUT")    //appointment/{id}

	// FORUM
	forum := router.PathPrefix("/forum").Subrouter()
	forum.Use(mid.JWTAuthorization)
	// create topic forum
	forum.HandleFunc("", handler.CreateTopic).Methods("POST") //forum
	// get suatu topic forum
	forum.HandleFunc("/{topicid}", handler.GetTopic).Methods("GET") //forum/topicid
	// edit pertanyaan topic forum
	forum.HandleFunc("/{topicid}", handler.UpdateTopic).Methods("PUT") //forum/topicid
	// delete pertanyaan topic forum
	forum.HandleFunc("/{topicid}", handler.DeleteTopic).Methods("DELETE") //forum/topicid
	// reply suatu topic
	forum.HandleFunc("/{topicid}/reply", handler.ReplyTopic).Methods("POST") //forum/topicid/reply
	// hapus sebuah reply
	forum.HandleFunc("/{topicid}/reply/{replyid}", handler.DeleteReply).Methods("DELETE")
	// update sebuah reply
	forum.HandleFunc("/{topicid}/reply/{replyid}", handler.UpdateReply).Methods("PUT")
	// get sebuah reply
	forum.HandleFunc("/{topicid}/reply/{replyid}", handler.GetReply).Methods("GET")
	// list semua topic yg ditanyakan user
	// forum.HandleFunc("/{userid}", handler.ListTopicByUserID).Methods("GET") //forum/userid
	// list semua topic dengan kategori tertentu
	// forum.HandleFunc("/{category}", handler.ListTopicByCategory).Methods("GET")

	// ARTICLE

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"*"},
	})

	handler := c.Handler(router)

	log.Println("Server listen at :8000")
	log.Fatal(http.ListenAndServe(":8000", handler))

}

func main() {
	handleRequest()
}
