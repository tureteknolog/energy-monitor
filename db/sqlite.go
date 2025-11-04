package db

import (
    "database/sql"
    "log"
    
    _ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(dataSourceName string) error {
    var err error
    DB, err = sql.Open("sqlite3", dataSourceName)
    if err != nil {
        return err
    }
    
    if err = DB.Ping(); err != nil {
        return err
    }
    
    log.Println("Database connected")
    return createTables()
}

func createTables() error {
    schema := `
    CREATE TABLE IF NOT EXISTS sites (
        id TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );
    
    CREATE TABLE IF NOT EXISTS energy_data (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        site_id TEXT NOT NULL,
        date TEXT NOT NULL,
        hour INTEGER NOT NULL,
        consumption_kwh REAL NOT NULL,
        outdoor_temp REAL NOT NULL,
        wind_speed REAL NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        UNIQUE(site_id, date, hour),
        FOREIGN KEY(site_id) REFERENCES sites(id)
    );
    
    CREATE TABLE IF NOT EXISTS holidays (
        site_id TEXT NOT NULL,
        date TEXT NOT NULL,
        name TEXT,
        PRIMARY KEY(site_id, date),
        FOREIGN KEY(site_id) REFERENCES sites(id)
    );
    
    CREATE INDEX IF NOT EXISTS idx_energy_data_site_date 
        ON energy_data(site_id, date);
    CREATE INDEX IF NOT EXISTS idx_energy_data_consumption 
        ON energy_data(consumption_kwh DESC);
    `
    
    _, err := DB.Exec(schema)
    if err != nil {
        return err
    }
    
    log.Println("Database tables created")
    return nil
}

func Close() {
    if DB != nil {
        DB.Close()
    }
}
