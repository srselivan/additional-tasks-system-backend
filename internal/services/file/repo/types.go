package repo

type (
	FilesRepoGetListOpts struct {
		GroupId int64
		Limit   int64
		Offset  int64
	}
	FilesRepoCreateOpts struct {
		Name     string
		Filename string
		Filepath string
		TaskId   *int64
		AnswerId *int64
	}
)
