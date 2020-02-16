-- Database schema for serve.rothskeller.net/portal.

-- The approval table is a single-row, single-column table containing a BLOB.
-- The BLOB is the protocol buffer encoding of model.Approvals, which contains a
-- list of all outstanding requests needing approval.
CREATE TABLE approval (data BLOB);

-- The authorizer table is a single-row, single-column table containing a BLOB.
-- The BLOB is the protocol buffer encoding of authz.Authorizer, which contains
-- all of the groups, roles, and privileges for the SERV portal.
CREATE TABLE authorizer (data BLOB);

-- The email_list table tracks all email distribution lists.
CREATE TABLE email_list (
    id   text PRIMARY KEY,
    data blob NOT NULL
);

-- The email_message table tracks all email messages handled by the portal.
CREATE TABLE email_message (
    id         integer PRIMARY KEY,
    message_id text    NOT NULL UNIQUE,
    timestamp  text    NOT NULL,
    data       blob    NOT NULL
);
CREATE INDEX email_message_timestamp_index ON email_message (timestamp DESC);

-- The email_message_body table contains the actual body of each email message,
-- including headers, in transfer-encoded form exactly as received.
CREATE TABLE email_message_body (
    id   integer PRIMARY KEY REFERENCES email_message ON DELETE CASCADE,
    body blob    NOT NULL
);

-- The folder table tracks all document folders.  The data column contains most
-- of the folder data, in protocol buffer encoding of model.Folder.
CREATE TABLE folder (
    id   integer PRIMARY KEY,
    data blob    NOT NULL
);

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
    unsubscribe   text    NOT NULL UNIQUE,
    data          blob    NOT NULL
);

-- The person_email table maps email addresses to people; it is used when
-- receiving an email to determine which of our people (if any) it is from.
-- Note that people can have multiple email addresses, and that people can share
-- email addresses.
CREATE TABLE person_email (
    email  text    NOT NULL COLLATE NOCASE,
    person integer NOT NULL REFERENCES person ON DELETE CASCADE,
    UNIQUE (email, person)
);
CREATE INDEX person_email_person_index ON person_email (person);

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

-- The text_number table maps inbound phone numbers to the text_message ID of
-- the text message most recently sent to that number.
CREATE TABLE text_number (
    number text    PRIMARY KEY,
    mid    integer NOT NULL REFERENCES text_message ON DELETE CASCADE
);
CREATE INDEX text_number_mid_index ON text_number (mid);
