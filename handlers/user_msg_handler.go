/*
 * @Descripttion: 
 * @Author: LongWeiYi
 * @Date: 2022-12-23 21:17:53
 * @LastEditors: LWYð
 * @LastEditTime: 2022-12-23 21:44:13
 * @FilePath: \wechatbot\handlers\user_msg_handler.go
 */
package handlers

import (
	"github.com/869413421/wechatbot/gtp"
	"github.com/eatmoreapple/openwechat"
	"log"
	"strings"
)

var _ MessageHandlerInterface = (*UserMessageHandler)(nil)

// UserMessageHandler ç§èæ¶æ¯å¤ç
type UserMessageHandler struct {
}

// handle å¤çæ¶æ¯
func (g *UserMessageHandler) handle(msg *openwechat.Message) error {
	if msg.IsText() {
		return g.ReplyText(msg)
	}
	return nil
}

// NewUserMessageHandler åå»ºç§èå¤çå¨
func NewUserMessageHandler() MessageHandlerInterface {
	return &UserMessageHandler{}
}

// ReplyText åéææ¬æ¶æ¯å°ç¾¤
func (g *UserMessageHandler) ReplyText(msg *openwechat.Message) error {
	// æ¥æ¶ç§èæ¶æ¯
	sender, err := msg.Sender()
	log.Printf("Received User %v Text Msg : %v", sender.NickName, msg.Content)
	if UserService.ClearUserSessionContext(sender.ID(), msg.Content) {
		_, err = msg.ReplyText("ä¸ä¸æå·²ç»æ¸ç©ºäºï¼ä½ å¯ä»¥é®ä¸ä¸ä¸ªé®é¢å¦ã")
		if err != nil {
			log.Printf("response user error: %v \n", err)
		}
		return nil
	}

	// è·åä¸ä¸æï¼åGPTåèµ·è¯·æ±
	requestText := strings.TrimSpace(msg.Content)
	requestText = strings.Trim(msg.Content, "\n")

	requestText = UserService.GetUserSessionContext(sender.ID()) + requestText
	reply, err := gtp.Completions(requestText)
	if err != nil {
		log.Printf("gtp request error: %v \n", err)
		msg.ReplyText("æºå¨äººç¥äºï¼æä¸ä¼åç°äºå°±å»ä¿®ã")
		return err
	}
	if reply == "" {
		return nil
	}

	// è®¾ç½®ä¸ä¸æï¼åå¤ç¨æ·
	reply = strings.TrimSpace(reply)
	reply = strings.Trim(reply, "\n")
	UserService.SetUserSessionContext(sender.ID(), requestText, reply)
	reply = reply
	_, err = msg.ReplyText(reply)
	if err != nil {
		log.Printf("response user error: %v \n", err)
	}
	return err
}
