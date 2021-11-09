package server

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

type ChessColor int // 棋子颜色
const (
	White ChessColor = 1
	Black ChessColor = 2

	max = 18
)

type Room struct {
	ID           string               // 房间号
	Checkerboard [max][max]ChessColor // 这个房间的棋盘
	PlayerA      *Player              // 玩家A
	PlayerB      *Player              // 玩家A
	CanJoin      bool                 // 还能否加入玩家
}

var (
	connectPlayer = make(map[string]*Player, 0) // 当前连接的玩家
	listener      *net.UDPConn
	RoomMap       = make(map[string]*Room, 0) // 当前服务器上的所有房间
)

func generateRoomID() string {
	return fmt.Sprintf("%04v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000)) //这里面前面的04v是和后面的1000相对应的
}

// 创建房间完成
func (r *Room) createRoomDone() {
	r.PlayerA.joinRoomSuccess()
}

// 创建房间
func newRoom(p *Player) *Room {
	var room *Room
	if p.Room != nil {
		return nil
	} else {
		room = &Room{
			ID:      generateRoomID(),
			PlayerA: p,
			CanJoin: true,
		}
		p.Room = room
		RoomMap[room.ID] = room
	}

	return room
}

// 把一个玩家加入房间
func (r *Room) joinPlayerToRoom(p *Player) {
	// 判断是不是这个人就加入了自己创建的房间

	if r.CanJoin && (r.PlayerA == nil || r.PlayerB == nil) {
		if r.PlayerA == p || r.PlayerB == p {
			p.joinRoomFail("已经在房间里了，不能重复加入")
			return
		}
		if r.PlayerA == nil {
			r.PlayerA = p
		} else if r.PlayerB == nil {
			r.PlayerB = p
		}
		p.Room = r
		if r.PlayerA != nil && r.PlayerB != nil {
			p.joinRoomSuccess()

			r.PlayerA.StartGame(White)
			r.PlayerB.StartGame(Black)
			r.PlayerA.OtherPlayer = r.PlayerB
			r.PlayerB.OtherPlayer = r.PlayerA
		} else {
			p.joinRoomSuccess()
		}
		r.CanJoin = false
		return
	} else {
		// 房间已满无法加入
		p.joinRoomFail("房间玩家人数已满")
		return
	}
}
