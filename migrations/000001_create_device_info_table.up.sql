CREATE TABLE device_info (
    serial_number VARCHAR(100) NOT NULL PRIMARY KEY,
    model VARCHAR(50),
    firmware VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    last_seen TIMESTAMPTZ
);
