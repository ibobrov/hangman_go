package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
	"unicode"
)

const DictionaryPath = "resources/dictionary.json"
const StagesPath = "resources/stages.json"

func getDictionary() []string {
	file, err := os.Open(DictionaryPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	var words []string
	err = json.Unmarshal(data, &words)
	if err != nil {
		log.Fatal(err)
	}

	return words
}

func getHangmanStages() []string {
	file, err := os.Open(StagesPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	var stages []string
	err = json.Unmarshal(data, &stages)
	if err != nil {
		log.Fatal(err)
	}

	return stages
}

func getRandomWord(words *[]string) string {
	generator := rand.New(rand.NewSource(time.Now().UnixNano()))
	var randomNumber = generator.Intn(len(*words) - 1)

	return (*words)[randomNumber]
}

func printGreeting() {
	fmt.Println("The game of 'Hangman' started...")
	fmt.Println("The word is hidden. Guess it!")
	fmt.Println()
}

func printWord(mask []rune) {
	for i := range mask {
		if mask[i] != '_' {
			char := mask[i]
			fmt.Printf("%c ", char)
			fmt.Print("")
		} else {
			fmt.Printf("_ ")
		}
	}
	fmt.Println()
}

func printGoodbye(word string, win bool) {
	fmt.Printf("Hidden word %s", word)
	fmt.Println()
	if win {
		fmt.Println("Game over, you win! =)")
	} else {
		fmt.Println("Game over, you lose! =(")
	}
	fmt.Println()
}

func printHangman(tries int) {
	stages := getHangmanStages()

	if tries < 0 || tries >= len(stages) {
		fmt.Println("Invalid stage number")
		return
	}

	fmt.Println(stages[tries])
}

func askNextGame() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Want to play a game? (yes/no)")
	fmt.Println()
	answer, _ := reader.ReadString('\n')

	for {
		if answer == "yes\n" {
			return true
		} else if answer == "no\n" {
			return false
		} else {
			fmt.Println("Enter the correct answer (yes/no): ")
			answer, _ = reader.ReadString('\n')
		}
	}
}

func askChar() rune {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter a letter: ")
	char, _, _ := reader.ReadRune()
	return unicode.ToLower(char)
}

func charInWord(char rune, word string) bool {
	for _, c := range word {
		if c == char {
			return true
		}
	}
	return false
}

func updateMask(word string, mask []rune, char rune) []rune {
	for i, c := range word {
		if c == char {
			mask[i] = char
		}
	}

	return mask
}

func startGame(words []string) {
	printGreeting()

	word := getRandomWord(&words)
	mask := []rune(word)

	for i := range mask {
		mask[i] = '_'
	}

	stage := 0
	maxTries := 6
	win := false

	printWord(mask)

	for maxTries >= 0 {
		char := askChar()
		if charInWord(char, word) {
			fmt.Println("There is such a letter")
			mask = updateMask(word, mask, char)
			printWord(mask)
		} else {
			fmt.Println("There is no such letter")
			printHangman(stage)
			stage += 1
			maxTries -= 1
			printWord(mask)
		}

		if !charInWord('_', string(mask)) {
			win = true
			break
		}
	}

	printGoodbye(word, win)
}

func main() {
	words := getDictionary()

	for askNextGame() {
		startGame(words)
	}
}
