package mcts

// MCTS主要有四個階段：選擇(Selection)、擴展(Expansion)、模擬(Rollout)和反向傳播(Backpropagation)

import (
	"math"
	"math/rand"
	"time"

	tictactoe "mcts/tictactoe"
)

// 節點資料(目前遊戲狀態、父節點、子節點、勝利次數、訪問次數和未探索位置)
type TreeNode struct {
	state          *tictactoe.GameState
	parent         *TreeNode
	children       []*TreeNode
	wins           float64
	visits         float64
	unexploredPosz []int
}

// 傳入目前狀態、玩家、迭代次數取得最佳動作
func MonteCarloTreeSearch(game *tictactoe.GameState, iterations int) int {
	rand.Seed(time.Now().UnixNano())

	root := &TreeNode{
		state:          game,
		unexploredPosz: game.GetLegalPosz(),
	}

	for i := 0; i < iterations; i++ {
		node := root.selectNode()
		winner, finalNode := node.rollout()
		//fmt.Println("開始反向傳播:", node.state)
		finalNode.backpropagation(winner)
	}

	bestChild := root.bestChild()
	return bestChild.state.LastPlaced
}

// 選擇(Selection)-選擇最佳UTC值得節點
func (t *TreeNode) selectNode() *TreeNode {
	// 目前節點如果是第一次訪問 或是 還有未探索的位置 或是沒有任何子節點 返回目前節點
	if t.visits == 0 || len(t.unexploredPosz) > 0 || len(t.children) == 0 {
		return t
	}
	var bestChild *TreeNode
	maxUCT := math.Inf(-1)
	for _, child := range t.children {

		// UCT計算公式為：UCT = (w / n) + C * sqrt(ln(N) / n)
		// 前半段的(w / n)代表利用 後半段的C * sqrt(ln(N) / n)代表探索  C探索率，通常是sqrt(2)，可以根據類型調整。C越大會使AI更注重探索，C越小會使AI注重利用
		uct := 0.0
		if child.visits != 0 {
			uct = child.wins/child.visits + math.Sqrt(2*math.Log(t.visits)/child.visits)
			//fmt.Println("uct=", uct)
		}

		if uct > maxUCT {
			maxUCT = uct
			bestChild = child
		}
	}
	return bestChild.selectNode()
}

// 擴展(Expansion)-優先探索尚未探索的位置，如果都探索了就跑selectNode
func (t *TreeNode) expand() *TreeNode {
	pos := 0
	// 如果此節點已探索完成(len(t.unexploredPosz)==0)
	if len(t.unexploredPosz) == 0 {
		return t.selectNode()
	} else {
		// 隨機選擇一個未探索過的位置
		posIndex := rand.Intn(len(t.unexploredPosz))
		pos = t.unexploredPosz[posIndex]
		// 移除選中的位置
		t.unexploredPosz = append(t.unexploredPosz[:posIndex], t.unexploredPosz[posIndex+1:]...)
	}

	// 設定新狀態為目前狀態並放置棋子
	newState := *t.state
	newState.Board[pos] = newState.CurrentPlayer()
	newState.LastPlaced = pos
	// 建立子節點
	child := &TreeNode{
		state:          &newState,
		parent:         t,
		unexploredPosz: newState.GetLegalPosz(),
	}
	t.children = append(t.children, child)

	return child
}

// 模擬(Rollout)-進行一次隨機模擬，並返回模擬結果中的贏家
func (t *TreeNode) rollout() (int, *TreeNode) {
	if t.state.GetGameState().IsTerminal {
		return t.state.GetGameState().Winner, t
	}

	return t.expand().rollout()
}

// 反向傳播(Backpropagation)-每次模擬(Rollout)結束時，會根據模擬結果更新從根節點到該模擬結束的節點之間的所有節點資料
func (t *TreeNode) backpropagation(winner int) {
	t.visits++
	if winner == 0 {
		t.wins += 0.1
	} else if t.state.CurrentPlayer() != winner { // 贏家不等於目前玩家就勝利次數就+1 代表上個行動的玩家最後是贏棋的
		t.wins++
	}
	if t.parent != nil {
		t.parent.backpropagation(winner)
	}
}

// 找出目前節點中訪問次數最高的子節點並返回
func (t *TreeNode) bestChild() *TreeNode {
	var bestChild *TreeNode
	maxVisits := math.Inf(-1)

	for _, child := range t.children {
		if child.visits > maxVisits {
			maxVisits = child.visits
			bestChild = child
		}
	}
	return bestChild
}
