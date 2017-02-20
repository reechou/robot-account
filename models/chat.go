package models

import (
	"fmt"
	"strconv"
	"time"

	"github.com/reechou/holmes"
)

const (
	ROBOT_CHAT_SOURCE_FROM_USER  = "来自用户"
	ROBOT_CHAT_SOURCE_FROM_WEB   = "来自web cms"
	ROBOT_CHAT_SOURCE_FROM_PHONE = "来自手机"
)

type RobotChat struct {
	ID           int64  `xorm:"pk autoincr" json:"id"`
	RobotId      int64  `xorm:"not null default 0 int index" json:"robotId"`
	RobotWx      string `xorm:"not null default '' varchar(128) index" json:"robotWx"`
	AccountId    int64  `xorm:"not null default 0 int index" json:"accountId"`
	FromName     string `xorm:"not null default '' varchar(128)" json:"fromName"`
	ToName       string `xorm:"not null default '' varchar(128)" json:"toName"`
	MsgType      string `xorm:"not null default '' varchar(16)" json:"msgType"`
	Content      string `xorm:"not null default '' varchar(768)" json:"content"`
	MediaTempUrl string `xorm:"not null default '' varchar(256)" json:"mediaTempUrl,omitempty"`
	Source       string `xorm:"not null default '' varchar(64)" json:"source"`
	CreatedAt    int64  `xorm:"not null default 0 int index" json:"createdAt"`
}

func (self *RobotChat) TableName() string {
	return "robot_chat_" + strconv.Itoa(int(self.RobotId)%ROBOT_CHAT_TABLE_NUM)
}

func CreateRobotChat(info *RobotChat) error {
	if info.RobotWx == "" {
		return fmt.Errorf("wx robot wx[%s] cannot be nil.", info.RobotWx)
	}

	now := time.Now().Unix()
	info.CreatedAt = now

	_, err := x.Insert(info)
	if err != nil {
		holmes.Error("create robot chat error: %v", err)
		return err
	}
	holmes.Info("create robot chat robot[%s] fromUser[%s] toUser[%s] content[%s] success.", info.RobotWx, info.FromName, info.ToName, info.Content)

	return nil
}

func GetRobotChatListCount(robotId, accountId int64) (int64, error) {
	count, err := x.Where("account_id = ?", accountId).Count(&RobotChat{RobotId: robotId})
	if err != nil {
		holmes.Error("account_id[%d] get robot chat list count error: %v", accountId, err)
		return 0, err
	}
	return count, nil
}

func GetRobotChatList(robotId, accountId, offset, num int64) ([]RobotChat, error) {
	var list []RobotChat
	err := x.Table(&RobotChat{RobotId: robotId}).Where("account_id = ?", accountId).Desc("created_at").Limit(int(num), int(offset)).Find(&list)
	if err != nil {
		holmes.Error("account_id[%d] get robot chat list error: %v", accountId, err)
		return nil, err
	}
	return list, nil
}

func GetRobotNewChatList(robotId, timestamp int64) ([]RobotChat, error) {
	var list []RobotChat
	err := x.Table(&RobotChat{RobotId: robotId}).Where("created_at > ?", timestamp).Desc("created_at").Limit(20).Find(&list)
	if err != nil {
		holmes.Error("robot_id[%d] get robot new chat list error: %v", robotId, err)
		return nil, err
	}
	return list, nil
}
