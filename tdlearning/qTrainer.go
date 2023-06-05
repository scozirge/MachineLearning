package main

import (
	"fmt"
	"math"
	"math/rand"
	"tdlearning/ticTacToe"
	"time"
)

// Q-學習(Q-learning)是強化學習的一種方法。Q-學習就是要記錄下學習過的策略，因而告訴智能體什麼情況下採取什麼行動會有最大的獎勵值
// 「Q」這個字母在強化學習中表示一個動作的期望獎勵
// Q:S(環境狀態)xA(行為)->R(獎勵)

const (
	learningRate         = 0.5    // 學習率 大於,小於等於1 較高的學習率會較快學習新策略，反之agent會比較傾向已經學到的策略
	discountFactor       = 0.7    // 折扣係數 0~1  當discountFactor數值越大時agent更加重視未來獲得的長期獎勵，discountFactor數值越小時，更加短視近利，只在乎目前可獲得的獎勵
	explorationRate      = 1.0    // 探索率(貪婪策略) 也就是agent選擇要探索還是利用的機率 範圍0~1 當探索率越高時，agent會更多嘗試新的行動(探索) 而不僅僅是依賴已知的策略(即從Q表中選擇最佳策略) 0代表不學習了只依賴目前Q表中的最佳策略(利用)
	explorationDecayRate = 0.9993 // 探索綠衰減 每次遊戲結束時 explorationRate會乘上此值來降低下一局的探索率 以便在訓練過程中適應已學習的策略
	trainTimes           = 100000 // 訓練次數(遊戲次數)
	AgentToken           = 1      // 表示代表agent的棋子(0:空格 1:圈圈 2:叉叉)
	PlayerToken          = 2      // 表示代表玩家的棋子(0:空格 1:圈圈 2:叉叉)
	checkWinRateInterval = 100    // 每X局訓練遊戲後報告一次智能體勝率
	learnFromRealPlayer  = false  //是否從跟玩家對戰中繼續學習
	learningMode         = true   //true時為訓練 false時為跟玩家對戰
)

type GameState int //遊戲狀態
const (
	notFinish GameState = iota
	win
	lose
	draw
)

func (gs GameState) String() string {
	names := [...]string{
		"未結束的棋局",
		"贏",
		"輸",
		"平手",
	}

	if gs < notFinish || gs > draw {
		return "unknown"
	}

	return names[gs]
}

var agentWins = 0
var agentLoses = 0

func main() {
	if learningMode {
		TrainAgent()
		return
	}
	rand.Seed(time.Now().UnixNano())
	// 加載訓練好的Q表
	qTable, err := ticTacToe.LoadQTableFromGob("qtable.gob")
	if err != nil {
		fmt.Printf("讀取Q表失敗：%v\n", err)
		return
	}

	// 初始化遊戲狀態
	state := ticTacToe.State{}
	gameFinished, _ := ticTacToe.IsGameFinished(state)

	// 遊戲循環
	for !gameFinished {
		// AI行動
		// 選擇行動
		action := ChooseAction(state, qTable, 0) // 將探索率設為0
		// 執行行動，並獲得新狀態
		state, _ = DoAction(AgentToken, state, action)

		// 檢查遊戲是否結束
		gameFinished, _ = ticTacToe.IsGameFinished(state)
		if gameFinished {
			break
		}

		// 玩家行動
		pAction := getPlayerInput(state) // 自行實現此函數，根據玩家輸入選擇行動
		// 執行行動 並獲得新狀態
		playerDoneState, playerDoneReward := DoAction(PlayerToken, state, pAction)
		//更新Q表
		updateQTable(PlayerToken, qTable, state, playerDoneState, pAction, playerDoneReward)
		state = playerDoneState
		// 檢查遊戲是否結束
		gameFinished, _ = ticTacToe.IsGameFinished(state)
	}
	fmt.Println(state.DrawTable())
	// 輸出遊戲結果
	fmt.Println("遊戲結束！結果:", checkGameState(2, state))
	if learnFromRealPlayer {
		err := ticTacToe.SaveQTableToGob(qTable, "qtable.gob")
		if err != nil {
			fmt.Printf("寫入Q表失敗：%v\n", err)
		} else {
		}

		err2 := ticTacToe.SaveQTableToJson(qTable, "qtable.json")
		if err2 != nil {
			fmt.Printf("寫入Q表失敗：%v\n", err2)
		} else {
		}
	}
}

// 取得玩家輸入
func getPlayerInput(state ticTacToe.State) int {
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
		if state[playerInput] != 0 {
			fmt.Println("該位置已被佔用，請選擇其他位置")
			continue
		}
		break
	}
	fmt.Println("放置玩家旗子到位置:", playerInput)
	return playerInput
}

