package forfile

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
)

// METOd APP func(a App) GetOpeFile(FN string)... save Del
func GetOpenFile(FN string) ([]byte, error) {

	filename := fmt.Sprintf("home/ubuntu/infile/%s", FN) //config
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err) //error handling
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err) //error handling
	}

	return b, nil
}

func SaveNewFile(FN string, file multipart.File) {
	f, err := os.OpenFile("/home/ubuntu/infile/"+FN, os.O_WRONLY|os.O_CREATE, 0666) //config
	if err != nil {
		log.Fatal(err) //error handling
	}
	defer f.Close()

	io.Copy(f, file)
}

func DeleteOldFile(FN string) {

	err := os.Remove("/home/ubuntu/infile/" + FN) //config
	if err != nil {
		log.Fatal(err) //error handling
	}
}
