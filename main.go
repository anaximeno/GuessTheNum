package main

import (
    "fmt"
)

const GAME_MAX_LEVEL = 100

const MENU_DEFAULT_OUTPUT_FORMAT = `
 Chose one option:
    1. Play the game
    2. View the game history
    0. Exit
 => `

const CORRECT_GUESS_DEFAULT_OUTPUT_FORMAT = `
               Correct Guess
             Next Level %d -> %d

  [Click 'q' to quit or another to continue]
`

const MAIN_SECTION_DEFAULT_OUTPUT_FORMAT = `
 Level %d: Range from %d to %d
         %2d more tentatives
 -----------------------------
 Hint: %s
 Guess => `

func main() {
    fmt.Println(MENU_DEFAULT_OUTPUT_FORMAT)
    fmt.Printf(MAIN_SECTION_DEFAULT_OUTPUT_FORMAT, 1, 0, 50, 3, "lower")
    fmt.Printf(CORRECT_GUESS_DEFAULT_OUTPUT_FORMAT, 1, 2)
}

