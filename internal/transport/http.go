package transport

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/g3orge/FIOApi/internal/db"
	"github.com/g3orge/FIOApi/internal/model"
	"github.com/gorilla/mux"
)

func GetF(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s\n", r.Method, r.URL.Path)

	queryName := r.URL.Query().Get("name")
	querySurname := r.URL.Query().Get("surname")
	if queryName == "" && querySurname == "" {
		log.Print("queryArgs are empty")
		return
	}

	// f, err := db.GetNames()
	// if err != nil {
	// 	log.Println(err)
	// 	b := "didnt find"
	// 	w.Write([]byte(b))
	// 	return
	// }

	datab := db.GetDB()
	query := datab.Table("names").Model(&model.OutF{})

	if queryName != "" {
		query = query.Where("name = ?", queryName)
	}

	if querySurname != "" {
		query = query.Where("name = ?", querySurname)
	}

	filteredF := model.OutF{}
	query.Find(&filteredF)

	national := "https://api.nationalize.io/?name=" + filteredF.Name
	gender := "https://api.genderize.io/?name=" + filteredF.Name
	age := "https://api.agify.io/?name=" + filteredF.Name

	rNat, err := http.Get(national)
	if err != nil {
		log.Println(err)
	}

	body, err := io.ReadAll(rNat.Body)
	if err != nil {
		log.Println(err)
	}

	var f2 model.F
	json.Unmarshal(body, &f2)

	var max float64
	max = 0
	for k, v := range f2.Country {
		if v.Probability > max {
			max = v.Probability
			filteredF.Country = f2.Country[k].Countryid
		}
	}

	rAge, err := http.Get(age)
	if err != nil {
		log.Println(err)
	}

	body, err = io.ReadAll(rAge.Body)
	if err != nil {
		log.Println(err)
	}
	json.Unmarshal(body, &f2)
	filteredF.Age = f2.Age

	rGend, err := http.Get(gender)
	if err != nil {
		log.Println(err)
	}

	body, err = io.ReadAll(rGend.Body)
	if err != nil {
		log.Println(err)
	}
	json.Unmarshal(body, &f2)

	filteredF.Gender = f2.Gender

	db.GetDB().Table("names").Where("name = ?", filteredF.Name).Save(&filteredF)

	json.NewEncoder(w).Encode(filteredF)
}

func UpdateF(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s\n", r.Method, r.URL.Path)

	id := mux.Vars(r)["id"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	var f model.OutF
	err = json.Unmarshal(body, &f)
	if err != nil {
		log.Println("in unmarsh json ", err)
	}
	fmt.Println(f)
	db.UpdateName(&f, id)
}

func DeleteF(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s\n", r.Method, r.URL.Path)

	id := mux.Vars(r)["id"]
	db.DeleteName(id)
}

func AddF(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s\n", r.Method, r.URL.Path)

	var f model.OutF

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("in body json ", err)
	}

	err = json.Unmarshal(body, &f)
	if err != nil {
		log.Println("in unmarsh json ", err)
	}
	db.AddName(&f)

	b, _ := json.Marshal(f)
	w.Write(b)
}

func GetAll(w http.ResponseWriter, r *http.Request) { //пагинация
	log.Printf("Received request: %s %s\n", r.Method, r.URL.Path)

	queryVal := r.URL.Query()
	page := strToInt(queryVal.Get("page"))
	limit := strToInt(queryVal.Get("limit"))

	f, err := db.GetNames()
	if err != nil {
		log.Println(err)
	}

	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = len(f)
	}

	startInd := (page - 1)
	endInd := startInd + limit

	if startInd >= len(f) {
		json.NewEncoder(w).Encode(f)
		return
	}
	if endInd > len(f) {
		endInd = len(f)
	}

	pageName := f[startInd:endInd]

	json.NewEncoder(w).Encode(pageName)
}

func strToInt(str string) int {
	result, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return result
}
