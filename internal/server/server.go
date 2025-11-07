package server

import (
    "database/sql"
    "embed"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"

    "vault-cli/internal/auth"
    "vault-cli/internal/aws"
    "vault-cli/internal/config"
    "vault-cli/internal/core"
    "vault-cli/internal/db"
    "vault-cli/internal/secrets"
    "vault-cli/internal/session"
)

//go:embed static/*
var staticFS embed.FS

type Server struct {
    cfg *config.Config
    db  *sql.DB
}

func New(cfg *config.Config, database *sql.DB) *Server {
    return &Server{cfg: cfg, db: database}
}

func (s *Server) Start(addr string) error {
    srv := &http.Server{
        Addr:              addr,
        Handler:           s.routes(),
        ReadHeaderTimeout: 10 * time.Second,
    }
    fmt.Printf("🌐 Vault UI listening on http://%s\n", addr)
    return srv.ListenAndServe()
}

func (s *Server) routes() http.Handler {
    mux := http.NewServeMux()

    mux.HandleFunc("/api/health", s.handleHealth)
    mux.HandleFunc("/api/login", s.handleLogin)
    mux.HandleFunc("/api/files", s.wrapAuth(s.handleListFiles))
    mux.HandleFunc("/api/audit", s.wrapAuth(s.handleListAudit))
    mux.HandleFunc("/api/upload", s.wrapAuth(s.handleUpload))
    mux.HandleFunc("/api/download", s.wrapAuth(s.handleDownload))
    mux.HandleFunc("/api/secrets", s.wrapAuth(s.handleSecrets))
    mux.HandleFunc("/api/secrets/value", s.wrapAuth(s.handleSecretValue))

    fileServer := http.FileServer(http.FS(staticFS))
    mux.Handle("/", s.serveIndex(fileServer))

    return withCORS(mux)
}

func withCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        w.Header().Set("Access-Control-Allow-Methods", "GET,POST,DELETE,OPTIONS")
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }
        next.ServeHTTP(w, r)
    })
}

func (s *Server) serveIndex(fs http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if strings.HasPrefix(r.URL.Path, "/api/") {
            http.NotFound(w, r)
            return
        }
        if r.URL.Path == "/" {
            f, err := staticFS.Open("static/index.html")
            if err != nil {
                http.Error(w, "index missing", http.StatusInternalServerError)
                return
            }
            defer f.Close()
            w.Header().Set("Content-Type", "text/html; charset=utf-8")
            _, _ = io.Copy(w, f)
            return
        }
        fs.ServeHTTP(w, r)
    })
}

func (s *Server) wrapAuth(handler http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if s.cfg.RequirePassword {
            if err := session.Require(); err != nil {
                s.writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
                return
            }
        }
        handler(w, r)
    }
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
    s.writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        s.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
        return
    }

    if !s.cfg.RequirePassword {
        s.writeJSON(w, http.StatusOK, map[string]any{"ok": true, "requiresPassword": false})
        return
    }

    var req struct {
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        s.writeError(w, http.StatusBadRequest, "invalid json")
        return
    }
    if strings.TrimSpace(req.Password) == "" {
        s.writeError(w, http.StatusBadRequest, "password required")
        return
    }

    ok, err := auth.CheckPassword(s.cfg.PasswordFile, req.Password)
    if err != nil {
        s.writeError(w, http.StatusInternalServerError, err.Error())
        return
    }
    if !ok {
        s.writeError(w, http.StatusUnauthorized, "invalid credentials")
        return
    }

    if err := session.Save("web", 15*time.Minute); err != nil {
        s.writeError(w, http.StatusInternalServerError, err.Error())
        return
    }
    if sess, err := session.Load(); err == nil {
        s.writeJSON(w, http.StatusOK, map[string]any{
            "ok":        true,
            "expiresAt": sess.ExpiresAt,
        })
        return
    }
    s.writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (s *Server) handleListFiles(w http.ResponseWriter, r *http.Request) {
    items, err := db.ListFiles(s.db)
    if err != nil {
        s.writeError(w, http.StatusInternalServerError, err.Error())
        return
    }
    s.writeJSON(w, http.StatusOK, items)
}

