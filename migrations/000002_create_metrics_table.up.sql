CREATE TABLE IF NOT EXISTS metrics (
  id bigserial PRIMARY KEY,
  instance_id bigserial NOT NULL,
  cpu_usage real NOT NULL,
  memory_usage real NOT NULL,
  uptime interval NOT NULL,
  recorded_at TIMESTAMPTZ DEFAULT NOW(),
  CONSTRAINT fk_instance FOREIGN KEY(instance_id) REFERENCES instances(id)
)