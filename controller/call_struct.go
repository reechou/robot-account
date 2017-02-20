package controller

import (
	"github.com/reechou/robot-account/models"
)

const (
	RESPONSE_OK = iota
	RESPONSE_ERR
)

type Response struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

type RobotSaveFriendsReq struct {
	RobotWx string       `json:"robotWx"`
	Friends []UserFriend `json:"friends"`
}

type GetRobotGroupInfoListReq struct {
	RobotGroupId int64 `json:"robotGroupId"`
}

type AddRobotGroupInfoReq struct {
	RobotGroupId int64          `json:"robotGroupId"`
	RobotList    []models.Robot `json:"robotList"`
}

type DelRobotGroupReq struct {
	RobotGroupId int64 `json:"robotGroupId"`
}

type DelRobotGroupInfoReq struct {
	ID int64 `json:"id"`
}

type UpdateRobotGroupReq struct {
	RobotGroupId int64  `json:"robotGroupId"`
	GroupName    string `json:"groupName"`
}

type GetRobotFriendReq struct {
	RobotId int64 `json:"robotId"`
	Offset  int64 `json:"offset"`
	Num     int64 `json:"num"`
}

type GetRobotLowerFriendReq struct {
	RobotId      int64  `json:"robotId"`
	SourceWechat string `json:"sourceWechat"`
	Offset       int64  `json:"offset"`
	Num          int64  `json:"num"`
}

type UpdateRobotFriendRemarkReq struct {
	RobotId int64  `json:"robotId"`
	UserId  int64  `json:"userId"`
	Remark  string `json:"remark"`
}

type GetRobotFriendChatReq struct {
	RobotId   int64 `json:"robotId"`
	AccountId int64 `json:"accountId"`
	Offset    int64 `json:"offset"`
	Num       int64 `json:"num"`
}

type GetRobotNewChatReq struct {
	RobotId   int64 `json:"robotId"`
	Timestamp int64 `json:"lastTime"`
}

type CreateAccountTagReq struct {
	RobotId     int64  `json:"robotId"`
	AccountId   int64  `json:"accountId"`
	AccountName string `json:"accountName"`
	TagName     string `json:"tagName"`
}

type GetAccountListFromTagReq struct {
	TagName string `json:"tagName"`
}

type DelAccountTagReq struct {
	ID int64 `json:"id"`
}

type GetTagListReq struct {
	Offset int64 `json:"offset"`
	Num    int64 `json:"num"`
}

type WxCircle struct {
	Id           int64    `json:"id"`
	Type         int64    `json:"type"`
	Words        string   `json:"words"`
	MaterialUrls []string `json:"materialUrls"`
}

type CreateWxCircleReq struct {
	Type         int64    `json:"type"`
	Words        string   `json:"words"`
	MaterialUrls []string `json:"materialUrls"`
}

type DelWxCircleReq struct {
	ID int64 `json:"id"`
}

type UpdateWxCircleReq struct {
	Id           int64    `json:"id"`
	Type         int64    `json:"type"`
	Words        string   `json:"words"`
	MaterialUrls []string `json:"materialUrls"`
}

type GetWxCircleListReq struct {
	Offset int64 `json:"offset"`
	Num    int64 `json:"num"`
}

type CreateWxCircleSettingReq struct {
	WxCircleId       int64   `json:"wxCircleId"`
	RobotGroupIdList []int64 `json:"robotGroupIdList"`
	SendTime         int64   `json:"sendTime"`
}

type DelWxCircleSettingReq struct {
	ID int64 `json:"id"`
}

type GetWxCircleSettingListReq struct {
	WxCircleId int64 `json:"wxCircleId"`
}
