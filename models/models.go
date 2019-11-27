package models

type Appointment struct {
	AppointmentID int64
	// Time_Posted         time.Time
	Time_Appointment string
	Doctor_name      string
	Pet_Owner_Name   string
	Pet_owner_id     int
	Pet_Type         string
	IsPaid           bool
	Complaint        string
	CreatedAt        string
	UpdatedAt        string
}

type Article struct {
	ID        int64
	Title     string
	Author    string
	AuthorID  int
	Content   string
	CreatedAt string
	UpdatedAt string
	PhotoPath string
}

type Account struct {
	ID       int
	Email    string
	Name     string
	Password string
	Role     int //1:PetOwner, 2:Doctor, 3:Admin
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

type ForumCategory struct {
	ID            int
	CategoryTitle string
	CategoryDesc  string
}

type ForumTopic struct {
	TopicID   int64
	Category  ForumCategory //apa category dari forum ini
	Title     string        //judul dari topic ini
	Author    string        //yang ngepost topic
	AuthorID  int
	Content   string //pertanyaan topic
	CreatedAt string
	UpdatedAt string
	// views       int
	Replies []ForumReply
}

type ForumReply struct {
	ReplyID int64
	// CategoryID  	int //apa category dari forum ini
	TopicID   int    //apa topic dari forum ini
	Author    string //yang ngepost reply
	AuthorID  int
	Content   string //isi dari reply
	CreatedAt string //kapan post reply
	UpdatedAt string
}

type ForumTopicResponse struct {
	Status  bool
	Message string
	Data    map[string]interface{}
}

type Response struct {
	Status  bool
	Message string
	Data    map[string]interface{}
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
