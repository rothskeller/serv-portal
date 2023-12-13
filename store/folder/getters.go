package folder

import "sunnyvaleserv.org/portal/store/enum"

// Fields returns the set of fields that have been retrieved for this folder.
func (f *Folder) Fields() Fields {
	return f.fields
}

// ID is the unique identifier of the Folder.
func (f *Folder) ID() ID {
	if f == nil {
		return 0
	}
	if f.fields&FID == 0 {
		panic("Folder.ID called without having fetched FID")
	}
	return f.id
}

// Parent is the unique identifier of the parent folder containing this folder.
// Note that the root folder (ID 1) is its own parent.
func (f *Folder) Parent() ID {
	if f.fields&FParent == 0 {
		panic("Folder.Parent called without having fetched FParent")
	}
	return f.parent
}

// Name is the name of the folder, as shown in the UI.
func (f *Folder) Name() string {
	if f.fields&FName == 0 {
		panic("Folder.Name called without having fetched FName")
	}
	return f.name
}

// URLName is the name of the folder, as it appears in the folder URL.  It is
// normally the kebab-case version of the folder name.
func (f *Folder) URLName() string {
	if f.fields&FURLName == 0 {
		panic("Folder.URLName called without having fetched FURLName")
	}
	return f.urlName
}

// Viewer is the organization and privilege level needed to view the folder.  If
// the organization is zero, having the requisite privilege level in any
// organization suffices.  If the privilege level is zero, the folder is visible
// to the general public without login.  Note that in the data model, the
// ability to view a folder is not dependent on being able to view its ancestor
// folders.  However, higher level code does enforce that restriction.
func (f *Folder) Viewer() (enum.Org, enum.PrivLevel) {
	if f.fields&FViewer == 0 {
		panic("Folder.Viewer called without having fetched FViewer")
	}
	return f.viewOrg, f.viewPriv
}

// Editor is the organization and privilege level needed to edit the folder
// contents, i.e., add, remove, and change files or folders within it.  If the
// organization is zero, having the requisite privilege level in any
// organization suffices.  Note that the ability to edit a folder is not
// dependent on being able to edit its ancestor folders.
func (f *Folder) Editor() (enum.Org, enum.PrivLevel) {
	if f.fields&FEditor == 0 {
		panic("Folder.Editor called without having fetched FEditor")
	}
	return f.editOrg, f.editPriv
}
