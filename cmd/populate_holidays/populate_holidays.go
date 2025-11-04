package main

import (
    "energy-monitor/db"
    "log"
)

func main() {
    // Initiera databas
    if err := db.InitDB("./data/energy.db"); err != nil {
        log.Fatal("Failed to initialize database:", err)
    }
    defer db.Close()
    
    // Svenska helgdagar 2025
    holidays := []struct {
        Date string
        Name string
    }{
        // Fasta helgdagar
        {"2025-01-01", "Nyårsdagen"},
        {"2025-01-06", "Trettondedag jul"},
        {"2025-05-01", "Första maj"},
        {"2025-06-06", "Nationaldagen"},
        {"2025-12-24", "Julafton"},
        {"2025-12-25", "Juldagen"},
        {"2025-12-26", "Annandag jul"},
        {"2025-12-31", "Nyårsafton"},
        
        // Rörliga helgdagar 2025
        {"2025-04-18", "Långfredagen"},
        {"2025-04-20", "Påskdagen"},
        {"2025-04-21", "Annandag påsk"},
        {"2025-05-29", "Kristi himmelsfärdsdag"},
        {"2025-06-08", "Pingstdagen"},
        {"2025-06-20", "Midsommarafton"},
        {"2025-06-21", "Midsommardagen"},
        {"2025-11-01", "Alla helgons dag"},
    }
    
    // Hämta alla sajter
    sites, err := db.GetAllSites()
    if err != nil {
        log.Fatal("Failed to get sites:", err)
    }
    
    if len(sites) == 0 {
        log.Println("No sites found. Holidays will be added when sites are created.")
        log.Println("You can run this script again after creating sites.")
        return
    }
    
    // Lägg till helgdagar för varje sajt
    for _, site := range sites {
        log.Printf("Adding holidays for site: %s", site.ID)
        for _, holiday := range holidays {
            query := `INSERT OR IGNORE INTO holidays (site_id, date, name) VALUES (?, ?, ?)`
            _, err := db.DB.Exec(query, site.ID, holiday.Date, holiday.Name)
            if err != nil {
                log.Printf("Failed to add holiday %s: %v", holiday.Name, err)
            }
        }
    }
    
    log.Println("Holidays populated successfully!")
}
