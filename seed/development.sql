-- Development seed data
-- This file is safe to run repeatedly after a fresh database reset.

-- Countries
INSERT INTO countries (code, name)
VALUES
    ('HRV', 'Croatia'),
    ('USA', 'United States'),
    ('CHN', 'China')
ON CONFLICT (code) DO NOTHING;


-- Visa rules
INSERT INTO visa_rules (
    nationality_code,
    destination_code,
    visa_required,
    effective_from,
    effective_to,
    source_name,
    source_url
)
VALUES
(
    'HRV',
    'USA',
    false,
    '2026-01-01',
    NULL,
    'U.S. Department of State',
    'https://travel.state.gov/'
),
(
    'HRV',
    'CHN',
    false,
    '2026-01-01',
    '2027-02-28',
    'Chinese Visa Application Service',
    'https://www.visaforchina.cn/'
),
(
    'HRV',
    'CHN',
    true,
    '2027-03-01',
    NULL,
    'Chinese Visa Application Service',
    'https://www.visaforchina.cn/'
)
ON CONFLICT DO NOTHING;


-- ETA rules
INSERT INTO eta_rules (
    nationality_code,
    destination_code,
    eta_required,
    effective_from,
    effective_to,
    source_name,
    source_url
)
VALUES
(
    'HRV',
    'USA',
    true,
    '2026-01-01',
    NULL,
    'Development ETA Source',
    'https://example.com/eta'
),
(
    'HRV',
    'CHN',
    false,
    '2026-01-01',
    NULL,
    'Development ETA Source',
    'https://example.com/eta'
)
ON CONFLICT DO NOTHING;


-- Stay duration rules
INSERT INTO stay_rules (
    nationality_code,
    destination_code,
    max_stay_days,
    effective_from,
    effective_to,
    source_name,
    source_url
)
VALUES
(
    'HRV',
    'USA',
    90,
    '2026-01-01',
    NULL,
    'Development Stay Source',
    'https://example.com/stay'
),
(
    'HRV',
    'CHN',
    90,
    '2026-01-01',
    NULL,
    'Development Stay Source',
    'https://example.com/stay'
)
ON CONFLICT DO NOTHING;