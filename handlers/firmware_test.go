package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestCheckFirmware_MissingVersion(t *testing.T) {
	handler := &FirmwareHandler{}

	req := httptest.NewRequest(http.MethodGet, "/check-firmware?serial_number=SN-001", nil)
	w := httptest.NewRecorder()

	handler.CheckFirmware(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCheckFirmware_MissingSerialNumber(t *testing.T) {
	handler := &FirmwareHandler{}

	req := httptest.NewRequest(http.MethodGet, "/check-firmware?v=1.0.0", nil)
	w := httptest.NewRecorder()

	handler.CheckFirmware(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCheckFirmware_WrongMethod(t *testing.T) {
	handler := &FirmwareHandler{}

	req := httptest.NewRequest(http.MethodPost, "/check-firmware", nil)
	w := httptest.NewRecorder()

	handler.CheckFirmware(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

func TestCheckFirmware_UpdateAvailable(t *testing.T) {
	// Create a mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "postgres")
	handler := &FirmwareHandler{DB: db}

	// Mock DB returns firmware version "2.0.0" for SN-001
	rows := sqlmock.NewRows([]string{"firmware"}).AddRow("2.0.0")
	mock.ExpectQuery("SELECT firmware FROM device_info").
		WithArgs("SN-001").
		WillReturnRows(rows)

	req := httptest.NewRequest(http.MethodGet, "/check-firmware?v=1.0.0&serial_number=SN-001", nil)
	w := httptest.NewRecorder()

	handler.CheckFirmware(w, req)

	if w.Code != http.StatusAccepted {
		t.Errorf("expected 202, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)

	if resp["status"] != false {
		t.Errorf("expected status false, got %v", resp["status"])
	}
}

func TestCheckFirmware_UpToDate(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "postgres")
	handler := &FirmwareHandler{DB: db}

	rows := sqlmock.NewRows([]string{"firmware"}).AddRow("1.0.0")
	mock.ExpectQuery("SELECT firmware FROM device_info").
		WithArgs("SN-001").
		WillReturnRows(rows)

	req := httptest.NewRequest(http.MethodGet, "/check-firmware?v=1.0.0&serial_number=SN-001", nil)
	w := httptest.NewRecorder()

	handler.CheckFirmware(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)

	if resp["status"] != true {
		t.Errorf("expected status true, got %v", resp["status"])
	}
}