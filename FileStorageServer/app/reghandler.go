package app

import (
	"FileStorageServer/database"
	"FileStorageServer/token"

	"github.com/sirupsen/logrus"

	"encoding/json"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"net/http"
	"strconv"
)

type UserData struct {
	Name     string      `json:"user_name"`
	Password interface{} `json:"user_password"` //if password != string?
}

var HashingPassword = func(Password interface{}, keyPassword string) string {
	PasswordString := fmt.Sprintf("%v", Password)
	crcH := crc32.ChecksumIEEE([]byte(PasswordString + keyPassword))
	hash := strconv.FormatUint(uint64(crcH), 10)
	return hash
}

func (app *App) SignUp(w http.ResponseWriter, r *http.Request) {

	UD := &UserData{}

	UserD, err := ioutil.ReadAll(r.Body) //What if the r.Body was EMPTY or it wasn't there
	if err != nil {
		app.MyLogger.WithFields(logrus.Fields{
			"func":    "ReadAll(r.Body) in SignUp",
			"package": "app",
		}).Info(err)
		http.Error(w, "Error when reading the request body", 500)
		return //error handling
	}
	defer r.Body.Close()

	err = json.Unmarshal(UserD, UD)
	if err != nil {
		app.MyLogger.WithFields(logrus.Fields{
			"func":    "Unmarshal(UserD, UD) in SignUp",
			"package": "app",
		}).Info(err)
		http.Error(w, "Error when Unmarshali r.Body", 500)
		return //error handling
	}

	PasswordString := fmt.Sprintf("%v", UD.Password)

	if UD.Name == "" || PasswordString == "" {
		app.MyLogger.Info("INCORRECT FILLING IN OF REGISTRATION DATA")
		http.Error(w, "Invalid request format, missing user Name or Password", 400)
		return //error handling
	}

	flag, err := database.CheckUserByName(app.Db, UD.Name) //If flag == true => The user is already registered
	if err != nil {
		app.MyLogger.WithFields(logrus.Fields{
			"func":    "database.CheckUserByName in SignUp",
			"package": "app",
		}).Info(err)
		http.Error(w, "Error when accessing the database", 500)
		return //error handling
	}

	if *flag == true {
		w.WriteHeader(400)
		response := "A USER WITH THIS NAME ALREADY EXISTS, COME UP WITH A NEW ONE!"
		w.Write([]byte(response))
		return //error handling
	} else {
		HashPassword := HashingPassword(UD.Password, app.Config.KeyPassword)

		err = database.AddUserInDB(app.Db, UD.Name, HashPassword)
		if err != nil {
			app.MyLogger.WithFields(logrus.Fields{
				"func":    "database.AddUserInDB in SignUp",
				"package": "app",
			}).Info(err)
			http.Error(w, "Error when accessing the database", 500)
			return //error handling
		}

		TokenString, err := token.CreateToken(UD.Name, app.Config.KeyToken)
		if err != nil {
			app.MyLogger.WithFields(logrus.Fields{
				"func":    "token.CreateToken in SignUp",
				"package": "app",
			}).Info(err)
			http.Error(w, "Error when creating a token", 500)
			return //error handling
		}

		//Token in body

		w.WriteHeader(201)
		err = json.NewEncoder(w).Encode(TokenString)
		if err != nil {
			app.MyLogger.WithFields(logrus.Fields{
				"func":    "json.NewEncoder(w).Encode(TokenString) in SignUp",
				"package": "app",
			}).Info(err)
			http.Error(w, http.StatusText(500), 500)
			return //error handling
		}
		//response := "THE USER HAS BEEN CREATED, REGISTRATION WAS SUCCESSFUL!"
		//w.Write([]byte(response))
	}
}

func (app *App) SignIn(w http.ResponseWriter, r *http.Request) {

	UD := &UserData{}

	UserD, err := ioutil.ReadAll(r.Body) //What if the r.Body was EMPTY or it wasn't there
	if err != nil {
		app.MyLogger.WithFields(logrus.Fields{
			"func":    "ReadAll(r.Body) in SignIn",
			"package": "app",
		}).Info(err)
		http.Error(w, "Error when reading the request body", 500)
		return //error handling
	}
	defer r.Body.Close()

	err = json.Unmarshal(UserD, UD)
	if err != nil {
		app.MyLogger.WithFields(logrus.Fields{
			"func":    "Unmarshal(UserD, UD) in SignIn",
			"package": "app",
		}).Info(err)
		http.Error(w, "Error when Unmarshali r.Body", 500)
		return //error handling
	}

	PasswordString := fmt.Sprintf("%v", UD.Password)

	if UD.Name == "" || PasswordString == "" {
		app.MyLogger.Info("INCORRECT FILLING IN OF REGISTRATION DATA")
		http.Error(w, "Invalid request format, missing user Name or Password", 400)
		return //error handling
	}

	flag, HashPasswordInDB, err := database.CheckUserByNameAndPassword(app.Db, UD.Name) //If flag == true => The user is already registered
	if err != nil {
		app.MyLogger.WithFields(logrus.Fields{
			"func":    "database.CheckUserByNameAndPassword in SignIn",
			"package": "app",
		}).Info(err)
		http.Error(w, "Error when accessing the database", 500)
		return //error handling
	}

	if *flag == true {
		HashPassword := HashingPassword(UD.Password, app.Config.KeyPassword)
		if HashPassword != HashPasswordInDB {
			http.Error(w, "INVALID PASSWORD", 400)
		}

		TokenString, err := token.CreateToken(UD.Name, app.Config.KeyToken)
		if err != nil {
			app.MyLogger.WithFields(logrus.Fields{
				"func":    "token.CreateToken in SignIn",
				"package": "app",
			}).Info(err)
			http.Error(w, "Error when creating a token", 500)
			return //error handling
		}

		//Token in body

		w.WriteHeader(201)
		err = json.NewEncoder(w).Encode(TokenString)
		if err != nil {
			app.MyLogger.WithFields(logrus.Fields{
				"func":    "json.NewEncoder(w).Encode(TokenString) in SignIn",
				"package": "app",
			}).Info(err)
			http.Error(w, http.StatusText(500), 500)
			return //error handling
		}
	} else {
		http.Error(w, "THE USER IS NOT REGISTERED", 401)
	}
}
