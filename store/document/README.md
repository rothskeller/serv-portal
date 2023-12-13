# Document Storage

Documents come in two types:  files and URLs.  Both have a row in the 'document'
table, and are identified by the primary key ("id") of that row.

```sql
CREATE TABLE document (
  id       integer PRIMARY KEY,
  folder   integer NOT NULL REFERENCES folder,
  name     text    NOT NULL,
  url      text,
  archived boolean NOT NULL DEFAULT false
);
CREATE UNIQUE INDEX document_name_idx ON document (folder, name) WHERE NOT archived;
```

For URL documents, the `name` field is the display name of the link and the
`url` field is the URL.  There is no other associated storage for URL documents.

For file documents, the `name` field is the filename with extension, and the
`url` field is not used.  The file contents are stored in the file system, in a
file named documents/XX/YY, where XX is the `id` divided by 100 and YY is the
`id` modulo 100.

In order to be able to undo an accidental or malicious file change (or removal),
file documents are never changed removed through the UI.  Instead, their
`archived` flag is set, and a new file document with a new ID is created in its
place.  Periodically when needed, offline tools can prune the archived files.
The same is true when a file document is changed to a URL document:  the old
document is archived and a new document is created.

Changes and removals of URL documents are easily reversed without this
mechanism, using the log or database backups.  So the `archived` flag is never
set on a URL document.  URL documents are deleted when no longer needed, and are
changed in place.
