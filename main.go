package main

import (
	handler "Halovet/handler/http"
	mid "Halovet/middleware"
	"fmt"

	// "path/filepath"

	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func index(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: homepage")
	fmt.Fprintf(w, "Welcome to Homepage")
}

// func servePic(w http.ResponseWriter, r *http.Request) {
// 	// querymap := r.URL.Query()
// 	// namadokter := querymap["doctor"][0]
// 	// petownername := querymap["petowner"][0]
// 	// timeappointment := querymap["time"][0]

// 	folder := ""
// 	fileserved := filepath.Join(folder, filename, ".jpg")
// 	http.ServeFile(w, r, fileserved)
// }

func handleRequest() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)

	//serve static file
	router.
		PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("public"))))

	// ACCOUNT
	account := router.PathPrefix("/account").Subrouter()
	account.HandleFunc("/register", handler.Register).Methods("POST") //account/register
	account.HandleFunc("/login", handler.Login).Methods("POST")       //account/login
	account.HandleFunc("/logout", handler.Logout)                     //account/logout

	// APPOINTMENT
	appointment := router.PathPrefix("/appointment").Subrouter()
	appointment.Use(mid.JWTAuthorization)
	//kalo mau apply lebih dari satu middleware tambahin aja di dalam kurung
	appointment.HandleFunc("", handler.CreateAppointment).Methods("POST")                //appointment
	appointment.HandleFunc("/{id}/uploadPayment", handler.UploadPayment).Methods("POST") //appointment/{id}/uploadPayment
	// appointment.HandleFunc("", handler.GetAllAppointments).Methods("GET")				//appointment
	appointment.HandleFunc("/{id}", handler.GetAppointmentByID).Methods("GET")   //appointment/{id}
	appointment.HandleFunc("/{id}", handler.DeleteAppointment).Methods("DELETE") //appointment/{id}
	appointment.HandleFunc("/{id}", handler.UpdateAppointment).Methods("PUT")    //appointment/{id}
	appointment.HandleFunc("/user/{userid}", handler.GetAppointmentByUserID).Methods("GET")

	// FORUM
	forum := router.PathPrefix("/forum").Subrouter()
	forums := router.PathPrefix("/forums").Subrouter()
	forum.Use(mid.JWTAuthorization)
	forums.Use(mid.JWTAuthorization)
	// create topic forum
	forum.HandleFunc("", handler.CreateTopic).Methods("POST") //forum
	// get suatu topic forum
	forum.HandleFunc("/{topicid}", handler.GetTopic).Methods("GET") //forum/topicid
	// edit pertanyaan topic forum
	forum.HandleFunc("/{topicid}", handler.UpdateTopic).Methods("PUT") //forum/topicid
	// delete pertanyaan topic forum
	forum.HandleFunc("/{topicid}", handler.DeleteTopic).Methods("DELETE") //forum/topicid
	// get topic forum by id
	forum.HandleFunc("/user/{userid}", handler.GetTopicByUserID).Methods("GET")
	//get all forum
	forums.HandleFunc("", handler.GetAllForum).Methods("GET")

	reply := router.PathPrefix("/forum").Subrouter()
	reply.Use(mid.JWTAuthorization)
	// reply suatu topic
	reply.HandleFunc("/{topicid}/reply", handler.ReplyTopic).Methods("POST") //forum/topicid/reply
	// hapus sebuah reply
	reply.HandleFunc("/{topicid}/reply/{replyid}", handler.DeleteReply).Methods("DELETE")
	// update sebuah reply
	reply.HandleFunc("/{topicid}/reply/{replyid}", handler.UpdateReply).Methods("PUT")
	// get sebuah reply
	reply.HandleFunc("/{topicid}/reply/{replyid}", handler.GetReply).Methods("GET")
	// list semua topic yg ditanyakan user
	// forum.HandleFunc("/{userid}", handler.ListTopicByUserID).Methods("GET") //forum/userid
	// list semua topic dengan kategori tertentu
	// forum.HandleFunc("/{category}", handler.ListTopicByCategory).Methods("GET")

	// ARTICLE
	// adminarticle := router.PathPrefix("/article").Subrouter()
	article := router.PathPrefix("/article").Subrouter()
	articles := router.PathPrefix("/articles").Subrouter()
	article.Use(mid.JWTAuthorization)
	articles.Use(mid.JWTAuthorization)
	article.HandleFunc("/{articleid}", handler.GetArticle).Methods("GET")
	// article.Handle("/static/{photopath}/", servePic)
	articles.HandleFunc("", handler.GetAllArticle).Methods("GET")

	//KHUSUS ADMIN
	admin := router.PathPrefix("/admin").Subrouter()
	admin.Use(mid.JWTAuthorization)
	admin.HandleFunc("/appointment", handler.GetAllAppointment).Methods("GET")
	//VALIDASI BUKTI PEMBAYARAN
	admin.HandleFunc("/appointment/{appointmentid}/validatePayment", handler.ValidatePay).Methods("PUT")
	// get all articles
	admin.HandleFunc("/article", handler.CreateArticle).Methods("POST")
	admin.HandleFunc("/article/{articleid}", handler.DeleteArticle).Methods("DELETE")
	admin.HandleFunc("/article/{articleid}", handler.UpdateArticle).Methods("PUT")

	// admin.HandleFunc("/forum", handler.GetAllForum).Methods("GET")

	//PENCARIAN

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
