package models

type Holiday struct {
    SiteID string `json:"site_id"`
    Date   string `json:"date"` // "2025-12-24" (lokal svensk tid)
    Name   string `json:"name"`
}
