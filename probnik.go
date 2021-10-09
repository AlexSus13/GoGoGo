package main

import (
	"database/sql"
	_"github.com/lib/pq"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"fmt"
	"log"
	"os"
	"io"
	"time"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

type headers struct {
	filename       string
	content_type   string
	content_lenght string
	id             int
}

type DB struct {
	name *sql.DB
}

type DBtwo struct {
	name *sql.DB
}

type DBthree struct {
	name *sql.DB
}

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)   //обработка ошибок
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime) //обработка ошибок

	db, err := sql.Open("postgres", "user=alex host=localhost password=12345 port=5432 dbname=dbforgolang sslmode=disable")
	if err != nil {
		errorLog.Fatal(err) //обработка ошибок
		return
	}
	//defer db.Close()

	errr := db.Ping()
	if errr != nil {
		errorLog.Fatal(errr) //обработка ошибок
		return
	}

	r := mux.NewRouter()

	DbHandlerGET := &DB{name: db}
	DbHandlerPOST := &DBtwo{name: db}
	DbHandlerFILE := &DBthree{name: db}

	r.HandleFunc("/", Hello).Methods("GET")
	r.Handle("/getfiles", DbHandlerGET).Methods("GET")
	r.Handle("/postfile", DbHandlerPOST).Methods("POST")
	r.Handle("/getfile/{filename}", DbHandlerFILE).Methods("GET")

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	infoLog.Printf("starting server at :8080 port") //обработка ошибок
	errorLog.Fatal(srv.ListenAndServe())            //обработка ошибок
}

func Hello(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello, this is the start page of the server"))
}

func (p *DB) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)
	db := p.name

	rows, err := db.Query("SELECT * FROM headers")
	if err != nil {
		errorLog.Fatal(err) //обработка ошибок
	}
	defer rows.Close()

	headersSlice := []headers{}

	for rows.Next() {
		h := headers{}
		err := rows.Scan(&h.filename, &h.content_type, &h.content_lenght, &h.id)
		if err != nil {
			errorLog.Fatal(err) //обработка ошибок
		}
		headersSlice = append(headersSlice, h)
	}
	for _, h := range headersSlice {
		result := fmt.Sprintf("filename=%s, content_type=%s, content_lenght=%s, id=%d\n", h.filename, h.content_type, h.content_lenght, h.id)
		w.Write([]byte(result))
	}
}

func (p *DBthree) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)
	db := p.name

	//ПОЛУЧЕНИЕ ИМЕНИ ФАЙЛА ИЗ УРЛА
	params := mux.Vars(r)
	FN := params["filename"]

	rows, err := db.Query("SELECT * FROM headers WHERE filename LIKE '%' || $1 || '%'", FN)
	if err != nil {
		errorLog.Fatal(err) //обработка ошибок
	}
	defer rows.Close()

	headersSlice := []headers{}

	for rows.Next() {
		h := headers{}
		err := rows.Scan(&h.filename, &h.content_type, &h.content_lenght, &h.id)
		if err != nil {
			errorLog.Fatal(err) //обработка ошибок
		}
		headersSlice = append(headersSlice, h)
	}
	for _, h := range headersSlice { //исправить на отправку файлов во всех get запросах
		result := fmt.Sprintf("filename=%s, content_type=%s, content_lenght=%s, id=%d\n", h.filename, h.content_type, h.content_lenght, h.id)
		w.Write([]byte(result))
	}
}

func (p *DBtwo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)
	db := p.name

	FN, CT, CL, err := FileUpload(r)
	if err != nil {
		errorLog.Fatal(err) //обработка ошибок
	}

	if FN == "" || CT == "" || CL == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	h := headers{}
	err = db.QueryRow("SELECT filename FROM headers WHERE filename=$1", FN).Scan(&h.filename)
	if err != nil {
		if err == sql.ErrNoRows { //значит строка с таким значением не найдена
			Insetr(db, w, FN, CT, CL)
		} else {
			errorLog.Fatal(err) //обработка ошибок
		}
	}

	//ЕСЛИ СОВПАДЕНИ ПЕРЕЗАТИРАЕМ
	_, err = db.Exec("UPDATE headers SET content_type=$1, content_lenght=$2 WHERE filename=$3", CT, CL, FN)
	if err != nil {
		errorLog.Fatal(err) //обработка ошибок
	}

	w.Write([]byte("Sending was successful"))
}

func Insetr(db *sql.DB, w http.ResponseWriter, FN string, CT string, CL string){
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)
	_, err := db.Exec("INSERT INTO headers (filename, content_type, content_lenght) VALUES ($1, $2, $3)", FN, CT, CL)
	if err != nil {
		errorLog.Fatal(err) //обработка ошибок
	}

	w.Write([]byte("Sending was successful"))
}

func FileUpload(r *http.Request) (string, string, string, error) {
        errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime) //обработка ошибок

	r.ParseMultipartForm(32 << 20)

	file, handler, err := r.FormFile("fileforgolang")
	if err != nil {
		errorLog.Fatal(err) //обработка ошибок
		//return "", err
	}
	defer file.Close()

	f, err := os.OpenFile("/home/ubuntu/infile/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		errorLog.Fatal(err) //обработка ошибок
		//return "", err
	}
	defer f.Close()

	io.Copy(f, file)

	ct := r.Header.Get("Content-Type")
	if r.Header.Get("Content-Type") == "" { //нужна ли эта проверка
		ct = handler.Header.Get("Content-Type")
	}

	cl := strconv.FormatInt(handler.Size, 10)
	return handler.Filename, ct, cl, nil
}
