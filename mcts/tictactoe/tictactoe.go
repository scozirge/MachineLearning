package tictactoe

const (
	None    = 0
	Player1 = 1
	Player2 = 2
)

// 棋局結果
type GameResult struct {
	IsTerminal bool
	Winner     int
}

// 定義贏線
var WinLines = [][3]int{
	{0, 1, 2},
	{3, 4, 5},
	{6, 7, 8},
	{0, 3, 6},
	{1, 4, 7},
	{2, 5, 8},
	{0, 4, 8},
	{2, 4, 6},
}

// 棋況
type GameState struct {
	Board      [9]int
	LastPlaced int
}

// 建立新的一局棋況
func New() *GameState {
	return &GameState{
		Board:      [9]int{},
		LastPlaced: -1,
	}
}

// 傳入棋況取得目前結果
func (t GameState) GetGameState() GameResult {
	for _, line := range WinLines {
		if t.Board[line[0]] != None &&
			t.Board[line[0]] == t.Board[line[1]] &&
			t.Board[line[0]] == t.Board[line[2]] {
			return GameResult{true, t.Board[line[0]]}
		}
	}

	return GameResult{IsFull(t), None}
}

// 確認是否有空格
func IsFull(t GameState) bool {
	for i := 0; i < 9; i++ {
		if t.Board[i] == None {
			return false
		}
	}
	return true
}

// 傳入棋盤取得仍可放置的空格位置
func (t *GameState) GetLegalPosz() []int {
	var emptyPosz []int

	for i := 0; i < 9; i++ {
		if t.Board[i] == None {
			emptyPosz = append(emptyPosz, i)
		}
	}

	return emptyPosz
}

// 取消動作
func (t *GameState) UndoAction(pos int) {
	t.Board[pos] = None
}

// 傳入棋況取得目前換哪位玩家行動
func (t *GameState) CurrentPlayer() int {
	p1Count := 0
	p2Count := 0

	for i := 0; i < 9; i++ {
		if t.Board[i] == Player1 {
			p1Count++
		} else if t.Board[i] == Player2 {
			p2Count++
		}
	}

	if p2Count == p1Count {
		return Player1
	}
	return Player2
}

// 畫出棋況結果圖
func (state GameState) DrawTable() string {
	symbols := []rune{' ', 'O', 'X'}
	gameStr := ""
	for i, index := range state.Board {
		gameStr += string(symbols[index])
		if i%3 == 2 {
			gameStr += "\n"
		} else {
			gameStr += "|"
		}
	}
	return gameStr
}
