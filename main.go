package main

import (
	"fmt"
	"math/rand"
	"time"
)


const CORRECT_GUESS_DEFAULT_OUTPUT_FORMAT string = `
                  <<< CORRECT GUESS >>>
           Continue to the next level %d -> %d ?

     Write 'q' to quit or click 'enter' to continue

==> `

const NO_MORE_TENTATIVES_DEFAULT_OUTPUT_FORMAT string = `
                    -> 0 MORE ATTEMPTS <-
         You losed the game! The correct guess was %d

      Write 'q' to quit or click 'enter' to play again.

==> `

const MAIN_SECTION_DEFAULT_OUTPUT_FORMAT string = `
 Level %d: Guess a number between %d and %d
 %d more attempts

 Hint: %s
 ----------------------------
 Guess => `

const WIN_GAME_OUTPUT string = `
           CONGRATULATIONS! YOU WON THE GAME.
       Now you are "officially" the guessing master.

`

const NO_GUESS int = -1

type GameLevel struct {
    id int
    nTries int
    minRange int
    maxRange int
    numberToGuess int
    next *GameLevel
    prev *GameLevel
}


type Player struct {
    triesLeft int
    lastGuess int
}


type GameState struct {
    nLevels int
    level *GameLevel
    player Player
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
    guessingGame := GameState{ player: Player {} }

    /* Levels of the game */
    guessingGame.addLevel(0, 2, 2)
    guessingGame.addLevel(0, 2, 2)
    // guessingGame.addLevel(0, 70, 5)
    // guessingGame.addLevel(0, 100, 5)
    // guessingGame.addLevel(0, 120, 5)
    // guessingGame.addLevel(0, 150, 6)
    // guessingGame.addLevel(0, 200, 6)
    // guessingGame.addLevel(0, 500, 7)
    // guessingGame.addLevel(0, 1000, 7)
    // guessingGame.addLevel(0, 2500, 8)
    // guessingGame.addLevel(0, 10000, 12)
    // guessingGame.addLevel(2, 13, 1)

    guessingGame.initPlayer()

    for {
        guessingGame.run()
        if g := &guessingGame; g.level.id == g.nLevels && g.checkGuess() {
            fmt.Print(WIN_GAME_OUTPUT)
            break
        } else if g.player.triesLeft == 0 && !g.checkGuess() {
            fmt.Printf(
                NO_MORE_TENTATIVES_DEFAULT_OUTPUT_FORMAT,
                g.level.numberToGuess,
            )

            if enterpoint() {
                g.reboot()
            } else {
                break
            }
        } else {
            break
        }
    }
}


func (game *GameState) initPlayer() {
    game.player.triesLeft = game.level.nTries
    game.player.lastGuess = NO_GUESS
}


func (game *GameState) reboot() {
    var level *GameLevel

    for level = game.level ; level.prev != nil ; level = level.prev {
        level.changeGuessingNumber()
    }

    game.level = level
    game.level.changeGuessingNumber()
    game.initPlayer()
}


func (game *GameState) run() {
    clear()
    for game.player.consumeTry() {
        fmt.Printf(
            MAIN_SECTION_DEFAULT_OUTPUT_FORMAT,
            game.level.id,
            game.level.minRange,
            game.level.maxRange,
            game.player.triesLeft + 1,
            game.giveHint(),
        )

        fmt.Scan(&game.player.lastGuess)
        clear()

        if game.checkGuess() {
            if game.transitLevel() {
                fmt.Printf(
                    CORRECT_GUESS_DEFAULT_OUTPUT_FORMAT,
                    game.level.prev.id,
                    game.level.id,
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


func (level *GameLevel) changeGuessingNumber() {
    max, min := level.maxRange, level.minRange
    level.numberToGuess = rand.Intn(max - 1) + min + 1
}


func (game *GameState) transitLevel() bool {
    if game.level != nil && game.level.next != nil {
        game.level = game.level.next
        game.initPlayer()
        return true
    } else {
        return false
    }
}


func (game *GameState) addLevel(minRange, maxRange, nTries int) {
    defer func () {
        game.nLevels += 1
    }()

    if game.level == nil {
        game.level = &GameLevel{
            id:             1,
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
            id:             level.id + 1,
            nTries:         nTries,
            minRange:       minRange,
            maxRange:       maxRange,
            prev:           level,
            next:           nil,
        }
        level.next.changeGuessingNumber()
    }
}


func (player *Player) consumeTry() bool {
    if player.triesLeft > 0 {
        player.triesLeft -= 1
        return true
    } else {
        return false
    }
}


func (game GameState) giveHint() string {
    if lastGuess := game.player.lastGuess; lastGuess == NO_GUESS {
        return "no hints yet"
    } else if game.level.numberToGuess > lastGuess {
        return fmt.Sprintf("greater than %d", lastGuess)
    } else {
        return fmt.Sprintf("lower than %d", lastGuess)
    }
}


func (game GameState) checkGuess() bool {
    return game.level.numberToGuess == game.player.lastGuess
}
