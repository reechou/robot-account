package models

import (
	"fmt"
	"strconv"
	"time"

	"github.com/reechou/holmes"
)

type RobotGroup struct {
	ID        int64  `xorm:"pk autoincr" json:"id"`
	GroupName string `xorm:"not null default '' varchar(128)" json:"groupName"`
	CreatedAt int64  `xorm:"not null default 0 int"  json:"createAt"`
}

func CreateRobotGroup(info *RobotGroup) error {
	if info.GroupName == "" {
		return fmt.Errorf("robot group[%s] cannot be nil.", info.GroupName)
	}

	now := time.Now().Unix()
	info.CreatedAt = now

	_, err := x.Insert(info)
	if err != nil {
		holmes.Error("create robot group error: %v", err)
		return err
	}
	holmes.Info("create robot group name[%s] success.", info.GroupName)

	return nil
}

func DelRobotGroup(id int64) error {
	rg := &RobotGroup{ID: id}
	_, err := x.Where("id = ?", id).Delete(rg)
	if err != nil {
		holmes.Error("id[%d] robot group delete error: %v", id, err)
		return err
	}
	return nil
}

func UpdateRobotGroupName(info *RobotGroup) error {
	_, err := x.ID(info.ID).Cols("group_name").Update(info)
	return err
}

func GetRobotGroupList() ([]RobotGroup, error) {
	var list []RobotGroup
	err := x.Find(&list)
	if err != nil {
		holmes.Error("get robot group list error: %v", err)
		return nil, err
	}
	return list, nil
}

type RobotGroupInfo struct {
	ID        int64  `xorm:"pk autoincr" json:"id"`
	GroupId   int64  `xorm:"not null default 0 int unique(group_robot_id)" json:"groupId"`
	RobotId   int64  `xorm:"not null default 0 int unique(group_robot_id)" json:"robotId"`
	RobotWx   string `xorm:"not null default '' varchar(128) index" json:"robotWx"`
	CreatedAt int64  `xorm:"not null default 0 int" json:"-"`
}

func CreateRobotGroupInfo(info *RobotGroupInfo) error {
	if info.GroupId == 0 {
		return fmt.Errorf("robot group[%d] cannot be nil.", info.GroupId)
	}

	now := time.Now().Unix()
	info.CreatedAt = now

	_, err := x.Insert(info)
	if err != nil {
		holmes.Error("create robot group info error: %v", err)
		return err
	}
	holmes.Info("create robot group info name[%v] success.", info)

	return nil
}

func CreateRobotGroupInfoList(list []RobotGroupInfo) error {
	if len(list) == 0 {
		return nil
	}
	_, err := x.Insert(&list)
	if err != nil {
		holmes.Error("create robot friend list error: %v", err)
		return err
	}
	return nil
}

func DelRobotGroupInfo(id int64) error {
	info := &RobotGroupInfo{ID: id}
	_, err := x.Where("id = ?", id).Delete(info)
	if err != nil {
		holmes.Error("id[%d] robot group info delete error: %v", id, err)
		return err
	}
	return nil
}

func GetRobotGroupInfoList(groupId int64) ([]RobotGroupInfo, error) {
	var list []RobotGroupInfo
	err := x.Where("group_id = ?", groupId).Find(&list)
	if err != nil {
		holmes.Error("get robot group info list error: %v", err)
		return nil, err
	}
	return list, nil
}

type Robot struct {
	ID            int64  `xorm:"pk autoincr" json:"id"`
	RobotWx       string `xorm:"not null default '' varchar(128)" json:"robotWx"`
	IfSaveFriend  int64  `xorm:"not null default 0 int" json:"-"`
	Tag           string `xorm:"not null default '' varchar(128)" json:"-"`
	Ip            string `xorm:"not null default '' varchar(64)"`
	OfPort        string `xorm:"not null default '' varchar(64)"`
	LastLoginTime int64  `xorm:"not null default 0 int" json:"lastLoginTime"`
	CreatedAt     int64  `xorm:"not null default 0 int" json:"createAt"`
	UpdatedAt     int64  `xorm:"not null default 0 int" json:"-"`
}

