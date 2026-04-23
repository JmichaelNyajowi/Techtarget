package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestSubmitTelemetry_WrongMethod(t *testing.T) {
	handler := &TelemetryHandler{}

	req := httptest.NewRequest(http.MethodGet, "/telemetry", nil)
	w := httptest.NewRecorder()

	handler.Submit(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

func TestSubmitTelemetry_MissingSerial(t *testing.T) {
	handler := &TelemetryHandler{}

	body := `{"info": {"model": "M1", "firmware": "v2.34"}, "telemetry": {}}`
	req := httptest.NewRequest(http.MethodPost, "/telemetry", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Submit(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestSubmitTelemetry_Success(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "postgres")
	handler := &TelemetryHandler{DB: db}

	// Mock the upsert into device_info
	mock.ExpectExec("INSERT INTO device_info").
		WithArgs("WB234A2", "M1", "v2.34").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock the insert into telemetry
	mock.ExpectExec("INSERT INTO telemetry").
		WithArgs("WB234A2", 234.3, 24.2, 243.1, 67.4).
		WillReturnResult(sqlmock.NewResult(1, 1))

	payload := map[string]interface{}{
		"info": map[string]string{
			"serial":   "WB234A2",
			"model":    "M1",
			"firmware": "v2.34",
		},
		"telemetry": map[string]float64{
			"vibration": 234.3,
			"x_accel":   24.2,
			"y_accel":   243.1,
			"z_accel":   67.4,
		},
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/telemetry", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Submit(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)

	if resp["status"] != true {
		t.Errorf("expected status true, got %v", resp["status"])
	}
}