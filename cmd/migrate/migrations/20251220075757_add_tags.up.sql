CREATE TABLE IF NOT EXISTS tags (
  id bigserial PRIMARY KEY,
  title text NOT NULL,
  updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);