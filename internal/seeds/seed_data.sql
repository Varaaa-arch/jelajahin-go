-- Seed Airlines
INSERT INTO airlines (id, code, name, headquarters_city, is_active, created_at, updated_at)
VALUES 
    (gen_random_uuid(), 'GA', 'Garuda Indonesia', 'Jakarta', true, NOW(), NOW()),
    (gen_random_uuid(), 'QG', 'Citilink', 'Jakarta', true, NOW(), NOW()),
    (gen_random_uuid(), 'IW', 'Indonesia AirAsia', 'Tangerang', true, NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Seed Airports
INSERT INTO airports (id, code, name, city, country, timezone, is_active, created_at, updated_at)
VALUES 
    (gen_random_uuid(), 'CGK', 'Soekarno-Hatta International', 'Jakarta', 'Indonesia', 'Asia/Jakarta', true, NOW(), NOW()),
    (gen_random_uuid(), 'DPS', 'Ngurah Rai International', 'Bali', 'Indonesia', 'Asia/Jakarta', true, NOW(), NOW()),
    (gen_random_uuid(), 'SUB', 'Juanda International', 'Surabaya', 'Indonesia', 'Asia/Jakarta', true, NOW(), NOW()),
    (gen_random_uuid(), 'SYD', 'Sydney Kingsford Smith', 'Sydney', 'Australia', 'Australia/Sydney', true, NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Seed Aircraft Types
INSERT INTO aircraft_types (id, name, manufacturer, model, total_seats, created_at)
VALUES 
    (gen_random_uuid(), 'Boeing 737-800', 'Boeing', '737-800', 189, NOW()),
    (gen_random_uuid(), 'Airbus A320', 'Airbus', 'A320', 194, NOW())
ON CONFLICT (name) DO NOTHING;

-- Seed Seat Classes
INSERT INTO seat_classes (id, name, display_name, baggage_allowance_kg, created_at)
VALUES 
    (gen_random_uuid(), 'economy', 'Economy', 20, NOW()),
    (gen_random_uuid(), 'premium_economy', 'Premium Economy', 30, NOW()),
    (gen_random_uuid(), 'business', 'Business', 40, NOW()),
    (gen_random_uuid(), 'first', 'First', 50, NOW())
ON CONFLICT (name) DO NOTHING;

-- Seed Aircraft Instances
INSERT INTO aircraft (id, aircraft_type_id, airline_id, registration_number, manufacture_year, is_active, created_at, updated_at)
SELECT 
    gen_random_uuid(),
    (SELECT id FROM aircraft_types WHERE name = 'Boeing 737-800' LIMIT 1),
    (SELECT id FROM airlines WHERE code = 'GA' LIMIT 1),
    'PK-GIA',
    2015,
    true,
    NOW(),
    NOW()
WHERE NOT EXISTS (SELECT 1 FROM aircraft WHERE registration_number = 'PK-GIA')
UNION ALL
SELECT 
    gen_random_uuid(),
    (SELECT id FROM aircraft_types WHERE name = 'Airbus A320' LIMIT 1),
    (SELECT id FROM airlines WHERE code = 'QG' LIMIT 1),
    'PK-CIJ',
    2018,
    true,
    NOW(),
    NOW()
WHERE NOT EXISTS (SELECT 1 FROM aircraft WHERE registration_number = 'PK-CIJ')
UNION ALL
SELECT 
    gen_random_uuid(),
    (SELECT id FROM aircraft_types WHERE name = 'Boeing 737-800' LIMIT 1),
    (SELECT id FROM airlines WHERE code = 'IW' LIMIT 1),
    'PK-AWU',
    2016,
    true,
    NOW(),
    NOW()
WHERE NOT EXISTS (SELECT 1 FROM aircraft WHERE registration_number = 'PK-AWU');

-- Seed Aircraft Seats untuk Boeing 737-800 (189 seats)
-- Economy: Rows 1-20 (A-F) = 120 seats
INSERT INTO aircraft_seats (id, aircraft_type_id, seat_class_id, seat_number, row_number, column_letter, is_exit_row)
SELECT 
    gen_random_uuid(),
    (SELECT id FROM aircraft_types WHERE name = 'Boeing 737-800' LIMIT 1),
    (SELECT id FROM seat_classes WHERE name = 'economy' LIMIT 1),
    row_num || column_letter,
    row_num,
    column_letter,
    false
FROM (
    SELECT generate_series(1, 20) as row_num
) rows
CROSS JOIN (
    SELECT unnest(ARRAY['A', 'B', 'C', 'D', 'E', 'F']) as column_letter
) columns
ON CONFLICT (aircraft_type_id, seat_number) DO NOTHING;

-- Business: Rows 21-24 (A-F) = 24 seats
INSERT INTO aircraft_seats (id, aircraft_type_id, seat_class_id, seat_number, row_number, column_letter, is_exit_row)
SELECT 
    gen_random_uuid(),
    (SELECT id FROM aircraft_types WHERE name = 'Boeing 737-800' LIMIT 1),
    (SELECT id FROM seat_classes WHERE name = 'business' LIMIT 1),
    row_num || column_letter,
    row_num,
    column_letter,
    false
FROM (
    SELECT generate_series(21, 24) as row_num
) rows
CROSS JOIN (
    SELECT unnest(ARRAY['A', 'B', 'C', 'D', 'E', 'F']) as column_letter
) columns
ON CONFLICT (aircraft_type_id, seat_number) DO NOTHING;

-- First: Rows 25-26 (A-F) = 12 seats
INSERT INTO aircraft_seats (id, aircraft_type_id, seat_class_id, seat_number, row_number, column_letter, is_exit_row)
SELECT 
    gen_random_uuid(),
    (SELECT id FROM aircraft_types WHERE name = 'Boeing 737-800' LIMIT 1),
    (SELECT id FROM seat_classes WHERE name = 'first' LIMIT 1),
    row_num || column_letter,
    row_num,
    column_letter,
    false
FROM (
    SELECT generate_series(25, 26) as row_num
) rows
CROSS JOIN (
    SELECT unnest(ARRAY['A', 'B', 'C', 'D', 'E', 'F']) as column_letter
) columns
ON CONFLICT (aircraft_type_id, seat_number) DO NOTHING;

-- Seed Routes
INSERT INTO routes (id, airline_id, origin_airport_id, destination_airport_id, flight_number_prefix, distance_km, estimated_duration_minutes, is_active, created_at, updated_at)
SELECT 
    gen_random_uuid(),
    (SELECT id FROM airlines WHERE code = 'GA' LIMIT 1),
    (SELECT id FROM airports WHERE code = 'CGK' LIMIT 1),
    (SELECT id FROM airports WHERE code = 'DPS' LIMIT 1),
    'GA',
    1400,
    210,
    true,
    NOW(),
    NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM routes 
    WHERE airline_id = (SELECT id FROM airlines WHERE code = 'GA' LIMIT 1)
    AND origin_airport_id = (SELECT id FROM airports WHERE code = 'CGK' LIMIT 1)
    AND destination_airport_id = (SELECT id FROM airports WHERE code = 'DPS' LIMIT 1)
)
UNION ALL
SELECT 
    gen_random_uuid(),
    (SELECT id FROM airlines WHERE code = 'GA' LIMIT 1),
    (SELECT id FROM airports WHERE code = 'DPS' LIMIT 1),
    (SELECT id FROM airports WHERE code = 'SUB' LIMIT 1),
    'GA',
    600,
    90,
    true,
    NOW(),
    NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM routes 
    WHERE airline_id = (SELECT id FROM airlines WHERE code = 'GA' LIMIT 1)
    AND origin_airport_id = (SELECT id FROM airports WHERE code = 'DPS' LIMIT 1)
    AND destination_airport_id = (SELECT id FROM airports WHERE code = 'SUB' LIMIT 1)
)
UNION ALL
SELECT 
    gen_random_uuid(),
    (SELECT id FROM airlines WHERE code = 'GA' LIMIT 1),
    (SELECT id FROM airports WHERE code = 'CGK' LIMIT 1),
    (SELECT id FROM airports WHERE code = 'SUB' LIMIT 1),
    'GA',
    1850,
    270,
    true,
    NOW(),
    NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM routes 
    WHERE airline_id = (SELECT id FROM airlines WHERE code = 'GA' LIMIT 1)
    AND origin_airport_id = (SELECT id FROM airports WHERE code = 'CGK' LIMIT 1)
    AND destination_airport_id = (SELECT id FROM airports WHERE code = 'SUB' LIMIT 1)
);

