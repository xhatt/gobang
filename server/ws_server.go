package server

import (
	"encoding/json"
	"fmt"
	"net"
)

//数据缓冲区
const BUFFSIZE = 1024

var buff = make([]byte, BUFFSIZE)

type Msg struct {
	Action int    `json:"action"`  // 这个消息的类型是干嘛的 0-连接 1-创建房间 2-加入房间 3-退出房间 4-落子
	RoomID string `json:"room_id"` // 房间号
	X      int    `json:"x"`
	Y      int    `json:"y"`
}

//消息解析，[]byte -> []string
func AnalyzeMessage(buff []byte, len int) string {

	return string(buff[len:])
}

// 处理消息
func handleMsg() {
	n, addr, err := listener.ReadFromUDP(buff) // 接收数据
	if err != nil {
		fmt.Println("read udp failed, err: ", err)
		return
	}
	if n == 0 {
		return
	}
	//msgStr := AnalyzeMessage(buff, n)
	var msg = Msg{}
	err = json.Unmarshal(buff[:n], &msg)
	if err != nil {
		fmt.Println("Umarshal failed:", err)
		return
	}
	fmt.Println(addr.String())
	switch msg.Action {
	case 0:
		// 无动作
		p := getPlayer(addr)
		p.connectSuccess()
	case 1:
		// 创建房间 并加入
		p := getPlayer(addr)
		room := newRoom(p)
		if room == nil {
			// 说明他已经加入的有房间了
			fmt.Printf("玩家%s已经加入了房间%s\n", addr.String(), p.Room.ID)
			p.joinRoomSuccess()
			return
		}
		room.createRoomDone()
		fmt.Printf("玩家%s创建并加入了房间%s\n", addr.String(), room.ID)
	case 2:
		// 加入房间
		room, ok := RoomMap[msg.RoomID]
		p := getPlayer(addr)
		if ok {
			room.joinPlayerToRoom(p)
		} else {
			p.joinRoomFail("未找到房间")
		}
	case 3:
		// 退出
	case 4:
		fmt.Println("收到玩家走子消息", msg.X, msg.Y)
		// 判断这个位置能不能走
		p := getPlayer(addr)
		p.confirm(msg.X, msg.Y)

	}

}

// UDP Server端
func RunServer() {
	var err error
	listener, err = net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 30000,
	})
	if err != nil {
		fmt.Println("Listen failed, err: ", err)
		return
	}
	defer listener.Close()
	for {
		handleMsg()
	}
}
