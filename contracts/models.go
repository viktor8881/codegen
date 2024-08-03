package contracts

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ListUserResponse struct {
	Users []UserResponse `json:"users"`
}

type ListUserByEmailRequest struct {
	Email string `json:"email" form:"email" valid:"required"`
}

type GetUserRequest struct {
	ID string `json:"id" form:"id" valid:"string,required"`
}

type GetUserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateUserRequest struct {
	Name  string `json:"name" form:"name" valid:"string,required"`
	Email string `json:"email" form:"email" valid:"email,required"`
}

type UpdateUserRequest struct {
	ID    string `json:"id" form:"id" valid:"string,required"`
	Name  string `json:"name" form:"name" valid:"string,required"`
	Email string `json:"email" form:"email" valid:"email,required"`
}

type DeleteUserRequest struct {
	ID string `json:"id" form:"id" valid:"int,required"`
}

type ProxyResponse any

type EmptyRequest struct{}

type EmptyResponse struct{}
