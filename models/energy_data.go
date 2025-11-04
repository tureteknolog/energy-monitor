package models

import "time"

type EnergyData struct {
    ID             int       `json:"id"`
    SiteID         string    `json:"site_id"`
    Date           string    `json:"date"`            // "2025-11-04" (UTC)
    Hour           int       `json:"hour"`            // 0-23 (UTC)
    ConsumptionKWh float64   `json:"consumption_kwh"`
    OutdoorTemp    float64   `json:"outdoor_temp"`
    WindSpeed      float64   `json:"wind_speed"`
    CreatedAt      time.Time `json:"created_at"`
}

type ConsumptionRequest struct {
    Timestamp      string  `json:"timestamp"`       // "2025-11-04T13:00:00Z"
    ConsumptionKWh float64 `json:"consumption_kwh"`
    OutdoorTemp    float64 `json:"outdoor_temp"`
    WindSpeed      float64 `json:"wind_speed"`
}
