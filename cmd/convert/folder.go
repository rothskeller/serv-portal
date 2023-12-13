package main

import (
	"path/filepath"

	"sunnyvaleserv.org/portal/model"
	ostore "sunnyvaleserv.org/portal/ostore"
	nstore "sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/store/document"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/folder"
)

func convertFolders(tx *ostore.Tx, st *nstore.Store) {
	for _, of := range tx.FetchFolders() {
		var nf folder.Updater

		if of.URL == "" {
			continue
		}
		parentPath, urlName := filepath.Dir(of.URL), filepath.Base(of.URL)
		npath, _ := folder.WithPath(st, filepath.Join("/files", parentPath), folder.FID|folder.FName)
		nf.Parent = npath[len(npath)-1]
		nf.Name = of.Name
		nf.URLName = urlName
		switch of.Visibility {
		case model.FolderVisibleToPublic:
			nf.ViewOrg, nf.ViewPriv = 0, 0
			nf.EditOrg, nf.EditPriv = enum.OrgAdmin, enum.PrivLeader
		case model.FolderVisibleToSERV:
			nf.ViewOrg, nf.ViewPriv = 0, enum.PrivMember
			nf.EditOrg, nf.EditPriv = enum.OrgAdmin, enum.PrivMember
		case model.FolderVisibleToOrg:
			nf.ViewOrg, nf.ViewPriv = enum.Org(of.Org), enum.PrivMember
			if nf.ViewOrg == enum.OrgAdmin {
				nf.EditOrg, nf.EditPriv = enum.OrgAdmin, enum.PrivMember
			} else {
				nf.EditOrg, nf.EditPriv = enum.Org(of.Org), enum.PrivLeader
			}
		default:
			panic("unknown folder visibility")
		}
		n := folder.Create(st, &nf)
		for _, od := range tx.FetchDocuments(of) {
			nd := document.Updater{
				Folder: n,
				Name:   od.Name,
				URL:    od.URL,
			}
			if nd.URL == "" {
				nd.LinkTo = filepath.Join("folders", of.URL, od.Name)
			}
			document.Create(st, &nd)
		}
	}
}
