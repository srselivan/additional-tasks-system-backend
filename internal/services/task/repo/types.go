package repo

import "time"

type (
	TasksRepoGetListOpts struct {
		GroupId int64
		Limit   int64
		Offset  int64
	}
	TasksRepoCreateOpts struct {
		GroupId       int64
		Title         string
		Text          string
		EffectiveFrom time.Time
		EffectiveTill time.Time
		Cost          int64
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
