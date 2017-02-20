package controller

import (
	"net/http"
	"encoding/json"
	"bytes"
	"fmt"
	"io/ioutil"
	
	"github.com/reechou/holmes"
	"github.com/reechou/robot-account/models"
)

type RobotExt struct {
	client *http.Client
}

func NewRobotExt() *RobotExt {
	return &RobotExt{
		client: &http.Client{},
	}
}

func (self *RobotExt) SendMsgs(robotWx string, msg *SendMsgInfo) error {
	holmes.Debug("msg: %v", msg)
	// get robot
	robot := &models.Robot{
		RobotWx: robotWx,
	}
	has, err := models.GetRobot(robot)
	if err != nil {
		holmes.Error("get robot error: %v", err)
		return err
	}
	if !has {
		return fmt.Errorf("cannot found this robot[%s]", robot)
	}
	
	reqBytes, err := json.Marshal(msg)
	if err != nil {
		holmes.Error("json encode error: %v", err)
		return err
	}
	
	url := "http://" + robot.Ip + robot.OfPort + "/sendmsgs"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBytes))
	if err != nil {
		holmes.Error("http new request error: %v", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := self.client.Do(req)
	if err != nil {
		holmes.Error("http do request error: %v", err)
		return err
	}
	defer resp.Body.Close()
	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		holmes.Error("ioutil ReadAll error: %v", err)
		return err
	}
	var response SendMsgResponse
	err = json.Unmarshal(rspBody, &response)
	if err != nil {
		holmes.Error("json decode error: %v [%s]", err, string(rspBody))
		return err
	}
	if response.Code != 0 {
		holmes.Error("send msg[%v] result code error: %d %s", msg, response.Code, response.Msg)
		return fmt.Errorf("send msg result error.")
	}
	
	return nil
}
