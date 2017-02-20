package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/reechou/holmes"
	"github.com/reechou/robot-account/models"
)

const (
	URL_DELIMITER = "_$$_"
)

func (self *Logic) CreateWxCircle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}

	rsp := &Response{Code: RESPONSE_OK}

	req := &CreateWxCircleReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("CreateWxCircle json decode error: %v", err)
		return
	}
	
	wxc := &models.WxCircle{
		Type:         req.Type,
		Words:        req.Words,
		MaterialUrls: strings.Join(req.MaterialUrls, URL_DELIMITER),
	}
	err := models.CreateWxCircle(wxc)
	if err != nil {
		holmes.Error("create wx circle error: %v", err)
		rsp.Code = RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("create wx circle error: %v", err)
	} else {
		rsp.Data = wxc.ID
	}

	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) DeleteWxCircle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}

	rsp := &Response{Code: RESPONSE_OK}

	req := &DelWxCircleReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("DeleteWxCircle json decode error: %v", err)
		return
	}

	err := models.DelWxCircle(req.ID)
	if err != nil {
		holmes.Error("delete wx circle error: %v", err)
		rsp.Code = RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("delete wx circle error: %v", err)
	}

	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) UpdateWxCircle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}

	rsp := &Response{Code: RESPONSE_OK}

	req := &UpdateWxCircleReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("UpdateWxCircle json decode error: %v", err)
		return
	}

	wxc := &models.WxCircle{
		ID:           req.Id,
		Type:         req.Type,
		Words:        req.Words,
		MaterialUrls: strings.Join(req.MaterialUrls, URL_DELIMITER),
	}
	_, err := models.UpdateWxCircle(wxc)
	if err != nil {
		holmes.Error("update wx circle error: %v", err)
		rsp.Code = RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("update wx circle error: %v", err)
	}

	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) GetWxCircleList(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}

	rsp := &Response{Code: RESPONSE_OK}

	req := &GetWxCircleListReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("GetWxCircleList json decode error: %v", err)
		return
	}

	type WxCircleList struct {
		Count int64      `json:"count"`
		List  []WxCircle `json:"list"`
	}
	count, err := models.GetWxCircleListCount()
	if err != nil {
		holmes.Error("get wx circle count error: %v", err)
		rsp.Code = RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("get wx circle count error: %v", err)
	} else {
		list, err := models.GetWxCircleList(req.Offset, req.Num)
		if err != nil {
			holmes.Error("get wx circle list error: %v", err)
			rsp.Code = RESPONSE_ERR
			rsp.Msg = fmt.Sprintf("get wx circle list error: %v", err)
		}
		var wxList []WxCircle
		for _, v := range list {
			wxList = append(wxList, WxCircle{
				Id:           v.ID,
				Type:         v.Type,
				Words:        v.Words,
				MaterialUrls: strings.Split(v.MaterialUrls, URL_DELIMITER),
			})
		}
		result := &WxCircleList{
			Count: count,
			List:  wxList,
		}
		rsp.Data = result
	}

	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) CreateWxCircleSetting(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}

	rsp := &Response{Code: RESPONSE_OK}

	req := &CreateWxCircleSettingReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("CreateWxCircleSetting json decode error: %v", err)
		return
	}

	var settingList []models.WxCircleSetting
	for _, v := range req.RobotGroupIdList {
		settingList = append(settingList, models.WxCircleSetting{
			WxCircleId:   req.WxCircleId,
			RobotGroupId: v,
			SendTime:     req.SendTime,
		})
	}
	err := models.CreateWxCircleSettingList(settingList)
	if err != nil {
		holmes.Error("create wc circle setting list error: %v", err)
		rsp.Code = RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("create wc circle setting list error: %v", err)
	}

	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) DeleteWxCircleSetting(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}

	rsp := &Response{Code: RESPONSE_OK}

	req := &DelWxCircleSettingReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("DeleteWxCircleSetting json decode error: %v", err)
		return
	}

	err := models.DelWxCircleSetting(req.ID)
	if err != nil {
		holmes.Error("delete wx circle setting error: %v", err)
		rsp.Code = RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("delete wx circle setting error: %v", err)
	}

	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) GetWxCircleSettingList(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}

	rsp := &Response{Code: RESPONSE_OK}

	req := &GetWxCircleSettingListReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("GetWxCircleSettingList json decode error: %v", err)
		return
	}

	list, err := models.GetWxCircleSettingList(req.WxCircleId)
	if err != nil {
		holmes.Error("get wx circle setting list error: %v", err)
		rsp.Code = RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("get wx circle setting list error: %v", err)
	} else {
		rsp.Data = list
	}

	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) GetRobotNewWxCircle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" || r.Method != "GET" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}

	r.ParseForm()
	robot := ""
	if len(r.Form["robot"]) > 0 {
		robot = r.Form["robot"][0]
	}

	rsp := &Response{Code: RESPONSE_OK}

	rwc := &models.RobotWxCircle{
		Robot: robot,
	}
	has, err := models.GetRobotWxCircle(rwc)
	if err != nil {
		holmes.Error("get robot wx circle error: %v", err)
		rsp.Code = RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("get robot wx circle error: %v", err)
		WriteJSON(w, http.StatusOK, rsp)
		return
	}
	var lastWxCircleId int64
	if has {
		lastWxCircleId = rwc.LastWxCircleId
	}
	wxCircle, err := models.GetRobotNewWxCircle(robot, lastWxCircleId)
	if err != nil {
		holmes.Error("get robot new wx circle error: %v", err)
		rsp.Code = RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("get robot new wx circle error: %v", err)
	} else {
		rsp.Data = &WxCircle{
			Id:           wxCircle.ID,
			Type:         wxCircle.Type,
			Words:        wxCircle.Words,
			MaterialUrls: strings.Split(wxCircle.MaterialUrls, URL_DELIMITER),
		}
		
		rwc.LastWxCircleId = wxCircle.ID
		if has {
			_, err = models.UpdateRobotWxCircle(rwc)
			if err != nil {
				holmes.Error("update robot wx circle error: %v", err)
			}
		} else {
			err = models.CreateRobotWxCircle(rwc)
			if err != nil {
				holmes.Error("create robot wx circle error: %v", err)
			}
		}
	}

	WriteJSON(w, http.StatusOK, rsp)
}
