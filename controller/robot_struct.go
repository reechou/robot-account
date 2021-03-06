package controller

const (
	RECEIVE_EVENT_MSG        = "receivemsg"
	RECEIVE_EVENT_ADD_FRIEND = "addfriend"
	RECEIVE_EVENT_ADD        = "receiveadd"

	CHAT_TYPE_PEOPLE = "people"
	CHAT_TYPE_GROUP  = "group"
	MSG_TYPE_TEXT    = "text"
	MSG_TYPE_IMG     = "img"
)

type UserFriend struct {
	Alias       string `json:"alias"`
	City        string `json:"city"`
	VerifyFlag  int    `json:"verifyFlag"`
	ContactFlag int    `json:"contactFlag"`
	NickName    string `json:"nickName"`
	RemarkName  string `json:"remarkName"`
	Sex         int    `json:"sex"`
	UserName    string `json:"userName"`
}

type BaseInfo struct {
	Uin           string `json:"uin"`
	UserName      string `json:"userName,omitempty"`   // 机器人username
	WechatNick    string `json:"wechatNick,omitempty"` // 微信昵称
	ReceiveEvent  string `json:"receiveEvent,omitempty"`
	FromType      string `json:"fromType,omitempty"`
	FromUserName  string `json:"fromUserName,omitempty"`
	FromNickName  string `json:"fromNickName,omitempty"`
	FromGroupName string `json:"fromGroupName,omitempty"`
}

type BaseToUserInfo struct {
	ToUserName  string `json:"toUserName,omitempty"`
	ToNickName  string `json:"toNickName,omitempty"`
	ToGroupName string `json:"toGroupName,omitempty"`
}

type SendBaseInfo struct {
	WechatNick string `json:"wechatNick,omitempty"` // 微信昵称
	ChatType   string `json:"chatType,omitempty"`
	NickName   string `json:"nickName,omitempty"`
	UserName   string `json:"userName,omitempty"`
	MsgType    string `json:"msgType,omitempty"`
	Msg        string `json:"msg,omitempty"`
}

type RetResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
}

type AddFriend struct {
	SourceWechat string `json:"sourceWechat,omitempty"`
	SourceNick   string `json:"sourceNick,omitempty"`
	UserWxid     string `json:"userWxid,omitempty"`
	UserWechat   string `json:"userWechat,omitempty"`
	UserNick     string `json:"userNick,omitempty"`
	UserCity     string `json:"userCity,omitempty"`
	UserSex      int    `json:"userSex,omitempty"`
	Ticket       string `json:"-"` // for verify
}

type ReceiveMsgInfo struct {
	BaseInfo       `json:"baseInfo,omitempty"`
	BaseToUserInfo `json:"baseToUserIno,omitempty"`
	AddFriend      `json:"addFriend,omitempty"`
	
	MsgType      string `json:"msgType,omitempty"`
	Msg          string `json:"msg,omitempty"`
	MediaTempUrl string `json:"mediaTempUrl,omitempty"`
}

type CallbackMsgInfo struct {
	RetResponse `json:"retResponse,omitempty"`
	BaseInfo    `json:"baseInfo,omitempty"`

	CallbackMsgs []SendBaseInfo `json:"msg,omitempty"`
}

type SendMsgInfo struct {
	SendMsgs []SendBaseInfo `json:"sendBaseInfo,omitempty"`
}

type SendMsgResponse struct {
	RetResponse `json:"retResponse,omitempty"`
}
