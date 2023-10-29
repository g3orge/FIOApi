package db

import (
	"errors"
	"fmt"
	"log"

	"github.com/g3orge/FIOApi/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB
var err error

func InitPostgres(db_user, db_pass, db_port, db_name, db_ssl string) {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=%s", db_user, db_pass, db_name, db_port, db_ssl)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal(err)
	}
	// db.Migrator().DropTable("names")

	db.Table("names").AutoMigrate(&model.OutF{})

	insertF := &[]model.OutF{{Name: "Dmitry", Surname: "Ushakov", Patronymic: "Ivanovich"},
		{Name: "Alexey", Surname: "Udakov", Patronymic: "Olegovich"},
		{Name: "Irina", Surname: "Ivanova", Patronymic: "Ivanovna"}}
	db.Table("names").Create(insertF)
}

func GetName(name string) (*model.OutF, error) {
	var f model.OutF
	res := db.Table("names").First(&f, "name = ?", name)
	if res.RowsAffected == 0 {
		return nil, errors.New(fmt.Sprintf("cannot find %s", name))
	}

	return &f, nil
}

func GetNames() ([]model.OutF, error) {
	var f []model.OutF
	res := db.Table("names").Find(&f)
	if res.RowsAffected == 0 {
		return nil, errors.New(fmt.Sprintf("some error"))
	}

	return f, nil
}

func AddName(f *model.OutF) {
	db.Table("names").Create(&f)
}

func DeleteName(name string) {
	var f model.F
	f.Name = name
	db.Table("names").Find(&f, "name = ?", name).Delete(&f)
}

func UpdateName(f *model.OutF, id string) {
	db.Table("names").Where("name = ?", id).Save(&f)
}

func GetDB() *gorm.DB {
	return db
}
