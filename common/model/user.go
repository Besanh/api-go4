package model

import "time"

type User struct {
	Id        string `json:"id"	gorm:"primaryKey;column:id;type:int(11) AUTO_INCREMENT;not null"`
	Username  string `json:"username" gorm:"unique;column:username;type:varchar(20);not null"`
	Level     string `json:"level" gorm:"column:level;type:int(4);not null"`
	Password  string `json:"password" gorm:"column:password;type:varchar(255);not null"`
	ApiKey    string `json:"api_key" gorm:"column:api_key;type:varchar(25);null"`
	CreatedAt string `json:"created_at" gorm:"column:created_at;type:TIMESTAMP"`
	UpdatedAt string `json:"updated_at" gorm:"column:updated_at;type:TIMESTAMP"`
}

func (User) TableName() string {
	return "vicidial_api_user"
}

// Type data tra ve
type UserAuthRes struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	ApiKey   string `json:"api_key"`
	Level    string `json:"level"`
}

// Nhan data tu body goi api
type AccessToken struct {
	ClientID     string    `json:"client_id"`
	UserID       string    `json:"user_id"`
	Token        string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ExpiredTime  int       `json:"expired_at"`
	Scope        string    `json:"scope"`
	TokenType    string    `json:"token_type"`
	JWT          string    `json:"jwt"`
}
