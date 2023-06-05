package ticTacToe

import "fmt"

// 表示棋盤的狀態(0:空格 1:圈圈 2:叉叉)
type State [9]int

// 將state轉成唯一字串作為QTable要輸出成json所用的的的key值
func (state State) ToKeyString() string {
	gameStr := ""
	for i, v := range state {
		gameStr += fmt.Sprint(v)
		if i != len(state)-1 {
			gameStr += "|"
		}
	}
	return gameStr
}

// 畫出棋況結果圖
func (state State) DrawTable() string {
	symbols := []rune{' ', 'O', 'X'}
	gameStr := ""
	for i, index := range state {
		gameStr += string(symbols[index])
		if i%3 == 2 {
			gameStr += "\n"
		} else {
			gameStr += "|"
		}
	}
	return gameStr
}

//定義每條贏線
type WinLine [3]int

var WinLines = []WinLine{
	{0, 1, 2},
	{3, 4, 5},
	{6, 7, 8},
	{0, 3, 6},
	{1, 4, 7},
	{2, 5, 8},
	{0, 4, 8},
	{2, 4, 6},
}

// 計算該盤的O/X數量 token傳入1就是計算O的數量 傳入2就是計算X的數量
func Count(token int, state State) int {
	cnt := 0
	for _, v := range state {
		if v == token {
			cnt++
		}
	}
	return cnt
}

// 判斷傳入棋局是否結束了回傳 是否結束,贏家toekn
func IsGameFinished(state State) (bool, int) {
	//根據每條贏線判斷是否有一方勝利了
	for _, line := range WinLines {
		token := state[line[0]]
		if token != 0 && token == state[line[1]] && token == state[line[2]] {
			return true, token
		}
	}

	// 判斷是否還有空位，如果有，則遊戲尚未結束
	for _, v := range state {
		if v == 0 {
			return false, 0
		}
	}

	// 遊戲結束，並且沒有贏家(平手)
	return true, 0
}
