package main

import (
	"fmt"
	"math/rand"
    "encoding/csv"
    "strconv"
	"os"
	"time"
)

const NO_GUESS int = -1

type Assets struct {
    correctGuessOutStr string
    gameWonOutStr string
    mainSectionOutStr string
    zeroAttempsOutStr string
}

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
    wasWon bool
    assets Assets
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


func check(err error) {
    if err != nil {
        panic(err)
    }
}


func main() {
    rand.Seed(time.Now().UnixNano())
    guessingGame := GameState{}
    guessingGame.init()

    for {
        guessingGame.run()

        if guessingGame.wasWon {
            // If the player wins show this output
            fmt.Print(guessingGame.assets.gameWonOutStr)
            break
        } else if guessingGame.player.triesLeft == 0 {
            // If the player doesn't have more tries
            // show the output bellow
            fmt.Printf(
                guessingGame.assets.zeroAttempsOutStr,
                guessingGame.level.numberToGuess,
            )

            if enterpoint() {
                // If the player clicks 'enter' or another keyword than 'q'
                // continue the game execution and reboot the game.
                guessingGame.reboot()
            } else {
                // If the player choses to quit break the execution here
                break
            }
        } else {
            // If this section is reach means that the player opted to 
            // to quit the game instead of play the next level
            break
        }
    }
}


func (asset *Assets) load() {
    var dat []byte
    var err error

    dat, err = os.ReadFile("./assets/correct_guess.txt")
    check(err)
    asset.correctGuessOutStr = string(dat)

    dat, err = os.ReadFile("./assets/no_more_attempts.txt")
    asset.zeroAttempsOutStr = string(dat)
    check(err)

    dat, err = os.ReadFile("./assets/main_section.txt")
    asset.mainSectionOutStr = string(dat)
    check(err)

    dat, err = os.ReadFile("./assets/game_won.txt")
    asset.gameWonOutStr = string(dat)
    check(err)
}


func (level *GameLevel) changeGuessingNumber() {
    max, min := level.maxRange, level.minRange
    level.numberToGuess = rand.Intn(max - 1) + min + 1
}


func (player *Player) consumeTry() bool {
    if player.triesLeft > 0 {
        player.triesLeft -= 1
        return true
    } else {
        return false
    }
}


func (game *GameState) init() {
    game.player = Player{}
    game.assets = Assets{}

    file, err := os.Open("./assets/levels.csv")
    check(err)

    defer file.Close()

    csvReader := csv.NewReader(file)
    data, err := csvReader.ReadAll()
    check(err)

    // Reads the info of the levels of the game 
    // from the file and insert new levels.
    for _, level := range data {
        min, err := strconv.Atoi(level[0])
        check(err)
        max, err := strconv.Atoi(level[1])
        check(err)
        attempts, err := strconv.Atoi(level[2])
        check(err)

        game.addLevel(min, max, attempts)
    }

    game.assets.load()
    game.initPlayer()
}


func (game *GameState) initPlayer() {
    game.player.triesLeft = game.level.nTries
    game.player.lastGuess = NO_GUESS
}


func (game *GameState) reboot() {
    var level *GameLevel

    // Traversing the list from the next to the previous 
    // node to reinitialize the number to be guessed.
    for level = game.level ; level.prev != nil ; level = level.prev {
        level.changeGuessingNumber()
    }

    // Referencing the current level of the game to the first level
    // and changing the number to predicted in this level.
    game.level = level
    game.level.changeGuessingNumber()
    game.initPlayer()
}


func (game *GameState) run() {
    clear()

    // Executes while the player has attempts left
    for game.player.consumeTry() {
        fmt.Printf(
            game.assets.mainSectionOutStr,
            game.level.id,
            game.level.minRange,
            game.level.maxRange,
            game.player.triesLeft + 1,
            game.giveHint(),
        )

        // Get the player's guess
        fmt.Scan(&game.player.lastGuess)
        clear()

        // Comparing the guess, if they match
        // the player will transit to the next level
        // else it will use another attempt
        if game.checkGuess() {
            // If there are any level to transit continue the game
            // to the next level, else it means that the player had
            // transited to all levels which also means that the 
            // player won the game.
            if game.transitLevel() {
                fmt.Printf(
                    game.assets.correctGuessOutStr,
                    game.level.prev.id,
                    game.level.id,
                )

                if !enterpoint() {
                    // If the player chose to quit,
                    // break here.
                    break
                } else {
                    // Else clear the screen and
                    // continues to the next level
                    clear()
                }
            } else {
                // If this section is reach, means that there are no more level,
                // that is, the game was won.
                game.wasWon = true
                break
            }
        }
    }
}


func (game *GameState) transitLevel() bool {
    // Changes the current level to the next level
    // in the list of levels.
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
        // Increases the number of level in the state
        // at the end of the scope of this function.
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


func (game GameState) giveHint() string {
    // analyses the last guess of the player and returns a hint,
    // informing if the number to be guessed it lower or greater
    // than the one the player gave.
    if lastGuess := game.player.lastGuess; lastGuess == NO_GUESS {
        return "no hints yet"
    } else if game.level.numberToGuess > lastGuess {
        return fmt.Sprintf("greater than %d", lastGuess)
    } else {
        return fmt.Sprintf("lower than %d", lastGuess)
    }
}


func (game GameState) checkGuess() bool {
    // compares the last guess of the player to the 
    // number to be guessed.
    return game.level.numberToGuess == game.player.lastGuess
}
