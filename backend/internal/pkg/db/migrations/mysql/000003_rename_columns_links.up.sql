BEGIN;
  ALTER TABLE `Links`
  RENAME COLUMN Title TO Description,
  RENAME COLUMN Address TO Url;
COMMIT;
