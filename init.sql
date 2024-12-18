-- init.sql

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS clients (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS scooters (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    occupied BOOLEAN DEFAULT FALSE NOT NULL
);

CREATE TABLE IF NOT EXISTS locations (
    id SERIAL PRIMARY KEY,
    latitude DECIMAL(9, 6) NOT NULL,
    longitude DECIMAL(9, 6) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    scooter_uuid UUID NOT NULL REFERENCES scooters(uuid) ON DELETE CASCADE
);

-- Insert sample data
INSERT INTO clients (uuid, name) VALUES 
    ('f50e26c5-592b-4da3-b4eb-387d68bececa', 'Aiden'),
    ('84eb800d-a7e6-4f46-a5b8-0ea6f2c32184', 'Sophie'),
    ('82e2a4a2-f562-43a6-baf5-8cdaa1433a98', 'Lucas');

INSERT INTO scooters (name, occupied) VALUES 
    ('Ottawa Scooter 1', FALSE),
    ('Ottawa Scooter 2', FALSE),
    ('Ottawa Scooter 3', FALSE),
    ('Montreal Scooter 1', FALSE),
    ('Montreal Scooter 2', FALSE);


INSERT INTO locations (latitude, longitude, scooter_uuid) VALUES 
    (45.4215, -75.6972, (SELECT uuid FROM scooters WHERE name = 'Ottawa Scooter 1')),
    (45.4245, -75.6950, (SELECT uuid FROM scooters WHERE name = 'Ottawa Scooter 2')),
    (45.4290, -75.6880, (SELECT uuid FROM scooters WHERE name = 'Ottawa Scooter 3')),
    
    (45.5088, -73.5540, (SELECT uuid FROM scooters WHERE name = 'Montreal Scooter 1')),
    (45.5055, -73.5655, (SELECT uuid FROM scooters WHERE name = 'Montreal Scooter 2'));