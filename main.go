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
    levelID int
    nTries int
    minRange int
    maxRange int
    numberToGuess int
    next *GameLevel
    prev *GameLevel
}


type GameState struct {
    level *GameLevel
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
                guessingGame.level.numberToGuess,
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
    game.playerTries = game.level.nTries
}


func (game *GameState) reboot() {
    var level *GameLevel
    for level = game.level ; level.prev != nil ; level = level.prev {
        level.changeGuessingNumber()
    }
    game.level = level
    game.level.changeGuessingNumber()
    game.init()
}


func (game *GameState) run() {
    clear()
    for game.consumeTry() {
        fmt.Printf(
            MAIN_SECTION_DEFAULT_OUTPUT_FORMAT,
            game.level.levelID,
            game.level.minRange,
            game.level.maxRange,
            game.playerTries + 1,
            game.giveHint(),
        )

        fmt.Scan(&game.lastPlayerGuess)
        clear()

        if game.checkGuess() {
            if game.transitLevel() {
                fmt.Printf(
                    CORRECT_GUESS_DEFAULT_OUTPUT_FORMAT,
                    game.level.prev.levelID,
                    game.level.levelID,
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
    if game.level != nil && game.level.next != nil {
        game.level = game.level.next
        game.playerTries = game.level.nTries
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

    if game.level == nil {
        game.level = &GameLevel{
            levelID:       1,
            nTries:         nTries,
            minRange:       minRange,
            maxRange:       maxRange,
            prev:           nil,
            next:           nil,
        }
        game.level.changeGuessingNumber()
    } else {
        var level *GameLevel
        for level = game.level ; level.next != nil ; level = level.next {}
        level.next = &GameLevel{
            levelID:       level.levelID + 1,
            nTries:         nTries,
            minRange:       minRange,
            maxRange:       maxRange,
            prev:           level,
            next:           nil,
        }
        level.next.changeGuessingNumber()
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
    if lastGuess := game.lastPlayerGuess; lastGuess == -1 {
        return "no hints yet"
    } else if game.level.numberToGuess > game.lastPlayerGuess {
        return fmt.Sprintf("greater than %d", game.lastPlayerGuess)
    } else {
        return fmt.Sprintf("lower than %d", game.lastPlayerGuess)
    }
}


func (game GameState) checkGuess() bool {
    return game.level.numberToGuess == game.lastPlayerGuess
}
