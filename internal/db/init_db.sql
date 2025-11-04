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
RecordFileToDynamo(keyName, hash, info.Size(), cfg.Mode, "s3")
RecordAuditToDynamo("upload", keyName, "s3", true, "")