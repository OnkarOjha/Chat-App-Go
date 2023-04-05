package Models

import (
	"github.com/golang-jwt/jwt/v4"
	// socketio "github.com/googollee/go-socket.io"
	"github.com/jinzhu/gorm"
)

// User Information
type User struct {
	Token           string `json:"token"`
	User_Id         string `json:"user_id" gorm:"default:uuid_generate_v4();"` //PK
	Name            string `json:"name"`
	Phone           string `json:"phone"`
	Email           string `json:"email"`
	Join_date       string `json:"join_date"`
	Profile_picture string `json:"profile_picture"`
	Is_active       bool   `json:"is_active"`
	Bio             string `json:"bio"`
	Is_deleted bool `json:"is_deleted"`

}

// Room Topic Information
type Topic struct {
	Topic_id    string `json:"topic_id" gorm:"default:uuid_generate_v4();unique;primaryKey"` //PK
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Chat Room Information
type Room struct {
	gorm.Model
	Room_id    string `json:"room_id" gorm:"default:uuid_generate_v4();"` //PK
	Admin_id    string `json:"admin_id"`
	Name       string `json:"name"`
	Created_at string `json:"created_at"`
	Topic_id   string `json:"topic_id"` //  hatana pdega
	Topic_name string `json:"topic_name"`
	User_count int    `json:"user_count"`
	Is_deleted bool `json:"is_deleted"`
}

//Message InformationVerifictaion failed
type Message struct {
	gorm.Model
	Message_id   string `json:"message_id" gorm:"default:uuid_generate_v4();unique;primaryKey"` //PK
	User_id      string `json:"user_id"`
	Room_id      string `json:"room_id"`
	Text         string `json:"text"`
	Message_type string `json:"message_type"`
}

// Participant Information
type Participant struct {
	gorm.Model
	P_Id        string `json:"id" gorm:"default:uuid_generate_v4();unique;primaryKey"` //PK
	User_id   string `json:"user_id"`
	Room_id   string `json:"room_id"`
	Room_name string `json:"room_name"`
	Is_deleted bool `json:"is_deleted"`
}

type Claims struct {
	User_id string `json:"user_id"`
	Phone string `json:"phone"`
	jwt.RegisteredClaims
}

