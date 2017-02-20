package models

import (
	"fmt"
	"time"

	"github.com/reechou/holmes"
)

type Tag struct {
	ID        int64  `xorm:"pk autoincr" json:"id"`
	Tag       string `xorm:"not null default '' varchar(128) unique" json:"tag"`
	CreatedAt int64  `xorm:"not null default 0 int"  json:"createAt"`
}

func CreateTag(info *Tag) error {
	if info.Tag == "" {
		return fmt.Errorf("tag[%s] cannot be nil.", info.Tag)
	}

	now := time.Now().Unix()
	info.CreatedAt = now

	_, err := x.Insert(info)
	if err != nil {
		holmes.Error("create tag error: %v", err)
		return err
	}
	holmes.Info("create tag name[%s] success.", info.Tag)

	return nil
}

func GetTag(info *Tag) (bool, error) {
	has, err := x.Where("tag = ?", info.Tag).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		return false, nil
	}
	return true, nil
}

func GetTagCount() (int64, error) {
	count, err := x.Count(&Tag{})
	if err != nil {
		holmes.Error("get all tag count error: %v", err)
		return 0, err
	}
	return count, nil
}

func GetTagList(offset, num int64) ([]Tag, error) {
	var list []Tag
	err := x.Limit(int(num), int(offset)).Find(&list)
	if err != nil {
		holmes.Error("get tag list error: %v", err)
		return nil, err
	}
	return list, nil
}

type AccountTag struct {
	ID          int64  `xorm:"pk autoincr" json:"id"`
	RobotId     int64  `xorm:"not null default 0 int index"  json:"robotId"`
	AccountId   int64  `xorm:"not null default 0 int index"  json:"accountId"`
	AccountName string `xorm:"not null default '' varchar(128)" json:"accountName"`
	TagId       int64  `xorm:"not null default 0 int index"  json:"tagId"`
	TagName     string `xorm:"not null default '' varchar(128)" json:"tagName"`
	CreatedAt   int64  `xorm:"not null default 0 int"  json:"createAt"`
}

func CreateAccountTag(info *AccountTag) error {
	if info.TagName == "" {
		return fmt.Errorf("account tag[%s] cannot be nil.", info.TagName)
	}

	now := time.Now().Unix()
	info.CreatedAt = now

	_, err := x.Insert(info)
	if err != nil {
		holmes.Error("create account tag error: %v", err)
		return err
	}
	holmes.Info("create account[%s] tag name[%s] success.", info.AccountName, info.TagName)

	return nil
}

func DelAccountTag(id int64) error {
	rg := &AccountTag{ID: id}
	_, err := x.Where("id = ?", id).Delete(rg)
	if err != nil {
		holmes.Error("id[%d] account tag delete error: %v", id, err)
		return err
	}
	return nil
}

func GetAccountTagList(robotId int64, accountIdList []int64) ([]AccountTag, error) {
	var list []AccountTag
	err := x.Where("robot_id = ?", robotId).In("account_id", accountIdList).Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func GetAccountListFromTag(tagId int64) ([]AccountTag, error) {
	var list []AccountTag
	err := x.Where("tag_id = ?", tagId).Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}
