package models

import (
	"time"
)

type Appointment struct {
	AppointmentID int64
	// Time_Posted         time.Time
	Time_Appointment string
	Doctor_name      string
	Pet_owner_name   string
	Pet_name         string
	Complaint        string
}

type Account struct {
	Email    string
	Name     string
	Password string
	Role     int //1:PetOwner, 2:Doctor
}

type Pet_Owner struct {
	ID       int64
	Email    string
	Name     string
	Password string
}

// type Admin struct {
// }

type Docter struct {
	ID       int64
	Email    string
	Name     string
	Password string
}

type Forum_Category struct {
	ID             int
	Category_title string
	Category_Desc  string
}

type Forum_Topic struct {
	ID          int
	Category    Forum_Category //apa category dari forum ini
	Title       string         //judul dari topic ini
	Author      Pet_Owner      //yang ngepost topic
	Content     string         //pertanyaan topic
	Time_Posted time.Time
	// views       int
}

type Forum_Reply struct {
	ID          int
	Category    Forum_Category //apa category dari forum ini
	forum_topic Forum_Topic    //apa topic dari forum ini
	Author      string         //yang ngepost reply
	Comment     string         //isi dari reply
	Time_Posted time.Time      //kapan post reply
}

type AppointmentResponse struct {
	Status  bool
	Message string
	Data    map[string]interface{}
}

type PetOwnerResponse struct {
	Status  bool
	Message string
	Data    map[string]interface{}
}

// func DBMigrate(db *gorm.DB) *gorm.DB {
//
// }