-- Seed Sample Flights
INSERT INTO flights (id, route_id, aircraft_id, flight_number, departure_date, departure_time, arrival_time, base_price, tax_surcharge, fuel_surcharge, status, seats_available, created_at, updated_at)
SELECT 
    gen_random_uuid(),
    (SELECT id FROM routes 
     WHERE flight_number_prefix = 'GA' 
     AND origin_airport_id = (SELECT id FROM airports WHERE code = 'CGK')
     AND destination_airport_id = (SELECT id FROM airports WHERE code = 'DPS')
     LIMIT 1),
    (SELECT id FROM aircraft WHERE registration_number = 'PK-GIA' LIMIT 1),
    'GA100',
    '2024-09-15'::date,
    '08:00:00'::time,
    '10:30:00'::time,
    1500000,
    150000,
    50000,
    'scheduled',
    189,
    NOW(),
    NOW()
WHERE NOT EXISTS (SELECT 1 FROM flights WHERE flight_number = 'GA100' AND departure_date = '2024-09-15'::date)
UNION ALL
SELECT 
    gen_random_uuid(),
    (SELECT id FROM routes 
     WHERE flight_number_prefix = 'GA' 
     AND origin_airport_id = (SELECT id FROM airports WHERE code = 'CGK')
     AND destination_airport_id = (SELECT id FROM airports WHERE code = 'DPS')
     LIMIT 1),
    (SELECT id FROM aircraft WHERE registration_number = 'PK-GIA' LIMIT 1),
    'GA101',
    '2024-09-15'::date,
    '14:00:00'::time,
    '16:30:00'::time,
    1500000,
    150000,
    50000,
    'scheduled',
    189,
    NOW(),
    NOW()
