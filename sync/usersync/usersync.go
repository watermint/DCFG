package usersync

import (
	"github.com/watermint/dcfg/integration/connector"
	"github.com/watermint/dcfg/integration/context"
	"github.com/watermint/dcfg/integration/directory"
)

type UserSync struct {
	DropboxConnector connector.DropboxConnector
	DropboxAccounts  directory.AccountDirectory
	GoogleAccounts   directory.AccountDirectory
	GoogleGroups     directory.GroupResolver
	GoogleEmail      directory.EmailResolver // Resolver for lookup
	GoogleConfirm    directory.EmailResolver // Resolver for confirmation
}

func NewUserSync(context context.ExecutionContext) UserSync {
	gd := directory.NewGoogleDirectory(context)
	dd := directory.NewDropboxDirectory(context)
	dp := connector.CreateConnector(context)
	gc := directory.NewGoogleEmailResolver(context)

	return UserSync{
		DropboxConnector: dp,
		DropboxAccounts:  dd,
		GoogleAccounts:   gd,
		GoogleGroups:     gd,
		GoogleEmail:      gd,
		GoogleConfirm:    gc,
	}
}

func (d *UserSync) membersNotInDirectory(members map[string]directory.Account, ad directory.AccountDirectory) (notInDir []directory.Account) {
	for _, x := range members {
		if _, exist := ad.Accounts()[x.Email]; !exist {
			notInDir = append(notInDir, x)
		}
	}
	return
}

func (d *UserSync) memberInGroups(needle directory.Account, haystack []directory.Group) bool {
	for _, x := range haystack {
		if needle.Email == x.GroupEmail {
			return true
		}
	}
	return false
}

func (d *UserSync) membersNotInGroup(member []directory.Account, gd directory.GroupResolver) (notInDir []directory.Account) {
	for _, x := range member {
		_, exist := gd.Group(x.Email)
		if !exist {
			notInDir = append(notInDir, x)
		}
	}
	return
}
