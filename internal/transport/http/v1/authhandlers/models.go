package authhandlers

type signUpRequest struct {
	GroupId    int64   `json:"groupId"`
	Email      string  `json:"email"`
	Password   string  `json:"password"`
	FirstName  string  `json:"firstName"`
	LastName   string  `json:"lastName"`
	MiddleName *string `json:"middleName"`
}
