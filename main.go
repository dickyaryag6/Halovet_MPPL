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
	account.HandleFunc("/register", handler.Register).Methods("POST")
	account.HandleFunc("/login", handler.Login).Methods("POST")
	account.HandleFunc("/logout", handler.Logout)

	// APPOINTMENT
	appointment := router.PathPrefix("/appointment").Subrouter()
	appointment.Use(mid.JWTAuthorization)
	appointment.HandleFunc("", handler.CreateAppointment).Methods("POST")
	appointment.HandleFunc("/{id}/uploadPayment", handler.UploadPayment).Methods("POST")
	// appointment.HandleFunc("", handler.GetAllAppointments).Methods("GET")
	appointment.HandleFunc("/{id}", handler.GetAppointmentByID).Methods("GET")
	appointment.HandleFunc("/{id}", handler.DeleteAppointment).Methods("DELETE")
	appointment.HandleFunc("/{id}", handler.UpdateAppointment).Methods("PUT")

	// FORUM
	// appointment := router.PathPrefix("/forum").Subrouter()
	// appointment.Use(mid.MiddlewareJWTAuthorization)
	// appointment.HandleFunc("", handler.CreateTopic).Methods("POST")
	// appointment.HandleFunc("{topicid}/reply", handler.ReplyTopic).Methods("POST")

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