func GetRobot(info *Robot) (bool, error) {
	has, err := x.Where("robot_wx = ?", info.RobotWx).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		holmes.Debug("cannot find robot from robot_wx[%s]", info.RobotWx)
		return false, nil
	}
	return true, nil
}

func GetRobotList() ([]Robot, error) {
	var list []Robot
	err := x.Find(&list)
	if err != nil {
		holmes.Error("get robot list error: %v", err)
		return nil, err
	}
	return list, nil
}

const (
	FRIEND_SOURCE_OWNER     = iota
	FRIEND_SOURCE_USER_ADD  // 用户自己添加
	FRIEND_SOURCE_OWNER_ADD // 机器人添加好友:通讯录,群等
)

type RobotFriend struct {
	ID           int64  `xorm:"pk autoincr" json:"id"`
	RobotId      int64  `xorm:"not null default 0 int index" json:"robotId"`
	RobotWx      string `xorm:"not null default '' varchar(128) unique(robot_friend)" json:"robotWx"`
	Name         string `xorm:"not null default '' varchar(128) unique(robot_friend)" json:"name"`
	UserName     string `xorm:"not null default '' varchar(128)" json:"userName"`
	Wechat       string `xorm:"not null default '' varchar(128)" json:"wechat"`
	WxId         string `xorm:"not null default '' varchar(128)" json:"wxId"`
	City         string `xorm:"not null default '' varchar(64)" json:"city"`
	Sex          int    `xorm:"not null default 0 int" json:"sex"`
	SourceWechat string `xorm:"not null default '' varchar(128) index" json:"sourceWechat"`
	SourceNick   string `xorm:"not null default '' varchar(128)" json:"sourceNick"`
	Source       int64  `xorm:"not null default 0 int" json:"source"`
	Remark       string `xorm:"not null default '' varchar(768)" json:"remark"`
	LastChatTime int64  `xorm:"not null default 0 int index" json:"lastChatTime"`
	ChatNum      int64  `xorm:"not null default 0 int index" json:"chatNum"`
	CreatedAt    int64  `xorm:"not null default 0 int index" json:"-"`
	UpdatedAt    int64  `xorm:"not null default 0 int" json:"-"`
}

func (self *RobotFriend) TableName() string {
	return "robot_friend_" + strconv.Itoa(int(self.RobotId)%ROBOT_FRIEND_TABLE_NUM)
}

func CreateRobotFriend(info *RobotFriend) error {
	if info.RobotWx == "" {
		return fmt.Errorf("wx robot wx[%s] cannot be nil.", info.RobotWx)
	}

	now := time.Now().Unix()
	info.LastChatTime = now
	info.CreatedAt = now
	info.UpdatedAt = now

	_, err := x.Insert(info)
	if err != nil {
		holmes.Error("create robot friend error: %v", err)
		return err
	}
	holmes.Info("create robot friend robot[%s] name[%s] success.", info.RobotWx, info.Name)

	return nil
}

func CreateRobotFriendList(list []RobotFriend) error {
	if len(list) == 0 {
		return nil
	}
	_, err := x.Insert(&list)
	if err != nil {
		holmes.Error("create robot friend list error: %v", err)
		return err
	}
	return nil
}

func GetRobotFriendListCount(robotId int64) (int64, error) {
	count, err := x.Where("robot_id = ?", robotId).Count(&RobotFriend{RobotId: robotId})
	if err != nil {
		holmes.Error("robot_id[%d] get robot friend list count error: %v", robotId, err)
		return 0, err
	}
	return count, nil
}

func GetRobotFriendList(robotId, offset, num int64) ([]RobotFriend, error) {
	var list []RobotFriend
	err := x.Table(&RobotFriend{RobotId: robotId}).Where("robot_id = ?", robotId).Desc("created_at").Limit(int(num), int(offset)).Find(&list)
	if err != nil {
		holmes.Error("robot_id[%d] get robot friend list error: %v", robotId, err)
		return nil, err
	}
	return list, nil
}

