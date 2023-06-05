package main

import (
	"fmt"
	"math/rand"
	"time"

	mcts "mcts/mcts"
	tictactoe "mcts/tictactoe"
)

const (
	playTimes = 1000
	selfPlay  = false
)

func main() {
	winTimes := 0
	tieTimes := 0
	rand.Seed(time.Now().UnixNano())
	// aa := mcts.MonteCarloTreeSearch(game, 10)
	// fmt.Println(aa)
	// return
	if selfPlay {
		for i := 0; i < playTimes; i++ {
			game := aiSelfPlay()
			winner := game.GetGameState().Winner
			if winner == tictactoe.None {
				tieTimes++
				fmt.Println("平手!")
			} else {
				if winner == tictactoe.Player1 {
					winTimes++
				}
				fmt.Printf("玩家 %d 獲勝!\n", winner)
			}
		}
		fmt.Printf("在%d局對戰中 玩家1的勝率為%.1f%% 平手率為%.1f%%", playTimes, float64(winTimes)/float64(playTimes)*100, float64(tieTimes)/float64(playTimes)*100)
	} else {
		game := playWithAI()
		winner := game.GetGameState().Winner
		switch winner {
		case 0:
			fmt.Println("平手")
		case 1:
			fmt.Println("玩家勝利")
		case 2:
			fmt.Println("AI勝利")
		default:
			fmt.Println("未定義的遊戲結果")
		}
	}

}

func aiSelfPlay() *tictactoe.GameState {
	game := tictactoe.New()
	for !game.GetGameState().IsTerminal {
		if game.CurrentPlayer() == tictactoe.Player1 { // 玩家1行動
			pos := mcts.MonteCarloTreeSearch(game, 1000)
			game.Board[pos] = tictactoe.Player1
			// fmt.Println(game.DrawTable())
			// fmt.Println("玩家1 放置旗子在位置", pos)
		} else { //玩家2行動
			pos := mcts.MonteCarloTreeSearch(game, 1)
			game.Board[pos] = tictactoe.Player2
			// fmt.Println(game.DrawTable())
			// fmt.Println("玩家2 放置旗子在位置", pos)
		}
		//time.Sleep(100 * time.Millisecond)
	}

	return game
}
func playWithAI() *tictactoe.GameState {
	game := tictactoe.New()
	for !game.GetGameState().IsTerminal {
		if game.CurrentPlayer() == tictactoe.Player1 { // 玩家1行動
			pos := getPlayerInput(game) // 自行實現此函數，根據玩家輸入選擇行動
			game.Board[pos] = tictactoe.Player1
			fmt.Println(game.DrawTable())
			fmt.Println("玩家 放置旗子在位置", pos)
		} else { //玩家2行動
			pos := mcts.MonteCarloTreeSearch(game, 1000)
			game.Board[pos] = tictactoe.Player2
			fmt.Println(game.DrawTable())
			fmt.Println("AI 放置旗子在位置", pos)
		}
		//time.Sleep(100 * time.Millisecond)
	}
	return game
}

// 取得玩家輸入
func getPlayerInput(state *tictactoe.GameState) int {
	var playerInput int
	for {
		// 請求玩家輸入
		fmt.Println(state.DrawTable())
		fmt.Println("請輸入你想放置棋子的位置(0-8):")

		_, err := fmt.Scanf("%d", &playerInput)
		if err != nil {
			fmt.Println("輸入有誤，請重新輸入")
			continue
		}
		// 清除換行符
		var newline rune
		fmt.Scanf("%c", &newline)

		// 檢查玩家輸入是否合法
		if playerInput < 0 || playerInput > 8 {
			fmt.Println("輸入範圍有誤，請輸入0-8之間的數字")
			continue
		}

		// 檢查選擇的位置是否已被佔用
		if state.Board[playerInput] != 0 {
			fmt.Println("該位置已被佔用，請選擇其他位置")
			continue
		}
		break
	}
	fmt.Println("放置玩家旗子到位置:", playerInput)
	return playerInput
}
