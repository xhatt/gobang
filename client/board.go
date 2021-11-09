package client

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"os"
	"strings"
	"sync"
)

type ClassBoard interface {
	MoveCursor() // 移动光标
	Confirm()    // 确认落子
	Draw()       // 绘制棋盘
}

type ChessColor int // 棋子颜色
const (
	White      ChessColor = 1
	Black      ChessColor = 2
	WhitePiece            = "o️"
	BlackPiece            = "x"

	max   = 18
	piece = "."
)

var (
	pieceMap = map[ChessColor]string{
		White: WhitePiece,
		Black: BlackPiece,
	}
	connect *Connect
	board   *Board
)

type Board struct {
	x             int
	y             int
	chessColor    ChessColor           // 棋子颜色  1-白色 2-黑色  不可被更改，实例化时从服务端获取
	Checkerboard  [max][max]ChessColor // 当前本地棋盘的全貌
	CanUnder      bool                 // 判断是否可以下子
	mutex         *sync.Mutex          // 棋盘的锁，操作棋盘时需要加锁
	RoomID        string               // 当前的房间号 默认 ""
	CurrentPlayer ChessColor           // 当前谁在走子
	//EnemySurplusTimes int         // 对方剩余时间
	//WeSurplusTimes    int         // 我方剩余时间
}

// 绘制棋盘
func (b *Board) Draw() {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	Clear()
	ResetCursor()
	fmt.Printf("方向键移动光标，空格键落子, q|Q 退出程序\n我使用的符号是：%s\n正在走子的玩家是: %s\n", pieceMap[board.chessColor], pieceMap[b.CurrentPlayer])
	for yIndex, xItem := range b.Checkerboard {
		//fmt.Print(x)
		for xIndex, yItem := range xItem {
			var (
				p  = piece
				ok bool
			)
			if b.x == xIndex && b.y == yIndex {
				// 说明光标选中了这个元素
				if p, ok = pieceMap[yItem]; ok {
					p = BackgroundRedColor(pieceMap[yItem])
				} else {
					p = BackgroundRedColor(piece)
				}
			} else {
				if p, ok = pieceMap[yItem]; ok {
					p = pieceMap[yItem]
				} else {
					p = piece

				}
			}

			fmt.Printf("  %s", p)

		}
		fmt.Println()
	}

}

func (b *Board) Confirm() {
	// 判断是否轮到自己下，下子，获取当前坐标，向服务器发送
	//b.Checkerboard[b.y][b.x] = BlackPiece
	if b.CanUnder {
		connect.confirm(b.x, b.y)
		select {
		case msg := <-connect.boardChan:
			if msg.Code != 0 {
				return
			}
		}
		b.Checkerboard[b.y][b.x] = b.chessColor
		b.CanUnder = false // 下了之后等待服务器信号来重置
		b.changeCurrentPlayer()
	}
}

func (b *Board) changeCurrentPlayer() {
	if b.CanUnder {
		b.CurrentPlayer = b.chessColor
	} else {
		if b.chessColor == White {
			b.CurrentPlayer = Black
		} else {
			b.CurrentPlayer = White
		}
	}
}

// 对手落子，修改棋盘并渲染
func (b *Board) modifyBoard(x, y int, color ChessColor) {
	b.CurrentPlayer = b.chessColor
	b.Checkerboard[y][x] = color
	b.CanUnder = true // 下了之后等待服务器信号来重置
	b.x = x
	b.y = y
	b.Draw()
}

func NewCheckerboard() (bo [max][max]ChessColor) {
	for i := 0; i < max; i++ {
		var temp [18]ChessColor
		for j := 0; j < max; j++ {
			temp[j] = 0
		}
		bo[i] = temp
	}
	return

}

func NewBoard(chessColor ChessColor) *Board {
	return &Board{
		chessColor:   chessColor,
		Checkerboard: NewCheckerboard(),
		mutex:        &sync.Mutex{},
	}
}

func Run() {
	Clear()
	connect = NewConnect()
	connect.heartbeat()
	fmt.Println("1. 创建并加入一个房间")
	fmt.Println("2. 加入一个已有的房间")
	id := ""
	Input("请输入序号:", &id)
	switch id {
	case "1":
		if board != nil {
			fmt.Println("已经加入了房间：", board.RoomID)
		} else {
			connect.createRoom()
			select {
			case msg := <-connect.roomChan:
				if msg.Code == 0 {
					fmt.Printf("房间已创建并加入，等待其他玩家加入。\n房间号: %s", msg.RoomID)
					select {
					case msg = <-connect.startChan:
						board = NewBoard(msg.ChessColor)
						board.CanUnder = msg.CanUnder
						board.changeCurrentPlayer()
						board.Draw()

						board.HandleKeyboard()
					}
				}
			}
		}
	case "2":
		roomID := ""
		Input("请输入房间号:", &roomID)
		connect.joinRoom(roomID)
		select {
		case msg := <-connect.roomChan:
			if msg.Code == 0 {
				select {
				case msg = <-connect.startChan:
					board = NewBoard(msg.ChessColor)
					board.CanUnder = msg.CanUnder
					board.changeCurrentPlayer()
					board.Draw()
					board.HandleKeyboard()
				}
			}
		}

	}

}

// 把光标往上移动
func (b *Board) moveCursorUp() {
	if b.y > 0 {
		b.y--
	}
}

func (b *Board) moveCursorDown() {
	if b.y < max-1 {
		b.y++
	}
}

func (b *Board) moveCursorRight() {
	if b.x < max-1 {
		b.x++
	}
}
func (b *Board) moveCursorLeft() {
	if b.x > 0 {
		b.x--
	}
}

// 移动光标
func (b *Board) MoveCursor(key keyboard.Key) {
	switch key {
	case keyboard.KeyArrowUp:
		b.moveCursorUp()
	case keyboard.KeyArrowDown:
		b.moveCursorDown()
	case keyboard.KeyArrowLeft:
		b.moveCursorLeft()
	case keyboard.KeyArrowRight:
		b.moveCursorRight()
	}
}

// 处理键盘除字母键以外的按键
func (b *Board) handleKey(key keyboard.Key) {
	switch key {
	case keyboard.KeyArrowRight, keyboard.KeyArrowLeft, keyboard.KeyArrowDown, keyboard.KeyArrowUp:
		b.MoveCursor(key)
	case keyboard.KeySpace:
		b.Confirm()
	}
	b.Draw()
}

// 处理字母按键
func (b *Board) handleChar(char rune) {
	ch := strings.ToLower(string(char))
	if ch == "q" {
		ShowCursor()
		os.Exit(0)
	}
	//} else if ch == "d" {
	//	//global.Delete = true
	//	global.setDelete()
	//	//os.Exit(0)
	//} else if ch == "w" {
	//	global.showDetail()
	//	GetTotalLength()
	//} else if ch == "a" {
	//	AddServer()
	//}
	//Flush()
}

// 处理键盘事件
func (b *Board) HandleKeyboard() {
	for {
		char, key, err := keyboard.GetSingleKey()
		if err != nil {
			panic(err)
		}
		if char != 0 {
			// 说明按下的是字符
			b.handleChar(char)
		} else if key != 0 {
			b.handleKey(key)
		} else {

		}
	}
}
