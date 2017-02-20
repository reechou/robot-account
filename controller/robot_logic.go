package controller

import (
	"github.com/reechou/holmes"
	"github.com/reechou/robot-account/models"
)

func (self *Logic) HandleReceiveMsg(msg *ReceiveMsgInfo) {
	switch msg.BaseInfo.ReceiveEvent {
	case RECEIVE_EVENT_ADD_FRIEND:
		self.AddRobotFriend(msg, models.FRIEND_SOURCE_USER_ADD)
	case RECEIVE_EVENT_ADD:
		self.AddRobotFriend(msg, models.FRIEND_SOURCE_OWNER_ADD)
	case RECEIVE_EVENT_MSG:
		self.AddRobotChat(msg)
	}
}

func (self *Logic) AddRobotFriend(msg *ReceiveMsgInfo, source int64) {
	robot := &models.Robot{
		RobotWx: msg.BaseInfo.WechatNick,
	}
	_, err := models.GetRobot(robot)
	if err != nil {
		holmes.Error("get robot error: %v", err)
	}

	rf := &models.RobotFriend{
		RobotId:      robot.ID,
		RobotWx:      msg.BaseInfo.WechatNick,
		Name:         msg.AddFriend.UserNick,
		UserName:     msg.BaseInfo.FromUserName,
		Wechat:       msg.AddFriend.UserWechat,
		WxId:         msg.AddFriend.UserWxid,
		City:         msg.AddFriend.UserCity,
		Sex:          msg.AddFriend.UserSex,
		SourceWechat: msg.AddFriend.SourceWechat,
		SourceNick:   msg.AddFriend.SourceNick,
		Source:       source,
	}
	err = models.CreateRobotFriend(rf)
	if err != nil {
		holmes.Error("create robot friend error: %v", err)
	}
}

func (self *Logic) AddRobotChat(msg *ReceiveMsgInfo) {
	if msg.BaseInfo.UserName == msg.BaseInfo.FromUserName {
		holmes.Debug("robot[%s] send msg[%s] to friend[%s]", msg.BaseInfo.FromNickName, msg.Msg, msg.BaseToUserInfo.ToNickName)
		self.addRobotChatFromRobot(msg)
		return
	}
	rf, err := self.getRobotFriend(msg.BaseInfo.WechatNick, msg.BaseInfo.FromUserName, msg.BaseInfo.FromNickName)
	if err != nil {
		holmes.Error("get msg robot friend error: %v", err)
		return
	}
	rc := &models.RobotChat{
		RobotId:      rf.RobotId,
		RobotWx:      msg.BaseInfo.WechatNick,
		AccountId:    rf.ID,
		FromName:     msg.BaseInfo.FromNickName,
		ToName:       msg.BaseInfo.WechatNick,
		MsgType:      msg.MsgType,
		Content:      msg.Msg,
		MediaTempUrl: msg.MediaTempUrl,
		Source:       models.ROBOT_CHAT_SOURCE_FROM_USER,
	}
	err = models.CreateRobotChat(rc)
	if err != nil {
		holmes.Error("create robot chat error: %v", err)
	}
	err = models.UpdateRobotFriendChatInfo(rf)
	if err != nil {
		holmes.Error("update robot friend chat info error: %v", err)
	}
}

func (self *Logic) addRobotChatFromRobot(msg *ReceiveMsgInfo) {
	rf, err := self.getRobotFriend(msg.BaseInfo.WechatNick, msg.BaseToUserInfo.ToUserName, msg.BaseToUserInfo.ToNickName)
	if err != nil {
		holmes.Error("get msg robot friend error: %v", err)
		return
	}
	rc := &models.RobotChat{
		RobotId:      rf.RobotId,
		RobotWx:      msg.BaseInfo.WechatNick,
		AccountId:    rf.ID,
		FromName:     msg.BaseInfo.WechatNick,
		ToName:       msg.BaseToUserInfo.ToNickName,
		MsgType:      msg.MsgType,
		Content:      msg.Msg,
		MediaTempUrl: msg.MediaTempUrl,
		Source:       models.ROBOT_CHAT_SOURCE_FROM_PHONE,
	}
	err = models.CreateRobotChat(rc)
	if err != nil {
		holmes.Error("create robot chat error: %v", err)
	}
	//err = models.UpdateRobotFriendChatInfo(rf)
	//if err != nil {
	//	holmes.Error("update robot friend chat info error: %v", err)
	//}
}

func (self *Logic) getRobotFriend(robotWx, fromUserName, fromNickName string) (*models.RobotFriend, error) {
	robot := &models.Robot{
		RobotWx: robotWx,
	}
	_, err := models.GetRobot(robot)
	if err != nil {
		holmes.Error("get robot error: %v", err)
		return nil, err
	}
	account := &models.RobotFriend{
		RobotId:  robot.ID,
		RobotWx:  robotWx,
		UserName: fromUserName,
	}
	has, err := models.GetRobotFriendFromUserName(account)
	if err != nil {
		holmes.Error("get robot friend from username[%v] error: %v", account, err)
		return nil, err
	}
	if has {
		if fromNickName != account.Name && fromNickName != "" {
			account.Name = fromNickName
			err = models.UpdateRobotFriendName(account)
			if err != nil {
				holmes.Error("update robot friend name[%v] error: %v", account, err)
				return nil, err
			}
		}
	}
	account = &models.RobotFriend{
		RobotId: robot.ID,
		RobotWx: robotWx,
		Name:    fromNickName,
	}
	has, err = models.GetRobotFriend(account)
	if err != nil {
		holmes.Error("get robot friend from [%v] error: %v", account, err)
		return nil, err
	}
	if !has {
		holmes.Error("cannot found this[%s %s %s] account", robotWx, fromUserName, fromNickName)
		return nil, MSG_ACCOUNT_GET_NONE
	}
	account.UserName = fromUserName
	err = models.UpdateRobotFriendUserName(account)
	if err != nil {
		holmes.Error("update account[%v] username error: %v", account, err)
	}
	return account, nil
}

func (self *Logic) SendRobotMsg(msg *SendBaseInfo) error {
	var msgs []SendBaseInfo
	msgs = append(msgs, *msg)
	err := self.re.SendMsgs(msg.WechatNick, &SendMsgInfo{SendMsgs: msgs})
	if err != nil {
		holmes.Error("send msgs of robot[%s] error: %v", msg.WechatNick, err)
		return err
	}
	holmes.Debug("send robot[%s] msg[%v] success.", msg.WechatNick, msg)
	self.addRobotChatFromWebCms(msg)
	return nil
}

func (self *Logic) addRobotChatFromWebCms(msg *SendBaseInfo) {
	rf, err := self.getRobotFriend(msg.WechatNick, msg.UserName, msg.NickName)
	if err != nil {
		holmes.Error("get msg robot friend error: %v", err)
		return
	}
	rc := &models.RobotChat{
		RobotId:      rf.RobotId,
		RobotWx:      msg.WechatNick,
		AccountId:    rf.ID,
		FromName:     msg.WechatNick,
		ToName:       msg.NickName,
		MsgType:      msg.MsgType,
		Content:      msg.Msg,
		Source:       models.ROBOT_CHAT_SOURCE_FROM_WEB,
	}
	err = models.CreateRobotChat(rc)
	if err != nil {
		holmes.Error("create robot chat error: %v", err)
	}
}
