package Database

import (
	"fmt"
	models "main/Models"
	constant "main/Utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	fmt.Println("Connecting to database...")
	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", constant.Host, constant.Port, constant.User, constant.Password, constant.Dbname)
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		fmt.Println("Error connecting to database")
		return err
	}

	db.Exec("CREATE SCHEMA IF NOT EXISTS public")
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	db.AutoMigrate(&models.User{}, &models.Topic{}, &models.Room{}, &models.Participant{}, &models.Message{})
	DB = db
	fmt.Println("Successfully Connected to database")
	return nil

}
