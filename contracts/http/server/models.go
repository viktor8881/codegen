package server

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

type GetRoleRequest struct {
	Id int64 `json:"id" form:"id" valid:"int,required"`
}

type GetRoleResponse struct {
	Id   int64  `json:"id" form:"id" valid:"int,required"`
	Name string `json:"name" form:"name" valid:"string,required"`
}

type CreateRoleRequest struct {
	Name string `json:"name" form:"name" valid:"string,required"`
}

type CreateRoleResponse struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
type UpdateRoleRequest struct {
	Id   int64  `json:"id" form:"id" valid:"int,required"`
	Name string `json:"name" form:"name" valid:"string,required"`
}

type UpdateRoleResponse struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type DeleteRoleRequest struct {
	Id int64 `json:"id" form:"id" valid:"int,required"`
}

type GetEmployeeRequest struct {
	Id int64 `json:"id" form:"id" valid:"int,required"`
}

type GetEmployeeResponse any

type GetPostRequest struct {
	Id int64 `json:"id" form:"id" valid:"int,required"`
}

type GetPostResponse any

type EmptyRequest struct{}

type EmptyResponse struct{}
