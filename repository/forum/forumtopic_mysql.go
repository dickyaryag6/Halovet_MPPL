package repository

import (
	dbCon "Halovet/driver"
	models "Halovet/models"
	"database/sql"
	"errors"
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

// CRUD TOPIC FORUM
func FindTopicbyID(topicid string) (models.ForumTopic, bool) {
	//CEK APAKAH ID USER YG LOGIN SAMA DENGAN YANG DI DATABASE

	var Topic models.ForumTopic

	var Reply models.ForumReply
	var Replies []models.ForumReply
	var CategoryID int
	// Print(topicid)
	id, err := strconv.Atoi(topicid)
	// Println(id)
	if err != nil {
		Println("format ID salah")
		return Topic, false
	}

	//QUERY TOPIC
	sqlStatement := "select * from forum_topic where id = ?"
	err = db.QueryRow(sqlStatement, id).
		Scan(&Topic.TopicID,
			&CategoryID,
			&Topic.Title,
			&Topic.Author,
			&Topic.Content,
			&Topic.CreatedAt,
			&Topic.UpdatedAt,
			&Topic.AuthorID)
	if err != nil {
		Println(err.Error())
		return Topic, false
	}

	Topic.Category, _ = FindCategory(CategoryID)

	// QUERY REPLIES
	sqlStatement = "select * from forum_reply where topic_id = ?"
	results, err := db.Query(sqlStatement, topicid)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	for results.Next() {
		err = results.Scan(&Reply.ReplyID,
			&Reply.TopicID,
			&Reply.Author,
			&Reply.AuthorID,
			&Reply.Content,
			&Reply.CreatedAt,
			&Reply.UpdatedAt)
		if err != nil {
			panic(err.Error())
		} else {
			Replies = append(Replies, Reply)
		}

	}
	Topic.Replies = Replies
	return Topic, true

}

func InsertTopic(title string, author string, authorid int, content string, categoryid int) (models.ForumTopic, error) {
	var Topic models.ForumTopic
	// var Category models.ForumCategory

	sqlStatement := "insert into forum_topic (title, author, author_id, content, category_id) values (?,?,?,?,?)"
	row, err := db.Exec(sqlStatement, title, author, authorid, content, categoryid)

	if err != nil {
		log.Println(err.Error())
		return Topic, errors.New("Insert Gagal")
	}
	id, err := row.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		return Topic, errors.New("Insert Gagal")
	}

	Topic.TopicID = id

	Topic.Category, err = FindCategory(categoryid)
	if err != nil {
		return Topic, errors.New("Kategori tidak ditemukan")
	}

	Topic.Title = title
	Topic.Author = author
	Topic.AuthorID = authorid
	Topic.Content = content

	t := time.Now()
	Topic.CreatedAt = Sprintf(t.Format("2006-01-02 15:04:05"))
	Topic.UpdatedAt = Sprintf(t.Format("2006-01-02 15:04:05"))

	return Topic, nil

}

func UpdateTopic(id string, title string, content string, categoryid int) bool {
	appointmentid, err := strconv.Atoi(id)
	if err != nil {
		Println("format ID salah")
		return false
	}

	timeNow := Sprintf(time.Now().Format("2006-01-02 15:04:05"))

	sqlStatement := "update forum_topic set title = ?, content = ?, category_id = ?, updated_at = ? where id = ?"
	_, err = db.Exec(sqlStatement, title, content, categoryid, timeNow, appointmentid)

	if err != nil {
		Println(err.Error())
		return false
	}

	return true
}

func DeleteTopic(id string) bool {
	topicid, err := strconv.Atoi(id)

	if err != nil {
		Println("format ID salah")
		return false
	}

	sqlStatement := "delete from forum_topic where id = ?"
	_, err = db.Exec(sqlStatement, topicid)

	if err != nil {
		Println(err.Error())
		return false
	}
	return true
}

// ---------------------------------------------------------------------

// CRUD CATEGORY
func GetCategoryID(CategoryTitle string) (int, bool) {
	sqlStatement := "select id from forum_category where category_title = ?"
	var ID int

	err = db.QueryRow(sqlStatement, CategoryTitle).
		Scan(&ID)

	if err != nil {
		Println(err.Error())
		return ID, false
	}

	return ID, true

}

