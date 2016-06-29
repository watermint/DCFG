package dcfg

import (
	"github.com/watermint/dcfg/connector"
	"github.com/watermint/dcfg/directory"
	"github.com/watermint/dcfg/usersync"
	"testing"
)

func TestUserSyncRemoveUser(t *testing.T) {
	provision := connector.DropboxConnectorMock{}
	googleAccounts := directory.AccountDirectoryMock{
		MockData: []directory.Account{
			directory.Account{
				Email: "a@example.com",
			},
			directory.Account{
				Email: "b@example.com",
			},
		},
	}
	dropboxAccounts := directory.AccountDirectoryMock{
		MockData: []directory.Account{
			directory.Account{
				Email: "a@example.com",
			},
			directory.Account{
				Email: "b@example.com",
			},
			directory.Account{
				Email: "c@example.com",
			},
		},
	}
	userSync := usersync.UserSync{
		DropboxConnector: &provision,
		DropboxAccounts:  &dropboxAccounts,
		GoogleAccounts:   &googleAccounts,
	}
	userSync.SyncDeprovision()

	unexpected, missing, success := provision.AssertLogs([]string{
		provision.CreateOperationLog("MembersRemove", "c@example.com"),
	})
	if !success {
		t.Error("Sync failed", unexpected, missing, success)
	}
}

func TestUserSyncRemoveUser2(t *testing.T) {
	provision := connector.DropboxConnectorMock{}
	googleAccounts := directory.AccountDirectoryMock{
		MockData: []directory.Account{
			directory.Account{
				Email: "a@example.com",
			},
			directory.Account{
				Email: "b@example.com",
			},
		},
	}
	dropboxAccounts := directory.AccountDirectoryMock{
		MockData: []directory.Account{
			directory.Account{
				Email: "a@example.com",
			},
			directory.Account{
				Email: "c@example.com",
			},
			directory.Account{
				Email: "d@example.com",
			},
		},
	}
	userSync := usersync.UserSync{
		DropboxConnector: &provision,
		DropboxAccounts:  &dropboxAccounts,
		GoogleAccounts:   &googleAccounts,
	}
	userSync.SyncDeprovision()

	unexpected, missing, success := provision.AssertLogs([]string{
		provision.CreateOperationLog("MembersRemove", "c@example.com"),
		provision.CreateOperationLog("MembersRemove", "d@example.com"),
	})
	if !success {
		t.Error("Sync failed", unexpected, missing, success)
	}
}

func TestUserSyncEqual(t *testing.T) {
	provision := connector.DropboxConnectorMock{}
	googleAccounts := directory.AccountDirectoryMock{
		MockData: []directory.Account{
			directory.Account{
				Email: "a@example.com",
			},
			directory.Account{
				Email: "b@example.com",
			},
		},
	}
	dropboxAccounts := directory.AccountDirectoryMock{
		MockData: []directory.Account{
			directory.Account{
				Email: "a@example.com",
			},
			directory.Account{
				Email: "b@example.com",
			},
		},
	}
	userSync := usersync.UserSync{
		DropboxConnector: &provision,
		DropboxAccounts:  &dropboxAccounts,
		GoogleAccounts:   &googleAccounts,
	}
	userSync.SyncDeprovision()

	unexpected, missing, success := provision.AssertLogs([]string{})
	if !success {
		t.Error("Sync failed", unexpected, missing, success)
	}
}

func TestUserSyncGoogleHasMore(t *testing.T) {
	provision := connector.DropboxConnectorMock{}
	googleAccounts := directory.AccountDirectoryMock{
		MockData: []directory.Account{
			directory.Account{
				Email: "a@example.com",
			},
			directory.Account{
				Email: "b@example.com",
			},
			directory.Account{
				Email: "c@example.com",
			},
		},
	}
	dropboxAccounts := directory.AccountDirectoryMock{
		MockData: []directory.Account{
			directory.Account{
				Email: "a@example.com",
			},
			directory.Account{
				Email: "b@example.com",
			},
		},
	}
	userSync := usersync.UserSync{
		DropboxConnector: &provision,
		DropboxAccounts:  &dropboxAccounts,
		GoogleAccounts:   &googleAccounts,
	}
	userSync.SyncDeprovision()

	unexpected, missing, success := provision.AssertLogs([]string{})
	if !success {
		t.Error("Sync failed", unexpected, missing, success)
	}
}