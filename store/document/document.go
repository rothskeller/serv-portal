// Package document defines the Document type, which describes a file or URL.
package document

import "sunnyvaleserv.org/portal/store/folder"

// ID uniquely identifies a document.
type ID int

// Document describes a document (a file or URL) in a folder.
type Document struct {
	// ID is the unique identifier of the document.
	ID ID
	// Folder is the identifier of the folder containing the document.
	Folder folder.ID
	// Name is the name of the document.  For file documents, it is the
	// filename; for URL documents, it is the display label of the link.
	Name string
	// URL is the URL for a URL document.  For file documents, it is empty.
	URL string
	// Archived is a flag indicating that the document is archived and no
	// longer in use.
	Archived bool
}
