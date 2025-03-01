package db

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func Connect() *pg.DB {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	db := pg.Connect(&pg.Options{
		Addr:     os.Getenv("DB_ADDR"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME_ACCOUNTS"),
	})

	err = createSchema(db)
	if err != nil {
		log.Println("Error creating schema:", err)
	}

	return db
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{} {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
