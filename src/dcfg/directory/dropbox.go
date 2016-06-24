package directory

import (
	"github.com/dropbox/dropbox-sdk-go-unofficial/team"
	"dcfg/auth"
	"dcfg/explorer"
	"github.com/cihub/seelog"
)

type DropboxDirectory struct {
	// API raw data structure
	rawMembers        []*team.TeamMemberInfo
	rawGroupSummaries []*team.GroupSummary
	rawGroupFullInfo  map[string]*team.GroupFullInfo

	// Abstract data structure
	groups            []Group
	accounts          []Account
}

const (
	dropboxLoadChunkSize = 100
)

func (d *DropboxDirectory) loadMembers() {
	d.rawMembers = []*team.TeamMemberInfo{}
	client := auth.DropboxClient()

	seelog.Trace("Loading Dropbox Team Member Info")

	sel := team.MembersListArg{}
	sel.Limit = dropboxLoadChunkSize
	ms, err := client.MembersList(&sel)
	if err != nil {
		explorer.Fatal("Unable to load Dropbox Team Member", err)
	}
	for _, m := range ms.Members {
		d.rawMembers = append(d.rawMembers, m)
	}
	seelog.Tracef("Dropbox Team Member Chunk loaded: %d member(s)", len(ms.Members))
	if !ms.HasMore {
		return
	}

	cursor := ms.Cursor
	for {
		sel := team.MembersListContinueArg{}
		sel.Cursor = cursor
		seelog.Trace("Loading Dropbox Team Member Info (Continue)")

		ms, err := client.MembersListContinue(&sel)
		if err != nil {
			explorer.Fatal("Unable to load Dropbox Team Member (while continue)", err)
		}
		for _, m := range ms.Members {
			d.rawMembers = append(d.rawMembers, m)
		}
		seelog.Tracef("Dropbox Team Member Chunk (Continue) loaded: %d member(s), has more: %t", len(ms.Members), ms.HasMore)
		if !ms.HasMore {
			break
		}
		cursor = ms.Cursor
	}
}

func (d *DropboxDirectory) loadGroupSummaries() {
	client := auth.DropboxClient()

	seelog.Trace("Loading Dropbox Group Summaries")

	sel := team.GroupsListArg{}
	sel.Limit = dropboxLoadChunkSize
	gs, err := client.GroupsList(&sel)
	if err != nil {
		explorer.Fatal("Unable to load Dropbox Group Summary", err)
	}
	for _, g := range gs.Groups {
		d.rawGroupSummaries = append(d.rawGroupSummaries, g)
	}
	seelog.Tracef("Dropbox Group Summary Chunk loaded: %d group(s)", len(gs.Groups))
	if !gs.HasMore {
		return
	}
	cursor := gs.Cursor
	for {
		sel := team.GroupsListContinueArg{}
		sel.Cursor = cursor
		gs, err := client.GroupsListContinue(&sel)
		if err != nil {
			explorer.Fatal("Unable to load Dropbox Group Summary (while continue)", err)
		}
		seelog.Tracef("Dropbox Group Summary (Continue) Chunk loaded: %d group(s)", len(gs.Groups))
		for _, g := range gs.Groups {
			d.rawGroupSummaries = append(d.rawGroupSummaries, g)
		}
		if !gs.HasMore {
			break
		}
		cursor = gs.Cursor
	}
}

func (d *DropboxDirectory) loadGroups() {
	groups := make(map[string]*team.GroupFullInfo)
	client := auth.DropboxClient()

	for i, gs := range d.rawGroupSummaries {
		sel := team.GroupsSelector{}
		sel.Tag = "group_ids"
		sel.GroupIds = []string{gs.GroupId}
		seelog.Tracef("Loading Dropbox Group Full Info [%d of %d]: Group ID[%s] Group Name[%s]", i, len(d.rawGroupSummaries), gs.GroupId, gs.GroupName)
		results, err := client.GroupsGetInfo(&sel)

		if err != nil {
			explorer.Fatal("Failed to load Dropbox Group", gs.GroupId, gs.GroupName, err)
		}

		for _, gr := range results {
			groups[gr.GroupInfo.GroupId] = gr.GroupInfo
		}
	}
	d.rawGroupFullInfo = groups
}

func (d *DropboxDirectory) Load() {
	d.loadMembers()
	d.loadGroupSummaries()
	d.loadGroups()

	d.groups = d.createGroups()
	d.accounts = d.createAccounts()
}

func (d *DropboxDirectory) createAccounts() (members []Account) {
	for _, m := range d.rawMembers {
		members = append(members, Account{
			Email: m.Profile.Email,
			GivenName: m.Profile.Name.GivenName,
			Surname: m.Profile.Name.Surname,
		})
	}
	return
}

func (d *DropboxDirectory) createGroups() (groups []Group) {
	for gid, g := range d.rawGroupFullInfo {
		members := []Account{}
		for _, m := range g.Members {
			members = append(members, Account{
				Email: m.Profile.Email,
				GivenName: m.Profile.Name.GivenName,
				Surname: m.Profile.Name.Surname,
			})
		}
		group := Group{
			GroupId: gid,
			GroupName: g.GroupName,
			CorrelationId: g.GroupExternalId,
			Members: members,
		}
		groups = append(groups, group)
	}
	return
}

func (d *DropboxDirectory) Groups() []Group {
	return d.groups
}

func (d *DropboxDirectory) Accounts() []Account {
	return d.accounts
}
