CREATE TABLE IF NOT EXISTS "user" (
  id uuid PRIMARY KEY NOT NULL,
  email VARCHAR NOT NULL,
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp DEFAULT now(),
  active BOOLEAN DEFAULT true
);

CREATE TABLE IF NOT EXISTS slug (
  id uuid PRIMARY KEY NOT NULL,
  user_id uuid NOT NULL,
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp DEFAULT now(),
  active BOOLEAN DEFAULT true,
  lat DECIMAL NOT NULL,
  long DECIMAL NOT NULL
);

CREATE TABLE IF NOT EXISTS merchant (
  id uuid PRIMARY KEY NOT NULL,
  user_id uuid NOT NULL,
  slug_id uuid NOT NULL,
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now(),
  name VARCHAR NOT NULL,
  active BOOLEAN DEFAULT true,
  lat DECIMAL NOT NULL,
  long DECIMAL NOT NULL
);

CREATE TABLE IF NOT EXISTS campaing (
  id uuid PRIMARY KEY NOT NULL,
  user_id uuid NOT NULL,
  slug_id uuid NOT NULL,
  merchant_id uuid NOT NULL,
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now(),
  active BOOLEAN DEFAULT true,
  lat DECIMAL NOT NULL,
  long DECIMAL NOT NULL,
  clicks INTEGER,
  impressions INTEGER
);