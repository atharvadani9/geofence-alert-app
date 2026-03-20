CREATE TABLE geofences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    caregiver_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tracked_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    radius INT NOT NULL,
    latitude NUMERIC NOT NULL,
    longitude NUMERIC NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE locations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tracked_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    latitude NUMERIC NOT NULL,
    longitude NUMERIC NOT NULL,
    timestamp TIMESTAMPTZ DEFAULT NOW(),
    accuracy NUMERIC NOT NULL,
);
CREATE INDEX idx_locations_tracked_user_id ON locations(tracked_user_id, timestamp);

CREATE TABLE alerts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    geofence_id UUID NOT NULL REFERENCES geofences(id) ON DELETE CASCADE,
    tracked_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    alert_type VARCHAR(20) NOT NULL CHECK (alert_type IN ('entry', 'exit')),
    triggered_at TIMESTAMPTZ DEFAULT NOW(),
    latitude NUMERIC NOT NULL,
    longitude NUMERIC NOT NULL,
    is_acknowledged BOOLEAN NOT NULL DEFAULT FALSE,
    acknowledged_at TIMESTAMPTZ,
);
