-- Database schema for serve.rothskeller.net/portal.

-- The audit table contains a record of every change to site data (except minor
-- trivia like session expiration).
CREATE TABLE audit (
    timestamp datetime NOT NULL,
    username  text     NOT NULL,
    request   text     NOT NULL,
    type      text     NOT NULL,
    id        any      NOT NULL,
    data      blob
);

-- The authz table is a single-row, single-column table containing a BLOB.
-- The BLOB is the protocol buffer encoding of model.AuthzData, which contains
-- all of the groups, roles, and privileges for the SERV portal.
CREATE TABLE authz (data BLOB);

-- The person table tracks all people associated (or formerly associated) with
-- SERV.  There is one row in this table for each such person.  Since each such
-- person has a (potentially disabled) login to the SERV portal, this table also
-- tracks users.  The data column contains most of the person data, in protocol
-- buffer encoding of model.Person.  The other columns are those needed for
-- lookups or sorting.
CREATE TABLE person (
    id            integer PRIMARY KEY,
    username      text    UNIQUE COLLATE NOCASE,
    pwreset_token text    UNIQUE,
    cell_phone    text    UNIQUE,
    data          blob    NOT NULL
);

-- The session table tracks all logged-in sessions.
CREATE TABLE session (
    token   text    PRIMARY KEY,
    person  integer NOT NULL REFERENCES person ON DELETE CASCADE,
    expires text    NOT NULL -- RFC3339
);
CREATE INDEX session_person_index ON session (person);

-- The venue table is a single-row, single-column table containing a BLOB.  The
-- BLOB is the protocol buffer encoding of model.Venues, which contains all of
-- the venues for the SERV portal.
CREATE TABLE venue (data BLOB);

-- The event table tracks all SERV events at which volunteer attendance is
-- tracked.
CREATE TABLE event (
    id          integer PRIMARY KEY,
    date        text    NOT NULL,
    scc_ares_id text    UNIQUE,
    data        blob    NOT NULL
);
CREATE INDEX event_date_index ON event (date);

-- The attendance table tracks which people attended which events.
CREATE TABLE attendance (
    event   integer NOT NULL REFERENCES event ON DELETE CASCADE,
    person  integer NOT NULL REFERENCES person ON DELETE CASCADE,
    type    integer NOT NULL DEFAULT 0,
    minutes integer NOT NULL DEFAULT 0,
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
    serv integer NOT NULL
);

-- The scc_ares_event_type table maps types of events imported from the
-- scc-ares-races.org site into our event types.  Any event type in their
-- database which doesn't have an entry here is mapped according to the entry
-- here for scc='' (which must exist).
CREATE TABLE scc_ares_event_type (
    scc  text PRIMARY KEY,
    serv text NOT NULL
);

-- The text_message table tracks all outgoing text messages.
CREATE TABLE text_message (
    id   integer PRIMARY KEY,
    data blob    NOT NULL
);

-- The text_delivery table tracks the delivery of each outgoing text message to
-- each recipient.  This includes tracking their responses if any.
CREATE TABLE text_delivery (
    message   integer NOT NULL REFERENCES text_message ON DELETE CASCADE,
    recipient integer NOT NULL REFERENCES person ON DELETE CASCADE,
    data      blob    NOT NULL,
    PRIMARY KEY (message, recipient)
);
CREATE INDEX text_delivery_recipient_index ON text_delivery (recipient);
