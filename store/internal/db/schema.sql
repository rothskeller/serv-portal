-- Database schema for serve.rothskeller.net/portal.

-- The approval table is a single-row, single-column table containing a BLOB.
-- The BLOB is the protocol buffer encoding of model.Approvals, which contains a
-- list of all outstanding requests needing approval.
CREATE TABLE approval (data BLOB);

-- The authorizer table is a single-row, single-column table containing a BLOB.
-- The BLOB is the protocol buffer encoding of authz.Authorizer, which contains
-- all of the groups, roles, and privileges for the SERV portal.
CREATE TABLE authorizer (data BLOB);

-- The folder table tracks all document folders.  The data column contains most
-- of the folder data, in protocol buffer encoding of model.Folder.
CREATE TABLE folder (
    id   integer PRIMARY KEY,
    data blob    NOT NULL
);

-- The lists table is a single-row, single-column table containing a BLOB.  The
-- BLOB is the protocol buffer encoding of model.Lists, which contains all of
-- the lists for the SERV portal.
CREATE TABLE lists (data blob NOT NULL);

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
    hours_token   text    UNIQUE,
    data          blob    NOT NULL
);

-- The roles table is a single-row, single-column table containing a BLOB.  The
-- BLOB is the protocol buffer encoding of model.Roles, which contains all of
-- the roles for the SERV portal.
CREATE TABLE roles (data blob NOT NULL);

-- The session table tracks all logged-in sessions.
CREATE TABLE session (
    token   text    PRIMARY KEY,
    person  integer NOT NULL REFERENCES person ON DELETE CASCADE,
    expires text    NOT NULL, -- RFC3339
    csrf    text    NOT NULL
);
CREATE INDEX session_person_index ON session (person);

-- The venue table is a single-row, single-column table containing a BLOB.  The
-- BLOB is the protocol buffer encoding of model.Venues, which contains all of
-- the venues for the SERV portal.
CREATE TABLE venue (data BLOB);

-- The event table tracks all SERV events at which volunteer attendance is
-- tracked.
CREATE TABLE event (
    id   integer PRIMARY KEY,
    date text    NOT NULL,
    data blob    NOT NULL
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

-- The search table contains full-text search information for all objects.
CREATE VIRTUAL TABLE search USING fts5 (
    type UNINDEXED,
    id   UNINDEXED,
    id2  UNINDEXED,
    documentName,
    documentContents,
    eventName,
    eventDetails,
    eventDate,
    folderName,
    personInformalName,
    personFormalName,
    personCallSign,
    personEmail,
    personEmail2,
    personHomeAddress,
    personWorkAddress,
    personMailAddress,
    roleName,
    roleTitle,
    textMessage,
    tokenize = 'unicode61'
);
