package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

func OpenDB(path string) (*sql.DB, error) {
	return sql.Open("sqlite", path)
}

func InitDB(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS files (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			filename TEXT,
			uploaded_at TEXT,
			hash TEXT,
			size INTEGER,
			location TEXT,
			mode TEXT
		);`,

		`CREATE TABLE IF NOT EXISTS audit (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			action TEXT,
			filename TEXT,
			target TEXT,
			success INTEGER,
			err TEXT,
			ts TEXT
		);`,

		`CREATE TABLE IF NOT EXISTS secrets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			category TEXT NOT NULL,
			name TEXT NOT NULL,
			ciphertext BLOB NOT NULL,
			nonce BLOB NOT NULL,
			mode TEXT NOT NULL,
			hash TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			UNIQUE(category, name)
		);`,
	}

	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			return fmt.Errorf("InitDB failed: %v", err)
		}
	}

	fmt.Println("✅ Database initialized with tables: files, audit, secrets")
	return nil
}

func RecordFile(db *sql.DB, filename, hash string, size int64, location, mode string) error {
	_, err := db.Exec(`INSERT INTO files(filename, uploaded_at, hash, size, location, mode) VALUES(?,?,?,?,?,?)`,
		filename, time.Now().UTC().Format(time.RFC3339), hash, size, location, mode)
	return err
}

func RecordAudit(db *sql.DB, action, filename, target string, success bool, errMsg string) error {
	sc := 0
	if success {
		sc = 1
	}
	_, err := db.Exec(`INSERT INTO audit(action, filename, target, success, err, ts) VALUES(?,?,?,?,?,?)`,
		action, filename, target, sc, errMsg, time.Now().UTC().Format(time.RFC3339))
	return err
}

func PrintDBEntries(db *sql.DB) error {
	rows, err := db.Query(`SELECT id, filename, uploaded_at, hash, size, location, mode FROM files ORDER BY uploaded_at DESC`)
	if err != nil {
		return err
	}
	defer rows.Close()
	fmt.Println("Files:")
	for rows.Next() {
		var id int
		var fn, ts, hash, loc, mode string
		var size int64
		if err := rows.Scan(&id, &fn, &ts, &hash, &size, &loc, &mode); err != nil {
			return err
		}
		fmt.Printf("%d %s %s %s %d %s %s\n", id, fn, ts, hash, size, loc, mode)
	}
	return rows.Err()
}

func PrintAudit(db *sql.DB) error {
	rows, err := db.Query(`SELECT id, action, filename, target, success, err, ts FROM audit ORDER BY ts DESC LIMIT 200`)
	if err != nil {
		return err
	}
	defer rows.Close()
	fmt.Println("Audit:")
	for rows.Next() {
		var id int
		var action, filename, target, errMsg, ts string
		var success int
		if err := rows.Scan(&id, &action, &filename, &target, &success, &errMsg, &ts); err != nil {
			return err
		}
		fmt.Printf("%d %s %s -> %s success=%d err=%s at %s\n", id, action, filename, target, success, errMsg, ts)
	}
	return nil
}
