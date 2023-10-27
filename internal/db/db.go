package db

import (
	"errors"
	"fmt"
	"log"

	"github.com/g3orge/FIOApi/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

func InitPostgres() {
	dsn := "user=postgres password=root dbname=Names port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.Table("names").AutoMigrate(&model.OutF{})
	insertF := &[]model.OutF{{Name: "Dmitry", Surname: "Ushakov", Patronymic: "Ivanovich"},
		{Name: "Alexey", Surname: "Udakov", Patronymic: "Olegovich"},
		{Name: "Irina", Surname: "Ivanova", Patronymic: "Ivanovna"}}

	db.Table("names").Create(insertF)

	readF := &model.OutF{}
	db.Table("names").First(&readF, "name = ?", "Irina")
	fmt.Println(readF.Name, readF.Surname)
}

func GetName(name string) (*model.OutF, error) {
	var f model.OutF
	res := db.Table("names").First(&f, "name = ?", name)
	if res.RowsAffected == 0 {
		return nil, errors.New(fmt.Sprintf("cannot find %s", name))
	}

	return &f, nil
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
