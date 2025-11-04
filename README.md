# Energy Monitor

Elförbrukningsmonitor med multi-site support.

## Funktioner

- Ta emot timförbrukning från HomeAssistant
- Lagra data med utetemperatur och vindstyrka
- Visa topp 3 förbrukningar för vardagar 7-18 (svensk tid)
- Beräkna medelvärde av topp 3
- Multi-site support
- Automatisk hantering av svenska helgdagar
- Mobilanpassat webbgränssnitt

## Installation

### Kör med Docker Compose
```bash
docker-compose up -d
```

### Kör lokalt
```bash
# Bygg
go build -o energy-monitor

# Kör
./energy-monitor
```

Servern startar på http://localhost:8080

## API

### POST /api/{site_id}/consumption

Skicka förbrukningsdata:
```json
{
  "timestamp": "2025-11-04T13:00:00Z",
  "consumption_kwh": 12.5,
  "outdoor_temp": 5.2,
  "wind_speed": 3.8
}
```

### GET /api/{site_id}/top-weekdays

Hämta topp 3 för innevarande månad (vardagar 7-18).

### GET /api/sites

Lista alla sajter.

## Populera helgdagar

Efter att sajter skapats, bygg och kör verktyget:
```bash
go build -o populate_holidays ./cmd/populate_holidays
./populate_holidays
```

Detta lägger till svenska helgdagar 2025 för alla sajter.

## HomeAssistant Integration

Exempel på automation i HomeAssistant:
```yaml
automation:
  - alias: "Send hourly energy data"
    trigger:
      platform: time_pattern
      minutes: 1
    action:
      - service: rest_command.send_energy_data
        data:
          site_id: "strangnas"
          timestamp: "{{ now().isoformat() }}"
          consumption: "{{ states('sensor.energy_consumption') }}"
          temp: "{{ states('sensor.outdoor_temperature') }}"
          wind: "{{ states('sensor.wind_speed') }}"

rest_command:
  send_energy_data:
    url: "http://your-server:8080/api/{{ site_id }}/consumption"
    method: POST
    content_type: "application/json"
    payload: >
      {
        "timestamp": "{{ timestamp }}",
        "consumption_kwh": {{ consumption }},
        "outdoor_temp": {{ temp }},
        "wind_speed": {{ wind }}
      }
```

## Tidszoner

- All lagring sker i UTC
- Filtrering av vardagar 7-18 görs i svensk tid (Europe/Stockholm)
- Automatisk hantering av sommar/vintertid

## Databas

SQLite databas lagras i `./data/energy.db`
