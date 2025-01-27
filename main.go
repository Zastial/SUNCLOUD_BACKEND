package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"SUNCLOUD_BACKEND/model"
	"SUNCLOUD_BACKEND/routes"
)

type File struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Path      string    `gorm:"not null;unique" json:"path"`
	Size      int64     `gorm:"not null" json:"size"`
	Type      string    `gorm:"not null" json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func InitDatabase() error {
	db, err := gorm.Open(sqlite.Open("cloud.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	return db.AutoMigrate(&File{})
}

func NewDatabase() (*model.Database, error) {
	if _, err := os.Stat("cloud.db"); os.IsNotExist(err) {
		log.Println("Initialisation de la base de données...")
		if err := InitDatabase(); err != nil {
			log.Fatal("Erreur lors de l'initialisation de la DB:", err)
		}
	}

	db, err := gorm.Open(sqlite.Open("cloud.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&File{})
	if err != nil {
		return nil, err
	}

	return &model.Database{DB: db}, nil
}

func main() {
	db, err := NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	FileHandler := &routes.FileHandler{DB: db}

	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/files", FileHandler.GetAllFiles)
		api.GET("/files/:id", FileHandler.GetFile)
		api.POST("/files", FileHandler.CreateFile)
		api.DELETE("/files/:id", FileHandler.DeleteFile)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
