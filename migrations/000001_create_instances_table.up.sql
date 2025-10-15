CREATE TABLE IF NOT EXISTS instances (
  id bigserial PRIMARY KEY,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  name text NOT NULL,
  ip text NOT NULL,
  refresh_rate interval NOT NULL DEFAULT '5 minutes',
  version integer NOT NULL DEFAULT 1
);