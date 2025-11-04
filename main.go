package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"vault-cli/internal/auth"
	"vault-cli/internal/config"
	"vault-cli/internal/core"
	"vault-cli/internal/db"
	"vault-cli/internal/local"
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
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	database, err := db.OpenDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	defer database.Close()

	if err := db.InitDB(database); err != nil {
		log.Fatalf("db init: %v", err)
	}

	if cfg.RequirePassword {
		if ok := auth.VerifyPassword(cfg.PasswordFile); !ok {
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
		if err := core.UploadHandler(file, cfg, database); err != nil {
			log.Fatalf("upload: %v", err)
		}
		fmt.Println("uploaded")

	case "download":
		if len(os.Args) < 3 {
			log.Fatal("download requires a file name")
		}
		file := os.Args[2]
		if err := core.DownloadHandler(file, cfg, database); err != nil {
			log.Fatalf("download: %v", err)
		}
		fmt.Println("downloaded")

	case "local-upload":
		if len(os.Args) < 3 {
			log.Fatal("local-upload requires a file path")
		}
		file := os.Args[2]
		if err := local.LocalUpload(file); err != nil {
			log.Fatalf("local-upload: %v", err)
		}
		fmt.Println("local uploaded")

	case "local-download":
		if len(os.Args) < 3 {
			log.Fatal("local-download requires a file name")
		}
		file := os.Args[2]
		if err := local.LocalDownload(file); err != nil {
			log.Fatalf("local-download: %v", err)
		}
		fmt.Println("local downloaded")

	case "list":
		if err := db.PrintDBEntries(database); err != nil {
			log.Fatalf("list: %v", err)
		}

	case "audit":
		if err := db.PrintAudit(database); err != nil {
			log.Fatalf("audit: %v", err)
		}
	case "report":
		core.GenerateReport(database)

	default:
		usage()
	}
}
