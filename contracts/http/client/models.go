package client

type GetArticleRequest struct {
	Id int64 `json:"id" form:"id" valid:"int,required"`
}

type ProxyResponse any

type EmptyRequest struct{}

type EmptyResponse struct{}
