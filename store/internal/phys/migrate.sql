PRAGMA foreign_keys=OFF;
BEGIN;
CREATE TABLE new_person (
  id                 integer PRIMARY KEY,
  volgistics_id      integer,
  informal_name      text    NOT NULL,
  formal_name        text    NOT NULL,
  sort_name          text    NOT NULL,
  call_sign          text,
  pronouns           text,
  email              text,
  email2             text,
  cell_phone         text,
  home_phone         text,
  work_phone         text,
  password           text,
  bad_login_count    integer,
  bad_login_time     text,            -- YYYY-MM-DDTHH:MM:SS (local)
  pwreset_token      text    UNIQUE,
  pwreset_time       text,            -- YYYY-MM-DDTHH:MM:SS (local)
  unsubscribe_token  text    UNIQUE,
  hours_token        text    UNIQUE,
  identification     integer NOT NULL DEFAULT 0,
  birthdate          text,            -- YYYY-MM-DD
  flags              integer NOT NULL DEFAULT 0
);
INSERT INTO new_person SELECT id, volgistics_id, informal_name, formal_name, sort_name, call_sign, pronouns, email, email2, cell_phone, home_phone, work_phone, password, bad_login_count, bad_login_time, pwreset_token, pwreset_time, unsubscribe_token, hours_token, identification, birthdate, flags FROM person;
DROP TABLE person;
ALTER TABLE new_person RENAME TO person;
CREATE INDEX person_sort_name_idx ON person (sort_name);
COMMIT;
