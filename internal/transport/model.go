package transport

type F struct {
	Name       string   `json:"name" gorm:"type:text"`
	Surname    string   `gorm:"type:text"`
	Patronymic string   `gorm:"type:text"`
	Age        int      `json:"age" gorm:"type:integer"`
	Gender     string   `json:"gender" gorm:"type:text"`
	Country    []string `json:"country" gorm:"type:text"`
}
