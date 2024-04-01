package pg

import (
	"backend/internal/models"
	"time"
)

type file struct {
	Id        int64      `db:"id"`
	Name      string     `db:"name"`
	Filename  string     `db:"filename"`
	Filepath  string     `db:"filepath"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	TaskId    *int64     `db:"task_id"`
	AnswerId  *int64     `db:"answer_id"`
}

func (f file) toServiceModel() models.File {
	return models.File{
		Id:        f.Id,
		Name:      f.Name,
		Filename:  f.Filename,
		Filepath:  f.Filepath,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
		TaskId:    f.TaskId,
		AnswerId:  f.AnswerId,
	}
}
