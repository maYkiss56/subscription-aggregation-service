CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL CHECK (price > 0),
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NULL,

    CONSTRAINT valid_dates CHECK (end_date IS NULL OR end_date >= start_date)
);
