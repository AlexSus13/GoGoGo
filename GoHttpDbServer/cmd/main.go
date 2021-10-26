package main

import (
	"GoHttpDbServer/database"
	"GoHttpDbServer/config"
	"GoHttpDbServer/app"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
	"os/signal"
	"net/http"
	"context"
	"syscall"
	"time"
	"log"
	"os"
)

func main() {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)   //error handling
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime) //error handling

	config := config.GetConf()

	db, err := database.NewPostgresDB(database.Config{
		User:     config.DB.User,
		Host:     config.DB.Host,
		Password: config.DB.Password,//os.Getenv("DB_PASSWORD"),
		Port:     config.DB.Port,
		DBName:   config.DB.Dbname,
		SSLMode:  config.DB.Sslmode,
	})
	if err != nil {
		errorLog.Fatal(err) //error handling
	}

	app := app.NewApp(db, errorLog, infoLog, config)

	r := mux.NewRouter()

	r.HandleFunc("/", app.Hello).Methods("GET")
	r.HandleFunc("/getfiles/{filename}", app.HandlerGET).Methods("GET")
	r.HandleFunc("/postfile", app.HandlerPOST).Methods("POST")
	r.HandleFunc("/getfile/{filename}", app.HandlerFILE).Methods("GET")

	srv := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func () {
		errorLog.Fatal(srv.ListenAndServe()) //error handling
	}()

	infoLog.Printf("starting server at :8080 port") //error handling

	signalChanel := make(chan os.Signal, 1)

	signal.Notify(signalChanel, syscall.SIGTERM, syscall.SIGINT)

	<- signalChanel

	infoLog.Printf("server at :8080 port Shutting Down") //error handling

	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)

	err = srv.Shutdown(ctx)
	if err != nil {
		errorLog.Fatal(err) //error handling
	}

	err = db.Close()
	if err != nil {
		errorLog.Fatal(err) //error handling
	}

	cancel()
}

