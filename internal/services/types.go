package services

import "time"

type (
	TaskServiceGetListOpts struct {
		GroupId int64
		Limit   int64
		Offset  int64
	}
	TaskServiceCreateOpts struct {
		GroupId       int64
		Title         string
		Text          string
		EffectiveFrom time.Time
		EffectiveTill time.Time
		Cost          int64
	}
	TaskServiceUpdateOpts struct {
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
	AnswerServiceGetListOpts struct {
		GroupId int64
		Limit   int64
		Offset  int64
	}
	AnswerServiceCreateOpts struct {
		GroupId int64
		Comment string
	}
	AnswerServiceUpdateOpts struct {
		Id      int64
		GroupId int64
		Comment string
	}
)

type (
	GroupServiceGetListOpts struct {
		Limit  int64
		Offset int64
	}
	GroupServiceCreateOpts struct {
		Name string
	}
	GroupServiceUpdateOpts struct {
		Id   int64
		Name string
	}
)

type (
	FileServiceCreateOpts struct {
		Name     string
		Filename string
		Filepath string
		TaskId   *int64
		AnswerId *int64
	}
)

type (
	UserServiceCreateOpts struct {
		GroupId    int64
		Email      string
		Password   string
		FirstName  string
		LastName   string
		MiddleName *string
	}
	UserServiceUpdateOpts struct {
	}
	UserServiceGetByCredentialsOpts struct {
	}
)
