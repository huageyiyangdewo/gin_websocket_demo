package service

import (
	"chat/conf"
	"chat/pkg/xcode"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	logging "github.com/sirupsen/logrus"
)

func (cm *ClientManager) Start() {
	for {
		logging.Info("-----listen chan msg------")
		select {
		case conn := <-cm.Register: // 建立连接
			logging.Infof("create new conn: %v", conn.Id)

			Manager.Clients[conn.Id] = conn
			replyMsg := &ReplyMsg{
				Code:    xcode.WebSocketSuccess,
				Content: "connected service",
			}
			msg, _ := json.Marshal(replyMsg)
			_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)
		case conn := <-cm.UnRegister: // 断开连接
			logging.Info("disconnect: %v", conn.Id)
			if _, ok := Manager.Clients[conn.Id]; ok {
				replyMsg := &ReplyMsg{
					Code:    xcode.WebSocketEnd,
					Content: "connect has been disconnected",
				}

				msg, _ := json.Marshal(replyMsg)
				_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)

				close(conn.Send) // 关闭通道
				delete(Manager.Clients, conn.Id)
			}

		case broadcast := <-Manager.Broadcast: // 广播信息
			message := broadcast.Message
			sendId := broadcast.Client.SendId
			flag := false // 默认对方不在线
			for id, conn := range Manager.Clients {
				if id != sendId {
					continue
				}

				select {
				case conn.Send <- message:
					flag = true
				default: // select 的其他分支没有准备好时，才执行这个分支，所以这里应该是不会被调用
					close(conn.Send)
					delete(Manager.Clients, conn.Id)
				}
			}

			id := broadcast.Client.Id
			if flag {
				logging.Info("对方在线应答")

				replyMsg := &ReplyMsg{
					Code:    xcode.WebSocketOnLineReply,
					Content: "other online response",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)

				// read 1 表示已读，这里简单处理
				err := InsertMsg2Mongo(conf.MongoDBName, id, string(message), 1, int64(3*month))
				if err != nil {
					logging.Errorf("InsertMsg2Mongo err: %s", err)
				}
			} else {
				logging.Info("对方不在线应答")

				replyMsg := &ReplyMsg{
					Code:    xcode.WebSocketOffLineReply,
					Content: "other offline response",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)

				// read 0 表示未读，这里简单处理
				fmt.Println("xxxxxx", conf.MongoDBClient)
				err := InsertMsg2Mongo(conf.MongoDBName, id, string(message), 0, int64(3*month))
				if err != nil {
					logging.Errorf("InsertMsg2Mongo err: %s", err)
				}

			}
		}
	}
}