package handler

import (
	models "Halovet/models"
	method "Halovet/repository/forum"
	"encoding/json"
	. "fmt"
	"log"
	"net/http"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

// CreateTopic : membuat topic baru
func CreateTopic(w http.ResponseWriter, r *http.Request) {
	Println("Endpoint Hit: CreateTopic")

	var Topic models.ForumTopic
	var result []models.ForumTopic
	var response models.ForumTopicResponse

	//AMBIL DATA UNTUK INSERT DATABASE

	Topic.Title = r.FormValue("Title")
	Topic.Content = r.FormValue("Content")

	//get id dari database
	CategoryID, _ := method.GetCategoryID(r.FormValue("Category"))

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	user := userInfo["User"]
	userReal, _ := user.(map[string]interface{})
	Topic.Author = Sprintf("%v", userReal["Name"])
	Topic.AuthorID, _ = strconv.Atoi(Sprintf("%v", userReal["ID"]))

	realResult, err := method.InsertTopic(
		Topic.Title,
		Topic.Author,
		Topic.AuthorID,
		Topic.Content,
		CategoryID,
	)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		response.Status = false
		response.Message = "Failed to Post New Topic"
		json.NewEncoder(w).Encode(response)
	} else {
		result = append(result, realResult)

		data := map[string]interface{}{
			"ForumTopic": result,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		response.Status = true
		response.Message = "Succesfully Post New Topic"
		response.Data = data
		json.NewEncoder(w).Encode(response)

		s := Sprintf("%s Succesfully Created New Topic", Topic.Author)
		log.Println(s)
	}
}

// GetTopicByID : get topic dengan id tertentu
func GetTopic(w http.ResponseWriter, r *http.Request) {
	Println("Endpoint Hit: GetTopic")

	// var topic models.ForumTopic
	var result []models.ForumTopic
	var response models.Response

	vars := mux.Vars(r)
	topicResult, status := method.FindTopicbyID(vars["topicid"])

	if status == false {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		response.Status = false
		response.Message = "Failed to Get Topic"
		json.NewEncoder(w).Encode(response)
	} else {
		result = append(result, topicResult)

		data := map[string]interface{}{
			"ForumTopic": result,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(202)
		response.Status = true
		response.Message = "Successfully Get Topic"
		response.Data = data
		json.NewEncoder(w).Encode(response)
	}
}

// UpdateTopicByID : update topic dengan id tertentu
func UpdateTopic(w http.ResponseWriter, r *http.Request) {
	Println("Endpoint Hit: UpdateTopicByID")

	vars := mux.Vars(r)

	var Topic models.ForumTopic
	// var result []models.ForumTopic
	var response models.Response

	Topic.Title = r.FormValue("Title")
	Topic.Content = r.FormValue("Content")
	//get id dari database
	CategoryID, _ := method.GetCategoryID(r.FormValue("Category"))

	status := method.UpdateTopic(
		vars["topicid"],
		Topic.Title,
		Topic.Content,
		CategoryID,
	)

	if status == false {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		response.Status = false
		response.Message = "Failed to Update Topic"
		json.NewEncoder(w).Encode(response)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		response.Status = true
		response.Message = "Succesfully to Update Topic"
		json.NewEncoder(w).Encode(response)
	}
}

// DeleteTopicByID : delete topic dengan id tertentu
func DeleteTopic(w http.ResponseWriter, r *http.Request) {
	Println("Endpoint Hit: DeleteTopicByID")
	var response models.Response
	vars := mux.Vars(r)

	status := method.DeleteTopic(vars["topicid"])

	if status == false {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		response.Status = false
		response.Message = "Failed to Delete Topic"
		json.NewEncoder(w).Encode(response)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		response.Status = true
		response.Message = "Succesfully Delete Topic"
		// response.Data = result
		json.NewEncoder(w).Encode(response)
	}
}

func ReplyTopic(w http.ResponseWriter, r *http.Request) {
	Println("Endpoint Hit: ReplyTopic")

	vars := mux.Vars(r)
	_, status := method.FindTopicbyID(vars["topicid"])

	if status == false {
		// log.Println("Forum Topic not found")
		//PRINT JSON FORUM ID TIDAK DITEMUKAN
		json.NewEncoder(w).Encode("Forum Topic not found")
		return
	}

	var Reply models.ForumReply
	var result []models.ForumReply
	var response models.Response

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	user := userInfo["User"]
	userReal, _ := user.(map[string]interface{})
	Reply.Author = Sprintf("%v", userReal["Name"])
	Reply.AuthorID, _ = strconv.Atoi(Sprintf("%v", userReal["ID"]))

	Reply.Content = r.FormValue("Content")

	realResult, err := method.InsertReply(
		vars["topicid"],
		Reply.Author,
		Reply.AuthorID,
		Reply.Content,
	)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		response.Status = false
		response.Message = "Failed to Post Reply"
		json.NewEncoder(w).Encode(response)
	} else {

		result = append(result, realResult)

		data := map[string]interface{}{
			"ForumReply": result,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		response.Status = true
		response.Message = "Succesfully to Post Reply"
		response.Data = data
		json.NewEncoder(w).Encode(response)
	}
}

func DeleteReply(w http.ResponseWriter, r *http.Request) {
	Println("Endpoint Hit: DeleteReply")

	vars := mux.Vars(r)

	_, status := method.FindTopicbyID(vars["topicid"])

	if status == false {
		// log.Println("Forum Topic not found")
		//PRINT JSON FORUM ID TIDAK DITEMUKAN
		json.NewEncoder(w).Encode("Forum Topic not found")
		return
	}

	var response models.Response
	err := method.DeleteReply(
		vars["topicid"],
		vars["replyid"],
	)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		response.Status = false
		response.Message = "Failed to Delete Reply"
		json.NewEncoder(w).Encode(response)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		response.Status = true
		response.Message = "Succesfully Delete Reply"
		json.NewEncoder(w).Encode(response)

	}

}

func UpdateReply(w http.ResponseWriter, r *http.Request) {
	Println("Endpoint Hit: UpdateReply")
	vars := mux.Vars(r)

	_, status := method.FindTopicbyID(vars["topicid"])

	if status == false {
		// log.Println("Forum Topic not found")
		//PRINT JSON FORUM ID TIDAK DITEMUKAN
		json.NewEncoder(w).Encode("Forum Topic not found")
		return
	}

	var Reply models.ForumReply
	// var result []models.ForumReply
	var response models.Response

	Reply.Content = r.FormValue("Content")

	err := method.UpdateReply(
		vars["topicid"],
		vars["replyid"],
		Reply.Content,
	)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		response.Status = false
		response.Message = "Failed to Update Reply"
		json.NewEncoder(w).Encode(response)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		response.Status = true
		response.Message = "Succesfully to Update Reply"
		json.NewEncoder(w).Encode(response)
	}

}

