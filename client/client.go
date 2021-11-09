package client

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type Msg struct {
	Action     int        `json:"action"` // 这个消息的类型是干嘛的 0-连接 1-创建房间 2-加入房间 3-退出房间 4-落子
	Code       int        `json:"code"`
	Msg        string     `json:"msg"`
	RoomID     string     `json:"room_id"`
	X          int        `json:"x"` // 落子的x
	Y          int        `json:"y"` // Y
	ChessColor ChessColor `json:"chess_color"`
	CanUnder   bool       `json:"can_under"` // 判断是否可以下子

}

type Connect struct {
	socket    *net.UDPConn
	roomChan  chan *Msg // 关于房间的通道
	boardChan chan *Msg // 关于走子的通道
	startChan chan *Msg // 游戏开始与结束的通道
	buff      []byte
}

func (c *Connect) send(msg map[string]interface{}) {
	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.socket.Write(data)
}

//心跳包
func (c *Connect) heartbeat() {
	c.send(map[string]interface{}{"action": 0})
}

// 创建房间
func (c *Connect) createRoom() {
	go c.send(map[string]interface{}{"action": 1})
}

// 加入房间
func (c *Connect) joinRoom(roomID string) {
	go c.send(map[string]interface{}{"action": 2, "room_id": roomID})
}

// 落子
func (c *Connect) confirm(x, y int) {
	go c.send(map[string]interface{}{"action": 4, "x": x, "y": y})
}

// 处理服务器消息
func (c *Connect) handleMsg() {
	n, _, err := c.socket.ReadFromUDP(c.buff) // 接收数据
	if err != nil {
		fmt.Println("接收数据失败, err: ", err)
		return
	}
	if n == 0 {
		return
	}
	var msg = Msg{}
	err = json.Unmarshal(c.buff[:n], &msg)
	//fmt.Println(msg)
	if err != nil {
		fmt.Println("Umarshal failed:", err)
		return
	}

	switch msg.Action {
	case 0:

	case 1:
		c.roomChan <- &msg
	case 2:
		// 加入房间
		c.roomChan <- &msg
	case 3:
		// 退出

	case 4:
		// 走子
		//c.boardChan <- &msg
		board.modifyBoard(msg.X, msg.Y, msg.ChessColor)
	case 6: // 游戏开始
		if msg.Code == 0 {
			c.startChan <- &msg
		} else if msg.Code == 1 {
			// 游戏结束 打印赢家
			board.Draw()
			if msg.ChessColor == board.chessColor {
				fmt.Println("游戏结束。你赢了")
			} else {
				fmt.Println("游戏结束。你输了")
			}
			os.Exit(0)
		}
	case 7:
		// 自己走子的确认信息
		c.boardChan <- &msg
	}
}

func NewConnect() *Connect {
	socket, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 30000,
	})
	if err != nil {
		fmt.Println("连接UDP服务器失败，err: ", err)
		os.Exit(0)
	}
	c := &Connect{
		socket:    socket,
		buff:      make([]byte, 1024),
		roomChan:  make(chan *Msg),
		boardChan: make(chan *Msg),
		startChan: make(chan *Msg),
	}
	go func() {
		for {
			c.handleMsg()
		}
	}()
	return c
}
