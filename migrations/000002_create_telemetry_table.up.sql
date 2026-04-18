CREATE TABLE telemetry (
    id SERIAL PRIMARY KEY,
    device_serial VARCHAR(100) REFERENCES device_info(serial_number) ON UPDATE CASCADE,
    vibration FLOAT,
    x_accel FLOAT,
    y_accel FLOAT,
    z_accel FLOAT,
    recorded_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);