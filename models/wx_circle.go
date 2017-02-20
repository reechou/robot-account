package models

import (
	"fmt"
	"strconv"
	"time"

	"github.com/reechou/holmes"
)

const (
	WX_CIRCLE_TYPE_WORDS = 0
)

// type: 0-仅文字 1-图片 2-视频 3-链接 (1,2,3都可带文字)
type WxCircle struct {
	ID           int64  `xorm:"pk autoincr" json:"id"`
	Type         int64  `xorm:"not null default 0 int" json:"type"`
	Words        string `xorm:"not null default '' varchar(768)" json:"words"`
	MaterialUrls string `xorm:"not null default '' varchar(768)" json:"materialUrls"`
	CreatedAt    int64  `xorm:"not null default 0 int" json:"createdAt"`
	UpdatedAt    int64  `xorm:"not null default 0 int" json:"updatedAt"`
}

func CreateWxCircle(info *WxCircle) error {
	if info.Type == WX_CIRCLE_TYPE_WORDS {
		if info.Words == "" {
			return fmt.Errorf("robot wx circle argv error.")
		}
	} else {
		if info.MaterialUrls == "" {
			return fmt.Errorf("robot wx circle argv error.")
		}
	}

	now := time.Now().Unix()
	info.CreatedAt = now

	_, err := x.Insert(info)
	if err != nil {
		holmes.Error("create wx circle error: %v", err)
		return err
	}
	holmes.Info("create robot wx circle[%s] success.", info.Words)

	return nil
}

func DelWxCircle(id int64) error {
	wxc := &WxCircle{ID: id}
	_, err := x.Where("id = ?", id).Delete(wxc)
	if err != nil {
		holmes.Error("id[%d] wx circle delete error: %v", id, err)
		return err
	}
	return nil
}

func UpdateWxCircle(info *WxCircle) (int64, error) {
	info.UpdatedAt = time.Now().Unix()
	affected, err := x.Cols("type", "words", "material_urls", "updated_at").Update(info, &WxCircle{ID: info.ID})
	return affected, err
}

func GetWxCircleListCount() (int64, error) {
	count, err := x.Count(&WxCircle{})
	if err != nil {
		holmes.Error("get wx circle list count error: %v", err)
		return 0, err
	}
	return count, nil
}

func GetWxCircleList(offset, num int64) ([]WxCircle, error) {
	var list []WxCircle
	err := x.Desc("created_at").Limit(int(num), int(offset)).Find(&list)
	if err != nil {
		holmes.Error("get wx circle list error: %v", err)
		return nil, err
	}
	return list, nil
}

type WxCircleSetting struct {
	ID           int64 `xorm:"pk autoincr" json:"id"`
	WxCircleId   int64 `xorm:"not null default 0 int index" json:"wxCircleId"`
	RobotGroupId int64 `xorm:"not null default 0 int index" json:"robotGroupId"`
	SendTime     int64 `xorm:"not null default 0 int index" json:"sendTime"`
	CreatedAt    int64 `xorm:"not null default 0 int" json:"createAt"`
}

func CreateWxCircleSettingList(list []WxCircleSetting) error {
	if len(list) == 0 {
		return nil
	}

	_, err := x.Insert(&list)
	if err != nil {
		holmes.Error("create wx circle setting list error: %v", err)
		return err
	}

	return nil
}

func DelWxCircleSetting(id int64) error {
	wxc := &WxCircleSetting{ID: id}
	_, err := x.Where("id = ?", id).Delete(wxc)
	if err != nil {
		holmes.Error("id[%d] wx circle setting delete error: %v", id, err)
		return err
	}
	return nil
}

type WxCircleSettingGroup struct {
	WxCircleSetting  `xorm:"extends"`
	GroupName string `json:"groupName"`
}
func (WxCircleSettingGroup) TableName() string {
	return "wx_circle_setting"
}

func GetWxCircleSettingList(wxCircleId int64) ([]WxCircleSettingGroup, error) {
	var list []WxCircleSettingGroup
	err := x.Join("INNER", "robot_group", "robot_group.id = wx_circle_setting.robot_group_id").Where("wx_circle_setting.wx_circle_id = ?", wxCircleId).Find(&list)
	if err != nil {
		holmes.Error("get wx circle setting list error: %v", err)
		return nil, err
	}
	return list, nil
}

type RobotWxCircle struct {
	ID             int64  `xorm:"pk autoincr" json:"id"`
	Robot          string `xorm:"not null default '' varchar(128)" json:"robot"`
	LastWxCircleId int64  `xorm:"not null default 0 int index" json:"lastWxCircleId"`
	CreatedAt      int64  `xorm:"not null default 0 int" json:"createAt"`
	UpdatedAt      int64  `xorm:"not null default 0 int" json:"-"`
}

func CreateRobotWxCircle(info *RobotWxCircle) error {
	if info.Robot == "" {
		return fmt.Errorf("robot[%s] wx circle cannot be nil.", info.Robot)
	}

	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now

	_, err := x.Insert(info)
	if err != nil {
		holmes.Error("create robot wx circle error: %v", err)
		return err
	}
	holmes.Info("create robot[%s] robot wx circle success.", info.Robot)

	return nil
}

func UpdateRobotWxCircle(info *RobotWxCircle) (int64, error) {
	info.UpdatedAt = time.Now().Unix()
	affected, err := x.Cols("last_wx_circle_id", "updated_at").Update(info, &RobotWxCircle{Robot: info.Robot})
	return affected, err
}

func GetRobotWxCircle(info *RobotWxCircle) (bool, error) {
	has, err := x.Where("robot = ?", info.Robot).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		holmes.Debug("cannot find robot wx circle from robot[%s]", info.Robot)
		return false, nil
	}
	return true, nil
}

func GetRobotNewWxCircle(robot string, lastWxCircleId int64) (*WxCircle, error) {
	results, err := x.Query("select wc.id, wc.type, wc.words, wc.material_urls, wc.created_at from wx_circle as wc left join wx_circle_setting as wcs on wc.id = wcs.wx_circle_id where (wcs.robot_group_id = 0"+
		" or wcs.robot_group_id in (select group_id from robot_group_info where robot_wx = ?)) and wc.id > ? limit 1", robot, lastWxCircleId)
	if err != nil {
		holmes.Error("get robot new wx circle error: %v", err)
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}
	
	
	id, _ := strconv.ParseInt(string(results[0]["wc.id"]), 10, 0)
	t, _ := strconv.ParseInt(string(results[0]["wc.type"]), 10, 0)
	createdAt, _ := strconv.ParseInt(string(results[0]["wc.created_at"]), 10, 0)
	wxCircle := &WxCircle{
		ID:           id,
		Type:         t,
		Words:        string(results[0]["wc.words"]),
		MaterialUrls: string(results[0]["wc.material_urls"]),
		CreatedAt:    createdAt,
	}

	return wxCircle, nil
}
