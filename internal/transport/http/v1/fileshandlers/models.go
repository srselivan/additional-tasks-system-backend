package answershandlers

type createRequest struct {
	TaskId   *int64 `json:"taskId"`
	AnswerId *int64 `json:"answerId"`
}
