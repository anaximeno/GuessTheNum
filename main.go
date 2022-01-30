package main

import (
	"fmt"
	"math/rand"
	"time"
)


const CORRECT_GUESS_DEFAULT_OUTPUT_FORMAT string = `
                     Correct Guess
                   Next Level %d -> %d

     Write 'q' to quit or click 'enter' to continue

==> `

const NO_MORE_TENTATIVES_DEFAULT_OUTPUT_FORMAT string = `
                   NO MORE TENTATIVES
                The correct number was %d

    Write 'q' to quit or click 'enter' to play again.

==> `

const MAIN_SECTION_DEFAULT_OUTPUT_FORMAT string = `
 Level %d: Guess a number between %d and %d
 %d more tentatives

 Hint: %s
 ----------------------------
 Guess => `

const WIN_GAME_OUTPUT string = `
                   YOU WON THE GAME
            Now you're the guessing master
                !!! CONGRATULATIONS !!!
`


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
    playerTries int
    lastPlayerGuess int
}


func clear() {
    fmt.Print("\033[H\033[2J")
}


func enterpoint() bool {
    var input string
    fmt.Scanln(&input)
    if input == "q" {
        return false
    } else {
        return true
    }
}


func main() {
    rand.Seed(time.Now().UnixNano())
    guessingGame := GameState{}

    /* Levels of the game */
    guessingGame.addLevel(0, 50, 5)
    guessingGame.addLevel(0, 70, 5)
    guessingGame.addLevel(0, 100, 5)
    guessingGame.addLevel(0, 120, 5)
    guessingGame.addLevel(0, 150, 6)
    guessingGame.addLevel(0, 200, 6)
    guessingGame.addLevel(0, 500, 7)
    guessingGame.addLevel(0, 1000, 7)
    guessingGame.addLevel(0, 2500, 8)
    guessingGame.addLevel(0, 10000, 12)
    guessingGame.addLevel(2, 13, 1)

    guessingGame.init()

    for {
        guessingGame.run()
        if guessingGame.checkGuess() {
            fmt.Print(WIN_GAME_OUTPUT)
            break
        } else {
            fmt.Printf(
                NO_MORE_TENTATIVES_DEFAULT_OUTPUT_FORMAT,
                guessingGame.currentLevel.numberToGuess,
            )
            if enterpoint() {
                guessingGame.reboot()
            } else {
                break
            }
        }
    }
}


func (game *GameState) init() {
    game.lastPlayerGuess = -1
    game.playerTries = game.currentLevel.nTries
}


func (game *GameState) reboot() {
    var cur *GameLevel
    for cur = game.currentLevel ; cur.prev != nil ; cur = cur.prev {
        cur.changeGuessingNumber()
    }
    cur.changeGuessingNumber()
    game.currentLevel = cur
    game.init()
}


func (game *GameState) run() {
    clear()
    for game.consumeTry() {
        fmt.Printf(
            MAIN_SECTION_DEFAULT_OUTPUT_FORMAT,
            game.currentLevel.levelNum,
            game.currentLevel.minRange,
            game.currentLevel.maxRange,
            game.playerTries + 1,
            game.giveHint(),
        )

        fmt.Scan(&game.lastPlayerGuess)
        clear()

        if game.checkGuess() {
            if game.transitLevel() {
                fmt.Printf(
                    CORRECT_GUESS_DEFAULT_OUTPUT_FORMAT,
                    game.currentLevel.prev.levelNum,
                    game.currentLevel.levelNum,
                )
                if !enterpoint() {
                    break
                } else {
                    clear()
                }
            } else {
                break
            }
        }
    }
}


func (gameLevel *GameLevel) changeGuessingNumber() {
    gameLevel.numberToGuess = 1 + rand.Intn(gameLevel.maxRange - 1) + gameLevel.minRange
}


func (game *GameState) transitLevel() bool {
    if game.currentLevel != nil && game.currentLevel.next != nil {
        game.currentLevel = game.currentLevel.next
        game.playerTries = game.currentLevel.nTries
        game.lastPlayerGuess = -1
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
            levelNum:       1,
            nTries:         nTries,
            minRange:       minRange,
            maxRange:       maxRange,
            prev:           nil,
            next:           nil,
        }
        game.currentLevel.changeGuessingNumber()
    } else {
        var cur *GameLevel
        for cur = game.currentLevel ; cur.next != nil ; cur = cur.next {}
        cur.next = &GameLevel{
            levelNum:       cur.levelNum + 1,
            nTries:         nTries,
            minRange:       minRange,
            maxRange:       maxRange,
            prev:           cur,
            next:           nil,
        }
        cur.next.changeGuessingNumber()
    }
}


func (game *GameState) consumeTry() bool {
    if game.playerTries > 0 {
        game.playerTries -= 1
        return true
    } else {
        return false
    }
}


func (game GameState) giveHint() string {
    if game.lastPlayerGuess == -1 {
        return "no hints yet"
    } else if game.currentLevel.numberToGuess > game.lastPlayerGuess {
        return fmt.Sprintf("greater than %d", game.lastPlayerGuess)
    } else {
        return fmt.Sprintf("lower than %d", game.lastPlayerGuess)
    }
}


func (game GameState) checkGuess() bool {
    return game.currentLevel.numberToGuess == game.lastPlayerGuess
}