func GetReply(w http.ResponseWriter, r *http.Request) {
	Println("Endpoint Hit: GetReply")
	var result []models.ForumReply
	var response models.Response

	vars := mux.Vars(r)
	_, status := method.FindTopicbyID(vars["topicid"])

	if status == false {
		// log.Println("Forum Topic not found")
		//PRINT JSON FORUM ID TIDAK DITEMUKAN
		json.NewEncoder(w).Encode("Forum Topic not found")
		return
	}

	replyResult, err := method.FindReply(
		vars["topicid"],
		vars["replyid"])

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		response.Status = false
		response.Message = "Failed to Get Reply"
		json.NewEncoder(w).Encode(response)
	} else {
		result = append(result, replyResult)

		data := map[string]interface{}{
			"ForumReply": result,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(202)
		response.Status = true
		response.Message = "Successfully Get Reply"
		response.Data = data
		json.NewEncoder(w).Encode(response)
	}

}

// ListTopicByUserID :
// func ListTopicByUserID(w http.ResponseWriter, r *http.Request) {
// 	Println("Endpoint Hit: ListTopicByUserID")
// 	vars := mux.Vars(r)
// 	//DAPETIN ID USER
// 	status := method.FindAllTopicbyID(vars["id"])
// 	//GET SEMUA TOPIC YANG AUTHORNYA == USER ID

// }

// // ListTopicByCategory :
// func ListTopicByCategory(w http.ResponseWriter, r *http.Request) {
// 	Println("Endpoint Hit: ListTopicByCategory")
// 	vars := mux.Vars(r)
// 	status := method.FindTopicbyCategory(vars["id"])
// }
