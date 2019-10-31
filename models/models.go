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

type Pet_Owner struct {
	ID       int64
	Email    string
	Name     string
	Password string
}

type Admin struct {
  
}
// type Docter struct {
//
// }

type Forum_Category struct {
	ID             int
	Category_title string
}

type Forum_Topic struct {
	ID          int
	Category    Forum_Category
	Title       string
	Author      Pet_Owner
	Content     string
	Time_Posted time.Time
	views       int
}

type Forum_Reply struct {
	ID          int
	Category    Forum_Category
	forum_topic Forum_Topic
	Author      string
	Comment     string
	Time_Posted time.Time
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
