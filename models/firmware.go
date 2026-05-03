package models

// --- Firmware ---

type FirmwareUpdateResponse struct {
	Status              bool   `json:"status"`
	CurrentFirmware     string `json:"current_firmware"`
	NextFirmwareVersion string `json:"next_firmware_version"`
}

type FirmwareUpToDateResponse struct {
	Status          bool   `json:"status"`
	CurrentFirmware string `json:"current_firmware"`
}

// --- Provisioning ---

type ProvisionRequest struct {
	SerialNumber string `json:"serial_number"`
	Model        string `json:"model"`
	Firmware     string `json:"firmware"`
	UniversalKey string `json:"universal_key"`
}

type ProvisionResponse struct {
	Status bool   `json:"status"`
	APIKey string `json:"api_key"`
}

// --- Telemetry ---

type DeviceInfo struct {
	Serial string `json:"serial"`
	APIKey string `json:"api_key"`
}

type TelemetryData struct {
	Vibration float64 `json:"vibration"`
	XAccel    float64 `json:"x_accel"`
	YAccel    float64 `json:"y_accel"`
	ZAccel    float64 `json:"z_accel"`
}

type TelemetryRequest struct {
	Info      DeviceInfo    `json:"info"`
	Telemetry TelemetryData `json:"telemetry"`
}

type TelemetryResponse struct {
	Status       bool   `json:"status"`
	SerialNumber string `json:"serial_number"`
}