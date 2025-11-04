CREATE TABLE IF NOT EXISTS files (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  filename TEXT,
  uploaded_at TEXT,
  hash TEXT,
  size INTEGER,
  location TEXT,
  mode TEXT
);

CREATE TABLE IF NOT EXISTS audit (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  action TEXT,
  filename TEXT,
  target TEXT,
  success INTEGER,
  err TEXT,
  ts TEXT
);

CREATE TABLE IF NOT EXISTS secrets (
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
);


RecordFileToDynamo(keyName, hash, info.Size(), cfg.Mode, "s3")
RecordAuditToDynamo("upload", keyName, "s3", true, "")