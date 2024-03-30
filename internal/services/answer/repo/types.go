package repo

type (
	AnswersRepoGetListOpts struct {
		GroupId int64
		Limit   int64
		Offset  int64
	}
	AnswersRepoCreateOpts struct {
		GroupId int64
		Comment string
	}
	AnswersRepoUpdateOpts struct {
		Id      int64
		GroupId int64
		Comment string
	}
)
