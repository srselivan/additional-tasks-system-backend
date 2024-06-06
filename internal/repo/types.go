package repo

import (
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("not found")
)

type (
	UsersRepoGetListByRoleIdOpts struct {
		RoleId int64
		Limit  int64
		Offset int64
	}
	UsersRepoGetByCredentialsOpts struct {
		Email    string
		Password string
	}
	UsersRepoCreateOpts struct {
		GroupId    *int64
		RoleId     int64
		Email      string
		Password   string
		FirstName  string
		LastName   string
		MiddleName *string
	}
)

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

type (
	AnswersRepoGetListOpts struct {
		TaskId int64
		Limit  int64
		Offset int64
	}
	AnswersRepoCreateOpts struct {
		TaskId  int64
		UserId  int64
		Comment *string
	}
	AnswersRepoUpdateOpts struct {
		Id      int64
		Comment string
	}
)

type (
	MarksRepoGetListByUserIdOpts struct {
		UserId int64
		Limit  int64
		Offset int64
	}
	MarksRepoCreateOpts struct {
		AnswerId int64
		Mark     int64
		Comment  *string
	}
	MarksRepoUpdateOpts struct {
		Id      int64
		Mark    int64
		Comment *string
	}
)

type (
	TasksRepoGetListForCreatorOpts struct {
		CreatedBy int64
		Limit     int64
		Offset    int64
	}
	TasksRepoGetListForUserOpts struct {
		UserId  int64
		GroupId int64
		Limit   int64
		Offset  int64
	}
	TasksRepoCreateOpts struct {
		CreatedBy     int64
		Title         string
		Text          *string
		EffectiveFrom time.Time
		EffectiveTill time.Time
	}
	TasksRepoUpdateOpts struct {
		Id            int64
		GroupId       int64
		Title         string
		Text          string
		EffectiveFrom time.Time
		EffectiveTill time.Time
		Cost          int64
	}
)

type (
	TaskLinksRepoCreateOpts struct {
		TaskId  int64
		UserId  *int64
		GroupId *int64
	}
)

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

type (
	GetStatisticsOpts struct {
		Limit  int64
		Offset int64
		From   *time.Time
		To     *time.Time
	}
)
