package main

import (
  . "fmt"
  "log"
  "net/http"
  "github.com/gorilla/mux"
  handler "Halovet/handler/http"
  mid "Halovet/middleware"
  "github.com/rs/cors"

)

func index(w http.ResponseWriter, r *http.Request) {
  Println("Endpoint Hit: homepage")
  Fprintf(w, "Welcome to Homepage")
}

func handleRequest() {

  router := mux.NewRouter().StrictSlash(true)
  router.HandleFunc("/", index)

  account := router.PathPrefix("/account").Subrouter()


  account.HandleFunc("/register", handler.Register).Methods("POST")
  account.HandleFunc("/login", handler.Login).Methods("POST")
  account.HandleFunc("/logout", handler.Logout)

  appointment := router.PathPrefix("/appointment").Subrouter()
  appointment.Use(mid.MiddlewareJWTAuthorization)
  appointment.HandleFunc("", handler.CreateAppointment).Methods("POST")
  // appointment.HandleFunc("", handler.GetAllAppointments).Methods("GET")
  appointment.HandleFunc("/{id}", handler.GetAppointmentByID).Methods("GET")
  appointment.HandleFunc("/{id}", handler.DeleteAppointment).Methods("DELETE")
  appointment.HandleFunc("/{id}", handler.UpdateAppointment).Methods("PUT")

  // appointment := router.PathPrefix("/forum").Subrouter()
  // appointment.Use(mid.MiddlewareJWTAuthorization)
  // appointment.HandleFunc("", handler.CreateTopic).Methods("POST")
  // appointment.HandleFunc("{topicid}/reply", handler.ReplyTopic).Methods("POST")

  c := cors.New(cors.Options{
  AllowedOrigins  :   []string{"*"},
  AllowCredentials:   true,
  AllowedMethods  :   []string{"*"},
  })

  handler := c.Handler(router)

  Println("Server listen at :8000")
  log.Fatal(http.ListenAndServe(":8000", handler))

}

func main() {
  handleRequest()
}
