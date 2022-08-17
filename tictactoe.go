package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var emptyBoard = []string{
	"   |   |   ", 
	"---+---+---", 
	"   |   |   ", 
	"---+---+---", 
	"   |   |   ",
}

func main() {
	startScreen()
}

type Coordinate struct {
	x, y interface{}
}

func gameMode() {
	fmt.Println("Starting a new game...")
	boardPieces := map[Coordinate]string{
		Coordinate{0, 1}: "1",
		Coordinate{0, 5}: "2",
		Coordinate{0, 9}: "3",
		Coordinate{2, 1}: "4",
		Coordinate{2, 5}: "5",
		Coordinate{2, 9}: "6",
		Coordinate{4, 1}: "7",
		Coordinate{4, 5}: "8",
		Coordinate{4, 9}: "9",
	}
	scanner := bufio.NewScanner(os.Stdin)
	playerTurn := 0
	var player int
	for !checkGameOver(boardPieces) {
		playerMark := false
		for !playerMark {
			playerTurn += 1
			if (playerTurn + 2) % 2 == 0 {
				player = 2
			} else {
				player = 1
			}
			var input string
			printBoard(boardPieces)
			enterInputMessage := "Enter position to mark a space: "
			fmt.Printf("Player %d, %s\n", player, enterInputMessage)
			scanner.Scan()
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "reading standard input: ", err)
			}
			input = scanner.Text()
			fmt.Printf("Player %d entered %q\n", player, input)
			var mark string
			if playerTurn % 2 + 1 == 1 {
				mark = "O"
			} else {
				mark = "X"
			}

			var x, y int
			if input == "1" {
				x, y = 0, 1
			} else if input == "2" {
				x, y = 0, 5
			} else if input == "3" {
				x, y= 0, 9
			} else if input == "4" {
				x, y = 2, 1
			} else if input == "5" {
				x, y = 2, 5
			} else if input == "6" {
				x, y = 2, 9
			} else if input == "7" {
				x, y = 4, 1
			} else if input == "8" {
				x, y = 4, 5
			} else if input == "9" {
				x, y = 4, 9
			} else {
				fmt.Println("[ERROR]: Incorrect position selected. Try again.")
				continue
			}
			if !isMarkedAlready(boardPieces[Coordinate{x, y}]) {
				boardPieces[Coordinate{x, y}] = mark
				playerMark = true
			} else {
				fmt.Println("[ERROR]: Incorrect position selected. Try again.")
			}
		}
	}
	if checkWin(boardPieces) {
		fmt.Printf("Game Over! Player %d wins!\n", player)
	} else {
		fmt.Printf("Game Over! Draw game!")
	}
	startScreen()
}

func isMarkedAlready(s string) bool {
	return s == "X" || s == "O"
}

func printBoard(boardPieces map[Coordinate]string) {
	for row, line := range emptyBoard {
		var markedLine string = emptyBoard[row]
		for col:= range line {
			coord := Coordinate{row, col}
			if mark, found := boardPieces[coord]; found {
				markedLine = replaceAtIndex(markedLine, []rune(mark)[0], col)
			}
		}
		fmt.Println(markedLine)
	}
}

func replaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}

func checkWin(boardPieces map[Coordinate]string) bool {
	winningCombinations := [][][3]int{
		// Horizontal
		{{0, 0, 0}, {1, 5, 9}},
		{{2, 2, 2}, {1, 5, 9}},
		{{4, 4, 4}, {1, 5, 9}},

		// Vertical
		{{0, 2, 4}, {1, 1, 1}},
		{{0, 2, 4}, {5, 5, 5}},
		{{0, 2, 4}, {9, 9, 9}},

		// Diagonal
		{{0, 2, 4}, {1, 5, 9}},
		{{0, 2, 4}, {9, 5, 1}},
	} 
	for _, v := range winningCombinations {
		if isThreeInARow(v[0], v[1], boardPieces) {
			return true
		}
	}
	
	return false
}

func checkGameOver(boardPieces map[Coordinate]string) bool {
	return checkWin(boardPieces) || checkDraw(boardPieces)
}

func checkDraw(boardPieces map[Coordinate]string) bool {
	var boardFilled bool = true
	for _, v := range boardPieces {
		if _, err := strconv.Atoi(v); err == nil {
			boardFilled = false
			break
		}
	}
	if boardFilled {
		return !checkWin(boardPieces)
	}
	
	return false
}

func isThreeInARow(xCoords [3]int, yCoords[3]int, boardPieces map[Coordinate]string) bool {
	if mark1, found := boardPieces[Coordinate{xCoords[0], yCoords[0]}]; found {
		if mark2, found := boardPieces[Coordinate{xCoords[1],yCoords[1]}]; found {
			if mark3, found := boardPieces[Coordinate{xCoords[2],yCoords[2]}]; found {
				if mark1 == mark2 && mark1 == mark3 {
					return true
				}
			}
		}
	}
	return false
}

func startScreen() {
	gameTitle := "---------------------\nGame:\tREVERSE TIC-TAC-TOE\n---------------------\n"
	gameInstructions := "Rules:\tTwo players take turns marking the spaces in a three-by-three grid with X or O.\n" +
		"\tThe first player to get three of their marks in a horizontal, vertical, or diagonal row is the LOSER. The other player is the WINNER.\n---------------------\n"
	fmt.Println(gameTitle + gameInstructions)

	startScreenPrompt := "Commands:\n\t1: Start a new game.\n\t2: Exit.\n"
	startScreenCommands := map[string]string{
		"1": "Start a new game.",
		"2": "Exit.",
	}
	startScreenOutputs := map[string]func(){
		"1": gameMode,
		"2": exitGame,
	}
	commandPrompt(startScreenPrompt, startScreenCommands, startScreenOutputs)
}

func commandPrompt(prompt string, commands map[string]string, outputs map[string]func()) {
	// Prompt user to input a command
	fmt.Println(prompt)
	fmt.Print("Enter a command: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	var input string
	input = scanner.Text()

	// Repeat input prompt until a correct command is inputted
	for commands[input] == "" {
		fmt.Println("Error: Invalid command.")
		fmt.Println(prompt)
		fmt.Print("Enter a command: ")
		scanner.Scan()
		input = scanner.Text()
	}

	// Run the function for the given input
	outputs[input]()

}

func exitGame() {
	fmt.Println("Exiting game...")
	os.Exit(0)
}
