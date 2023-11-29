package types

type ThumbupMsg struct {
	BizId    string `json:"biz_id,omitempty"`    //业务 id
	ObjId    int64  `json:"obj_id,omitempty"`    //点赞对象id
	UserId   int64  `json:"user_id,omitempty"`   //用户id
	LikeType int32  `json:"like_type,omitempty"` //点赞类型
}
