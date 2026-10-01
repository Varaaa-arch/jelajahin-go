-- Seed Airlines
INSERT INTO airlines (id, code, name, created_at, updated_at)
VALUES
    (gen_random_uuid(), 'GA', 'Garuda Indonesia', NOW(), NOW()),
    (gen_random_uuid(), 'QG', 'Citilink', NOW(), NOW()),
    (gen_random_uuid(), 'IW', 'Indonesia AirAsia', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Seed Airports
INSERT INTO airports (id, code, name, city, country, timezone, created_at, updated_at)
VALUES
    (gen_random_uuid(), 'CGK', 'Soekarno-Hatta International', 'Jakarta', 'Indonesia', 'Asia/Jakarta', NOW(), NOW()),
    (gen_random_uuid(), 'DPS', 'Ngurah Rai International', 'Bali', 'Indonesia', 'Asia/Jakarta', NOW(), NOW()),
    (gen_random_uuid(), 'SUB', 'Juanda International', 'Surabaya', 'Indonesia', 'Asia/Jakarta', NOW(), NOW()),
    (gen_random_uuid(), 'SYD', 'Sydney Kingsford Smith', 'Sydney', 'Australia', 'Australia/Sydney', NOW(), NOW())
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
INSERT INTO routes (id, airline_id, origin_airport_id, destination_airport_id, flight_number_prefix, distance_km, estimated_duration_minutes, created_at, updated_at)
SELECT 
    gen_random_uuid(),
    (SELECT id FROM airlines WHERE code = 'GA' LIMIT 1),
    (SELECT id FROM airports WHERE code = 'CGK' LIMIT 1),
    (SELECT id FROM airports WHERE code = 'DPS' LIMIT 1),
    'GA',
    1400,
    210,
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

-- ============================================================
-- SEED BANYAK (domestik): bandara + maskapai + rute + penerbangan
-- Idempoten: ON CONFLICT / WHERE NOT EXISTS, aman dijalankan ulang
-- dan aman digabung dengan hasil `php artisan db:seed` Laravel
-- (nomor penerbangan Go di range 5xxxx, Laravel di 1xxx-12xxx).
-- ============================================================

-- Tambahan bandara domestik (sama dengan Laravel AirportSeeder)
INSERT INTO airports (id, code, name, city, country, timezone, created_at, updated_at)
VALUES
    (gen_random_uuid(), 'HLP', 'Halim Perdanakusuma International', 'Jakarta', 'Indonesia', 'Asia/Jakarta', NOW(), NOW()),
    (gen_random_uuid(), 'KNO', 'Kualanamu International', 'Medan', 'Indonesia', 'Asia/Jakarta', NOW(), NOW()),
    (gen_random_uuid(), 'UPG', 'Sultan Hasanuddin International', 'Makassar', 'Indonesia', 'Asia/Makassar', NOW(), NOW()),
    (gen_random_uuid(), 'YIA', 'Yogyakarta International', 'Yogyakarta', 'Indonesia', 'Asia/Jakarta', NOW(), NOW()),
    (gen_random_uuid(), 'SRG', 'Jenderal Ahmad Yani International', 'Semarang', 'Indonesia', 'Asia/Jakarta', NOW(), NOW()),
    (gen_random_uuid(), 'SOC', 'Adi Soemarmo International', 'Solo', 'Indonesia', 'Asia/Jakarta', NOW(), NOW()),
    (gen_random_uuid(), 'BDO', 'Husein Sastranegara International', 'Bandung', 'Indonesia', 'Asia/Jakarta', NOW(), NOW()),
    (gen_random_uuid(), 'PLM', 'Sultan Mahmud Badaruddin II International', 'Palembang', 'Indonesia', 'Asia/Jakarta', NOW(), NOW()),
    (gen_random_uuid(), 'BTH', 'Hang Nadim International', 'Batam', 'Indonesia', 'Asia/Jakarta', NOW(), NOW()),
    (gen_random_uuid(), 'PNK', 'Supadio International', 'Pontianak', 'Indonesia', 'Asia/Jakarta', NOW(), NOW()),
    (gen_random_uuid(), 'BPN', 'Sultan Aji Muhammad Sulaiman', 'Balikpapan', 'Indonesia', 'Asia/Makassar', NOW(), NOW()),
    (gen_random_uuid(), 'LOP', 'Zainuddin Abdul Madjid International', 'Lombok', 'Indonesia', 'Asia/Makassar', NOW(), NOW()),
    (gen_random_uuid(), 'MDC', 'Sam Ratulangi International', 'Manado', 'Indonesia', 'Asia/Makassar', NOW(), NOW()),
    (gen_random_uuid(), 'AMQ', 'Pattimura International', 'Ambon', 'Indonesia', 'Asia/Jayapura', NOW(), NOW()),
    (gen_random_uuid(), 'DJJ', 'Sentani International', 'Jayapura', 'Indonesia', 'Asia/Jayapura', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Tambahan maskapai domestik (sama dengan Laravel AirlineSeeder)
INSERT INTO airlines (id, code, name, created_at, updated_at)
VALUES
    (gen_random_uuid(), 'JT', 'Lion Air', NOW(), NOW()),
    (gen_random_uuid(), 'ID', 'Batik Air', NOW(), NOW()),
    (gen_random_uuid(), 'SJ', 'Sriwijaya Air', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Tambahan rute domestik (airline, origin, destination, jarak_km, durasi_mnt)
-- Ditulis set-based via CTE agar ringkas dan idempoten.
WITH tmp_new_routes (airline_code, origin_code, destination_code, distance_km, duration_min) AS (
    VALUES
        ('JT', 'CGK', 'SUB', 690, 90),
        ('JT', 'SUB', 'CGK', 690, 90),
        ('JT', 'CGK', 'DPS', 980, 110),
        ('JT', 'CGK', 'KNO', 1390, 150),
        ('JT', 'CGK', 'BPN', 1250, 140),
        ('JT', 'BPN', 'CGK', 1250, 140),
        ('JT', 'CGK', 'LOP', 1070, 120),
        ('ID', 'CGK', 'KNO', 1390, 150),
        ('ID', 'KNO', 'CGK', 1390, 150),
        ('ID', 'CGK', 'PLM', 430, 60),
        ('ID', 'PLM', 'CGK', 430, 60),
        ('ID', 'CGK', 'BTH', 880, 105),
        ('ID', 'CGK', 'SRG', 440, 60),
        ('ID', 'HLP', 'DPS', 980, 110),
        ('QG', 'CGK', 'UPG', 1390, 155),
        ('QG', 'CGK', 'YIA', 520, 70),
        ('QG', 'YIA', 'SUB', 260, 50),
        ('QG', 'DPS', 'LOP', 170, 40),
        ('QG', 'UPG', 'BPN', 500, 70),
        ('IW', 'DPS', 'CGK', 980, 110),
        ('IW', 'DPS', 'SUB', 310, 55),
        ('IW', 'DPS', 'UPG', 640, 80),
        ('SJ', 'CGK', 'PNK', 730, 95),
        ('SJ', 'UPG', 'MDC', 900, 110),
        ('SJ', 'CGK', 'SOC', 510, 70),
        ('GA', 'CGK', 'YIA', 520, 70),
        ('GA', 'CGK', 'KNO', 1390, 150),
        ('GA', 'CGK', 'UPG', 1390, 155)
)
INSERT INTO routes (id, airline_id, origin_airport_id, destination_airport_id, flight_number_prefix, distance_km, estimated_duration_minutes, created_at, updated_at)
SELECT
    gen_random_uuid(),
    al.id,
    o.id,
    d.id,
    t.airline_code,
    t.distance_km,
    t.duration_min,
    NOW(),
    NOW()
FROM tmp_new_routes t
JOIN airlines al ON al.code = t.airline_code
JOIN airports o ON o.code = t.origin_code
JOIN airports d ON d.code = t.destination_code
WHERE NOT EXISTS (
    SELECT 1 FROM routes r
    WHERE r.airline_id = al.id
    AND r.origin_airport_id = o.id
    AND r.destination_airport_id = d.id
);

-- Penerbangan massal: 14 hari ke depan × 2 slot/hari untuk tiap rute aktif.
-- Nomor prefix + 5xxxx (tidak tabrakan dengan seed Laravel 1xxx-12xxx).
-- Urutan rn deterministik (kode maskapai/bandara) agar idempoten.
WITH route_list AS (
    SELECT
        r.id AS route_id,
        r.flight_number_prefix AS prefix,
        ROW_NUMBER() OVER (ORDER BY al.code, o.code, dest.code) AS rn,
        COALESCE(
            (SELECT a2.id FROM aircraft a2 WHERE a2.airline_id = r.airline_id LIMIT 1),
            (SELECT a3.id FROM aircraft a3 LIMIT 1)
        ) AS aircraft_id,
        COALESCE(r.estimated_duration_minutes, 90) AS dur
    FROM routes r
    JOIN airlines al ON al.id = r.airline_id
    JOIN airports o ON o.id = r.origin_airport_id
    JOIN airports dest ON dest.id = r.destination_airport_id
),
days AS (
    SELECT generate_series(0, 13) AS d
),
slots AS (
    SELECT * FROM (VALUES (0, '07:00:00'::time), (1, '17:30:00'::time)) AS t(s, dep)
),
gen AS (
    SELECT
        rl.route_id,
        rl.aircraft_id,
        rl.prefix || (50000 + rl.rn * 30 + days.d * 2 + slots.s)::text AS fn,
        (CURRENT_DATE + days.d)::date AS dep_date,
        slots.dep AS dep_time,
        (slots.dep + (rl.dur + 15) * INTERVAL '1 minute')::time AS arr_time,
        (600000 + rl.rn * 20000 + days.d * 10000)::numeric(12,2) AS price
    FROM route_list rl
    CROSS JOIN days
    CROSS JOIN slots
)
INSERT INTO flights (id, route_id, aircraft_id, flight_number, departure_date, departure_time, arrival_time, base_price, tax_surcharge, fuel_surcharge, status, seats_available, created_at, updated_at)
SELECT
    gen_random_uuid(),
    gen.route_id,
    gen.aircraft_id,
    gen.fn,
    gen.dep_date,
    gen.dep_time,
    gen.arr_time,
    gen.price,
    50000,
    40000,
    'scheduled',
    60,
    NOW(),
    NOW()
FROM gen
WHERE NOT EXISTS (SELECT 1 FROM flights f WHERE f.flight_number = gen.fn);

-- Kursi untuk penerbangan baru (pola sama seperti blok flight_seats di atas)
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
-- Hanya penerbangan range Go (suffix angka >= 50000), bukan seed Laravel.
WHERE (substring(f.flight_number from '[0-9]+$'))::int >= 50000
AND NOT EXISTS (
    SELECT 1 FROM flight_seats
    WHERE flight_id = f.id AND aircraft_seat_id = aset.id
);
