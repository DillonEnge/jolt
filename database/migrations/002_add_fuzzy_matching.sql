-- Write your migrate up statements here

CREATE EXTENSION IF NOT EXISTS "fuzzystrmatch";

---- create above / drop below ----

DROP EXTENSION "fuzzystrmatch";

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
