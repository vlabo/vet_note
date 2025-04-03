CREATE TABLE IF NOT EXISTS patients (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type TEXT,
    name TEXT,
    gender TEXT,
    age NUMERIC,
    chip_id TEXT,
    weight REAL,
    castrated NUMERIC DEFAULT 0,
    note TEXT,
    owner TEXT,
    owner_phone TEXT,
    folder NUMERIC DEFAULT -1,
    index_folder NUMERIC DEFAULT -1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

CREATE TABLE IF NOT EXISTS procedures (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type TEXT,
    date TEXT,
    details TEXT,
    patient_id INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY(patient_id) REFERENCES patients(id)
);

CREATE TABLE IF NOT EXISTS settings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    value TEXT,
    type TEXT,
    idx INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

-- ALTER TABLE patients ADD COLUMN weight REAL;
-- ALTER TABLE patients ADD COLUMN folder NUMERIC DEFAULT -1;
-- ALTER TABLE patients ADD COLUMN index_folder NUMERIC DEFAULT -1;
