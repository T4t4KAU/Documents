package handlers

type UserParam struct {
	UserName string `json:"username" form:"username"`
	PassWord string `json:"password" form:"password"`
}

type RelationParam struct {
	ToUserId   int64
	ActionType int32
}
