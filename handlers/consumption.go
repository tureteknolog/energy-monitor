package handlers

import (
    "encoding/json"
    "energy-monitor/db"
    "energy-monitor/models"
    "net/http"
    "time"
    
    "github.com/gorilla/mux"
)

func PostConsumption(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    siteID := vars["site_id"]
    
    var req models.ConsumptionRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    // Parse timestamp
    timestamp, err := time.Parse(time.RFC3339, req.Timestamp)
    if err != nil {
        http.Error(w, "Invalid timestamp format, use RFC3339", http.StatusBadRequest)
        return
    }
    
    // Skapa sajt om den inte finns
    if err := db.CreateOrGetSite(siteID); err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }
    
    // Extrahera datum och timme (UTC)
    date := timestamp.Format("2006-01-02")
    hour := timestamp.Hour()
    
    // Skapa energidata
    energyData := &models.EnergyData{
        SiteID:         siteID,
        Date:           date,
        Hour:           hour,
        ConsumptionKWh: req.ConsumptionKWh,
        OutdoorTemp:    req.OutdoorTemp,
        WindSpeed:      req.WindSpeed,
    }
    
    // Spara i databas
    if err := db.InsertEnergyData(energyData); err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{
        "status": "success",
        "site_id": siteID,
        "date": date,
    })
}
