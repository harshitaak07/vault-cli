package core

import (
	"database/sql"
	"fmt"
)

func GenerateReport(db *sql.DB) {
	rows, err := db.Query("SELECT COUNT(*), SUM(size) FROM files")
	if err != nil {
		fmt.Println("report error:", err)
		return
	}
	defer rows.Close()

	var count int
	var totalSize int64
	rows.Next()
	rows.Scan(&count, &totalSize)

	fmt.Printf("Files Stored: %d\n", count)
	fmt.Printf("Total Size: %.2f MB\n", float64(totalSize)/1024/1024)

	fmt.Println("\nRecent Uploads:")
	rows2, err := db.Query("SELECT filename, uploaded_at FROM files ORDER BY uploaded_at DESC LIMIT 5")
	if err != nil {
		fmt.Println("report error:", err)
		return
	}
	defer rows2.Close()

	for rows2.Next() {
		var name, ts string
		rows2.Scan(&name, &ts)
		fmt.Printf("- %s (%s)\n", name, ts)
	}
}
