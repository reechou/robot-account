package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/reechou/holmes"
	"github.com/reechou/robot-account/models"
)

func (self *Logic) CreateAccountTag(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}

	rsp := &Response{Code: RESPONSE_OK}

	req := &CreateAccountTagReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("CreateAccountTag json decode error: %v", err)
		return
	}

	tag := &models.Tag{
		Tag: req.TagName,
	}
	has, err := models.GetTag(tag)
	if err != nil {
		holmes.Error("get tag error: %v", err)
		rsp.Code = RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("get tag error: %v", err)
		WriteJSON(w, http.StatusOK, rsp)
		return
	}

	if !has {
		err = models.CreateTag(tag)
		if err != nil {
			holmes.Error("create tag error: %v", err)
			rsp.Code = RESPONSE_ERR
			rsp.Msg = fmt.Sprintf("create tag error: %v", err)
			WriteJSON(w, http.StatusOK, rsp)
			return
		}
	}
	accountTag := &models.AccountTag{
		RobotId:     req.RobotId,
		AccountId:   req.AccountId,
		AccountName: req.AccountName,
		TagId:       tag.ID,
		TagName:     req.TagName,
	}
	err = models.CreateAccountTag(accountTag)
	if err != nil {
		holmes.Error("create account tag error: %v", err)
		rsp.Code = RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("create account tag error: %v", err)
		WriteJSON(w, http.StatusOK, rsp)
		return
	}
	rsp.Data = accountTag.ID

	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) GetAccountListFromTag(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}

	rsp := &Response{Code: RESPONSE_OK}

	req := &GetAccountListFromTagReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("GetAccountListFromTag json decode error: %v", err)
		return
	}

	tag := &models.Tag{
		Tag: req.TagName,
	}
	has, err := models.GetTag(tag)
	if err != nil {
		holmes.Error("get tag error: %v", err)
		rsp.Code = RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("get tag error: %v", err)
		WriteJSON(w, http.StatusOK, rsp)
		return
	}
	if has {
		list, err := models.GetAccountListFromTag(tag.ID)
		if err != nil {
			holmes.Error("get account list from tag error: %v", err)
			rsp.Code = RESPONSE_ERR
			rsp.Msg = fmt.Sprintf("get account list from tag error: %v", err)
			WriteJSON(w, http.StatusOK, rsp)
			return
		}
		rsp.Data = list
	}

	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) DeleteAccountTag(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}

	rsp := &Response{Code: RESPONSE_OK}

	req := &DelAccountTagReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("DeleteAccountTag json decode error: %v", err)
		return
	}

	err := models.DelAccountTag(req.ID)
	if err != nil {
		holmes.Error("delete tag error: %v", err)
		rsp.Code = RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("delete tag error: %v", err)
	}

	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) GetTagList(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}
	
	rsp := &Response{Code: RESPONSE_OK}
	
	req := &GetTagListReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("GetTagList json decode error: %v", err)
		return
	}
	count, err := models.GetTagCount()
	if err != nil {
		holmes.Error("get all tag count error: %v", err)
		rsp.Code = RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("get all tag count error: %v", err)
	} else {
		list, err := models.GetTagList(req.Offset, req.Num)
		if err != nil {
			holmes.Error("get robot lower friend list error: %v", err)
			rsp.Code = RESPONSE_ERR
			rsp.Msg = fmt.Sprintf("get robot lower friend list error: %v", err)
		}
		result := &TagList{
			Count: count,
			List:  list,
		}
		rsp.Data = result
	}
	
	WriteJSON(w, http.StatusOK, rsp)
}
