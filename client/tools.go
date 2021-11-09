package client

import (
	"bufio"
	"fmt"
	"os"
)

// 给字符串添加背景颜色
func BackgroundColor(s string) string { // d   b   f
	return fmt.Sprintf("%c[%d;%d;%dm%s%c[0m", 0x1B, 1, 40, 32, s, 0x1B)
}

// 给字符串添加红色背景颜色
func BackgroundRedColor(s string) string { // d   b   f
	return fmt.Sprintf("%c[%d;%d;%dm%s%c[0m", 0x1B, 5, 41, 32, s, 0x1B)
}

// 给字符串添加红色背景颜色
//func BackgroundRedColor(s string) string { // d   b   f
//	return fmt.Sprintf("\033[5m%s", s)
//}

// 打印（不换行）
// 字体颜色为默色
func Log(a ...interface{}) {
	fmt.Print(a...)
}

// 打印信息（不换行）
// 字体颜色为绿色
func Info(a ...interface{}) {

	fmt.Print("\033[32m")
	Log(a...)
	fmt.Print("\033[0m")
}

// 重置光标位置
func ResetCursor() {
	// 第0行，第0列
	fmt.Printf("\033[0;0H")
}

func Clear() {
	fmt.Printf("\033c")
	//fmt.Printf("\033[2J")
}

// 隐藏光标
func HideCursor() {
	fmt.Printf("\033[?25l")
}

// 显示光标
func ShowCursor() {
	fmt.Printf("\033[?25h")
}

func SetCursorSeek(w, h int) {
	fmt.Printf("\033[%d;%dH", h, w)
}


// GO自带的fmt.Scanln将空格也当作结束符，若需要读取含有空格的句子请使用该方法
func Scanln(a *string) {
	reader := bufio.NewReader(os.Stdin)
	data, _, err := reader.ReadLine()
	if err != nil {
		fmt.Print(err.Error())
	}
	if data == nil {
		os.Exit(0)
	}
	*a = string(data)
}

func Input(tips string, v *string) {
	Info(tips)
	Scanln(v)
}
