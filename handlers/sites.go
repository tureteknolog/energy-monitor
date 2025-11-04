package handlers

import (
    "encoding/json"
    "energy-monitor/db"
    "net/http"
)

func GetSites(w http.ResponseWriter, r *http.Request) {
    sites, err := db.GetAllSites()
    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(sites)
}
