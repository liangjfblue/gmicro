package models

type AddCommentRequest struct {
	Uid       string `json:"uid"`
	ArticleId int32  `json:"articleId"`
	Comment   string `json:"comment"`
	FromId    string `json:"fromId"`
	ToId      string `json:"toId"`
}

type AddCommentRespond struct {
	Code int32 `json:"code"`
}

type DelCommentRequest struct {
	Uid       string `json:"uid"`
	ArticleId int32  `json:"articleId"`
	CommentId int32  `json:"commentId"`
}

type DelCommentRespond struct {
	Code int32 `json:"code"`
}

type ListCommentRequest struct {
	Uid       string `json:"uid"`
	ArticleId int32  `json:"articleId"`
	Page      int32  `json:"page"`
	Size      int32  `json:"size"`
}

type Comment struct {
	Id      int32  `json:"id"`
	Comment string `json:"comment"`
	Time    string `json:"time"`
}

type ListCommentRespond struct {
	Count int32     `json:"count"`
	Lists []Comment `json:"lists"`
}
