package main
//IIIIIII GdEEEEEEEEE
import (
	//"fmt"
	"net/http"
	"log"
	"encoding/json"
	"github.com/gorilla/mux"
	"os"
)

type Animal struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

var animals []Animal

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)
	m := mux.NewRouter()
	m.HandleFunc("/zoo", GetZoo).Methods("GET")
	m.HandleFunc("/zoo", PostZoo).Methods("POST")
	m.HandleFunc("/zoo/{type}", GetAnimal).Methods("GET")
        m.HandleFunc("/zoo/{type}", PutAnimal).Methods("PUT")
	m.HandleFunc("/zoo/{type}", DelAnimal).Methods("DELETE")
	infoLog.Printf("Server 37.139.43.30 on :8080 port START")
	err := http.ListenAndServe(":8080", m)
	errorLog.Fatal(err)
}

func GetZoo(w http.ResponseWriter, r *http.Request) {
	infoLog := log.New(os.Stdout, "INFO\t",  log.Ldate|log.Ltime)
	w.Header().Set("Content-Type", "application/json")
	if len(animals) == 0 {
		infoLog.Printf("Zoo is empty,But you can add new animals")
		return
	}
	json.NewEncoder(w).Encode(animals)
}

func PostZoo(w http.ResponseWriter, r *http.Request) {
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)
	var animal Animal
	err := json.NewDecoder(r.Body).Decode(&animal)
	defer r.Body.Close()
	if err != nil {
		errorLog.Fatal(err)
	}
	animals = append(animals, animal)
	json.NewEncoder(w).Encode(animal)
}

func GetAnimal(w http.ResponseWriter, r *http.Request) {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, value := range animals {
		if value.Type == params["type"] {
			json.NewEncoder(w).Encode(value)
			return
		}
	}
	json.NewEncoder(w).Encode(&Animal{})
        infoLog.Printf("The zoo does not have this animal, but you can add it")
}

func PutAnimal(w http.ResponseWriter, r *http.Request) {
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)
        w.Header().Set("Content-Type", "application/json")
        params := mux.Vars(r)
        for key, value := range animals {
	        if value.Type == params["type"] {
			animals = append(animals[:key], animals[key+1:]...)
			var animal Animal
			err := json.NewDecoder(r.Body).Decode(&animal)
			defer r.Body.Close()
			if err != nil {
				errorLog.Fatal(err)
			}
			animals = append(animals, animal)
			json.NewEncoder(w).Encode(animals)
                        return
                }
        }
}

func DelAnimal(w http.ResponseWriter, r *http.Request) {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	w.Header().Set("Content-Type","application/json")
	params := mux.Vars(r)
	for key, value := range animals {
		if value.Type == params["type"] {
			animals = append(animals[:key], animals[key+1:]...)
			json.NewEncoder(w).Encode(animals)
			return
		}
	}
	infoLog.Printf("This animal is not in the zoo, see for yourself")
	json.NewEncoder(w).Encode(animals)
}
