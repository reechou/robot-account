package controller

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/reechou/holmes"
	"github.com/reechou/robot-account/config"
	"github.com/reechou/robot-account/models"
)

type Logic struct {
	sync.Mutex

	cfg *config.Config
	re  *RobotExt
}

func NewLogic(cfg *config.Config) *Logic {
	l := &Logic{
		cfg: cfg,
	}
	l.re = NewRobotExt()

	models.InitDB(cfg)
	l.init()

	return l
}

func (self *Logic) init() {
	http.HandleFunc("/robot/save_friends", self.RobotSaveFriends)
	http.HandleFunc("/robot/receive_msg", self.RobotReceiveMsg)

	http.HandleFunc("/robotgroup/get_list", self.RobotGroupList)
	http.HandleFunc("/robotgroup/get_info_list", self.RobotGroupInfoList)
	http.HandleFunc("/robotgroup/create", self.CreateRobotGroup)
	http.HandleFunc("/robotgroup/add_group_info_list", self.AddRobotGroupInfoList)
	http.HandleFunc("/robotgroup/del_robot_group_info", self.DelRobotGroupInfo)
	http.HandleFunc("/robotgroup/get_robot_list", self.GetRobotList)
	http.HandleFunc("/robotgroup/del_robot_group", self.DelRobotGroup)
	http.HandleFunc("/robotgroup/update_robot_group", self.UpdateRobotGroup)

	http.HandleFunc("/robotfriend/get_robot_friend_list", self.GetRobotFriendList)
	http.HandleFunc("/robotfriend/get_robot_lower_friend_list", self.GetRobotLowerFriendList)
	http.HandleFunc("/robotfriend/update_robot_friend_reamrk", self.UpdateRobotFriendRemark)
	http.HandleFunc("/robotfriend/get_robot_friend_list_7days_nochat", self.GetRobotFriendList7DaysNoChat)
	http.HandleFunc("/robotfriend/get_robot_friend_list_active", self.GetRobotFriendListActive)
	http.HandleFunc("/robotfriend/get_robot_friend_chat_list", self.GetRobotFriendChatList)
	http.HandleFunc("/robotfriend/get_robot_new_chat_list", self.GetRobotNewChatList)
	http.HandleFunc("/robotfriend/send_text_msg", self.SendTextMsg)

	http.HandleFunc("/robottag/create_account_tag", self.CreateAccountTag)
	http.HandleFunc("/robottag/get_account_list_from_tag", self.GetAccountListFromTag)
	http.HandleFunc("/robottag/delete_account_tag", self.DeleteAccountTag)
	http.HandleFunc("/robottag/get_tag_list", self.GetTagList)
	
	http.HandleFunc("/robotwxcircle/create_wx_circle", self.CreateWxCircle)
	http.HandleFunc("/robotwxcircle/delete_wx_circle", self.DeleteWxCircle)
	http.HandleFunc("/robotwxcircle/update_wx_circle", self.UpdateWxCircle)
	http.HandleFunc("/robotwxcircle/get_wx_circle_list", self.GetWxCircleList)
	http.HandleFunc("/robotwxcircle/create_wx_circle_setting", self.CreateWxCircleSetting)
	http.HandleFunc("/robotwxcircle/delete_wx_circle_setting", self.DeleteWxCircleSetting)
	http.HandleFunc("/robotwxcircle/get_wx_circle_setting_list", self.GetWxCircleSettingList)
	http.HandleFunc("/robotwxcircle/get_robot_new_wx_circle", self.GetRobotNewWxCircle)
}

func (self *Logic) Run() {
	defer holmes.Start(holmes.LogFilePath("./log"),
		holmes.EveryDay,
		holmes.AlsoStdout,
		holmes.DebugLevel).Stop()

	if self.cfg.Debug {
		EnableDebug()
	}

	holmes.Info("server starting on[%s]..", self.cfg.Host)
	holmes.Infoln(http.ListenAndServe(self.cfg.Host, nil))
}

func WriteJSON(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "x-requested-with,content-type")
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func WriteBytes(w http.ResponseWriter, code int, v []byte) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "x-requested-with,content-type")
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	w.Write(v)
}

func EnableDebug() {
}
