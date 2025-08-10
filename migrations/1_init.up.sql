CREATE TABLE IF NOT EXISTS user_subscriptions (
                                                  id SERIAL PRIMARY KEY,
                                                  service_name VARCHAR(255) NOT NULL,
    price INT NOT NULL,
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_subscription UNIQUE (user_id, service_name, start_date, end_date)
    );