package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/reechou/holmes"
	"github.com/reechou/robot-account/models"
)

func (self *Logic) GetRobotFriendList(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}

	rsp := &Response{Code: RESPONSE_OK}

	req := &GetRobotFriendReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("GetRobotFriendList json decode error: %v", err)
		return
	}
	//type RobotFriendList struct {
	//	Count int64                `json:"count"`
	//	List  []models.RobotFriend `json:"list"`
	//}
	count, err := models.GetRobotFriendListCount(req.RobotId)
	if err != nil {
		holmes.Error("get robot friend count error: %v", err)
		rsp.Code = RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("get robot friend count error: %v", err)
	} else {
		list, err := models.GetRobotFriendList(req.RobotId, req.Offset, req.Num)
		if err != nil {
			holmes.Error("get robot friend list error: %v", err)
			rsp.Code = RESPONSE_ERR
			rsp.Msg = fmt.Sprintf("get robot friend list error: %v", err)
			WriteJSON(w, http.StatusOK, rsp)
			return
		}
		var accountList []int64
		accountMap := make(map[int64]RobotFriend)
		for _, v := range list {
			accountList = append(accountList, v.ID)
			accountMap[v.ID] = RobotFriend{
				Friend: v,
			}
		}
		tList, err := models.GetAccountTagList(req.RobotId, accountList)
		if err != nil {
			holmes.Error("get robot friend tag list error: %v", err)
			rsp.Code = RESPONSE_ERR
			rsp.Msg = fmt.Sprintf("get robot friend tag list error: %v", err)
			WriteJSON(w, http.StatusOK, rsp)
			return
		}
		tagsMap := make(map[int64][]models.AccountTag)
		for _, v := range tList {
			tl := tagsMap[v.AccountId]
			tl = append(tl, v)
			tagsMap[v.AccountId] = tl
			//rf, ok := accountMap[v.AccountId]
			//if ok {
			//	rf.Tags = append(rf.Tags, v)
			//	rfList = append(rfList, rf)
			//	delete(accountMap, v.AccountId)
			//}
		}
		var rfList []RobotFriend
		for _, v := range accountMap {
			tl, ok := tagsMap[v.Friend.ID]
			if ok {
				v.Tags = tl
			}
			rfList = append(rfList, v)
		}
		result := &RobotFriendList{
			Count: count,
			List:  rfList,
		}
		rsp.Data = result
	}

	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) GetRobotLowerFriendList(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}

	rsp := &Response{Code: RESPONSE_OK}

	req := &GetRobotLowerFriendReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("GetRobotLowerFriendList json decode error: %v", err)
		return
	}
	//type RobotFriendList struct {
	//	Count int64                `json:"count"`
	//	List  []models.RobotFriend `json:"list"`
	//}
	count, err := models.GetRobotLowerFriendListCount(req.RobotId, req.SourceWechat)
	if err != nil {
		holmes.Error("get robot lower friend count error: %v", err)
		rsp.Code = RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("get robot lower friend count error: %v", err)
	} else {
		list, err := models.GetRobotLowerFriendList(req.RobotId, req.Offset, req.Num, req.SourceWechat)
		if err != nil {
			holmes.Error("get robot lower friend list error: %v", err)
			rsp.Code = RESPONSE_ERR
			rsp.Msg = fmt.Sprintf("get robot lower friend list error: %v", err)
		}
		var accountList []int64
		accountMap := make(map[int64]RobotFriend)
		for _, v := range list {
			accountList = append(accountList, v.ID)
			accountMap[v.ID] = RobotFriend{
				Friend: v,
			}
		}
		tList, err := models.GetAccountTagList(req.RobotId, accountList)
		if err != nil {
			holmes.Error("get robot friend tag list error: %v", err)
			rsp.Code = RESPONSE_ERR
			rsp.Msg = fmt.Sprintf("get robot friend tag list error: %v", err)
			WriteJSON(w, http.StatusOK, rsp)
			return
		}
		var rfList []RobotFriend
		for _, v := range tList {
			rf, ok := accountMap[v.AccountId]
			if ok {
				rf.Tags = append(rf.Tags, v)
				rfList = append(rfList, rf)
				delete(accountMap, v.AccountId)
			}
		}
		for _, v := range accountMap {
			rfList = append(rfList, v)
		}
		result := &RobotFriendList{
			Count: count,
			List:  rfList,
		}
		rsp.Data = result
	}

	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) UpdateRobotFriendRemark(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}

	rsp := &Response{Code: RESPONSE_OK}

	req := &UpdateRobotFriendRemarkReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("UpdateRobotFriendRemark json decode error: %v", err)
		return
	}
	rf := &models.RobotFriend{
		ID:      req.UserId,
		RobotId: req.RobotId,
		Remark:  req.Remark,
	}
	err := models.UpdateRobotFriendRemark(rf)
	if err != nil {
		holmes.Error("update robot friend remark error: %v", err)
		rsp.Code = RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("update robot friend remark error: %v", err)
	}

	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) GetRobotFriendList7DaysNoChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}

	rsp := &Response{Code: RESPONSE_OK}

	req := &GetRobotFriendReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("GetRobotFriendList7DaysNoChat json decode error: %v", err)
		return
	}
	//type RobotFriendList struct {
	//	Count int64                `json:"count"`
	//	List  []models.RobotFriend `json:"list"`
	//}
	count, err := models.GetRobotFriendListCountOf7DaysNoChat(req.RobotId)
	if err != nil {
		holmes.Error("get robot friend count of 7 days no chat error: %v", err)
		rsp.Code = RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("get robot friend count of 7 days no chat error: %v", err)
	} else {
		list, err := models.GetRobotFriendListOf7DaysNoChat(req.RobotId, req.Offset, req.Num)
		if err != nil {
			holmes.Error("get robot friend list of 7 days no chat error: %v", err)
			rsp.Code = RESPONSE_ERR
			rsp.Msg = fmt.Sprintf("get robot friend list of 7 days no chat error: %v", err)
		}
		var accountList []int64
		accountMap := make(map[int64]RobotFriend)
		for _, v := range list {
			accountList = append(accountList, v.ID)
			accountMap[v.ID] = RobotFriend{
				Friend: v,
			}
		}
		tList, err := models.GetAccountTagList(req.RobotId, accountList)
		if err != nil {
			holmes.Error("get robot friend tag list error: %v", err)
			rsp.Code = RESPONSE_ERR
			rsp.Msg = fmt.Sprintf("get robot friend tag list error: %v", err)
			WriteJSON(w, http.StatusOK, rsp)
			return
		}
		var rfList []RobotFriend
		for _, v := range tList {
			rf, ok := accountMap[v.AccountId]
			if ok {
				rf.Tags = append(rf.Tags, v)
				rfList = append(rfList, rf)
				delete(accountMap, v.AccountId)
			}
		}
		for _, v := range accountMap {
			rfList = append(rfList, v)
		}
		result := &RobotFriendList{
			Count: count,
			List:  rfList,
		}
		rsp.Data = result
	}

	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) GetRobotFriendListActive(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}

	rsp := &Response{Code: RESPONSE_OK}

	req := &GetRobotFriendReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("GetRobotFriendListActive json decode error: %v", err)
		return
	}
	//type RobotFriendList struct {
	//	Count int64                `json:"count"`
	//	List  []models.RobotFriend `json:"list"`
	//}
	count, err := models.GetRobotFriendListCount(req.RobotId)
	if err != nil {
		holmes.Error("get robot friend count of active error: %v", err)
		rsp.Code = RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("get robot friend count of active error: %v", err)
	} else {
		list, err := models.GetRobotFriendListOfActive(req.RobotId, req.Offset, req.Num)
		if err != nil {
			holmes.Error("get robot friend list of active error: %v", err)
			rsp.Code = RESPONSE_ERR
			rsp.Msg = fmt.Sprintf("get robot friend list of active error: %v", err)
		}
		var accountList []int64
		accountMap := make(map[int64]RobotFriend)
		for _, v := range list {
			accountList = append(accountList, v.ID)
			accountMap[v.ID] = RobotFriend{
				Friend: v,
			}
		}
		tList, err := models.GetAccountTagList(req.RobotId, accountList)
		if err != nil {
			holmes.Error("get robot friend tag list error: %v", err)
			rsp.Code = RESPONSE_ERR
			rsp.Msg = fmt.Sprintf("get robot friend tag list error: %v", err)
			WriteJSON(w, http.StatusOK, rsp)
			return
		}
		var rfList []RobotFriend
		for _, v := range tList {
			rf, ok := accountMap[v.AccountId]
			if ok {
				rf.Tags = append(rf.Tags, v)
				rfList = append(rfList, rf)
				delete(accountMap, v.AccountId)
			}
		}
		for _, v := range accountMap {
			rfList = append(rfList, v)
		}
		result := &RobotFriendList{
			Count: count,
			List:  rfList,
		}
		rsp.Data = result
	}

	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) GetRobotFriendChatList(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}

	rsp := &Response{Code: RESPONSE_OK}

	req := &GetRobotFriendChatReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("GetRobotFriendChatListActive json decode error: %v", err)
		return
	}
	type RobotFriendChatList struct {
		Count int64              `json:"count"`
		List  []models.RobotChat `json:"list"`
	}
	count, err := models.GetRobotChatListCount(req.RobotId, req.AccountId)
	if err != nil {
		holmes.Error("get robot friend chat count error: %v", err)
		rsp.Code = RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("get robot friend chat count error: %v", err)
	} else {
		list, err := models.GetRobotChatList(req.RobotId, req.AccountId, req.Offset, req.Num)
		if err != nil {
			holmes.Error("get robot friend chat list error: %v", err)
			rsp.Code = RESPONSE_ERR
			rsp.Msg = fmt.Sprintf("get robot friend chat list error: %v", err)
		}
		result := &RobotFriendChatList{
			Count: count,
			List:  list,
		}
		rsp.Data = result
	}

	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) GetRobotNewChatList(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}
	
	rsp := &Response{Code: RESPONSE_OK}
	
	req := &GetRobotNewChatReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("GetRobotNewChatList json decode error: %v", err)
		return
	}
	list, err := models.GetRobotNewChatList(req.RobotId, req.Timestamp)
	if err != nil {
		holmes.Error("get robot new chat list error: %v", err)
		rsp.Code = RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("get robot new chat list error: %v", err)
	} else {
		rsp.Data = list
	}
	
	WriteJSON(w, http.StatusOK, rsp)
}

func (self *Logic) SendTextMsg(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteJSON(w, http.StatusOK, nil)
		return
	}
	
	rsp := &Response{Code: RESPONSE_OK}
	
	req := &SendBaseInfo{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("SendTextMsg json decode error: %v", err)
		return
	}
	req.ChatType = CHAT_TYPE_PEOPLE
	req.MsgType = MSG_TYPE_TEXT
	
	err := self.SendRobotMsg(req)
	if err != nil {
		holmes.Error("send robot msg error: %v", err)
		rsp.Code = RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("send robot msg error: %v", err)
	}
	
	WriteJSON(w, http.StatusOK, rsp)
}
