package service

import (
	"chat/cache"
	"chat/conf"
	"chat/pkg/xcode"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"time"

	logging "github.com/sirupsen/logrus"
)

const month = 60*60*24*30  // 按照30天一个月

// SendMsg 发消息的类型
type SendMsg struct {
	Type int `json:"type"`
	Content string `json:"content"`
}


// ReplyMsg 回复的消息
type ReplyMsg struct {
	From string `json:"from"`
	Code int `json:"code"`
	Content string `json:"content"`
}

// Client 用户类
type Client struct {
	Id string
	SendId string
	Socket *websocket.Conn
	Send chan []byte
}

// Broadcast 广播类，包括广播内容和源用户
type Broadcast struct {
	Client *Client
	Message []byte
	Type int
}

// ClientManager 用户管理
type ClientManager struct {
	Clients  map[string]*Client
	Broadcast chan *Broadcast
	Reply chan *Client
	Register chan *Client
	UnRegister chan *Client
}

// Message 信息转json 包括发送者，接收者，内容
type Message struct {
	Sender string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content string `json:"content,omitempty"`
}

var Manager = ClientManager{
	Clients: make(map[string]*Client),  // 参与连接的用户
	Broadcast: make(chan *Broadcast),
	Reply: make(chan *Client),
	Register: make(chan *Client),
	UnRegister: make(chan *Client),

}

func CreateId(uid, toUid string) string {
	return uid + "->" + toUid
}

func Handler (c *gin.Context)  {
	uid := c.Query("uid")
	toUid := c.Query("toUid")

	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(c.Writer, c.Request, nil)  // 升级ws协议
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}

	// 创建一个用户实例
	client := &Client{
		Id: CreateId(uid, toUid),
		SendId: CreateId(toUid, uid),
		Socket: conn,
		Send: make(chan []byte),
	}

	// 用户注册到用户管理上
	Manager.Register <- client
	go client.Read()
	go client.Write()
}

func (c *Client) Read() {
	defer func() {
		Manager.UnRegister <- c
		_ = c.Socket.Close()
	}()

	for  {
		c.Socket.PongHandler()

		sendMsg := new(SendMsg)
		//c.Socket.ReadMessage()
		err := c.Socket.ReadJSON(sendMsg)
		if err != nil {
			logging.Errorf("websocket readjson failed, err:%v \n", err)
			Manager.UnRegister <- c
			_ = c.Socket.Close()
			break
		}

		if sendMsg.Type == 1 { // 1->2 发送消息
			r1, _ := cache.RedisClient.Get(context.Background(), c.Id).Result()
			r2, _ := cache.RedisClient.Get(context.Background(), c.SendId).Result()

			if r1 >= "3" && r2 == "" {  // 限制单聊次数
				replyMsg := ReplyMsg{
					Code: xcode.WebSocketLimit,
					Content: "number of times reached",
				}

				msg, _ := json.Marshal(replyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
				// 一个月以后才可以重新发消息，防止骚扰
				_, _ = cache.RedisClient.Expire(context.Background(), c.Id, time.Hour*24*30).Result()
				continue
			} else {
				cache.RedisClient.Incr(context.Background(), c.Id)
				// 防止过快 丢失 连接数据，建立连接3个月过期
				_, _ = cache.RedisClient.Expire(context.Background(), c.Id, time.Hour*24*30*3).Result()
			}

			logging.Info(c.Id, "----send msg:----", sendMsg.Content)
			Manager.Broadcast <- &Broadcast{
				Client: c,
				Message: []byte(sendMsg.Content),
			}
		} else if sendMsg.Type == 2 { //拉取历史消息
			timeT, err := strconv.Atoi(sendMsg.Content) // 传送来时间
			if err != nil {
				timeT = 999999999
			}
			results, _ := FindMany(conf.MongoDBName, c.SendId, c.Id, int64(timeT), 10)
			if len(results) > 10 {
				results = results[:10]
			} else if len(results) == 0 {
				replyMsg := ReplyMsg{
					Code:    xcode.WebSocketEnd,
					Content: "到底了",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)

				continue
			}
			for _, result := range results {
				replyMsg := ReplyMsg{
					From:    result.From,
					Content: fmt.Sprintf("%s", result.Msg),
				}
				msg, _ := json.Marshal(replyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)

			}
		} else if sendMsg.Type == 3 {
			results, err := FirsFindMsg(conf.MongoDBName, c.SendId, c.Id)
			fmt.Printf("-----%v\n", results)
			if err != nil {
				log.Println(err)
			}
			for _, result := range results {
				replyMsg := ReplyMsg{
					From:    result.From,
					Content: fmt.Sprintf("%s", result.Msg),
				}
				msg, _ := json.Marshal(replyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
			}

		}
	}

}

func (c *Client) Write() {
	defer func() {
		_ = c.Socket.Close()
	}()

	for  {
		select {
		case msg, ok := <- c.Send:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			logging.Info(c.Id, "receive msg: ", string(msg))

			replyMsg := &ReplyMsg{
				Code: xcode.WebSocketSuccessMsg,
				Content: fmt.Sprintf("%s", string(msg)),
			}

			msg2, _ := json.Marshal(replyMsg)
			_ = c.Socket.WriteMessage(websocket.TextMessage, msg2)
		}
	}

}