func GetRobotFriendListOfActive(robotId, offset, num int64) ([]RobotFriend, error) {
	var list []RobotFriend
	err := x.Table(&RobotFriend{RobotId: robotId}).Where("robot_id = ?", robotId).Desc("chat_num").Limit(int(num), int(offset)).Find(&list)
	if err != nil {
		holmes.Error("robot_id[%d] get robot friend list error: %v", robotId, err)
		return nil, err
	}
	return list, nil
}

func GetRobotFriendListCountOf7DaysNoChat(robotId int64) (int64, error) {
	now := time.Now().Unix()
	checkTime := now - 7*86400
	count, err := x.Where("robot_id = ?", robotId).And("last_chat_time < ?", checkTime).Count(&RobotFriend{RobotId: robotId})
	if err != nil {
		holmes.Error("robot_id[%d] get robot friend list count error: %v", robotId, err)
		return 0, err
	}
	return count, nil
}

func GetRobotFriendListOf7DaysNoChat(robotId, offset, num int64) ([]RobotFriend, error) {
	now := time.Now().Unix()
	checkTime := now - 7*86400
	var list []RobotFriend
	err := x.Table(&RobotFriend{RobotId: robotId}).Where("robot_id = ?", robotId).And("last_chat_time < ?", checkTime).Desc("created_at").Limit(int(num), int(offset)).Find(&list)
	if err != nil {
		holmes.Error("robot_id[%d] get robot friend list error: %v", robotId, err)
		return nil, err
	}
	return list, nil
}

func GetRobotLowerFriendListCount(robotId int64, sourceWechat string) (int64, error) {
	count, err := x.Where("robot_id = ?", robotId).And("source_wechat = ?", sourceWechat).Count(&RobotFriend{RobotId: robotId})
	if err != nil {
		holmes.Error("robot_id[%d] source_wechat[%s] get robot lower friend list count error: %v", robotId, sourceWechat, err)
		return 0, err
	}
	return count, nil
}

func GetRobotLowerFriendList(robotId, offset, num int64, sourceWechat string) ([]RobotFriend, error) {
	var list []RobotFriend
	err := x.Table(&RobotFriend{RobotId: robotId}).Where("robot_id = ?", robotId).And("source_wechat = ?", sourceWechat).Desc("created_at").Limit(int(num), int(offset)).Find(&list, &RobotFriend{RobotId: robotId})
	if err != nil {
		holmes.Error("robot_id[%d] source_wechat[%s] get robot friend list error: %v", robotId, sourceWechat, err)
		return nil, err
	}
	return list, nil
}

func GetRobotFriend(info *RobotFriend) (bool, error) {
	has, err := x.Where("robot_wx = ?", info.RobotWx).And("name = ?", info.Name).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		holmes.Debug("cannot find robot frined from unionid[%s-%s]", info.RobotWx, info.Name)
		return false, nil
	}
	return true, nil
}

func GetRobotFriendFromUserName(info *RobotFriend) (bool, error) {
	has, err := x.Where("robot_wx = ?", info.RobotWx).And("user_name = ?", info.UserName).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		holmes.Debug("cannot find fx account from user_name[%s]", info.UserName)
		return false, nil
	}
	return true, nil
}

func UpdateRobotFriendName(info *RobotFriend) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.Cols("name", "updated_at").Update(info, &RobotFriend{RobotWx: info.RobotWx, UserName: info.UserName})
	return err
}

func UpdateRobotFriendUserName(info *RobotFriend) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.Cols("user_name", "updated_at").Update(info, &RobotFriend{RobotWx: info.RobotWx, Name: info.Name})
	return err
}

func UpdateRobotFriendRemark(info *RobotFriend) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.ID(info.ID).Cols("remark", "updated_at").Update(info)
	return err
}

func UpdateRobotFriendChatInfo(info *RobotFriend) error {
	info.UpdatedAt = time.Now().Unix()

	sql := fmt.Sprintf("update robot_friend_%d set last_chat_time=?, chat_num=chat_num+1, updated_at=? where id=?", info.RobotId%ROBOT_FRIEND_TABLE_NUM)
	var err error
	_, err = x.Exec(sql, info.UpdatedAt, info.UpdatedAt, info.ID)
	if err != nil {
		return err
	}
	return nil
}
