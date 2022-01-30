package main

import (
	"fmt"
	"math/rand"
	"time"
)


const GAME_MAX_LEVEL int = 100


const CORRECT_GUESS_DEFAULT_OUTPUT_FORMAT string = `
                     Correct Guess
                  Next Level %d -> %d

     Write 'q' to quit or click 'enter' to continue

==> `

const NO_MORE_TENTATIVES_DEFAULT_OUTPUT_FORMAT string = `
                   NO MORE TENTATIVES
                The correct guess was %d

    Write 'q' to quit or click 'enter' to play again.

==> `

const MAIN_SECTION_DEFAULT_OUTPUT_FORMAT string = `
 Level %d: A number between %d and %d
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
    playerTries int
    totalLevels int
    lastGuessTentative int
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
    guessingGame.addLevel(0, 50, 4)
    guessingGame.addLevel(0, 70, 4)
    guessingGame.addLevel(0, 100, 4)
    guessingGame.addLevel(0, 120, 4)
    guessingGame.addLevel(0, 150, 5)
    guessingGame.addLevel(0, 200, 6)
    guessingGame.addLevel(0, 500, 7)

    guessingGame.init()

    for true {
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
    game.lastGuessTentative = -1
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

        fmt.Scan(&game.lastGuessTentative)
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
    if game.lastGuessTentative == -1 {
        return "no hints yet"
    } else if game.currentLevel.numberToGuess > game.lastGuessTentative {
        return fmt.Sprintf("greater than %d", game.lastGuessTentative)
    } else {
        return fmt.Sprintf("lower than %d", game.lastGuessTentative)
    }
}


func (game GameState) checkGuess() bool {
    return game.currentLevel.numberToGuess == game.lastGuessTentative
}
