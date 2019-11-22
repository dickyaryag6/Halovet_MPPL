package handler

import (
	models "Halovet/models"
	method "Halovet/repository/article"
	"encoding/json"
	. "fmt"
	"net/http"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

func CreateArticle(w http.ResponseWriter, r *http.Request) {
	Println("Endpoint Hit: CreateArticle")
	var article models.Article
	var result []models.Article
	var response models.Response

	article.Title = r.FormValue("Title")
	article.Content = r.FormValue("Content")
	if len(article.Title) == 0 || len(article.Content) == 0 {
		json.NewEncoder(w).Encode("Title atau Content tidak boleh Kosong")
		return
	}

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	user := userInfo["User"]
	userReal, _ := user.(map[string]interface{})
	article.Author = Sprintf("%v", userReal["Name"])

	article.AuthorID, err = strconv.Atoi(Sprintf("%v", userReal["ID"]))
	if err != nil {
		Println("format ID salah")
	}

	realResult, err := method.InsertArticle(
		article.Title,
		article.Content,
		article.Author,
		article.AuthorID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		response.Status = false
		response.Message = "Failed to Create Article"
		json.NewEncoder(w).Encode(response)
	} else {

		result = append(result, realResult)

		data := map[string]interface{}{
			"Article": result,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		response.Status = true
		response.Message = "Succesfully Create Article"
		response.Data = data
		json.NewEncoder(w).Encode(response)
	}
}

func GetArticle(w http.ResponseWriter, r *http.Request) {
	Println("Endpoint Hit: GetArticle")

	// var article models.Article
	var result []models.Article
	var response models.Response

	vars := mux.Vars(r)

	realResult, err := method.FindArticle(vars["articleid"])
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		response.Status = false
		response.Message = "Failed to Get Article"
		json.NewEncoder(w).Encode(response)
	} else {
		result = append(result, realResult)

		data := map[string]interface{}{
			"Article": result,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		response.Status = true
		response.Message = "Succesfully Create Article"
		response.Data = data
		json.NewEncoder(w).Encode(response)

	}

}

func UpdateArticle(w http.ResponseWriter, r *http.Request) {
	Println("Endpoint Hit: UpdateArticle")

	vars := mux.Vars(r)

	var article models.Article
	// var result []models.Article
	var response models.Response

	article.Title = r.FormValue("Title")
	article.Content = r.FormValue("Content")

	err := method.UpdateArticle(
		vars["articleid"],
		article.Title,
		article.Content,
	)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		response.Status = false
		response.Message = "Article Failed to Update"
		json.NewEncoder(w).Encode(response)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(202)
		response.Status = true
		response.Message = "Article Succesfully Updated"
		// response.Data = result
		json.NewEncoder(w).Encode(response)
	}

}

func DeleteArticle(w http.ResponseWriter, r *http.Request) {
	Println("Endpoint Hit: DeleteArticle")

	var response models.Response

	vars := mux.Vars(r)
	err := method.DeleteArticle(vars["articleid"])
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		response.Status = false
		response.Message = "Article Failed to Delete"
		json.NewEncoder(w).Encode(response)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(202)
		response.Status = true
		response.Message = "Article Succesfully Delete"
		json.NewEncoder(w).Encode(response)
	}

}
