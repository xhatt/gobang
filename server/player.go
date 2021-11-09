package server

import (
	"encoding/json"
	"fmt"
	"net"
)

var (
	OkMap = map[string]interface{}{"code": 0, "msg": "加入房间成功"}
)

type Player struct {
	addr        *net.UDPAddr
	Room        *Room   // 玩家当前加入的房间
	OtherPlayer *Player // 另外一个玩家
	ChessColor  ChessColor
}

func (p *Player) send(msg map[string]interface{}) {
	// 同步其他人数据
	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	listener.WriteToUDP(data, p.addr)
}

func (p *Player) StartGame(color ChessColor) {
	data := p.getOkMap()
	data["action"] = 6
	data["room_id"] = p.Room.ID
	data["msg"] = "开始游戏"
	data["chess_color"] = color
	if color == White {
		data["can_under"] = true
	}
	p.ChessColor = color
	go p.send(data)
}

func (p *Player) Win() {
	// 该玩家获得游戏胜利
	data := p.getOkMap()
	data["action"] = 6
	data["code"] = 1
	data["chess_color"] = p.ChessColor
	go p.send(data)
	go p.OtherPlayer.send(data)


}

func (p *Player) joinRoomSuccess() {
	// 玩家加入房间成功
	data := p.getOkMap()
	data["action"] = 1
	data["room_id"] = p.Room.ID
	data["msg"] = "加入房间成功"
	p.send(data)
}

func (p *Player) joinRoomFail(msg string) {
	// 玩家加入房间失败
	data := p.getOkMap()
	data["msg"] = msg
	data["code"] = 1
	data["action"] = 1
	p.send(data)
}

func (p *Player) check(x, y int) bool {
	if p.Room.Checkerboard[y][x] == 0 {
		return true
	} else {
		return false
	}

}

func (p *Player) confirm(x, y int) {

	data := p.getOkMap()
	data["action"] = 4
	data["x"] = x
	data["y"] = y
	data["chess_color"] = p.ChessColor
	if !p.check(x, y) {
		fmt.Println("走子无效，这个位置已经有其他棋子了")
		data["code"] = 1
		data["action"] = 7
		p.send(data)
		return
	} else {
		data1 := p.getOkMap()
		data1["action"] = 7
		p.send(data1)
		p.OtherPlayer.send(data)
	}

	p.Room.Checkerboard[y][x] = p.ChessColor
	// 走完了之后，判断输赢
	// 玩家加入房间失败
	if CheckWin(y, x, p.ChessColor, p.Room.Checkerboard) {
		fmt.Println("游戏结束")
		// 发送游戏结束通知
		p.Win()
		return
	}

}

func (p *Player) connectSuccess() {
	// 玩家连接服务器成功
	data := p.getOkMap()
	data["msg"] = "连接成功"
	p.send(data)
}

func getPlayer(addr *net.UDPAddr) *Player {
	if _, ok := connectPlayer[addr.String()]; !ok {
		p := &Player{addr, nil, nil, 0}
		connectPlayer[addr.String()] = p
	}
	return connectPlayer[addr.String()]
}

func (p *Player) getOkMap() map[string]interface{} {
	return map[string]interface{}{"code": 0, "msg": "", "action": 0}
}

func luozi(x, y int, color ChessColor, board [max][max]ChessColor) bool {
	board[x][y] = color
	for i := range direct {
		count := 0
		_x := x
		_y := y
		for j := 0; j < 5; j++ {
			if judgeValid(_x, _y, color, board) {
				_x += direct[i][0]
				_y += direct[i][1]
				count++
				continue
			}
			break
		}
		_x = x - direct[i][0]
		_y = y - direct[i][1]
		for j := 0; j < 5; j++ {
			if judgeValid(_x, _y, color, board) {
				_x -= direct[i][0]
				_y -= direct[i][1]
				count++
				continue
			}
			break
		}

		if count >= 5 {
			return true
		}
	}

	return false
}

var direct = [][]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}

func judgeValid(x, y int, color ChessColor, board [max][max]ChessColor) bool {
	if x < 0 || x >= len(board) || y < 0 || y >= len(board[0]) {
		return false
	}

	if board[x][y] == color {
		return true
	}

	return false
}

func CheckWin(x, y int, color ChessColor, board [max][max]ChessColor) bool {
	return luozi(x, y, color, board)
}
