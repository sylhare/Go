-- migrate:up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

ALTER TABLE authors
ADD COLUMN uuid UUID DEFAULT uuid_generate_v4();

-- migrate:down
ALTER TABLE authors
DROP
COLUMN uuid;