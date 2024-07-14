package contracts

type UserResponse struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ListUserResponse struct {
	Users []UserResponse `json:"users"`
}

type GetUserRequest struct {
	Id int64 `json:"id" form:"id" valid:"int,required"`
}

type GetUserResponse struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateUserRequest struct {
	Name  string `json:"name" form:"name" valid:"string,required"`
	Email string `json:"email" form:"email" valid:"email,required"`
}

type CreateUserResponse struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserRequest struct {
	Id    int64  `json:"id"`
	Name  string `json:"name" form:"name" valid:"string,required"`
	Email string `json:"email" form:"email" valid:"email,required"`
}

type UpdateUserResponse struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type DeleteUserRequest struct {
	Id int64 `json:"id"`
}

type EmptyRequest struct{}

type EmptyResponse struct{}
