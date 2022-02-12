package app

import (
	"FileStorageServer/database"
	"FileStorageServer/filesoperation"

	"github.com/sirupsen/logrus"

	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (app *App) Hello(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello, this is the start page of the server, for information on interacting with the server see the README.md"))
}

func (app *App) ListFileHeaders(w http.ResponseWriter, r *http.Request) {

	var FN string

	param := r.URL.Query()

	if _, ok := param["filename"]; ok {
		FN = param.Get("filename")
	}

	headersSlice, err := database.ListFilesHeaders(app.Db, FN)
	if err != nil {
		//app.MyLogger.Fatal(err)
		app.MyLogger.WithFields(logrus.Fields{
			"func":    "database.ListFilesHeaders",
			"package": "app",
		}).Info(err)
		http.Error(w, "Error when accessing the database", 500)
		return //error handling
	}

	headersSliceJSON, err := json.Marshal(headersSlice)
	if err != nil {
		app.MyLogger.WithFields(logrus.Fields{
			"func":    "json.Marshal",
			"package": "app",
		}).Info(err)
		http.Error(w, http.StatusText(500), 500)
		return //error handling
	}
	w.Write(headersSliceJSON)

}

func (app *App) GetFileAndHeaders(w http.ResponseWriter, r *http.Request) {

	var FN string

	param := r.URL.Query()

	if _, ok := param["filename"]; ok {
		FN = param.Get("filename")
	} else {
		http.Error(w, "Request parameter not specified", 400)
		return //error handling
	}

	headers, err := database.GetFileHeaders(app.Db, FN)
	if err != nil {
		app.MyLogger.WithFields(logrus.Fields{
			"func":    "database.GetFileHeaders",
			"package": "app",
		}).Info(err)
		http.Error(w, "Error when accessing the database", 500)
		return //error handling
	}

	File, err := filesoperation.Get(app.Config.PathToFile, FN)
	if err != nil {
		app.MyLogger.WithFields(logrus.Fields{
			"func":    "filesoperation.Get",
			"package": "app",
		}).Info(err)
		http.Error(w, "Error when interacting with files", 500)
		return //error handling
	}

	w.Header().Add("Content-Type", headers.ContentType)
	w.Header().Add("Content-Length", headers.ContentLength)

	w.Write(File)
}

func (app *App) SaveFileAndHeaders(w http.ResponseWriter, r *http.Request) {

	FileContents, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.MyLogger.WithFields(logrus.Fields{
			"func":    "ReadAll(r.Body)",
			"package": "app",
		}).Info(err)
		http.Error(w, "Error when reading the request body", 500)
		return //error handling
	}
	defer r.Body.Close()

	FN := r.URL.Query().Get("filename")
	CT := r.Header.Get("Content-Type")
	CL := r.Header.Get("Content-Length")

	if FN == "" || CT == "" || CL == "" {
		app.MyLogger.Info("NO HEADERS OR FILE NAME IN RESPONSE")
		http.Error(w, "Invalid request format, missing headers or file name", 400)
		return //error handling
	}

	flag, err := database.CheckFileByName(app.Db, FN) //If flag == true => a file with this name exists
	if err != nil {
		app.MyLogger.WithFields(logrus.Fields{
			"func":    "database.CheckFileByName",
			"package": "app",
		}).Info(err)
		http.Error(w, "Error when accessing the database", 500)
		return //error handling
	}

	if *flag == true {

		err := filesoperation.DeleteOld(app.Config.PathToFile, FN)
		if err != nil {
			app.MyLogger.WithFields(logrus.Fields{
				"func":    "filesoperation.DeleteOld",
				"package": "app",
			}).Info(err)
			http.Error(w, "Error when interacting with files", 500)
			return //error handling
		}

		err = filesoperation.SaveNew(app.Config.PathToFile, FN, FileContents)
		if err != nil {
			app.MyLogger.WithFields(logrus.Fields{
				"func":    "filesoperation.SaveNew",
				"package": "app",
			}).Info(err)
			http.Error(w, "Error when interacting with files", 500)
			return //error handling
		}

		err = database.UpdateTable(app.Db, CT, CL, FN)
		if err != nil {
			app.MyLogger.WithFields(logrus.Fields{
				"func":    "database.UpdateTable",
				"package": "app",
			}).Info(err)
			http.Error(w, "Error when accessing the database", 500)
			return //error handling
		}

	} else {

		err = filesoperation.SaveNew(app.Config.PathToFile, FN, FileContents)
		if err != nil {
			app.MyLogger.WithFields(logrus.Fields{
				"func":    "filesoperation.SaveNew",
				"package": "app",
			}).Info(err)
			http.Error(w, "Error when interacting with files", 500)
			return //error handling
		}

		err := database.PostFileHeaders(app.Db, FN, CT, CL)
		if err != nil {
			app.MyLogger.WithFields(logrus.Fields{
				"func":    "database.PostFileHeaders",
				"package": "app",
			}).Info(err)
			http.Error(w, "Error when accessing the database", 500)
			return //error handling
		}
	}

	w.WriteHeader(201)
	response := "SUCCESSFUL SAVING OF THE FILE AND INFORMATION IN THE DATABASE!"
	w.Write([]byte(response))
}
