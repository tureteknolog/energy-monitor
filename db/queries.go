package db

import (
    "database/sql"
    "energy-monitor/models"
    "fmt"
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

// GetHolidays hämtar alla helgdagar för en sajt
func GetHolidays(siteID string) (map[string]bool, error) {
    query := `SELECT date FROM holidays WHERE site_id = ?`
    rows, err := DB.Query(query, siteID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    holidays := make(map[string]bool)
    for rows.Next() {
        var date string
        if err := rows.Scan(&date); err != nil {
            return nil, err
        }
        holidays[date] = true
    }
    
    return holidays, nil
}

// GetMonthEnergyData hämtar all energidata för en månad
func GetMonthEnergyData(siteID string, year int, month int) ([]models.EnergyData, error) {
    query := `
        SELECT id, site_id, date, hour, consumption_kwh, outdoor_temp, wind_speed, created_at
        FROM energy_data
        WHERE site_id = ?
        AND substr(date, 1, 7) = ?
        ORDER BY consumption_kwh DESC
    `
    
    yearMonth := fmt.Sprintf("%04d-%02d", year, month)
    rows, err := DB.Query(query, siteID, yearMonth)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var data []models.EnergyData
    for rows.Next() {
        var ed models.EnergyData
        if err := rows.Scan(&ed.ID, &ed.SiteID, &ed.Date, &ed.Hour, 
            &ed.ConsumptionKWh, &ed.OutdoorTemp, &ed.WindSpeed, &ed.CreatedAt); err != nil {
            return nil, err
        }
        data = append(data, ed)
    }
    
    return data, nil
}