//訓練Agent
func TrainAgent() {
	rand.Seed(time.Now().UnixNano())
	agentQTable := ticTacToe.InitQTable()      //初始化AgentQ表
	curAgentExplorationRate := explorationRate //目前agent探索率
	curPlayer := 1

	for trainNO := 0; trainNO < trainTimes; trainNO++ {
		// 初始化遊戲狀態
		state := ticTacToe.State{}

		gameFinished, _ := ticTacToe.IsGameFinished(state)
		for !gameFinished { //行動迴圈
			if curPlayer == 1 {
				// agnet行動
				agentDoneState := agentAction(state, agentQTable, curAgentExplorationRate)
				// 設定新狀態為當前狀態
				state = agentDoneState
			} else {
				//玩家行動
				if finished, _ := ticTacToe.IsGameFinished(state); !finished {
					playerDoneState := playerAction(state, agentQTable, curAgentExplorationRate)
					// 設定新狀態為當前狀態
					state = playerDoneState
				}
			}
			gameFinished, _ = ticTacToe.IsGameFinished(state)
			if curPlayer == 1 {
				curPlayer = 2
			} else {
				curPlayer = 1
			}
		}

		gameState := checkGameState(AgentToken, state)
		if gameState == win {
			agentWins++
		} else if gameState == lose {
			agentLoses++
		}

		if trainNO != 0 && (trainNO+1)%checkWinRateInterval == 0 {
			winRate := float64(agentWins) / float64(checkWinRateInterval) * 100
			loseRate := float64(agentLoses) / float64(checkWinRateInterval) * 100
			fmt.Printf("在第%d-%d局訓練遊戲中，agent失敗率為 %.2f%% 勝率為%.2f%%：\n", trainNO-checkWinRateInterval+2, trainNO+1, loseRate, winRate)
			agentWins = 0
			agentLoses = 0
		}
		curAgentExplorationRate *= explorationDecayRate //獎低探索率
	}
	fmt.Println(agentQTable[[9]int{0, 0, 0, 0, 0, 0, 0, 0, 0}])
	fmt.Println("探索率:", curAgentExplorationRate)
	fmt.Println("訓練完成!")

	err := ticTacToe.SaveQTableToGob(agentQTable, "qtable.gob")
	if err != nil {
		fmt.Printf("寫入Q表失敗：%v\n", err)
	} else {
		fmt.Println("寫入Q表成功")
	}

}

func agentAction(state ticTacToe.State, agentQTable ticTacToe.QTable, curAgentExplorationRate float64) ticTacToe.State {
	// 選擇行動
	action := ChooseAction(state, agentQTable, curAgentExplorationRate)
	// 執行行動，並獲得新狀態和獎勵值
	agentDoneState, agentDoneReward := DoAction(AgentToken, state, action)
	// 更新Q表
	updateQTable(AgentToken, agentQTable, state, agentDoneState, action, agentDoneReward)
	return agentDoneState
}
func playerAction(state ticTacToe.State, agentQTable ticTacToe.QTable, curAgentExplorationRate float64) ticTacToe.State {
	pAction := playerChooseRandomAction(state)
	// 執行行動，並獲得新狀態和獎勵值
	playerDoneState, playerDoneReward := DoAction(PlayerToken, state, pAction)
	updateQTable(PlayerToken, agentQTable, state, playerDoneState, pAction, -playerDoneReward)
	return playerDoneState
}

//傳入目前棋況並依據Q表與探索率來行動
func ChooseAction(state ticTacToe.State, qTable ticTacToe.QTable, explorationRate float64) int {

	//從Q表中獲取當前棋況的行動值
	actionValues := qTable[state]
	//隨機值如果小於探索率，則進行探索(隨機選擇一個合法行動)
	if rand.Float64() < explorationRate {
		//找到所有合法行動(未被佔據的位置)
		legalActions := make([]int, 0)
		for i, value := range state {
			if value == 0 {
				legalActions = append(legalActions, i)
			}
		}
		//隨機選擇一個合法行動
		return legalActions[rand.Intn(len(legalActions))]
	}

	//否則，選擇最大Q值的行動(利用)
	myAction := -1
	bestValue := math.Inf(-1)
	for action, value := range actionValues {
		if value > bestValue && state[action] == 0 {
			bestValue = value
			myAction = action
		}
	}
	return myAction
}

// 無策略下棋方法
func playerChooseRandomAction(state ticTacToe.State) int {
	legalActions := make([]int, 0)
	for i, value := range state {
		if value == 0 {
			legalActions = append(legalActions, i)
		}
	}

	return legalActions[rand.Intn(len(legalActions))]
}

// 執行選擇的動作
func DoAction(token int, state ticTacToe.State, action int) (newState ticTacToe.State, reward float64) {
	newState = state //陣列可以這樣 但如果ticTacToe.State是宣告為slice就要用深複製
	//執行行動，將棋子放置在選定的位置
	newState[action] = token

	//檢查遊戲結果

	result := checkGameState(token, newState)
	//根據遊戲結果設定獎勵值
	switch result {
	case notFinish: //遊戲尚未結束
		reward = 0
	case win: //勝利獎勵1
		reward = 1
	case draw: //平局獎勵
		reward = 0
	}

	return newState, reward
}

// 依據棋況返回目前遊戲狀態GameState
func checkGameState(token int, state ticTacToe.State) GameState {

	// 棋局尚未結束判斷
	isGameFinished, winningToken := ticTacToe.IsGameFinished(state)
	if !isGameFinished {
		return notFinish
	}
	// 有任一方贏了
	if winningToken == 0 {
		return draw
	} else if winningToken == token {
		return win
	} else {
		return lose
	}
}

// 更新Q表
func updateQTable(token int, qTable ticTacToe.QTable, state, nextState ticTacToe.State, action int, reward float64) {

	//時序差分學習(Temporal-Difference Learning，簡稱TD Learning)
	updatedValue := qTable[state][action] + learningRate*(reward+discountFactor*maxQ(nextState, qTable)-qTable[state][action])
	qTable[state][action] = updatedValue
	// fmt.Println("/////////////////////////////////////////////")
	// fmt.Println("updatedValue=", updatedValue)
	// fmt.Println(state.DrawTable())
	// fmt.Println(qTable[state])
	// fmt.Println("result=", checkGameState(1, nextState))
}

// 依照Q表中取得目前棋況最高價值的行動價值
func maxQ(state ticTacToe.State, qTable ticTacToe.QTable) float64 {
	actionQ, ok := qTable[state]
	if !ok || len(actionQ) == 0 { //state在填滿的狀態下不會有下一步的Q值資料，此時返回0
		return 0
	}

	maxQ := math.Inf(-1)
	for _, q := range qTable[state] {
		if q > maxQ {
			maxQ = q
		}
	}
	return maxQ
}