func (s *Server) handleListAudit(w http.ResponseWriter, r *http.Request) {
    limit := 100
    if v := r.URL.Query().Get("limit"); v != "" {
        if n, err := strconv.Atoi(v); err == nil && n > 0 {
            limit = n
        }
    }
    items, err := db.ListAudit(s.db, limit)
    if err != nil {
        s.writeError(w, http.StatusInternalServerError, err.Error())
        return
    }
    s.writeJSON(w, http.StatusOK, items)
}

func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        s.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
        return
    }

    if err := r.ParseMultipartForm(128 << 20); err != nil {
        s.writeError(w, http.StatusBadRequest, fmt.Sprintf("parse form: %v", err))
        return
    }
    file, header, err := r.FormFile("file")
    if err != nil {
        s.writeError(w, http.StatusBadRequest, "file field required")
        return
    }
    defer file.Close()

    filename := sanitizeFilename(header.Filename)
    if filename == "" {
        s.writeError(w, http.StatusBadRequest, "invalid filename")
        return
    }

    tempDir, err := os.MkdirTemp("", "vault-upload-")
    if err != nil {
        s.writeError(w, http.StatusInternalServerError, err.Error())
        return
    }
    defer os.RemoveAll(tempDir)

    tempPath := filepath.Join(tempDir, filename)
    tempFile, err := os.Create(tempPath)
    if err != nil {
        s.writeError(w, http.StatusInternalServerError, err.Error())
        return
    }
    if _, err := io.Copy(tempFile, file); err != nil {
        tempFile.Close()
        s.writeError(w, http.StatusInternalServerError, err.Error())
        return
    }
    tempFile.Close()

    if s.cfg.Mode == "local" {
        if err := s.handleLocalUpload(tempPath, filename); err != nil {
            s.writeError(w, http.StatusInternalServerError, err.Error())
            return
        }
    } else {
        if err := aws.EncryptAndUpload(tempPath, s.cfg, s.db); err != nil {
            s.writeError(w, http.StatusInternalServerError, err.Error())
            return
        }
    }

    s.writeJSON(w, http.StatusCreated, map[string]string{"message": "upload complete"})
}

func (s *Server) handleLocalUpload(srcPath, name string) error {
    destDir := s.cfg.LocalPath
    if destDir == "" {
        destDir = os.Getenv("VAULT_REMOTE_PATH")
    }
    if destDir == "" {
        return errors.New("local path not configured")
    }
    if err := os.MkdirAll(destDir, 0755); err != nil {
        return err
    }
    destPath := filepath.Join(destDir, filepath.Base(name))

    if err := copyFile(srcPath, destPath); err != nil {
        return err
    }

    info, err := os.Stat(srcPath)
    if err != nil {
        return err
    }
    hash, err := core.FileSHA256(srcPath)
    if err != nil {
        return err
    }
    if s.db != nil {
        _ = db.RecordFile(s.db, name, hash, info.Size(), destDir, "local")
        _ = db.RecordAudit(s.db, "upload", name, destDir, true, "")
    }
    return nil
}

func copyFile(src, dest string) error {
    in, err := os.Open(src)
    if err != nil {
        return err
    }
    defer in.Close()

    out, err := os.Create(dest)
    if err != nil {
        return err
    }
    defer out.Close()

    if _, err := io.Copy(out, in); err != nil {
        return err
    }
    return out.Sync()
}

