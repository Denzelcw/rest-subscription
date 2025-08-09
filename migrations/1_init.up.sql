CREATE TABLE IF NOT EXISTS user_subscriptions (
                                                  id SERIAL PRIMARY KEY,
                                                  service_name VARCHAR(255) NOT NULL,
    price INT NOT NULL,
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE OR REPLACE FUNCTION enforce_unique_active_subscription()
RETURNS TRIGGER AS $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM user_subscriptions
        WHERE user_id = NEW.user_id
          AND service_name = NEW.service_name
          AND (end_date IS NULL OR end_date > CURRENT_DATE)
          AND (TG_OP = 'INSERT' OR id <> NEW.id)
    ) AND (NEW.end_date IS NULL OR NEW.end_date > CURRENT_DATE) THEN
        RAISE EXCEPTION 'user % already has active subscription on %', NEW.user_id, NEW.service_name;
END IF;

RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_unique_active_subscription
    BEFORE INSERT OR UPDATE ON user_subscriptions
                         FOR EACH ROW
                         EXECUTE FUNCTION enforce_unique_active_subscription();