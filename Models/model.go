package Models

import (
	"github.com/golang-jwt/jwt/v4"
	// socketio "github.com/googollee/go-socket.io"
	"github.com/jinzhu/gorm"
	
)

// User Information
type User struct {
	Token           string `json:"token"`
	User_Id         string `json:"userId" gorm:"default:uuid_generate_v4();"` //PK
	Name            string `json:"name" validate:"required"`
	Phone           string `json:"phone" validate:"required,e164"`
	Email           string `json:"email" validate:"required,email"`
	Join_date       string `json:"joinDate"`
	Profile_picture string `json:"profilePicture"`
	Is_active       bool   `json:"isActive"`
	Bio             string `json:"bio" validate:"required"`
	Is_deleted      bool   `json:"isDeleted"`
	Is_verified     bool   `json:"isVerified"`
}

// Room Topic Information
type Topic struct {
	Topic_id    string `json:"topicId" gorm:"default:uuid_generate_v4();unique;primaryKey"` //PK
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

// Chat Room Information
type Room struct {
	gorm.Model
	Room_id    string `json:"roomId" gorm:"default:uuid_generate_v4();"` //PK
	Admin_id   string `json:"adminId"`
	Name       string `json:"name" validate:"required"`
	Created_at string `json:"createdAt"`
	Topic_id   string `json:"topicId"` //  hatana pdega
	Topic_name string `json:"topicName" validate:"required"`
	User_count int64    `json:"userCount"`
	Is_deleted bool   `json:"isDeleted"`
}

//Message InformationVerifictaion failed
type Message struct {
	gorm.Model
	Message_id   string `json:"messageId" gorm:"default:uuid_generate_v4();unique;primaryKey"` //PK
	User_id      string `json:"userId"`
	Room_id      string `json:"roomId"`
	Text         string `json:"text" validate:"required"`
	Message_type string `json:"messageType"`
}

// Participant Information
type Participant struct {
	gorm.Model
	P_Id       string `json:"id" gorm:"default:uuid_generate_v4();unique;primaryKey"` //PK
	User_id    string `json:"userId"`
	Room_id    string `json:"roomId"`
	Room_name  string `json:"roomName"`
	Is_deleted bool   `json:"isDeleted"`
}

type Claims struct {
	User_id string `json:"userId"`
	Phone   string `json:"phone"`
	jwt.RegisteredClaims
}
