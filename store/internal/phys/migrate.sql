PRAGMA foreign_keys=off;
BEGIN;
CREATE TABLE new_list (
  id         integer PRIMARY KEY,
  type       integer NOT NULL CHECK (type IN (1, 2)),
  name       text    NOT NULL,
  moderators text
);
INSERT INTO new_list SELECT id, type, name, NULL from list;
UPDATE new_list SET moderators='rothskeller@gmail.com,sroth@sunnyvale.ca.gov' WHERE type=1;
DROP TABLE list;
ALTER TABLE new_list RENAME TO list;
COMMIT;
