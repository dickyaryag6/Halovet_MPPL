package handler

import (
	models "Halovet/models"
	method "Halovet/repository/article"
	"encoding/json"
	"fmt"
	. "fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

func GetAllArticle(w http.ResponseWriter, r *http.Request) {
	Println("Endpoint Hit: GetAllArticle")

	var response models.Response
	querymap := r.URL.Query()
	limitstart := querymap["limitstart"][0]
	// Printf("%T\n", limitstart)
	limit := querymap["limit"][0]

	realResult, rowcount, err := method.FindAllArticles(limitstart, limit)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		response.Status = false
		response.Message = "Failed to Get Articles"
		json.NewEncoder(w).Encode(response)
	} else {

		data := map[string]interface{}{
			"Articles":  realResult,
			"Row_Count": rowcount,
		}

		w.Header().Set("Content-Type", "application/json")
		message := "Articles Get Succesfully"
		w.WriteHeader(302)
		response.Status = true
		response.Message = message
		response.Data = data
		json.NewEncoder(w).Encode(response)
	}
}

func savePhoto() {

}

func CreateArticle(w http.ResponseWriter, r *http.Request) {
	Println("Endpoint Hit: CreateArticle")
	var article models.Article
	var result []models.Article
	var response models.Response

	article.Title = r.FormValue("title")
	article.Content = r.FormValue("content")
	if len(article.Title) == 0 || len(article.Content) == 0 {
		json.NewEncoder(w).Encode("Title atau Content tidak boleh Kosong")

	}
	uploadedFile, handler, err := r.FormFile("photo")
	if err != nil {
		// Println(err.Error())
		json.NewEncoder(w).Encode(err.Error())
	}
	defer uploadedFile.Close()

	dir, err := os.Getwd()
	// dir == folder Project
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
	}

	timeNow := Sprintf(time.Now().Format("2006-01-02"))
	article.PhotoPath = fmt.Sprintf("%s-%s%s",
		article.Title,
		timeNow,
		filepath.Ext(handler.Filename))

	fileLocation := filepath.Join(dir, "public/articlephotos", article.PhotoPath)
	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())

	}

	defer targetFile.Close()
	// Printf("%T\n", article.PhotoPath)
	if _, err := io.Copy(targetFile, uploadedFile); err != nil {
		json.NewEncoder(w).Encode(err.Error())

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
		article.AuthorID,
		article.PhotoPath,
	)
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
		w.WriteHeader(200)
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
		w.WriteHeader(302)
		response.Status = true
		response.Message = "Succesfully Get Article"
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

	article.Title = r.FormValue("title")
	article.Content = r.FormValue("content")

	result := method.UpdateArticle(
		vars["articleid"],
		article.Title,
		article.Content,
	)

	if result == false {
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
		w.WriteHeader(301)
		response.Status = true
		response.Message = "Article Succesfully Delete"
		json.NewEncoder(w).Encode(response)
	}

}
