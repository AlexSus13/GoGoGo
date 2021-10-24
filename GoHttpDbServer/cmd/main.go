package main

import (
	"GoHttpDbServer/app"
	"GoHttpDbServer/config"
	"GoHttpDbServer/database"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	//"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)   //обработка ошибок
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime) //обработка ошибок

	//initialization of configuration files
	err := initConfig()
	if err != nil {
		log.Fatal(err) //error handling
	}

	//loading environment variables from an '.env' file
	err = godotenv.Load()
	if err != nil {
		log.Fatal(err) //error handling
	}

	config := config.GetConf()

	db, err := database.NewPostgresDB(database.Config{
		User:     config.db.user,
		Host:     config.db.host,
		Password: os.Getenv("DB_PASSWORD"),
		Port:     config.db.port,
		DBName:   config.db.dbname,
		SSLMode:  config.db.sslmode,
	})
	if err != nil {
		log.Fatal(err) //error handling
	}

	app := app.NewApp(db, errorLog, infoLog, config)

	r := mux.NewRouter()

	r.HandleFunc("/", app.Hello).Methods("GET")
	r.Handle("/getfiles/{filename}", app.HandlerGET).Methods("GET")
	r.Handle("/postfile", app.HandlerPOST).Methods("POST")
	r.Handle("/getfile/{filename}", app.HandlerFILE).Methods("GET")

	srv := &http.Server{
		Addr:         ":" + config.port,
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	infoLog.Printf("starting server at :8080 port") //error handling
	errorLog.Fatal(srv.ListenAndServe())            //error handling

}

func initConfig() error {
	viper.AddConfigPath("/home/ubuntu/zooProject/GoGoGo/GoHttpDbServer/etc")
	viper.SetConfigName("etc")
	return viper.ReadInConfig()
}
