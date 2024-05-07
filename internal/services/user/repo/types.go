package repo

type (
	UsersRepoGetByCredentialsOpts struct {
		Email    string
		Password string
	}
	UsersRepoCreateOpts struct {
		GroupId    int64
		Email      string
		Password   string
		FirstName  string
		LastName   string
		MiddleName *string
	}
	UsersRepoUpdateOpts struct {
	}
)
