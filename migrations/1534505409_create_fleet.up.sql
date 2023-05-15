CREATE TABLE driver
(
    id TEXT PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    created_at INTEGER
);

CREATE TABLE vehicle
(
    id TEXT PRIMARY KEY,
    vin TEXT NOT NULL,
    type TEXT NOT NULL,
    driver_id TEXT,
    created_at INTEGER,
    FOREIGN KEY (driver_id) REFERENCES driver (id) ON DELETE CASCADE
);

CREATE TABLE fleet
(
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    capacity INTEGER,
    vehicle_id TEXT,
    FOREIGN KEY (vehicle_id) REFERENCES vehicle (id) ON DELETE CASCADE
);