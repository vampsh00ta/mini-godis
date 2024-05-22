package request

type Login struct {
	Username string `json:"username" validate:"required"`
}
