package transport

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/g3orge/FIOApi/internal/db"
	"github.com/g3orge/FIOApi/internal/model"
	"github.com/gorilla/mux"
)

func GetF(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	f, err := db.GetName(id)
	if err != nil {
		log.Println(err)
		b := "didnt find"
		w.Write([]byte(b))
		return
	}

	national := "https://api.nationalize.io/?name=" + f.Name
	gender := "https://api.genderize.io/?name=" + f.Name
	age := "https://api.agify.io/?name=" + f.Name

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
			f.Country = f2.Country[k].Countryid
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
	f.Age = f2.Age

	rGend, err := http.Get(gender)
	if err != nil {
		log.Println(err)
	}

	body, err = io.ReadAll(rGend.Body)
	if err != nil {
		log.Println(err)
	}
	json.Unmarshal(body, &f2)

	f.Gender = f2.Gender
	b, _ := json.Marshal(f)

	w.Write(b)
}

func UpdateF(w http.ResponseWriter, r *http.Request) {
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
	id := mux.Vars(r)["id"]
	db.DeleteName(id)
}

func AddF(w http.ResponseWriter, r *http.Request) {
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
