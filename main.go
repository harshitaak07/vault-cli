package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func usage() {
	fmt.Println(`Usage:
  vault [command] [args...]

Commands:
  upload <file>            Upload file to cloud (or local if mode=local)
  download <file>          Download file from cloud and decrypt
  local-upload <file>      Copy file to local vault
  local-download <file>    Retrieve file from local vault
  list                     List local DB entries
  audit                    Show audit log (DB)
  report                   Show summary report (file count, total size)

`)
}

func main() {
	godotenv.Load(".env")
	_ = os.Setenv("GOCACHE", os.Getenv("GOCACHE"))
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := OpenDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	defer db.Close()

	if err := InitDB(db); err != nil {
		log.Fatalf("db init: %v", err)
	}

	if cfg.RequirePassword {
		if ok := VerifyPassword(cfg.PasswordFile); !ok {
			log.Fatal("access denied: wrong password")
		}
	}

	if len(os.Args) < 2 {
		usage()
		return
	}
	cmd := os.Args[1]

	switch cmd {
	case "upload":
		if len(os.Args) < 3 {
			log.Fatal("upload requires a file path")
		}
		file := os.Args[2]
		if err := UploadHandler(file, cfg, db); err != nil {
			log.Fatalf("upload: %v", err)
		}
		fmt.Println("uploaded")

	case "download":
		if len(os.Args) < 3 {
			log.Fatal("download requires a file name")
		}
		file := os.Args[2]
		if err := DownloadHandler(file, cfg, db); err != nil {
			log.Fatalf("download: %v", err)
		}
		fmt.Println("downloaded")

	case "local-upload":
		if len(os.Args) < 3 {
			log.Fatal("local-upload requires a file path")
		}
		file := os.Args[2]
		if err := LocalUpload(file); err != nil {
			log.Fatalf("local-upload: %v", err)
		}
		fmt.Println("local uploaded")

	case "local-download":
		if len(os.Args) < 3 {
			log.Fatal("local-download requires a file name")
		}
		file := os.Args[2]
		if err := LocalDownload(file); err != nil {
			log.Fatalf("local-download: %v", err)
		}
		fmt.Println("local downloaded")

	case "list":
		if err := PrintDBEntries(db); err != nil {
			log.Fatalf("list: %v", err)
		}

	case "audit":
		if err := PrintAudit(db); err != nil {
			log.Fatalf("audit: %v", err)
		}
	case "report":
		GenerateReport(db)

	default:
		usage()
	}
}
