package db

import (
	"fmt"
	"log"
	"stream-video/pkg/config"
	"stream-video/pkg/domain"


	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", cfg.Db_host, cfg.Db_username, cfg.Db_password, cfg.Db_name, cfg.Db_port)
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	db.AutoMigrate(&domain.Video{})
	return db, err
}
