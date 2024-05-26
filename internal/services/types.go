package services

import "time"

type (
	TaskServiceGetListForCreatorOpts struct {
		CreatedBy int64
		Limit     int64
		Offset    int64
	}
	TaskServiceGetListForUserOpts struct {
		UserId  int64
		GroupId int64
		Limit   int64
		Offset  int64
	}
	TaskServiceCreateOpts struct {
		GroupIds      []int64
		UserIds       []int64
		Title         string
		Text          *string
		CreatedBy     int64
		EffectiveFrom *time.Time
		EffectiveTill time.Time
		FileIds       []int64
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
		TaskId int64
		Limit  int64
		Offset int64
	}
	AnswerServiceCreateOpts struct {
		TaskId  int64
		UserId  int64
		Comment *string
		Files   []int64
	}
	AnswerServiceUpdateOpts struct {
		Id      int64
		Comment string
		Files   []int64
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
	MarkServiceGetListByUserIdOpts struct {
		UserId int64
		Limit  int64
		Offset int64
	}
	MarkServiceCreateOpts struct {
		AnswerId int64
		Mark     int64
		Comment  *string
	}
	MarkServiceUpdateOpts struct {
		Id      int64
		Mark    int64
		Comment *string
	}
)

type (
	UserServiceGetListByRoleIdOpts struct {
		RoleId int64
		Limit  int64
		Offset int64
	}
)

type (
	TaskLinksServiceCreateOpts struct {
		TaskId  int64
		UserId  *int64
		GroupId *int64
	}
)
