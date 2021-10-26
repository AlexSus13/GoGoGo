package app

import (
	"GoHttpDbServer/database"
	"GoHttpDbServer/forfile"
	"github.com/gorilla/mux"
	"encoding/json"
	"net/http"
	"strconv"
)

func (a *App) Hello(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello, this is the start page of the server"))
}

func (a *App) HandlerGET(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	FN := params["filename"]

	headersSlice, err := database.GetAllDb(FN, a.Db)
	if err != nil {
		a.errorLog.Fatal(err) //error handling
	}

	headersSliceJson, err := json.Marshal(headersSlice)
	if err != nil {
		a.errorLog.Fatal(err) //error handling
	}

	//set headers and status response
	w.Write(headersSliceJson)

}

func (a *App) HandlerFILE(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	FN := params["filename"]

	headers, err := database.GetFileDb(FN, a.Db)
	if err != nil {
		a.errorLog.Fatal(err) //error handling
	}

	File, err := forfile.GetOpenFile(FN)
	if err != nil {
		a.errorLog.Fatal(err) //error handling
	}

	w.Header().Add("Content-Type", headers.Content_Type)
	w.Header().Add("Content-Length", headers.Content_Length)

	//set headers and status response
	w.Write(File)
}

func (a *App) HandlerPOST(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(32 << 20)

	file, handler, err := r.FormFile("fileforgolang") //config?//
	if err != nil {
		a.errorLog.Fatal(err) //error handling
	}
	defer file.Close()

	FN := handler.Filename
	CT := handler.Header.Get("Content-Type")
	CL := strconv.FormatInt(handler.Size, 10)

	if FN == "" || CT == "" || CL == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	flag, err := database.CheckFileInDb(FN, a.Db) //If flag == true => a file with this name exists
	if err != nil {
		a.errorLog.Fatal(err) //error handling
	}

	if *flag == true {

		err := forfile.DeleteOldFile(FN)
		if err != nil {
			a.errorLog.Fatal(err) //error handling
		}

		database.UpdateDb(CT, CL, FN, a.Db)

		err = forfile.SaveNewFile(FN, file)
		if err != nil {
			a.errorLog.Fatal(err) //error handling
		}

	} else {

		database.PostFileDb(FN, CT, CL, a.Db)

		err := forfile.SaveNewFile(FN, file)
		if err != nil {
			 a.errorLog.Fatal(err) //error handling
		}
	}

	response := "SUCCESSFUL SAVING OF THE FILE AND INFORMATION IN THE DATABASE!"
	w.Write([]byte(response))
}
