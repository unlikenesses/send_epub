package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"

	"gopkg.in/gomail.v2"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	attachment := GetAttachmentFilename()
	SendEmail(attachment)
}

func GetAttachmentFilename() string {
	path := "./"
	files, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer files.Close()

	names, err := files.Readdirnames(0)
	for _, file := range names {
		if strings.Contains(file, "epub") {
			return RenameFile(file)
		}
	}

	log.Fatal("No file found")

	return ""
}

// Renaming an epub to a png currently makes Amazon convert it to a file readable by Kindle ¯\_(ツ)_/¯
func RenameFile(filename string) string {
	newFilename := strings.Replace(filename, "epub", "png", 1)
	err := os.Rename(filename, newFilename)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found file %v and renamed to %v, now sending...", filename, newFilename)

	return newFilename
}

func SendEmail(attachment string) {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	password := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("FROM_ADDRESS")
	to := os.Getenv("TO_ADDRESS")
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetBody("text/html", "Hi Kindle")
	m.Attach(attachment)

	d := gomail.NewDialer(smtpHost, smtpPort, from, password)

	err := d.DialAndSend(m)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Email sent!")
}
