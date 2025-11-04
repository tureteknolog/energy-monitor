package main

import (
    "energy-monitor/db"
    "energy-monitor/handlers"
    "log"
    "net/http"
    "os"
    
    "github.com/gorilla/mux"
)

func main() {
    // Initiera databas
    dbPath := os.Getenv("DB_PATH")
    if dbPath == "" {
        dbPath = "./data/energy.db"
    }
    
    if err := db.InitDB(dbPath); err != nil {
        log.Fatal("Failed to initialize database:", err)
    }
    defer db.Close()
    
    // Setup router
    r := mux.NewRouter()
    
    // API endpoints
    r.HandleFunc("/api/{site_id}/consumption", handlers.PostConsumption).Methods("POST")
    r.HandleFunc("/api/{site_id}/top-weekdays", handlers.GetTopWeekdays).Methods("GET")
    r.HandleFunc("/api/sites", handlers.GetSites).Methods("GET")
    r.HandleFunc("/health", healthCheck).Methods("GET")
    
    // Serve static files
    r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
    
    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8081"
    }
    
    log.Printf("Server starting on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}
