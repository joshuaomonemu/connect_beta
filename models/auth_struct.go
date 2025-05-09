package models

type AuthResp struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Token   string      `json:"token"`
}

type Role string

const (
	RoleCreator  Role = "creator"
	RoleConsumer Role = "consumer"
)

type User struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Email       string `json:"email" gorm:"uniqueIndex"`
	Fullname    string `json:"fullname"`
	Password    string `json:"-"` // hashed
	DisplayName string `json:"display_name"`
	Role        Role   `json:"role" gorm:"type:varchar(20)"`
}

type Photo struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	FileURL  string `json:"file_url"`
	UserID   string `json:"user_id"`
	Uploaded string `json:"uploaded,omitempty"` // if you have a timestamp column
}
