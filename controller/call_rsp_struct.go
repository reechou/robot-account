package controller

import (
	"github.com/reechou/robot-account/models"
)

type RobotFriendList struct {
	Count int64         `json:"count"`
	List  []RobotFriend `json:"list"`
}

type RobotFriend struct {
	Friend models.RobotFriend  `json:"friend"`
	Tags   []models.AccountTag `json:"tags"`
}

type TagList struct {
	Count int64        `json:"count"`
	List  []models.Tag `json:"list"`
}
