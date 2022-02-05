package main

import (
	"FileStorageServer/app"
	"FileStorageServer/config"
	"FileStorageServer/database"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	MyLogger := logrus.New()

	MyLogger.Formatter = &logrus.JSONFormatter{
		DisableTimestamp: false,
		PrettyPrint:      true,
	}

	config, err := config.Get()
	if err != nil {
		MyLogger.WithFields(logrus.Fields{
			"func":    "config.Get",
			"package": "main",
		}).Fatal(err) //error handling
	}

	db, err := database.NewPostgresDB(database.Config{
		User:     config.DB.User,
		Host:     config.DB.Host,
		Password: config.DB.Password, //os.Getenv("DB_PASSWORD"),
		Port:     config.DB.Port,
		DBName:   config.DB.Dbname,
		SSLMode:  config.DB.Sslmode,
	})
	if err != nil {
		MyLogger.WithFields(logrus.Fields{
			"func":    "database.NewPostgresDB",
			"package": "main",
		}).Fatal(err) //error handling
	}

	app := app.NewApp(db, MyLogger, config)

	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api").Subrouter()

	subrouter.Use(app.CheckAuthMiddleware)

	subrouter.HandleFunc("/", app.Hello).Methods("GET")
	subrouter.HandleFunc("/listfiles", app.ListFileHeaders).Methods("GET")
	subrouter.HandleFunc("/postfile", app.SaveFileAndHeaders).Methods("POST")
	subrouter.HandleFunc("/getfile", app.GetFileAndHeaders).Methods("GET")

	router.HandleFunc("/reg", app.SignUp).Methods("POST")
	router.HandleFunc("/auth", app.SignIn).Methods("POST")

	MWrouter := app.LogMiddleware(router)

	srv := &http.Server{
		Addr:         config.Host + ":" + config.Port,
		Handler:      MWrouter,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		err := srv.ListenAndServe()
		switch err {
		case http.ErrServerClosed:
			MyLogger.Info("Server at :8080 port Stopped") //error handling
		default:
			MyLogger.WithFields(logrus.Fields{
				"func":    "srv.ListenAndServe",
				"package": "main",
			}).Fatal(err) //error handling
		}
	}()

	MyLogger.Info("Server at :8080 port Start") //error handling

	signalChanel := make(chan os.Signal, 1)

	signal.Notify(signalChanel, syscall.SIGTERM, syscall.SIGINT) //os.interrupt==syscall.SIGINT?

	<-signalChanel

	MyLogger.Info("server at :8080 port Shutting Down") //error handling

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	err = srv.Shutdown(ctx)
	if err != nil {
		MyLogger.WithFields(logrus.Fields{
			"func":    "srv.Shutdown",
			"package": "main",
		}).Fatal(err) //error handling
	}

	err = db.Close()
	if err != nil {
		MyLogger.WithFields(logrus.Fields{
			"func":    "db.Close",
			"package": "main",
		}).Fatal(err) //error handling
	}

	cancel()
}
