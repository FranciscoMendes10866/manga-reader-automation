package config

import (
	"github.com/FranciscoMendes10866/queues/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Database *gorm.DB
var DATABASE_URI string = "host=localhost user=pg password=pg dbname=dango port=5432 sslmode=disable"

func Connect() error {
	var err error

	Database, err = gorm.Open(postgres.Open(DATABASE_URI), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		panic(err)
	}

	Database.AutoMigrate(&entities.MangaEntity{}, &entities.ChapterEntity{}, &entities.PagesEntity{})

	return nil
}
