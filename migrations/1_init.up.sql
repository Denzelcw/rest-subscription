CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE IF NOT EXISTS user_subscriptions (
    id SERIAL PRIMARY KEY,
    service_name VARCHAR(255) NOT NULL,
    price INT NOT NULL,
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unique_subscription UNIQUE (user_id, service_name, start_date, end_date),

    CONSTRAINT valid_period CHECK (
                                      end_date IS NULL OR start_date < end_date
                                  )
    );

ALTER TABLE user_subscriptions
    ADD CONSTRAINT no_overlap
    EXCLUDE USING gist (
        user_id WITH =,
        service_name WITH =,
        daterange(
            start_date,
            COALESCE(end_date, DATE '9999-12-31'),
            '[]'
        ) WITH &&
    );