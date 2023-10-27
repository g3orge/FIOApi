package model

type F struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
	Age        int    `json:"age"`
	Gender     string `json:"gender"`
	Country    []Ctry `json:"country"`
}

type Ctry struct {
	Countryid   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

type OutF struct {
	Name       string `json:"name" gorm:"type:text"`
	Surname    string `json:"surname" gorm:"type:text"`
	Patronymic string `json:"patronymic,omitempty" gorm:"type:text"`
	Age        int    `json:"age" gorm:"type:integer"`
	Gender     string `json:"gender" gorm:"type:text"`
	Country    string `json:"country" gorm:"type:text"`
}
