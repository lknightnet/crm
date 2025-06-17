package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

type PostgreSQL struct {
	DB *gorm.DB
}

func NewConnection(nameDatabase string, attempt, maxAttempt int, models ...any) (*PostgreSQL, error) {
	if attempt >= maxAttempt {
		return nil, fmt.Errorf("attempt exceeds max attempt: %d", attempt)
	}
	lite := &PostgreSQL{}
	db, err := gorm.Open(postgres.Open(nameDatabase), &gorm.Config{
		AllowGlobalUpdate: true,
	})
	if err != nil {
		log.Printf("Attempt %d: Failed to connect to database. Error: %v\n", attempt, err)
		time.Sleep(2 * time.Second)
		return NewConnection(nameDatabase, attempt+1, maxAttempt, models...)
	}

	for _, v := range models {
		err := db.AutoMigrate(v)
		if err != nil {
			log.Println(err)
		}
	}

	lite.DB = db
	return lite, nil
}
