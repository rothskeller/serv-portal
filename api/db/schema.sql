-- Database schema for serve.rothskeller.net/portal.

-- The team table tracks all SERV teams (except for the implicit "All People"
-- team).
CREATE TABLE team (
    id          integer PRIMARY KEY, -- autoincrement
    parent      integer NOT NULL REFERENCES team ON DELETE CASCADE,
    tag         text    UNIQUE,
    type        integer NOT NULL,
    name        text    NOT NULL,
    email       text    NOT NULL
);
CREATE INDEX team_parent_index ON team (parent);

-- The team_privilege table tracks which teams have privileges on which other
-- teams.  Note that Webmasters implicitly has all privileges on all teams.
CREATE TABLE team_privilege (
    actor      integer NOT NULL REFERENCES team ON DELETE CASCADE,
    target     integer NOT NULL REFERENCES team ON DELETE CASCADE,
    privileges integer NOT NULL,
    PRIMARY KEY (actor, target)
) WITHOUT ROWID;
CREATE INDEX team_privilege_target_index ON team_privilege (target);

-- The role table tracks all roles within the teams.
CREATE TABLE role (
    id    integer PRIMARY KEY, -- autoincrement
    team  integer NOT NULL REFERENCES team ON DELETE CASCADE,
    name  text    NOT NULL
);
CREATE INDEX role_team_index ON role (team);
CREATE INDEX role_name_index ON role (name); -- for sorting

-- The role_privilege table tracks which roles have which privileges on which
-- teams.  Note that Webmaster roles implicitly have all privileges on all
-- teams.
CREATE TABLE role_privilege (
    role       integer NOT NULL REFERENCES role ON DELETE CASCADE,
    team       integer NOT NULL REFERENCES team ON DELETE CASCADE,
    privileges integer NOT NULL,
    PRIMARY KEY (role, team)
) WITHOUT ROWID;
CREATE INDEX role_privilege_team_index ON role_privilege (team);

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

-- The person_role table records which people have which roles.
CREATE TABLE person_role (
    person integer NOT NULL REFERENCES person ON DELETE CASCADE,
    role   integer NOT NULL REFERENCES role ON DELETE CASCADE,
    PRIMARY KEY (person, role)
) WITHOUT ROWID;
CREATE INDEX person_role_role_index ON person_role (role);

-- The event table tracks all SERV events at which volunteer attendance is
-- tracked.
CREATE TABLE event (
    id    integer PRIMARY KEY, -- autoincrement
    date  text    NOT NULL,
    name  text    NOT NULL,
    hours float   NOT NULL,
    type  text    NOT NULL,
    UNIQUE (date, name)
);

-- The event_team table tracks which teams are invited to which events.
CREATE TABLE event_team (
    event integer NOT NULL REFERENCES event ON DELETE CASCADE,
    team  integer NOT NULL REFERENCES team ON DELETE CASCADE,
    PRIMARY KEY (event, team)
) WITHOUT ROWID;
CREATE INDEX event_team_team_index ON event_team (team);

-- The attendance table tracks which people attended which events.
CREATE TABLE attendance (
    event  integer NOT NULL REFERENCES event ON DELETE CASCADE,
    person integer NOT NULL REFERENCES person ON DELETE CASCADE,
    PRIMARY KEY (event, person)
) WITHOUT ROWID;
CREATE INDEX attendance_person_index ON attendance (person);