WHERE NOT EXISTS (SELECT 1 FROM flights WHERE flight_number = 'GA101' AND departure_date = '2024-09-15'::date)
UNION ALL
SELECT 
    gen_random_uuid(),
    (SELECT id FROM routes 
     WHERE flight_number_prefix = 'GA' 
     AND origin_airport_id = (SELECT id FROM airports WHERE code = 'CGK')
     AND destination_airport_id = (SELECT id FROM airports WHERE code = 'DPS')
     LIMIT 1),
    (SELECT id FROM aircraft WHERE registration_number = 'PK-GIA' LIMIT 1),
    'GA102',
    '2024-09-16'::date,
    '08:00:00'::time,
    '10:30:00'::time,
    1500000,
    150000,
    50000,
    'scheduled',
    189,
    NOW(),
    NOW()
WHERE NOT EXISTS (SELECT 1 FROM flights WHERE flight_number = 'GA102' AND departure_date = '2024-09-16'::date);

-- Seed Flight Seats (untuk semua flights)
INSERT INTO flight_seats (id, flight_id, aircraft_seat_id, current_price, is_available, created_at, updated_at)
SELECT 
    gen_random_uuid(),
    f.id,
    aset.id,
    CASE 
        WHEN sc.name = 'first' THEN 3000000
        WHEN sc.name = 'business' THEN 2000000
        WHEN sc.name = 'premium_economy' THEN 1800000
        ELSE 1500000
    END,
    true,
    NOW(),
    NOW()
FROM flights f
JOIN aircraft a ON f.aircraft_id = a.id
JOIN aircraft_seats aset ON aset.aircraft_type_id = a.aircraft_type_id
JOIN seat_classes sc ON aset.seat_class_id = sc.id
WHERE NOT EXISTS (
    SELECT 1 FROM flight_seats 
    WHERE flight_id = f.id AND aircraft_seat_id = aset.id
);

-- Verify counts
SELECT 'Airlines' as entity, COUNT(*) as count FROM airlines
UNION ALL
SELECT 'Airports', COUNT(*) FROM airports
UNION ALL
SELECT 'Aircraft Types', COUNT(*) FROM aircraft_types
UNION ALL
SELECT 'Aircraft', COUNT(*) FROM aircraft
UNION ALL
SELECT 'Routes', COUNT(*) FROM routes
UNION ALL
SELECT 'Flights', COUNT(*) FROM flights
UNION ALL
SELECT 'Aircraft Seats', COUNT(*) FROM aircraft_seats
UNION ALL
SELECT 'Flight Seats', COUNT(*) FROM flight_seats;
