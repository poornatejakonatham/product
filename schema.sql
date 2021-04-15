-- Create DB.
DROP DATABASE IF EXISTS product;
CREATE DATABASE product;

-- Connect to DB.
\c product;

DROP TABLE IF EXISTS product;
CREATE TABLE product (
    id SERIAL,
    name TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
    createdOn TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_product_id PRIMARY KEY (id),
    CONSTRAINT uk_product_name UNIQUE(name)
);

-- Create user.
-- DROP OWNED BY productuser;
-- DROP USER productuser;
-- CREATE USER productuser WITH PASSWORD 'Pr0ductUser@123';
-- GRANT ALL PRIVILEGES ON DATABASE product to productuser;