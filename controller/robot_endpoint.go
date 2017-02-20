package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/reechou/holmes"
	"github.com/reechou/robot-account/models"
)

func (self *Logic) RobotSaveFriends(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}
	req := &RobotSaveFriendsReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("RobotSaveFriends json decode error: %v", err)
		return
	}

	rsp := &Response{Code: RESPONSE_OK}

	robot := &models.Robot{
		RobotWx: req.RobotWx,
	}
	_, err := models.GetRobot(robot)
	if err != nil {
		holmes.Error("get robot error: %v", err)
	}

	now := time.Now().Unix()
	var list []models.RobotFriend
	for _, v := range req.Friends {
		list = append(list, models.RobotFriend{
			RobotId:      robot.ID,
			RobotWx:      req.RobotWx,
			Name:         v.RemarkName,
			UserName:     v.UserName,
			Wechat:       v.Alias,
			City:         v.City,
			Sex:          v.Sex,
			Source:       models.FRIEND_SOURCE_OWNER,
			LastChatTime: now,
			CreatedAt:    now,
			UpdatedAt:    now,
		})
	}
	err = models.CreateRobotFriendList(list)
	if err != nil {
		holmes.Error("Error robot save friends error: %v", err)
		rsp.Code = RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("Error robot save friends error: %v", err)
	}

	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) RobotReceiveMsg(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}
	req := &ReceiveMsgInfo{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("RobotReceiveMsg json decode error: %v", err)
		return
	}

	self.HandleReceiveMsg(req)

	rsp := &Response{Code: RESPONSE_OK}

	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) RobotGroupList(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}

	rsp := &Response{Code: RESPONSE_OK}
	list, err := models.GetRobotGroupList()
	if err != nil {
		holmes.Error("get robot group list error: %v", err)
		rsp.Code = RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("get robot group list error: %v", err)
	} else {
		rsp.Data = list
	}

	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) RobotGroupInfoList(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}

	rsp := &Response{Code: RESPONSE_OK}

	req := &GetRobotGroupInfoListReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("RobotGroupInfoList json decode error: %v", err)
		return
	}
	list, err := models.GetRobotGroupInfoList(req.RobotGroupId)
	if err != nil {
		holmes.Error("get robot group info list error: %v", err)
		rsp.Code = RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("get robot group info list error: %v", err)
	} else {
		rsp.Data = list
	}

	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) CreateRobotGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}

	rsp := &Response{Code: RESPONSE_OK}

	req := &models.RobotGroup{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("CreateRobotGroup json decode error: %v", err)
		return
	}
	err := models.CreateRobotGroup(req)
	if err != nil {
		holmes.Error("create robot group error: %v", err)
		rsp.Code = RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("create robot group error: %v", err)
	} else {
		rsp.Data = req.ID
	}

	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) AddRobotGroupInfoList(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}

	rsp := &Response{Code: RESPONSE_OK}

	req := &AddRobotGroupInfoReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("AddRobotGroupInfo json decode error: %v", err)
		return
	}
	var list []models.RobotGroupInfo
	now := time.Now().Unix()
	for _, v := range req.RobotList {
		list = append(list, models.RobotGroupInfo{
			GroupId:   req.RobotGroupId,
			RobotId:   v.ID,
			RobotWx:   v.RobotWx,
			CreatedAt: now,
		})
	}
	err := models.CreateRobotGroupInfoList(list)
	if err != nil {
		holmes.Error("create robot group list error: %v", err)
		rsp.Code = RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("create robot group list error: %v", err)
	}

	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) DelRobotGroupInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}

	rsp := &Response{Code: RESPONSE_OK}

	req := &DelRobotGroupInfoReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("DelRobotGroupInfo json decode error: %v", err)
		return
	}
	err := models.DelRobotGroupInfo(req.ID)
	if err != nil {
		holmes.Error("delete robot group info error: %v", err)
		rsp.Code = RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("delete robot group info error: %v", err)
	}

	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) DelRobotGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}

	rsp := &Response{Code: RESPONSE_OK}

	req := &DelRobotGroupReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("DelRobotGroup json decode error: %v", err)
		return
	}
	err := models.DelRobotGroup(req.RobotGroupId)
	if err != nil {
		holmes.Error("delete robot group error: %v", err)
		rsp.Code = RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("delete robot group error: %v", err)
	}

	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) UpdateRobotGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}

	rsp := &Response{Code: RESPONSE_OK}

	req := &UpdateRobotGroupReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("UpdateRobotGroup json decode error: %v", err)
		return
	}
	err := models.UpdateRobotGroupName(&models.RobotGroup{ID: req.RobotGroupId, GroupName: req.GroupName})
	if err != nil {
		holmes.Error("update robot group error: %v", err)
		rsp.Code = RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("update robot group error: %v", err)
	}

	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) GetRobotList(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}

	rsp := &Response{Code: RESPONSE_OK}

	list, err := models.GetRobotList()
	if err != nil {
		holmes.Error("get robot list error: %v", err)
		rsp.Code = RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("get robot list error: %v", err)
	} else {
		rsp.Data = list
	}

	WriteJSON(w, http.StatusOK, rsp)
}
