package main

import (
    "fmt"
    "math/rand"
    "time"
)


const GAME_MAX_LEVEL int = 100


const MENU_DEFAULT_OUTPUT_FORMAT string = `
 Chose one option:
    1. Play the game
    2. View the game history
    0. Exit
 => `


const CORRECT_GUESS_DEFAULT_OUTPUT_FORMAT string = `
               Correct Guess
             Next Level %d -> %d

  [Click 'q' to quit or another to continue]
`

const MAIN_SECTION_DEFAULT_OUTPUT_FORMAT string = `
 Level %d: A number between %d to %d
 %d+ tentatives

 Hint: %s
 ----------------------------
 Guess =>

`

/* Game Level is a double linked list of levels */
type GameLevel struct {
    levelNum int
    nTries int
    minRange int
    maxRange int
    numberToGuess int
    next *GameLevel
    prev *GameLevel
}


type GameState struct {
    currentLevel *GameLevel
    totalLevels int
    lastGuessTentative int
}


func (gameLevel *GameLevel) changeGuessingNumber() {
    gameLevel.numberToGuess = rand.Intn(gameLevel.maxRange) + gameLevel.minRange;
}


func (game *GameState) transitLevel() bool {
    if game.currentLevel != nil && game.currentLevel.next != nil {
        game.currentLevel = game.currentLevel.next
        game.lastGuessTentative = -1
        return true
    } else {
        return false
    }
}


func (game *GameState) addLevel(minRange, maxRange, nTries int) {
    defer func () {
        game.totalLevels += 1
    }()

    if game.currentLevel == nil {
        game.currentLevel = &GameLevel {
            levelNum:   1,
            nTries:     nTries,
            minRange:   minRange,
            maxRange:   maxRange,
            prev:       nil,
            next:       nil,
        }
        game.currentLevel.changeGuessingNumber()
    } else {
        var cur *GameLevel
        for cur = game.currentLevel ; cur.next != nil ; cur = cur.next {}
        cur.next = &GameLevel{
            levelNum:   cur.levelNum + 1,
            nTries:     nTries,
            minRange:   minRange,
            maxRange:   maxRange,
            prev:       cur,
            next:       nil,
        }
        cur.next.changeGuessingNumber()
    }
}


func (game *GameState) consumeTry() bool {
    if game.currentLevel.nTries > 0 {
        game.currentLevel.nTries -= 1
        return true
    } else {
        return false
    }
}


func (game GameState) giveHint() string {
    if game.lastGuessTentative == -1 {
        return "No hints"
    } else if game.currentLevel.numberToGuess > game.lastGuessTentative {
        return fmt.Sprintf("greater than %d", game.lastGuessTentative)
    } else {
        return fmt.Sprintf("lower than %d", game.lastGuessTentative)
    }
}


func (game GameState) checkGuess(guess int) bool {
    return game.currentLevel.numberToGuess == guess
}


func main() {
    rand.Seed(time.Now().UnixNano())
    guessingGame := GameState{lastGuessTentative: -1}
    guessingGame.addLevel(0, 50, 3)
    guessingGame.addLevel(0, 100, 5)

    fmt.Printf(
        MAIN_SECTION_DEFAULT_OUTPUT_FORMAT,
        guessingGame.currentLevel.levelNum,
        guessingGame.currentLevel.minRange,
        guessingGame.currentLevel.maxRange,
        guessingGame.currentLevel.nTries,
        guessingGame.giveHint(),
    )

}
