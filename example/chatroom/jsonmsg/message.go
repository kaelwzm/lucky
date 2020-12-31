package jsonmsg

import (
	"github.com/helloh2o/lucky/core/iduck"
	"github.com/helloh2o/lucky/core/iproto"
	"github.com/helloh2o/lucky/example/chatroom/chatnode"
	"github.com/helloh2o/lucky/log"
)

const (
	// EnterRoomCode NO.
	EnterRoomCode = 1001
	// LeaveRoomCode NO.
	ChatMessageCode = 1002
	// LeaveRoomCode NO.
	LeaveRoomCode = 1003
	// JoinSuccessCode NO.
	JoinSuccessCode = 2001
)

// EnterRoom msg
type EnterRoom struct {
}

// JoinRoomSuccess msg
type JoinRoomSuccess struct {
}

// ChatMessage msg
type ChatMessage struct {
	FromName string
	Content  string
}

// LeaveRoom msg
type LeaveRoom struct {
}

// Processor is message handler
var Processor = iproto.NewJSONProcessor()

func init() {
	Processor.RegisterHandler(EnterRoomCode, &EnterRoom{}, func(args ...interface{}) {
		conn := args[iproto.Conn].(iduck.IConnection)
		conn.AfterClose(func() {
			chatnode.GetRoom().DelConn(conn.GetUuid())
		})
		chatnode.GetRoom().AddConn(conn)
		conn.SetNode(chatnode.GetRoom())
		conn.WriteMsg(&JoinRoomSuccess{})
		// 房间的最近20条历史消息
		msgs := <-chatnode.GetRoom().GetAllMessage()
		record := make([]interface{}, 0)
		if len(msgs) > 20 {
			record = append(record, msgs[:20]...)
		} else {
			record = msgs
			for _, m := range record {
				conn.WriteMsg(m)
			}
			log.Debug("write %d history message.", len(record))
		}
	})

	// 将聊天消息转发给节点
	Processor.RegisterHandler(ChatMessageCode, &ChatMessage{}, func(args ...interface{}) {
		conn := args[iproto.Conn].(iduck.IConnection)
		if nd := conn.GetNode(); nd != nil {
			nd.OnProtocolMessage(args[iproto.Msg].(*ChatMessage))
		}
	})

	Processor.RegisterHandler(LeaveRoomCode, &LeaveRoom{}, func(args ...interface{}) {
		conn := args[iproto.Conn].(iduck.IConnection)
		if nd := conn.GetNode(); nd != nil {
			nd.DelConn(conn.GetUuid())
		}
	})

	Processor.RegisterHandler(JoinSuccessCode, &JoinRoomSuccess{}, nil)
}
