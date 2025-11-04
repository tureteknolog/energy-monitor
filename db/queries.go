package db

import (
    "database/sql"
    "energy-monitor/models"
    "strings"
)

// CreateOrGetSite skapar sajt om den inte finns, returnerar site
func CreateOrGetSite(siteID string) error {
    // Skapa ett snyggt namn från ID (strangnas -> Strängnäs)
    name := strings.Title(siteID)
    
    query := `INSERT OR IGNORE INTO sites (id, name) VALUES (?, ?)`
    _, err := DB.Exec(query, siteID, name)
    return err
}

// InsertEnergyData sparar energidata
func InsertEnergyData(data *models.EnergyData) error {
    query := `
        INSERT INTO energy_data (site_id, date, hour, consumption_kwh, outdoor_temp, wind_speed)
        VALUES (?, ?, ?, ?, ?, ?)
        ON CONFLICT(site_id, date, hour) DO UPDATE SET
            consumption_kwh = excluded.consumption_kwh,
            outdoor_temp = excluded.outdoor_temp,
            wind_speed = excluded.wind_speed
    `
    
    _, err := DB.Exec(query, data.SiteID, data.Date, data.Hour, 
        data.ConsumptionKWh, data.OutdoorTemp, data.WindSpeed)
    return err
}

// GetAllSites hämtar alla sajter
func GetAllSites() ([]models.Site, error) {
    query := `SELECT id, name, created_at FROM sites ORDER BY name`
    rows, err := DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var sites []models.Site
    for rows.Next() {
        var site models.Site
        if err := rows.Scan(&site.ID, &site.Name, &site.CreatedAt); err != nil {
            return nil, err
        }
        sites = append(sites, site)
    }
    
    return sites, nil
}
