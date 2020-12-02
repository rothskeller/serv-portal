package folder

import (
	"errors"

	"sunnyvaleserv.org/portal/model"
)

type vo struct {
	v model.FolderVisibility
	o model.Org
}

func (vo vo) String() string {
	if vo.v == model.FolderVisibleToOrg {
		return vo.o.String()
	}
	return vo.v.String()
}

func folderVOString(f *model.Folder) string {
	return vo{f.Visibility, f.Org}.String()
}

func (vo vo) Label() string {
	if vo.v == model.FolderVisibleToOrg {
		return vo.o.Label()
	}
	return vo.v.Label()
}

func folderVOLabel(f *model.Folder) string {
	return vo{f.Visibility, f.Org}.Label()
}

func parseVO(s string) (vo vo, err error) {
	if vo.v, err = model.ParseFolderVisibility(s); err == nil {
		return
	}
	if vo.o, err = model.ParseOrg(s); err == nil {
		vo.v = model.FolderVisibleToOrg
		return
	}
	return vo, errors.New("invalid visibility")
}

var allVOs []vo

func init() {
	allVOs = make([]vo, 0, len(model.AllFolderVisibilities)+len(model.AllOrgs)-1)
	for _, v := range model.AllFolderVisibilities {
		if v == model.FolderVisibleToOrg {
			for _, o := range model.AllOrgs {
				allVOs = append(allVOs, vo{v, o})
			}
		} else {
			allVOs = append(allVOs, vo{v, model.OrgNone})
		}
	}
}
