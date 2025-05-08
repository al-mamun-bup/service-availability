# 🛰️ Service Availability API

This project provides a RESTful API that checks **food service availability** for a city based on:

- ✅ Open hours defined per city
- ✅ Geofence logic using H3 Uber index

Built with **Go 1.23** using the **Echo** web framework, the project follows **clean architecture** principles and integrates with a remote registry API for dynamic city settings.

---

## 🚀 Quick Start (Docker Compose)

```bash
git clone https://github.com/al-mamun-bup/service-availability.git
cd service-availability
```

Once running, the API will be available at:

```
http://localhost:8080/v1/is_service_available
```

---

## 📡 API Endpoint

### `GET /v1/is_service_available`

Checks if the food service is available for a given location and time.

---

### 🔍 Query Parameters

| Name      | Type   | Required | Description                                           |
|-----------|--------|----------|-------------------------------------------------------|
| `check`   | string | ✅       | Service type to check (`food` supported)              |
| `city_id` | int    | ✅       | ID of the city to evaluate                            |
| `lat`     | float  | ✅       | Latitude of user's location                           |
| `long`    | float  | ✅       | Longitude of user's location                          |

---

### ✅ Example Success Response

```json
{
  "Success": true
}
```

### ❌ Example Failure Responses

**Outside geofence or closed hours:**

```json
{
  "Success": false
}
```

**Invalid `city_id`:**

```json
{
  "error": "Invalid city_id format"
}
```

**Missing query parameters:**

```json
{
  "error": "Missing or invalid query parameters"
}
```

---

## 🧠 Project Architecture

This project uses **Clean Architecture** principles. The structure looks like this:

```
service-availability/
├── cmd/                      # Main application entrypoint
│   └── main.go
├── internal/
│   ├── handlers/             # Echo handlers (availability, geo logic)
│   ├── services/             # Business logic (time, geofence)
│   ├── registry/             # External API calls
│   ├── models/               # Data models (CitySettings)
│   └── utils/                # Helper utilities (time parser, geo helpers)
├── tests/                    # Unit tests with mocks
├── go.mod / go.sum
└── README.md
```

---

## 🧪 Run Unit Tests

```bash
go test ./...
```

- Tests include: Time checks, polygon matching, full API handler with mocks.

---

## 👨‍💻 Author

**Abdullah Al Mamun**  
Backend Engineer @ Pathao  
[GitHub](https://github.com/al-mamun-bup)

---

## 📜 License

This project is licensed under the **MIT License**.
