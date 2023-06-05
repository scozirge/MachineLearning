package ticTacToe

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// 存每個行動的Q值
type ActionQ map[int]float64

// 存每個狀態和對應的行動價值
type QTable map[State]ActionQ

//輸出json時轉換換用類型
type ExportableQTable struct {
	States map[string]map[string]float64
}

//取絕對值函式 註:很意外math庫的Abs沒有傳入int的函式
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// 窮舉並產生所有可能的棋局(放滿棋盤才算一個棋局)
func generateStates(state [9]int, index int, states *[]State) {

	// 如果index是棋局中的第幾步，index==9，代表完成一盤棋了
	if index == 9 {
		oCount := Count(1, state)    //計算O的數量
		xCount := Count(2, state)    //計算X的數量
		if abs(oCount-xCount) <= 1 { //O和X的數量差不大於1才合法，才加到紀錄中
			*states = append(*states, state)
		}
		return
	}

	for i := 0; i < 3; i++ { // 空格可能的內容(0:空格 1:圈圈 2:叉叉)
		state[index] = i
		generateStates(state, index+1, states)
	}
}

// 初始化Q表
func InitQTable() QTable {
	qTable := make(QTable) //每個狀態對應的行動價值mpa

	var state [9]int
	states := []State{} //棋盤狀態slice用[9]int來標示每一格填入的是O還是X
	// 生成所有可能的棋盤狀態
	generateStates(state, 0, &states)
	for _, state := range states {
		actions := make(ActionQ)                //actions是map[int]float64用來存每個行動的Q值
		for action := 0; action < 9; action++ { //OOXX有9個格子所以每盤棋會有9次行動
			// 如果該位置為空格，則行動價值初始化為0
			if state[action] == 0 {
				actions[action] = 0.0
			}
		}
		// 將當前狀態的行動價值添加到Q表中
		qTable[state] = actions
	}

	return qTable
}

// 寫入Q表到本地(gob格式)
func SaveQTableToGob(qTable QTable, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(qTable)
	if err != nil {
		return err
	}

	return nil
}

// 從本地讀取Q表(gob格式)
func LoadQTableFromGob(filename string) (QTable, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var qTable QTable
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&qTable)
	if err != nil {
		return nil, err
	}

	return qTable, nil
}

// 寫入Q表到本地(json格式)
func SaveQTableToJson(qTable QTable, filename string) error {
	exportableQTable := ConvertQTableToExportable(qTable)

	jsonData, err := json.Marshal(exportableQTable)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}

// 從本地讀取Q表(json格式)並轉換回QTable
func LoadQTableFromJson(filename string) (QTable, error) {
	jsonData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var exportableQTable ExportableQTable
	err = json.Unmarshal(jsonData, &exportableQTable)
	if err != nil {
		return nil, err
	}

	qTable := ConvertExportableToQTable(exportableQTable)

	return qTable, nil
}

//將Q表輸出json時轉換成ExportableQTable類型
func ConvertQTableToExportable(qTable QTable) ExportableQTable {
	exportableQTable := ExportableQTable{
		States: make(map[string]map[string]float64),
	}

	for state, actionQ := range qTable {
		stateStr := state.ToKeyString()
		exportableQTable.States[stateStr] = make(map[string]float64)

		for action, qValue := range actionQ {
			actionStr := fmt.Sprintf("%d", action)
			exportableQTable.States[stateStr][actionStr] = qValue
		}
	}

	return exportableQTable
}

// 將ExportableQTable類型轉換回QTable
func ConvertExportableToQTable(exportableQTable ExportableQTable) QTable {
	qTable := make(QTable)

	for stateStr, actionQ := range exportableQTable.States {
		state := StateFromKeyString(stateStr)
		actions := make(ActionQ)

		for actionStr, qValue := range actionQ {
			action := 0
			fmt.Sscanf(actionStr, "%d", &action)
			actions[action] = qValue
		}

		qTable[state] = actions
	}

	return qTable
}

// 從狀態字串中建立State
func StateFromKeyString(stateStr string) State {
	var state State
	fmt.Sscanf(stateStr, "%d%d%d%d%d%d%d%d%d",
		&state[0], &state[1], &state[2],
		&state[3], &state[4], &state[5],
		&state[6], &state[7], &state[8],
	)
	return state
}
