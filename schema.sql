CREATE TYPE category AS enum(
  'phone',
  'laptop',
  'desktop',
  'audio',
  'tablet',
  'earphone'
);
CREATE TABLE products(
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE,
  description TEXT,
  types category,
  created_date TIMESTAMP,
  updated_date TIMESTAMP
);
CREATE TYPE seller_source AS enum(
  'PCHOME',
  'MOMO',
  'SHOPEE',
  'AMAZON',
  'YAHOO',
  'OTHER'
);
CREATE TABLE price(
  id SERIAL PRIMARY KEY,
  products_id INTEGER REFERENCES products(id) NOT NULL,
  price NUMERIC(10, 2),
  discount NUMERIC(3),
  source_url TEXT,
  source_type seller_source,
  updated_date TIMESTAMP,
  UNIQUE (id, products_id)
);
CREATE TABLE history_price(
  id SERIAL PRIMARY KEY,
  products_id INTEGER REFERENCES products(id) NOT NULL,
  price NUMERIC(10, 2),
  discount NUMERIC(3),
  source_url TEXT,
  source_type seller_source,
  created_date TIMESTAMP
);