func FindCategory(categoryid int) (models.ForumCategory, error) {

	var Category models.ForumCategory

	sqlStatement := "select * from forum_category where id = ?"
	err = db.QueryRow(sqlStatement, categoryid).
		Scan(&Category.ID, &Category.CategoryTitle, &Category.CategoryDesc)
	if err != nil {
		return Category, errors.New("Kategori Tidak ada")
	}
	return Category, nil
}

//-------------------------------------------------------------------

// CRUD TOPIC REPLY
func FindReply(topicid string, replyid string) (models.ForumReply, error) {
	var Reply models.ForumReply

	realtopicid, err := strconv.Atoi(topicid)
	if err != nil {
		Println("format ID salah")
		return Reply, errors.New("Get Reply Gagal")
	}
	realreplyid, err := strconv.Atoi(replyid)
	if err != nil {
		Println("format ID salah")
		return Reply, errors.New("Get Reply Gagal")
	}

	// var TopicID int

	sqlStatement := "select * from forum_reply where id = ?"
	err = db.QueryRow(sqlStatement, realreplyid).
		Scan(&Reply.ReplyID,
			&Reply.TopicID,
			&Reply.Author,
			&Reply.AuthorID,
			&Reply.Content,
			&Reply.CreatedAt,
			&Reply.UpdatedAt)
	if err != nil {
		Println(err.Error())
		return Reply, err
	}
	if realtopicid != Reply.TopicID {
		Println(err.Error())
		return Reply, errors.New("Get Reply Gagal")
	}
	return Reply, nil

}

func InsertReply(id string, author string, authorid int, content string) (models.ForumReply, error) {
	var Reply models.ForumReply

	realtopicid, err := strconv.Atoi(id)

	if err != nil {
		Println("format ID salah")
		return Reply, errors.New("Insert Reply Gagal")
	}

	sqlStatement := "insert into forum_reply (topic_id, author, author_id, content) values (?,?,?,?)"
	row, err := db.Exec(sqlStatement, realtopicid, author, authorid, content)

	if err != nil {
		log.Println(err.Error())
		return Reply, errors.New("Insert Gagal")
	}
	replyid, err := row.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		return Reply, errors.New("Insert Gagal")
	}

	Reply.ReplyID = replyid

	// Reply.Topic, _ = FindTopicbyID(id)
	Reply.TopicID = realtopicid
	Reply.Author = author
	Reply.AuthorID = authorid
	Reply.Content = content

	t := time.Now()
	Reply.CreatedAt = Sprintf(t.Format("2006-01-02 15:04:05"))
	Reply.UpdatedAt = Sprintf(t.Format("2006-01-02 15:04:05"))
	return Reply, nil
}

func UpdateReply(topicid string, replyid string, content string) error {
	// var Reply models.ForumReply

	// realtopicid, err := strconv.Atoi(topicid)
	// if err != nil {
	// 	Println("format ID salah")
	// 	return false
	// }
	realreplyid, err := strconv.Atoi(replyid)
	if err != nil {
		Println("format ID salah")
		return errors.New("Update Gagal")
	}

	//CEK APAKAH TOPICID ADA
	_, status := FindTopicbyID(topicid)
	if status == false {
		Println("Topic ID Salah")
		return errors.New("Update Gagal")
	}

	timeNow := Sprintf(time.Now().Format("2006-01-02 15:04:05"))

	sqlStatement := "update forum_reply set content = ?, updated_at = ? where id = ?"
	_, err = db.Exec(sqlStatement, content, timeNow, realreplyid)

	if err != nil {
		Println(err.Error())
		return errors.New("Insert Gagal")
	}

	return nil
}

func DeleteReply(topicid string, replyid string) error {
	// var Reply models.ForumReply

	// realtopicid, err := strconv.Atoi(topicid)
	// if err != nil {
	// 	Println("format ID salah")
	// 	return false
	// }
	realreplyid, err := strconv.Atoi(replyid)
	if err != nil {
		Println("format ID salah")
		return errors.New("Delete Gagal")
	}

	//CEK APAKAH TOPICID ADA
	_, status := FindTopicbyID(topicid)
	if status == false {
		Println("Topic ID Salah")
		return errors.New("Delete Gagal")
	}

	sqlStatement := "delete from forum_reply where id = ?"
	_, err = db.Exec(sqlStatement, realreplyid)

	if err != nil {
		Println(err.Error())
		return errors.New("Insert Gagal")
	}
	return nil
}
