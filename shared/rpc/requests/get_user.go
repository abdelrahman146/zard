package requests

type GetUserRequest struct {
	ID string `json:"id" validate:"required"`
}

func (r *GetUserRequest) Subject() string {
	return "identity.user.get"
}

func (r *GetUserRequest) Consumer(group string) string {
	return "identity_user_get_" + group
}

type GetUserResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
