-- Database schema for serve.rothskeller.net/portal.

-- The role table tracks all SERV roles.
CREATE TABLE role (
    id           integer PRIMARY KEY,
    tag          text    UNIQUE,
    name         text    NOT NULL,
    member_label text    NOT NULL,
    serv_group   text    NOT NULL,
    imply_only   boolean NOT NULL,
    individual   boolean NOT NULL,
    privileges   blob    NOT NULL
);
CREATE INDEX role_name_index ON role (name);

-- The role_privilege table tracks which actor roles have which privileges on
-- which target roles.  It is redundant with the privileges column of the role
-- table, but it is maintained so that the data are in a form usable in offline
-- SQL queries, and so that the privileges blobs can be recalculated if
-- necessary.  This table is written but never read by the portal server.  Note
-- that role implications are also encoded here, as the LSB of each privileges
-- bitmask.
CREATE TABLE role_privilege (
    actor      integer NOT NULL REFERENCES role ON DELETE CASCADE,
    target     integer NOT NULL REFERENCES role ON DELETE CASCADE,
    privileges integer NOT NULL,
    PRIMARY KEY (actor, target)
) WITHOUT ROWID;
CREATE INDEX role_privilege_target_index ON role_privilege (target);

-- The person table tracks all people associated (or formerly associated) with
-- SERV.  There is one row in this table for each such person.  Since each such
-- person has a (potentially disabled) login to the SERV portal, this table also
-- tracks users.
CREATE TABLE person (
    id              integer PRIMARY KEY, -- autoincrement
    first_name      text    NOT NULL,
    last_name       text    NOT NULL,
    email           text    NOT NULL UNIQUE,
    phone           text    NOT NULL,
    password        text    NOT NULL,
    bad_login_count integer NOT NULL DEFAULT 0,
    bad_login_time  text    NOT NULL DEFAULT '', -- RFC3339
    pwreset_token   text    UNIQUE,
    pwreset_time    text    NOT NULL DEFAULT '' -- RFC3339
);

-- The session table tracks all logged-in sessions.
CREATE TABLE session (
    token   text    PRIMARY KEY,
    person  integer NOT NULL REFERENCES person ON DELETE CASCADE,
    expires text    NOT NULL -- RFC3339
);
CREATE INDEX session_person_index ON session (person);

-- The person_role table records which people have which roles.  It includes
-- only direct role membership, not transitive memberships.
CREATE TABLE person_role (
    person integer NOT NULL REFERENCES person ON DELETE CASCADE,
    role   integer NOT NULL REFERENCES role ON DELETE CASCADE,
    PRIMARY KEY (person, role)
) WITHOUT ROWID;
CREATE INDEX person_role_role_index ON person_role (role);

-- The venue table tracks all venues at which SERV events are held.
CREATE TABLE venue (
    id      integer PRIMARY KEY, -- autoincrement
    name    text    NOT NULL,
    address text    NOT NULL,
    city    text    NOT NULL,
    url     text    NOT NULL
);

-- The event table tracks all SERV events at which volunteer attendance is
-- tracked.
CREATE TABLE event (
    id          integer PRIMARY KEY, -- autoincrement
    name        text    NOT NULL,
    date        text    NOT NULL,
    start       text    NOT NULL,
    end         text    NOT NULL,
    venue       integer REFERENCES venue ON DELETE SET NULL,
    details     text    NOT NULL,
    type        text    NOT NULL,
    scc_ares_id text    UNIQUE,
    UNIQUE (date, name)
);
CREATE INDEX event_venue_index ON event (venue);

-- The event_role table tracks which roles are invited to which events.
CREATE TABLE event_role (
    event integer NOT NULL REFERENCES event ON DELETE CASCADE,
    role  integer NOT NULL REFERENCES role ON DELETE CASCADE,
    PRIMARY KEY (event, role)
) WITHOUT ROWID;
CREATE INDEX event_role_role_index ON event_role (role);

-- The attendance table tracks which people attended which events.
CREATE TABLE attendance (
    event  integer NOT NULL REFERENCES event ON DELETE CASCADE,
    person integer NOT NULL REFERENCES person ON DELETE CASCADE,
    PRIMARY KEY (event, person)
) WITHOUT ROWID;
CREATE INDEX attendance_person_index ON attendance (person);

-- The scc_ares_event_name table allows renaming and deleting of events imported
-- from the scc-ares-races.org site.  The event name from that site is looked up
-- in the 'scc' column.  Those whose 'serv' column is empty are not imported.
-- Those whose 'serv' column is non-empty are renamed accordingly.  Those whose
-- name is not found in this table are imported unchanged.
CREATE TABLE scc_ares_event_name (
    scc  text PRIMARY KEY,
    serv text NOT NULL
);

-- The scc_ares_event_location table maps locations of events imported from the
-- scc-ares-races.org site into venues in our database.  Any location in their
-- database which doesn't have an entry here is mapped according to the entry
-- here for scc='' (which must exist).
CREATE TABLE scc_ares_event_location (
    scc  text    PRIMARY KEY,
    serv integer REFERENCES venue ON DELETE CASCADE
);
CREATE INDEX scc_ares_event_location_serv_index ON scc_ares_event_location (serv);

-- The scc_ares_event_type table maps types of events imported from the
-- scc-ares-races.org site into our event types.  Any event type in their
-- database which doesn't have an entry here is mapped according to the entry
-- here for scc='' (which must exist).
CREATE TABLE scc_ares_event_type (
    scc  text PRIMARY KEY,
    serv text NOT NULL
);