func (s *Server) handleDownload(w http.ResponseWriter, r *http.Request) {
    name := sanitizeFilename(r.URL.Query().Get("name"))
    if name == "" {
        s.writeError(w, http.StatusBadRequest, "name required")
        return
    }

    var data []byte
    var err error

    if s.cfg.Mode == "local" {
        destDir := s.cfg.LocalPath
        if destDir == "" {
            destDir = os.Getenv("VAULT_REMOTE_PATH")
        }
        if destDir == "" {
            s.writeError(w, http.StatusInternalServerError, "local path not configured")
            return
        }
        path := filepath.Join(destDir, name)
        data, err = os.ReadFile(path)
        if err != nil {
            s.writeError(w, http.StatusInternalServerError, err.Error())
            return
        }
        if s.db != nil {
            _ = db.RecordAudit(s.db, "download", name, destDir, true, "")
        }
    } else {
        if err := aws.DownloadAndDecrypt(name, s.cfg, s.db); err != nil {
            s.writeError(w, http.StatusInternalServerError, err.Error())
            return
        }
        localPath := "decrypted_" + filepath.Base(name)
        data, err = os.ReadFile(localPath)
        _ = os.Remove(localPath)
        if err != nil {
            s.writeError(w, http.StatusInternalServerError, err.Error())
            return
        }
    }

    w.Header().Set("Content-Type", "application/octet-stream")
    w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", name))
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write(data)
}

func (s *Server) handleSecrets(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        cat := r.URL.Query().Get("category")
        items, err := secrets.List(s.db, cat)
        if err != nil {
            s.writeError(w, http.StatusInternalServerError, err.Error())
            return
        }
        s.writeJSON(w, http.StatusOK, items)
    case http.MethodPost:
        var req struct {
            Category string `json:"category"`
            Name     string `json:"name"`
            Value    string `json:"value"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            s.writeError(w, http.StatusBadRequest, "invalid json")
            return
        }
        if strings.TrimSpace(req.Category) == "" || strings.TrimSpace(req.Name) == "" {
            s.writeError(w, http.StatusBadRequest, "category and name required")
            return
        }
        sec := secrets.Secret{Category: req.Category, Name: req.Name, Value: req.Value}
        if err := secrets.Add(s.db, s.cfg, sec); err != nil {
            s.writeError(w, http.StatusInternalServerError, err.Error())
            return
        }
        if s.db != nil {
            _ = db.RecordAudit(s.db, "secret:add", fmt.Sprintf("%s/%s", sec.Category, sec.Name), "secrets", true, "")
        }
        s.writeJSON(w, http.StatusCreated, map[string]string{"message": "secret stored"})
    case http.MethodDelete:
        cat := r.URL.Query().Get("category")
        name := r.URL.Query().Get("name")
        if strings.TrimSpace(cat) == "" || strings.TrimSpace(name) == "" {
            s.writeError(w, http.StatusBadRequest, "category and name required")
            return
        }
        if err := secrets.Delete(s.db, cat, name); err != nil {
            s.writeError(w, http.StatusInternalServerError, err.Error())
            return
        }
        if s.db != nil {
            _ = db.RecordAudit(s.db, "secret:delete", fmt.Sprintf("%s/%s", cat, name), "secrets", true, "")
        }
        s.writeJSON(w, http.StatusOK, map[string]string{"message": "secret deleted"})
    default:
        s.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
    }
}

func (s *Server) handleSecretValue(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        s.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
        return
    }
    cat := r.URL.Query().Get("category")
    name := r.URL.Query().Get("name")
    if strings.TrimSpace(cat) == "" || strings.TrimSpace(name) == "" {
        s.writeError(w, http.StatusBadRequest, "category and name required")
        return
    }
    val, err := secrets.Get(s.db, s.cfg, cat, name)
    if err != nil {
        s.writeError(w, http.StatusInternalServerError, err.Error())
        return
    }
    s.writeJSON(w, http.StatusOK, map[string]string{"value": val})
}

func (s *Server) writeError(w http.ResponseWriter, status int, msg string) {
    s.writeJSON(w, status, map[string]string{"error": msg})
}

func (s *Server) writeJSON(w http.ResponseWriter, status int, payload any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(payload)
}

func sanitizeFilename(name string) string {
    name = filepath.Base(strings.TrimSpace(name))
    name = strings.ReplaceAll(name, "\\", "")
    name = strings.ReplaceAll(name, "/", "")
    return name
}

