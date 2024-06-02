DROP TABLE IF EXISTS class;
CREATE TABLE class (
  id        integer PRIMARY KEY,
  type      integer NOT NULL,
  start     text    NOT NULL, -- YYYY-MM-DD
  en_desc   text    NOT NULL,
  es_desc   text    NOT NULL,
  elimit    integer NOT NULL CHECK (elimit >= 0),
  referrals integer NOT NULL
);
CREATE UNIQUE INDEX class_start_idx ON class (start, type);

DROP TABLE IF EXISTS classreg;
CREATE TABLE classreg (
  id            integer PRIMARY KEY,
  class         integer NOT NULL REFERENCES class,
  person        integer          REFERENCES person,
  registered_by integer NOT NULL REFERENCES person,
  first_name    text    NOT NULL,
  last_name     text    NOT NULL,
  email         text,
  cell_phone    text
);
CREATE INDEX classreg_class_index ON classreg (class);
CREATE INDEX classreg_person_index ON classreg (person);
CREATE INDEX classreg_regby_index ON classreg (registered_by);

DROP TABLE IF EXISTS document;
CREATE TABLE document (
  id       integer PRIMARY KEY,
  folder   integer NOT NULL REFERENCES folder,
  name     text    NOT NULL,
  url      text,
  archived boolean NOT NULL DEFAULT false
);
CREATE UNIQUE INDEX document_name_idx ON document (folder, name) WHERE NOT archived;

DROP TABLE IF EXISTS event;
CREATE TABLE event (
  id         integer PRIMARY KEY,
  name       text    NOT NULL,
  start      text    NOT NULL,                      -- YYYY-MM-DDTHH:MM (local)
  end        text    NOT NULL CHECK (end >= start), -- YYYY-MM-DDTHH:MM (local)
  venue      integer REFERENCES venue, -- *not* ON DELETE CASCADE
  venue_url  text,
  activation text,
  details    text,
  flags      integer NOT NULL DEFAULT 0
);
CREATE INDEX event_start_idx ON event (start);
CREATE INDEX event_venue_idx ON event (venue);

DROP TABLE IF EXISTS folder;
CREATE TABLE folder (
  id        integer PRIMARY KEY,
  parent    integer NOT NULL REFERENCES folder,
  name      text    NOT NULL,
  url_name  text    NOT NULL,
  view_org  integer NOT NULL,
  view_priv integer NOT NULL,
  edit_org  integer NOT NULL,
  edit_priv integer NOT NULL
);
CREATE UNIQUE INDEX folder_urlname_idx ON folder (parent, url_name);
INSERT INTO folder VALUES (1, 1, 'Files', '', 0, 0, 1, 3);

DROP TABLE IF EXISTS list;
CREATE TABLE list (
  id         integer PRIMARY KEY,
  type       integer NOT NULL CHECK (type IN (1, 2)),
  name       text    NOT NULL,
  moderators text
);

DROP TABLE IF EXISTS list_person;
CREATE TABLE list_person (
  list   integer NOT NULL REFERENCES list ON DELETE CASCADE,
  person integer NOT NULL REFERENCES person ON DELETE CASCADE,
  sender boolean NOT NULL DEFAULT false,
  sub    boolean NOT NULL DEFAULT false,
  unsub  boolean NOT NULL DEFAULT false,
  PRIMARY KEY (list, person)
) WITHOUT ROWID;
CREATE INDEX list_person_person_idx ON list_person (person);

DROP TABLE IF EXISTS list_role;
CREATE TABLE list_role (
  list     integer NOT NULL REFERENCES list ON DELETE CASCADE,
  role     integer NOT NULL REFERENCES role ON DELETE CASCADE,
  sender   boolean NOT NULL DEFAULT false,
  submodel integer NOT NULL CHECK (submodel IN (0, 1, 2, 3)),
  PRIMARY KEY (list, role)
) WITHOUT ROWID;
CREATE INDEX list_role_role_idx ON list_role (role);

