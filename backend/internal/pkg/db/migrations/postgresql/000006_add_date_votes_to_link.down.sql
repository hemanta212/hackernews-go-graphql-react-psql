BEGIN;

  ALTER TABLE Links

  DROP COLUMN CreatedAt;

COMMIT;