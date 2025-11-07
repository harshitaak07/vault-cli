package db

import "database/sql"

type FileRecord struct {
    ID        int    `json:"id"`
    Filename  string `json:"filename"`
    Uploaded  string `json:"uploaded_at"`
    Hash      string `json:"hash"`
    Size      int64  `json:"size"`
    Location  string `json:"location"`
    Mode      string `json:"mode"`
}

func ListFiles(db *sql.DB) ([]FileRecord, error) {
    rows, err := db.Query(`SELECT id, filename, uploaded_at, hash, size, location, mode FROM files ORDER BY uploaded_at DESC`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []FileRecord
    for rows.Next() {
        var r FileRecord
        if err := rows.Scan(&r.ID, &r.Filename, &r.Uploaded, &r.Hash, &r.Size, &r.Location, &r.Mode); err != nil {
            return nil, err
        }
        items = append(items, r)
    }
    return items, rows.Err()
}

type AuditRecord struct {
    ID       int    `json:"id"`
    Action   string `json:"action"`
    Filename string `json:"filename"`
    Target   string `json:"target"`
    Success  bool   `json:"success"`
    Error    string `json:"error"`
    TS       string `json:"timestamp"`
}

func ListAudit(db *sql.DB, limit int) ([]AuditRecord, error) {
    if limit <= 0 {
        limit = 100
    }
    rows, err := db.Query(`SELECT id, action, filename, target, success, err, ts FROM audit ORDER BY ts DESC LIMIT ?`, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []AuditRecord
    for rows.Next() {
        var (
            r       AuditRecord
            success int
        )
        if err := rows.Scan(&r.ID, &r.Action, &r.Filename, &r.Target, &success, &r.Error, &r.TS); err != nil {
            return nil, err
        }
        r.Success = success == 1
        items = append(items, r)
    }
    return items, rows.Err()
}


