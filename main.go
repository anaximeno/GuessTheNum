package main
// TODO: add feature to add the attempts that left from the
// previous level to this level
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


// Struct defining the levels of the
// game as a doubly-linked list.
type GameLevel struct {
    id int
    nAttempts int
    minRange int
    maxRange int
    numberToGuess int
    next *GameLevel
    prev *GameLevel
}

type Player struct {
    attempts int
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


// TODO: Add expect message on error
func check(err error, expect string) {
    if err != nil {
        fmt.Println(expect)
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
            output := guessingGame.assets.gameWonOutStr
            fmt.Print(output)
            break
        } else if guessingGame.player.attempts == 0 {
            // If the player doesn't have more attempts
            // show the output bellow
            output := guessingGame.assets.zeroAttempsOutStr
            num := guessingGame.level.numberToGuess
            fmt.Printf(output, num)
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
    check(err, "Asset file 'correct_guess.txt' not found!")
    asset.correctGuessOutStr = string(dat)

    dat, err = os.ReadFile("./assets/no_more_attempts.txt")
    check(err, "Asset file 'no_more_attemps.txt' not found!")
    asset.zeroAttempsOutStr = string(dat)

    dat, err = os.ReadFile("./assets/main_section.txt")
    check(err, "Asset file 'main_section.txt' not found!")
    asset.mainSectionOutStr = string(dat)

    dat, err = os.ReadFile("./assets/game_won.txt")
    check(err, "Asset file 'game_won.txt' not found!")
    asset.gameWonOutStr = string(dat)
}


func (level *GameLevel) changeGuessingNumber() {
    max, min := level.maxRange - 1, level.minRange + 1
    level.numberToGuess = rand.Intn(max) + min
}


func (player *Player) consumeAttempt() bool {
    if player.attempts > 0 {
        player.attempts -= 1
        return true
    } else {
        return false
    }
}


func (game *GameState) init() {
    game.player = Player{}
    game.assets = Assets{}

    defer game.assets.load()
    defer game.initPlayer()

    file, err := os.Open("./assets/levels.csv")
    check(err, "Asset file 'levels.csv' not found!")
    csvReader := csv.NewReader(file)
    data, err := csvReader.ReadAll()
    check(err, "Could not read the levels csv file!")
    file.Close()

    // Reads the info of the levels of the game
    // from the file and insert new levels.
    for i, level := range data {
        if i == 0 { continue } // Ignore the name of the columns
        min, err := strconv.Atoi(level[0])
        check(err, "Error converting value to integer!")
        max, err := strconv.Atoi(level[1])
        check(err, "Error converting value to integer!")
        attempts, err := strconv.Atoi(level[2])
        check(err, "Error converting value to integer!")

        game.addLevel(min, max, attempts)
    }
}


func (game *GameState) initPlayer() {
    game.player.attempts = game.level.nAttempts
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
    for game.player.consumeAttempt() {
        output := game.assets.mainSectionOutStr
        level := game.level.id
        min := game.level.minRange
        max := game.level.maxRange
        attempts := game.player.attempts + 1
        nLevels := game.nLevels
        hint := game.giveHint()

        fmt.Printf(output, level, nLevels, min, max, attempts, hint)
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
                output := game.assets.correctGuessOutStr
                prevLevel := game.level.prev.id
                level := game.level.id
                fmt.Printf(output, prevLevel, level)
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
                // If this section is reach,
                // it means that there are no more levels,
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


func (game *GameState) addLevel(minRange, maxRange, attempts int) {

    defer func () {
        // Increases the number of level in the state
        // at the end of the scope of this function.
        game.nLevels += 1
    }()

    if game.level == nil {
        game.level = &GameLevel{
            id:             1,
            nAttempts:      attempts,
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
            nAttempts:      attempts,
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
