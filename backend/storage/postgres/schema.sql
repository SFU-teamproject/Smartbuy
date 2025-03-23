DROP TABLE IF EXISTS smartphones CASCADE;
CREATE TABLE smartphones (
    id SERIAL PRIMARY KEY,
    model TEXT NOT NULL,
    producer TEXT NOT NULL,
    memory INTEGER,
    CHECK(memory > 0),
    ram INTEGER,
    CHECK(ram > 0),
    display_size NUMERIC(3,2),
    CHECK(display_size >= 3 AND display_size <= 9.99),
    price INTEGER,
    CHECK(price >= 0),
    ratings_sum INTEGER DEFAULT 0,
    ratings_count INTEGER DEFAULT 0,
    image_path TEXT NOT NULL,
    description TEXT NOT NULL
);

DROP TABLE IF EXISTS users cascade;
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    password TEXT NOT NULL,
    role VARCHAR(10) DEFAULT 'user',
    CHECK (role IN ('admin', 'user')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX ON users (LOWER(name));

DROP TABLE IF EXISTS reviews cascade;
CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,
    smartphone_id INTEGER REFERENCES smartphones ON DELETE CASCADE,
    user_id INTEGER REFERENCES users ON DELETE CASCADE,
    rating SMALLINT,
    CHECK (rating >= 1 and rating <= 5),
    comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX ON reviews (smartphone_id, user_id);

CREATE OR REPLACE FUNCTION update_smartphone_rating()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE smartphones
        SET ratings_sum = ratings_sum + NEW.rating,
            ratings_count = ratings_count + 1
        WHERE id = NEW.smartphone_id;

    ELSIF TG_OP = 'UPDATE' THEN
        UPDATE smartphones
        SET ratings_sum = ratings_sum - OLD.rating + NEW.rating
        WHERE id = NEW.smartphone_id;

    ELSIF TG_OP = 'DELETE' THEN
        UPDATE smartphones
        SET ratings_sum = ratings_sum - OLD.rating,
            ratings_count = ratings_count - 1
        WHERE id = OLD.smartphone_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_smarthpone_rating
AFTER INSERT OR UPDATE OF rating OR DELETE ON reviews
FOR EACH ROW
EXECUTE FUNCTION update_smartphone_rating();

DROP TABLE IF EXISTS carts cascade;
CREATE TABLE carts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX ON carts(user_id);

DROP TABLE IF EXISTS cart_items cascade;
CREATE TABLE cart_items (
    id SERIAL PRIMARY KEY,
    cart_id INT NOT NULL REFERENCES carts ON DELETE CASCADE,
    smartphone_id INT NOT NULL REFERENCES smartphones ON DELETE CASCADE,
    quantity INT NOT NULL DEFAULT 1 CHECK (quantity > 0)
);
CREATE UNIQUE INDEX ON cart_items(cart_id, smartphone_id);

CREATE OR REPLACE FUNCTION create_cart_for_new_user()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO carts (user_id)
    VALUES (NEW.id);
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_create_cart
AFTER INSERT ON users
FOR EACH ROW
EXECUTE FUNCTION create_cart_for_new_user();

CREATE OR REPLACE FUNCTION update_cart_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE carts
    SET updated_at = CURRENT_TIMESTAMP
    WHERE id = (
        CASE
            WHEN TG_OP = 'INSERT' THEN NEW.cart_id
            WHEN TG_OP = 'UPDATE' THEN NEW.cart_id
            WHEN TG_OP = 'DELETE' THEN OLD.cart_id
        END
    );
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_cart_updated_at
AFTER INSERT OR UPDATE OR DELETE ON cart_items
FOR EACH ROW
EXECUTE FUNCTION update_cart_updated_at();