DROP TABLE IF EXISTS person;
CREATE TABLE person (
  id                 integer PRIMARY KEY,
  volgistics_id      integer,
  informal_name      text    NOT NULL,
  formal_name        text    NOT NULL,
  sort_name          text    NOT NULL,
  call_sign          text,
  pronouns           text,
  email              text    COLLATE NOCASE,
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
CREATE INDEX person_email_idx ON person (email);
CREATE INDEX person_sort_name_idx ON person (sort_name);

DROP TABLE IF EXISTS person_address;
CREATE TABLE person_address (
  person        integer NOT NULL REFERENCES person ON DELETE CASCADE,
  type          integer NOT NULL CHECK (type IN (0, 1, 2)),
  same_as_home  boolean NOT NULL,
  address       text,
  latitude      float,
  longitude     float,
  fire_district integer,
  PRIMARY KEY (person, type)
) WITHOUT ROWID;

DROP TABLE IF EXISTS person_bgcheck;
CREATE TABLE person_bgcheck (
  person  integer NOT NULL REFERENCES person ON DELETE CASCADE,
  type    integer NOT NULL CHECK (type IN (0, 1, 2)),
  cleared text, -- YYYY-MM-DD
  nli     text, -- YYYY-MM-DD
  assumed boolean NOT NULL DEFAULT false,
  PRIMARY KEY (person, type)
) WITHOUT ROWID;

DROP TABLE IF EXISTS person_dswreg;
CREATE TABLE person_dswreg (
  person     integer NOT NULL REFERENCES person ON DELETE CASCADE,
  class      integer NOT NULL CHECK (class BETWEEN 1 AND 14),
  registered text    NOT NULL, -- YYYY-MM-DD
  expiration text,             -- YYYY-MM-DD
  PRIMARY KEY (person, class)
) WITHOUT ROWID;

DROP TABLE IF EXISTS person_emcontact;
CREATE TABLE person_emcontact (
  person       integer NOT NULL REFERENCES person ON DELETE CASCADE,
  name         text    NOT NULL,
  home_phone   text,
  cell_phone   text,
  relationship text    NOT NULL
);
CREATE INDEX person_emcontact_person_idx ON person_emcontact (person);

DROP TABLE IF EXISTS person_note;
CREATE TABLE person_note (
  person     integer NOT NULL REFERENCES person ON DELETE CASCADE,
  note       text    NOT NULL,
  date       text    NOT NULL, -- YYYY-MM-DD
  visibility integer NOT NULL
);
CREATE INDEX person_note_person_idx ON person_note (person, date);

DROP TABLE IF EXISTS person_privlevel;
CREATE TABLE person_privlevel (
  person    integer NOT NULL REFERENCES person ON DELETE CASCADE,
  org       integer NOT NULL,
  privlevel integer NOT NULL,
  PRIMARY KEY (person, org)
) WITHOUT ROWID;

DROP TABLE IF EXISTS person_role;
CREATE TABLE person_role (
  person   integer NOT NULL REFERENCES person ON DELETE CASCADE,
  role     integer NOT NULL REFERENCES role ON DELETE CASCADE,
  explicit boolean NOT NULL,
  PRIMARY KEY (person, role)
) WITHOUT ROWID;
CREATE INDEX person_role_role_idx ON person_role (role);

DROP TABLE IF EXISTS role;
CREATE TABLE role (
  id        integer PRIMARY KEY,
  name      text    NOT NULL UNIQUE,
  title     text    UNIQUE,
  priority  integer NOT NULL, -- UNIQUE except while Reorder is running
  org       integer NOT NULL,
  privlevel integer NOT NULL,
  flags     integer NOT NULL DEFAULT 0
);

DROP TABLE IF EXISTS role_implies;
CREATE TABLE role_implies (
  implier integer NOT NULL REFERENCES role ON DELETE CASCADE,
  implied integer NOT NULL REFERENCES role ON DELETE CASCADE,
  PRIMARY KEY (implier, implied)
) WITHOUT ROWID;
CREATE INDEX role_implies_implied_idx ON role_implies (implied);

DROP TABLE IF EXISTS session;
CREATE TABLE session (
  token   text    PRIMARY KEY,
  person  integer NOT NULL REFERENCES person ON DELETE CASCADE,
  expires text    NOT NULL, -- YYYY-MM-DDTHH:MM:SS (local)
  csrf    text    NOT NULL
);
CREATE INDEX session_person_index ON session (person);

DROP TABLE IF EXISTS shift;
CREATE TABLE shift (
  id      integer PRIMARY KEY,
  task    integer NOT NULL REFERENCES task ON DELETE CASCADE,
  venue   integer REFERENCES venue,              -- *not* ON DELETE CASCADE
  start   text    NOT NULL,                      -- YYYY-MM-DDTHH:MM (local)
  end     text    NOT NULL CHECK (end >= start), -- YYYY-MM-DDTHH:MM (local)
  min     integer NOT NULL CHECK (min >= 0),
  max     integer CHECK (max IS NULL OR (max > 0 AND max >= min))
);
CREATE INDEX shift_task_idx ON shift (task);
CREATE INDEX shift_start_idx ON shift (start);
CREATE INDEX shift_venue_idx ON shift (venue);

DROP TABLE IF EXISTS shift_person;
CREATE TABLE shift_person (
  shift     integer NOT NULL REFERENCES shift ON DELETE CASCADE,
  person    integer NOT NULL REFERENCES person,
  signed_up integer NOT NULL,
  PRIMARY KEY (shift, person)
) WITHOUT ROWID;
CREATE UNIQUE INDEX shift_person_person_idx ON shift_person (shift, signed_up);

DROP TABLE IF EXISTS task;
CREATE TABLE task (
  id      integer PRIMARY KEY,
  event   integer NOT NULL REFERENCES event ON DELETE CASCADE,
  sort    integer NOT NULL CHECK (sort > 0),
  name    text,
  org     integer NOT NULL,
  flags   integer NOT NULL DEFAULT 0,
  details text
);
CREATE UNIQUE INDEX task_event_idx ON task (event, sort);
CREATE UNIQUE INDEX task_name_idx ON task (event, name);

DROP TABLE IF EXISTS task_person;
CREATE TABLE task_person (
  task      integer NOT NULL REFERENCES task,
  person    integer NOT NULL REFERENCES person,
  minutes   integer,
  flags     integer NOT NULL DEFAULT 0,
  PRIMARY KEY (task, person)
) WITHOUT ROWID;
CREATE INDEX task_person_person_idx ON task_person (person);

DROP TABLE IF EXISTS task_role;
CREATE TABLE task_role (
  task  integer NOT NULL REFERENCES task ON DELETE CASCADE,
  role  integer NOT NULL REFERENCES role,
  PRIMARY KEY (task, role)
) WITHOUT ROWID;
CREATE INDEX task_role_role_idx ON task_role (role);

DROP TABLE IF EXISTS textmsg;
CREATE TABLE textmsg (
  id          integer PRIMARY KEY,
  sender      integer NOT NULL REFERENCES person, -- *NOT* DELETE CASCADE
  timestamp   text    NOT NULL, -- YYYY-MM-DDTHH:MM:SS (local)
  message     text    NOT NULL
);
CREATE INDEX textmsg_sender_idx ON textmsg (sender);
CREATE INDEX textmsg_timestamp_idx ON textmsg (timestamp DESC);

DROP TABLE IF EXISTS textmsg_list;
CREATE TABLE textmsg_list (
  textmsg integer NOT NULL REFERENCES textmsg ON DELETE CASCADE,
  list    integer REFERENCES list ON DELETE SET NULL,
  name    text    NOT NULL,
  PRIMARY KEY (textmsg, list)
) WITHOUT ROWID;
CREATE INDEX textmsg_list_list_idx ON textmsg_list (list);

DROP TABLE IF EXISTS textmsg_number;
CREATE TABLE textmsg_number (
  number  text    PRIMARY KEY,
  textmsg integer NOT NULL REFERENCES textmsg ON DELETE CASCADE
);
CREATE INDEX textmsg_number_textmsg_idx ON textmsg_number (textmsg);

DROP TABLE IF EXISTS textmsg_recipient;
CREATE TABLE textmsg_recipient (
  textmsg   integer NOT NULL REFERENCES textmsg ON DELETE CASCADE,
  recipient integer NOT NULL REFERENCES person, -- *NOT* DELETE CASCADE,
  number    text,
  status    text,
  timestamp text    NOT NULL, -- YYYY-MM-DDTHH:MM:SS (local)
  PRIMARY KEY (textmsg, recipient)
) WITHOUT ROWID;
CREATE INDEX textmsg_recipient_recipient_idx ON textmsg_recipient (recipient);

DROP TABLE IF EXISTS textmsg_reply;
CREATE TABLE textmsg_reply (
  textmsg   integer NOT NULL REFERENCES textmsg ON DELETE CASCADE,
  recipient integer NOT NULL REFERENCES person, -- *NOT* DELETE CASCADE,
  reply     text    NOT NULL,
  timestamp text    NOT NULL -- YYYY-MM-DDTHH:MM:SS (local)
);
CREATE INDEX textmsg_reply_textmsg_idx ON textmsg_reply (textmsg, recipient, timestamp DESC);
CREATE INDEX textmsg_reply_recipient_idx ON textmsg_reply (recipient);

DROP TABLE IF EXISTS venue;
CREATE TABLE venue (
  id      integer PRIMARY KEY,
  name    text    NOT NULL UNIQUE,
  url     text,
  flags   integer NOT NULL
);
