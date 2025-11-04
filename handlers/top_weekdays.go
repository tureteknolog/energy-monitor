package handlers

import (
    "encoding/json"
    "energy-monitor/db"
    "energy-monitor/models"
    "net/http"
    "time"
    
    "github.com/gorilla/mux"
)

type TopWeekdaysResponse struct {
    SiteID  string                `json:"site_id"`
    Month   string                `json:"month"`
    Top3    []models.EnergyData   `json:"top3"`
    Average float64               `json:"average"`
}

func GetTopWeekdays(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    siteID := vars["site_id"]
    
    // Hämta aktuell månad
    now := time.Now()
    year := now.Year()
    month := int(now.Month())
    
    // Hämta all data för månaden
    allData, err := db.GetMonthEnergyData(siteID, year, month)
    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }
    
    // Hämta helgdagar
    holidays, err := db.GetHolidays(siteID)
    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }
    
    // Filtrera data: vardagar 7-18 svensk tid, exkludera helgdagar
    var filtered []models.EnergyData
    stockholmLoc, _ := time.LoadLocation("Europe/Stockholm")
    
    for _, data := range allData {
        // Parse datum
        date, err := time.Parse("2006-01-02", data.Date)
        if err != nil {
            continue
        }
        
        // Kolla om helgdag
        if holidays[data.Date] {
            continue
        }
        
        // Kolla om vardag (måndag=1, fredag=5)
        weekday := date.Weekday()
        if weekday == time.Sunday || weekday == time.Saturday {
            continue
        }
        
        // Konvertera UTC-timme till svensk timme
        utcTime := time.Date(date.Year(), date.Month(), date.Day(), data.Hour, 0, 0, 0, time.UTC)
        localTime := utcTime.In(stockholmLoc)
        localHour := localTime.Hour()
        
        // Endast timmar 7-18 (7:00-18:59)
        if localHour >= 7 && localHour <= 18 {
            filtered = append(filtered, data)
        }
    }
    
    // Ta topp 3
    top3 := filtered
    if len(top3) > 3 {
        top3 = filtered[:3]
    }
    
    // Beräkna medelvärde
    var sum float64
    for _, data := range top3 {
        sum += data.ConsumptionKWh
    }
    
    var average float64
    if len(top3) > 0 {
        average = sum / float64(len(top3))
    }
    
    response := TopWeekdaysResponse{
        SiteID:  siteID,
        Month:   time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).Format("2006-01"),
        Top3:    top3,
        Average: average,
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
