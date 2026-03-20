CREATE TABLE IF NOT EXISTS users(
   ID BIGSERIAL PRIMARY KEY,
   created_at TIMESTAMPTZ,
   updated_at TIMESTAMPTZ,
   deleted_at TIMESTAMPTZ,
   name TEXT NOT NULL,
   middle_name TEXT NOT NULL,
   email TEXT UNIQUE NOT NULL,
   password TEXT NOT NULL,
   phone TEXT NOT NULL,
   roles TEXT[] NOT NULL,
   active_rol TEXT NOT NULL
);