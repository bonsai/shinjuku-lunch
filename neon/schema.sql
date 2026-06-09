CREATE TABLE IF NOT EXISTS areas (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS genres (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS restaurants (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  area_id INTEGER REFERENCES areas(id),
  genre_id INTEGER REFERENCES genres(id),
  address TEXT,
  station TEXT,
  walk_min INTEGER,
  latitude DOUBLE PRECISION,
  longitude DOUBLE PRECISION,
  business_hours TEXT,
  url_tabelog TEXT,
  url_hotpepper TEXT,
  notes TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS lunch_logs (
  id SERIAL PRIMARY KEY,
  restaurant_id INTEGER REFERENCES restaurants(id),
  menu TEXT NOT NULL,
  price INTEGER NOT NULL,
  rating INTEGER CHECK (rating BETWEEN 1 AND 5),
  comment TEXT,
  revisit BOOLEAN DEFAULT FALSE,
  visited_date DATE DEFAULT CURRENT_DATE,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_restaurants_area ON restaurants(area_id);
CREATE INDEX IF NOT EXISTS idx_restaurants_genre ON restaurants(genre_id);
CREATE INDEX IF NOT EXISTS idx_lunch_logs_restaurant ON lunch_logs(restaurant_id);
