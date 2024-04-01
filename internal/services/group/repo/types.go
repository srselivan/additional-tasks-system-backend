package repo

type (
	GroupsRepoGetListOpts struct {
		Limit  int64
		Offset int64
	}
	GroupsRepoCreateOpts struct {
		Name string
	}
	GroupsRepoUpdateOpts struct {
		Id   int64
		Name string
	}
)
