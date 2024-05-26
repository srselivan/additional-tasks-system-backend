package authhandlers

type signUpRequest struct {
	GroupId    *int64  `json:"groupId"`
	RoleId     int64   `json:"roleId"`
	Email      string  `json:"email"`
	Password   string  `json:"password"`
	FirstName  string  `json:"firstName"`
	LastName   string  `json:"lastName"`
	MiddleName *string `json:"middleName"`
}

type signInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type jwtResponse struct {
	AccessToken string `json:"accessToken,omitempty"`
}
