package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Headers struct {
	FileName       string `json:"filename"`
	Content_Type   string `json:"content_type"`
	Content_Length string `json:"content_lenght"`
	Id             int    `json:"id"`
}

func GetAllDb(FN string, db *sql.DB) (*[]Headers, error) {

	var rows *sql.Rows
	var err error

	switch FN {
	case "allfiles":
		rows, err = db.Query("SELECT * FROM headers")
	default:
		rows, err = db.Query("SELECT * FROM headers WHERE filename LIKE '%' || $1 || '%'", FN)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	headersSlice := []Headers{}

	for rows.Next() {
		H := Headers{}
		err := rows.Scan(&H.FileName, &H.Content_Type, &H.Content_Length, &H.Id)
		if err != nil {
			return nil, err
		}
		headersSlice = append(headersSlice, H)
	}

	return &headersSlice, nil
}

func GetFileDb(FN string, db *sql.DB) (*Headers, error) {

	H := Headers{}

	err := db.QueryRow("SELECT * FROM headers WHERE filename=$1", FN).Scan(&H.FileName, &H.Content_Type, &H.Content_Length, &H.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		} else {
			return nil, err
		}
	}

	return &H, nil
}

func PostFileDb(FN, CT, CL string, db *sql.DB) error {

	_, err := db.Exec("INSERT INTO headers (filename, content_type, content_lenght) VALUES ($1, $2, $3)", FN, CT, CL)
	if err != nil {
		return err
	}

	return nil
}

func CheckFileInDb(FN string, db *sql.DB) (*bool, error) {

	var flag bool

	H := Headers{}

	err := db.QueryRow("SELECT filename FROM headers WHERE filename=$1", FN).Scan(&H.FileName)
	if err != nil {
		if err == sql.ErrNoRows {
			flag = false
			return &flag, nil
		} else {
			return nil, err
		}
	}
	flag = true
	return &flag, nil
}

func UpdateDb(CT, CL, FN string, db *sql.DB) error {

	_, err := db.Exec("UPDATE headers SET content_type=$1, content_lenght=$2 WHERE filename=$3", CT, CL, FN)
	if err != nil {
		return err
	}

	return nil
}
