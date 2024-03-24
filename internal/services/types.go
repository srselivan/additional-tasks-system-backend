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
