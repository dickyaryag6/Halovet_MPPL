package repository

import (
	dbCon "Halovet/driver"
	models "Halovet/models"
	"database/sql"
	. "fmt"
	"log"
	"strconv"
	"time"
)

var err error
var db *sql.DB

func init() {
	// KONEK KE DATABASE
	db, err = dbCon.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func InsertArticle(title string, content string, author string, authorid int) (models.Article, error) {
	var Article models.Article

	sqlStatement := "insert into articles (title, content, author, author_id) values (?,?,?,?)"
	row, err := db.Exec(sqlStatement, title, content, author, authorid)
	if err != nil {
		Print(err.Error())
		return Article, nil
	}

	// Println(title)
	// Println(content)

	id, err := row.LastInsertId()
	if err != nil {
		log.Fatal(err.Error())
	}

	Article.ID = id
	Article.Title = title
	Article.Author = author
	Article.AuthorID = authorid
	Article.Content = content
	Article.CreatedAt = Sprintf(time.Now().Format("2006-01-02 15:04:05"))
	Article.UpdatedAt = Sprintf(time.Now().Format("2006-01-02 15:04:05"))

	return Article, nil
}

func FindAllArticles(limitstart string, limit string) ([]models.Article, error) {
	var Article models.Article
	var Articles []models.Article

	realLimitStart, err := strconv.Atoi(limitstart)
	if err != nil {
		Println("format limit salah")
		return Articles, err
	}
	realLimit, err := strconv.Atoi(limit)
	if err != nil {
		Println("format limit salah")
		return Articles, err
	}

	sqlStatement := "select * from articles order by created_at limit ?, ?"
	results, err := db.Query(sqlStatement, realLimitStart, realLimit)
	if err != nil {
		panic(err.Error())
		return Articles, err
	}
	for results.Next() {
		err = results.Scan(&Article.ID,
			&Article.Title,
			&Article.Author,
			&Article.AuthorID,
			&Article.Content,
			&Article.CreatedAt,
			&Article.UpdatedAt)
		if err != nil {
			panic(err.Error())
		} else {
			Articles = append(Articles, Article)
		}

	}

	return Articles, nil
}

func FindArticle(articleid string) (models.Article, error) {
	var Article models.Article
	id, err := strconv.Atoi(articleid)
	if err != nil {
		Println("format ID salah")
	}

	sqlStatement := "select * from articles where id = ?"
	err = db.QueryRow(sqlStatement, id).
		Scan(&Article.ID,
			&Article.Title,
			&Article.Author,
			&Article.AuthorID,
			&Article.Content,
			&Article.CreatedAt,
			&Article.UpdatedAt)
	if err != nil {
		Println(err.Error())
		return Article, err
	}

	return Article, nil
}

func UpdateArticle(articleid string, title string, content string) bool {
	id, err := strconv.Atoi(articleid)
	if err != nil {
		Println("format ID salah")
		return false
	}

	sqlfind := "select count(*) from articles where id = ?"
	var i int
	err = db.QueryRow(sqlfind, id).
		Scan(&i)
	// Println(i)
	if i == 0 {
		// Println(err.Error())
		return false
	} else {
		// Println("ha")
		timeNow := Sprintf(time.Now().Format("2006-01-02 15:04:05"))
		sqlStatement := "update articles set title = ?, content = ?, updated_at = ? where id = ?"
		_, err := db.Exec(sqlStatement, title, content, timeNow, id)
		if err != nil {
			Println(err.Error())
			return false
		}
		return true
	}

}

func DeleteArticle(articleid string) error {
	id, err := strconv.Atoi(articleid)
	if err != nil {
		Println("format ID salah")
		return err
	}

	sqlStatement := "delete from articles where id = ?"
	_, err = db.Exec(sqlStatement, id)
	if err != nil {
		Println(err.Error())
		return err
	}
	return nil

}
