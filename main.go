package main

import (
	"strconv"
	"syscall/js"
)

var signal = make(chan int)
var board = [][]uint8{
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0},
}
var (
	mouseX       = 0
	canClick     = true
	mouseY       = 0
	player1Score = 0
	player2Score = 0
	whosTurn     = 1
	window       = js.Global()
	doc          = window.Get("document")
	canvas       = doc.Call("getElementById", "canvas")
	context      = canvas.Call("getContext", "2d")
)

func keepAlive() {
	for {
		<-signal
	}
}

func draw() {

	for i := 1; i < 3; i++ {

		for j := 1; j < 3; j++ {
			context.Call("moveTo", j*100, 0)
			context.Call("lineTo", j*100, 300)
		}
		context.Call("moveTo", 0, i*50)
		context.Call("lineTo", 300, i*50)
	}
	context.Call("stroke")
}

func clicked(this js.Value, inputs []js.Value) interface{} {
	if canClick == true {
		mouseX = inputs[0].Get("clientX").Int()
		mouseY = inputs[0].Get("clientY").Int()
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				control1 := 510 + (j * 100)
				control2 := 510 + ((j + 1) * 100)
				control3 := 225 + (i * 100)
				control4 := 225 + ((i + 1) * 100)

				if mouseX >= control1 && mouseX < control2 && mouseY >= control3 && mouseY < control4 {
					if board[i][j] == 0 {
						if whosTurn == 1 {
							whosTurn = 2
						} else {
							whosTurn = 1
						}

						putImage(j, i)
						changePlayer()
						switch winControl() {
						case 1:
							player1Score++
							changeScore()
							nextButton := doc.Call("getElementById", "next")
							nextButton.Set("disabled", false)
							canClick = false
							changeInfo("PLAYER 1 WINS THE ROUND!PLEASE CLICK NEXT FOR CONTINUE")

						case 2:
							player2Score++
							changeScore()
							nextButton := doc.Call("getElementById", "next")
							nextButton.Set("disabled", false)
							canClick = false
							changeInfo("PLAYER 2 WINS THE ROUND!PLEASE CLICK NEXT FOR CONTINUE")

						case 3:
							nextButton := doc.Call("getElementById", "next")
							nextButton.Set("disabled", false)
							canClick = false
							changeInfo("DRAW!PLEASE CLICK NEXT FOR CONTINUE")

						}

					}
					return nil
				}
			}
		}
	}

	return nil

}

func changePlayer() {
	turn := doc.Call("getElementById", "turn")
	player := "Turn: PLAYER " + strconv.Itoa(whosTurn)
	turn.Set("innerHTML", player)
}
func changeScore() {
	player1 := doc.Call("getElementById", "player1")
	player2 := doc.Call("getElementById", "player2")
	score1 := "PLAYER 1: " + strconv.Itoa(player1Score)
	score2 := "PLAYER 2: " + strconv.Itoa(player2Score)
	player1.Set("innerHTML", score1)
	player2.Set("innerHTML", score2)
}
func putImage(x int, y int) {
	if whosTurn == 2 {
		board[y][x] = 2
		img := doc.Call("getElementById", "imageX")
		context.Call("drawImage", img, 1+(x*100), 1+(y*50), 99, 49)

	} else {
		board[y][x] = 1
		img := doc.Call("getElementById", "imageO")
		context.Call("drawImage", img, 1+(x*100), 1+(y*50), 99, 49)

	}
}
func winControl() int {
	var check0 = 0
	var check1 = 0
	var check2 = 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == 0 {
				check1 = 0
				check2 = 0
				check0++
				break
			} else if board[i][j] == 1 {
				check1++
				if check1 < 3 {
					continue
				} else {
					check1 = 0
					return 1
				}
			} else if board[i][j] == 2 {
				check2++
				if check2 < 3 {
					continue
				} else {
					check2 = 0
					return 2
				}
			}

		}
		check1 = 0
		check2 = 0
		for j := 0; j < 3; j++ {
			if board[j][i] == 0 {
				check1 = 0
				check2 = 0
				check0++
				break
			} else if board[j][i] == 1 {
				check1++
				if check1 < 3 {
					continue
				} else {
					check1 = 0
					return 1
				}
			} else if board[j][i] == 2 {
				check2++
				if check2 < 3 {
					continue
				} else {
					check2 = 0
					return 2
				}
			}

		}
		check1 = 0
		check2 = 0
		if board[0][0] == 1 && board[1][1] == 1 && board[2][2] == 1 {
			return 1
		}
		if board[0][0] == 2 && board[1][1] == 2 && board[2][2] == 2 {
			return 2
		}
		if board[0][2] == 1 && board[1][1] == 1 && board[2][0] == 1 {
			return 1
		}
		if board[0][2] == 2 && board[1][1] == 2 && board[2][0] == 2 {
			return 2
		}
	}
	if check0 == 0 {
		return 3
	} else {
		return 0
	}
}
func resetBoard() {
	canClick = true
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			board[i][j] = 0
		}
	}
	context.Call("clearRect", 0, 0, 300, 300)
	draw()
}

func next(this js.Value, inputs []js.Value) interface{} {
	resetBoard()
	nextButton := doc.Call("getElementById", "next")
	nextButton.Set("disabled", true)
	changeInfo("START!")
	return nil
}
func changeInfo(info string) {
	information := doc.Call("getElementById", "info")
	information.Set("innerHTML", info)
}

func reset(this js.Value, inputs []js.Value) interface{} {
	whosTurn = 1
	player1Score = 0
	player2Score = 0
	resetBoard()
	changePlayer()
	changeScore()
	changeInfo("START!")
	return nil
}

func registerCallbacks() {

	js.Global().Set("reset", js.FuncOf(reset))
	js.Global().Set("next", js.FuncOf(next))
}

func main() {
	doc.Call("addEventListener", "click", js.FuncOf(clicked))
	draw()
	registerCallbacks()
	keepAlive()
}
