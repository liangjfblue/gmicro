package models

type PostArticleRequest struct {
	Uid        string `json:"uid"`
	Title      string `json:"title"`
	Topic      string `json:"topic"`
	Author     string `json:"author"`
	IsOriginal int32  `json:"isOriginal"`
	Content    string `json:"content"`
}

type PostArticleRespond struct {
	ArticleId int32 `json:"articleId"`
}

type GetArticleRequest struct {
	Uid       string `json:"uid"`
	ArticleId int32  `json:"articleId"`
}

type GetArticleRespond struct {
	Title          string `json:"Title"`
	Topic          string `json:"Topic"`
	Author         string `json:"Author"`
	IsOriginal     int32  `json:"IsOriginal"`
	Content        string `json:"Content"`
	CreateTime     string `json:"CreateTime"`
	LastUpdateTime string `json:"LastUpdateTime"`
}

type DelArticleRequest struct {
	Uid       string `json:"uid"`
	ArticleId int32  `json:"articleId"`
}

type DelArticleRespond struct {
	Code int32 `json:"code"`
}
