package database

import (
	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"database/sql"
)

type HeadersTable struct {
	FileName      string `json:"filename"`
	ContentType   string `json:"content_type"`
	ContentLength string `json:"content_lenght"`
	ID            int    `json:"id"`
}

type UserData struct {
	Name         string `json:"user_name"`
	HashPassword string `json:"user_password"`
}

func ListFilesHeaders(db *sql.DB, FN string) (*[]HeadersTable, error) {

	var rows *sql.Rows
	var err error

	switch FN {
	case "":
		rows, err = db.Query("SELECT * FROM headers")
	default:
		rows, err = db.Query("SELECT * FROM headers WHERE filename LIKE '%' || $1 || '%'", FN)
	}
	if err != nil {
		return nil, errors.Wrap(err, "DB Query, func ListFilesHeaders") //error handling
	}
	defer rows.Close()

	HeadersSlice := []HeadersTable{}

	for rows.Next() {
		Headers := HeadersTable{}
		err := rows.Scan(&Headers.FileName, &Headers.ContentType, &Headers.ContentLength, &Headers.ID)
		if err != nil {
			return nil, errors.Wrap(err, "rows.Scan, func ListFilesHeaders") //error handling
		}
		HeadersSlice = append(HeadersSlice, Headers)
	}

	return &HeadersSlice, nil
}

func GetFileHeaders(db *sql.DB, FN string) (*HeadersTable, error) {

	Headers := HeadersTable{}

	err := db.QueryRow("SELECT * FROM headers WHERE filename=$1", FN).Scan(&Headers.FileName, &Headers.ContentType, &Headers.ContentLength, &Headers.ID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, err
		default:
			return nil, errors.Wrap(err, "DB Query, func GetFileHeaders") //error handling
		}
	}

	return &Headers, nil
}

func PostFileHeaders(db *sql.DB, FN, CT, CL string) error {

	_, err := db.Exec("INSERT INTO headers (filename, content_type, content_lenght) VALUES ($1, $2, $3)", FN, CT, CL)
	if err != nil {
		return errors.Wrap(err, "DB Query, func PostFileHeaders") //error handling
	}

	return nil
}

func CheckFileByName(db *sql.DB, FN string) (*bool, error) {

	var flag bool

	Headers := HeadersTable{}

	err := db.QueryRow("SELECT filename FROM headers WHERE filename=$1", FN).Scan(&Headers.FileName)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			flag = false
			return &flag, nil
		default:
			return nil, errors.Wrap(err, "DB Query, func CheckFileByName") //error handling
		}
	}
	flag = true
	return &flag, nil
}

func UpdateTable(db *sql.DB, CT, CL, FN string) error {

	_, err := db.Exec("UPDATE headers SET content_type=$1, content_lenght=$2 WHERE filename=$3", CT, CL, FN)
	if err != nil {
		return errors.Wrap(err, "DB Query, func UpdateTable") //error handling
	}

	return nil
}

func CheckUserByName(db *sql.DB, Name string) (*bool, error) {

	var flag bool

	User := UserData{}

	err := db.QueryRow("SELECT user_name FROM user_data WHERE user_name=$1", Name).Scan(&User.Name)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			flag = false
			return &flag, nil
		default:
			return nil, errors.Wrap(err, "DB Query, func CheckUserByName") //error handling
		}
	}
	flag = true
	return &flag, nil
}

func CheckUserByNameAndPassword(db *sql.DB, Name string) (*bool, string, error) {

	var flag bool

	User := UserData{}

	err := db.QueryRow("SELECT user_password FROM user_data WHERE user_name=$1", Name).Scan(&User.HashPassword)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			flag = false
			return &flag, "", nil
		default:
			return nil, "", errors.Wrap(err, "DB Query, func CheckUserByName") //error handling
		}
	}

	flag = true
	return &flag, User.HashPassword, nil

}

func AddUserInDB(db *sql.DB, UserName, HashPassword string) error {

	_, err := db.Exec("INSERT INTO user_data (user_name, user_password) VALUES ($1, $2)", UserName, HashPassword)
	if err != nil {
		return errors.Wrap(err, "DB Query, func AddUserInDB") //error handling
	}

	return nil
}
