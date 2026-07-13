# ONVIF Edge Controller

ONVIF Edge Controller is a Go-based REST API service for discovering and controlling ONVIF-compatible IP cameras.

The service provides APIs for:

- ONVIF camera discovery
- Camera connection and validation
- RTSP stream URL retrieval
- Snapshot URL retrieval
- PTZ camera control
- Network interface discovery

---

# Features

## Camera Discovery

Discover ONVIF-compatible cameras available on the local network.

## Camera Information

Connect to cameras and retrieve:

- Manufacturer
- Model
- Firmware version
- Hardware ID
- Serial number

## Stream APIs

Retrieve:

- RTSP stream URLs
- Snapshot URLs
- Media profiles

## PTZ Control

Supports PTZ camera actions:

- left
- right
- up
- down

---

# Requirements

- Go 1.23+
- ONVIF-compatible IP camera
- Network access to camera

---

# Installation

Clone repository:

```bash
git clone https://tvmgit.thinkpalm.info/netvire/platform/edge/onvif-edge-controller.git
```

Go to project:

```bash
cd onvif-edge-controller
```

Install dependencies:

```bash
go mod tidy
```

---

# Dependency Notes

This project uses a locally patched version of the ONVIF PTZ library:

```go
replace github.com/0x524a/onvif-go => ./internal/vendor/onvif-go
```

The patch improves PTZ compatibility with certain ONVIF cameras that require strict XML namespace validation for PTZ requests and otherwise may reject PTZ operations with errors similar to:

```text
Validation constraint violation:
tag name or namespace mismatch in element 'PanTilt'
```

### Important

The patched dependency is located at:

```text
internal/vendor/onvif-go
```

This directory must be included in source control.

Building the project without the patched dependency may cause PTZ operations to fail on certain ONVIF-compliant cameras.

---

# Run Application

```bash
go run cmd/server/main.go
```

Server starts on:

```text
http://localhost:8080
```

---
# APIs

## Authentication

All `/api/v1/*` APIs are protected using Basic Authentication.

Provide the following headers while calling the APIs:

```http
Authorization: Basic <base64(username:password)>
```

Example credentials are configured using environment variables:

```env
API_USERNAME=*****
API_PASSWORD=******
```

## APIs Without Authentication

The following endpoint does NOT require authentication:

```http
GET /health
```
---

## Available APIs

### Health Check
> No authentication required

```http
GET /health
```

---

### Discover Cameras
> Requires Basic Authentication

```http
GET /api/v1/discovery
```

Optional interface:

```http
GET /api/v1/discovery?interface=eth0
```

---

### List Network Interfaces
> Requires Basic Authentication

```http
GET /api/v1/interfaces
```

---

### Connect Camera
> Requires Basic Authentication

```http
POST /api/v1/camera/connect
```

Example request:

```json
{
  "ip":"192.168.0.10",
  "username":"admin",
  "password":"admin"
}
```

---

### Get Stream URLs
> Requires Basic Authentication

```http
POST /api/v1/camera/streams
```

Example request:

```json
{
  "ip":"192.168.0.10",
  "username":"admin",
  "password":"admin"
}
```

---

### Get PTZ Capabilities

> Requires Basic Authentication

```http
POST /api/v1/camera/ptz-capabilities
```
Example request:

```json
{
  "ip":"192.168.0.10",
  "username":"admin",
  "password":"admin"
}
```

---

### PTZ & Preset Actions
> Requires Basic Authentication

```http
POST /api/v1/action
```

Example request (PTZ movement):

```json
{
  "ip": "192.168.0.124",
  "username": "admin",
  "password": "admin",
  "action": "DIRECTION_LEFT",
  "speed": 0.2
}
```
Note: speed is optional. If omitted, the default PTZ speed configured by the service is used.

Example request (Set Preset):

```json
{
  "ip": "192.168.0.124",
  "username": "admin",
  "password": "admin",
  "action": "SET_PRESET",
  "presetToken": "1"
}
```

Example request (Go to Preset):

```json
{
  "ip": "192.168.0.124",
  "username": "admin",
  "password": "admin",
  "action": "GOTO_PRESET",
  "presetToken": "1"
}
```

---

## Supports actions:

### PTZ Movement

- DIRECTION_LEFT
- DIRECTION_RIGHT
- DIRECTION_UP
- DIRECTION_DOWN
- ZOOM_IN
- ZOOM_OUT

### Preset Operations
- SET_PRESET
- GOTO_PRESET

---

# Technology Stack

- Go
- Gin Web Framework
- ONVIF Go SDK

---
