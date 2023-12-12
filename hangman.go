package hangman

import (
	"fmt"
	"os"
	"os/exec"
	"math/rand"
	"unicode/utf8"
	"encoding/json"
	"github.com/nsf/termbox-go"
	"log"
	"bufio"
)

type HangManData struct {
	Word             []rune   // Word composed of '_', ex: H_ll_
	ToFind           string   // Final word chosen by the program at the beginning. It is the word to find
	Attempts         int      // Number of attempts left
	HangmanPositions int      // Positions parsed in "hangman.txt" are stored
	ListWord         []string // List of words suggested by the user
	ListLetter       []rune   // List of letter sugested by the user
	LastFail         bool     // Used to find out the status of the last input (used in the display).
}

type Game struct {
	save       bool   // True if the --startWith (-sw) argument is given
	classic    bool   // True if the --classic (-c) argument is given
	ascii      bool   // True if the --ascii (-a) argument is given
	letter     bool   // True if the --letter (-l) argument is given
	saveFile   string // Name of the file given after --startWith (-sw) where the backup is stored
	letterFile string // Name of the file given after --letter (-l) where the ascii art is stored
	dico       string // First argument given, contains the name of the file containing the desired dictionary
}


// Display a manual for the utilisation of argument
func Help() {
	filePath := "Ressources/doc.txt"

	pagerCommand := "less"

	cmd := exec.Command(pagerCommand, filePath)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error when executing %s: %v\n", pagerCommand, err)
		os.Exit(1)
	}
}

// Display the rule of the game
func Rules() {
	filePath := "Ressources/rules.txt"

	pagerCommand := "less"

	cmd := exec.Command(pagerCommand, filePath)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error when executing %s: %v\n", pagerCommand, err)
		os.Exit(1)
	}
}

// Set HangManData's first value
func (hangman *HangManData) setData() {
	hangman.Word = []rune{}
	hangman.ToFind = ""
	hangman.Attempts = 10
	hangman.HangmanPositions = -1
	hangman.ListLetter = []rune{}
	hangman.ListWord = []string{}
}

// Set Word and ToFind for HangManData
func (hang *HangManData) SetWord(dico []string) {
	// Find a random word
	randomIndex := rand.Intn(len(dico) - 1)
	hang.ToFind = dico[randomIndex]

	nbVisibleLetter := len(hang.ToFind)/2 - 1 // Set the number of letters that will be visible

	for range hang.ToFind { // Set Word
		hang.Word = append(hang.Word, '_')
	}

	again := false
	var place []int

	for nbVisibleLetter > 0 { // Reveal random letters in the word to find
		randomIndex = rand.Intn(len(hang.ToFind))
		for _, j := range place {
			if j == randomIndex {
				again = true
			}
		}
		if !again {
			place = append(place, randomIndex)
			nbVisibleLetter--
		} else {
			again = false
		}
	}

	WordRune := []rune(hang.ToFind)
	for _, index := range place { // Add the different letter into Word
		hang.Word[index] = WordRune[index]
		hang.LetterInWord(WordRune[index])
		if !hang.UsedVerif(string(WordRune[index])) {
			hang.UsedLetter(WordRune[index])
		}
	}
}

// Function to check whether the given letter is in ToFind
func (game *HangManData) LetterInWord(oneRune rune) {
	var place []int
	for index, letters := range game.ToFind { // Transforms letters into lower case if they are not already lower case
		if oneRune >= 'A' && oneRune <= 'Z' {
			oneRune = oneRune + 32
		}
		if letters == oneRune || letters == oneRune-32 {
			place = append(place, index) // Saves the index(es) of the position where the letter was found
		}
	}
	if len(place) != 0 { // If any letters have been found then replace the corresponding slots with the letters
		game.LastFail = false
		for _, index := range place {
			game.Word[index] = oneRune
		}
	} else { // If the letter is not found, an attempt is lost
		game.LastFail = true
		game.Attempts--
		game.HangmanPositions++
	}
}

// Function to check whether the given word is ToFind (return true if this is the case)
func (game *HangManData) IsThisTheWord(word string) bool {
	var oneRune rune
	if len(word) != len(game.ToFind) { // Check that words are the same size
		return false
	} else {
		wordRune := []rune(word)
		toFindRune := []rune(game.ToFind)
		for index, runes := range wordRune {
			oneRune = runes
			if runes >= 'A' && runes <= 'Z' { // Puts letters in lower case for easier comparison
				oneRune = oneRune + 32
			}
			if oneRune != toFindRune[index] && oneRune != toFindRune[index]+32 {
				return false
			}
		}
	}
	return true
}


// Adds the rune passed as a parameter to ListLetter if it's not already there
func (game *HangManData) UsedLetter(oneRune rune) {
	if oneRune < 'A' || oneRune > 'Z' { // Make the letter uppercase if it isn't already
		oneRune = oneRune - 32
	}
	if oneRune >= 'A' && oneRune <= 'Z' { // Add the letter into ListLetter (only if it's a letter)
		game.ListLetter = append(game.ListLetter, oneRune)
	}
}

// Adds the rune passed as a parameter to ListLetter if it's not already there
func (game *HangManData) UsedWord(word string) {
	runeWord := []rune(word)
	for index, runes := range runeWord { // Transforms letters into lower case if they are not already lower case
		if runes < 'a' || runes > 'z' {
			runeWord[index] = runes + 32
		}
	}
	game.ListWord = append(game.ListWord, string(runeWord)) // Add the word into ListLetter
}

// Checks if the input is already in one of the ListWord or ListLetter (return true if it's the case)
func (game *HangManData) UsedVerif(intput string) bool {
	runeWord := []rune(intput)
	oneRune := runeWord[0]
	if utf8.RuneCountInString(intput) > 1 {
		for index, runes := range runeWord { // Transforms letters into lower case if they are not already lower case
			if runes < 'a' || runes > 'z' {
				runeWord[index] = runes + 32
			}
		}
		for _, words := range game.ListWord { // Search the word into ListWord
			if string(runeWord) == words {
				return true
			}
		}
	} else {
		if oneRune < 'A' || oneRune > 'Z' { // Make the letter uppercase if it isn't already
			oneRune = oneRune - 32
		}
		for _, letter := range game.ListLetter { // Search the letter into ListLetter
			if letter == oneRune {
				return true
			}
		}
	}
	return false
}


// Saves the party's progress, which is stored in the HangManData structure
func (data HangManData) Save(filename string) error {
	file, err := os.Create(filename) // Create a file
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(data); err != nil { // Save HangManData
		return err
	}

	return nil
}

// Load the party's progress, which is stored in filename, return a HangManData struct
func Load(filename string) (HangManData, error) {
	var data HangManData

	file, err := os.Open(filename) // Open filename
	if err != nil {
		return data, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil { // Load HangManData
		return data, err
	}

	return data, nil
}


// Draw a box in x, y with size width/height, color borderColor with title in terminal
func drawBox(x, y, width, height int, borderColor termbox.Attribute, title string) {
	// Draw the box frame
	for i := x; i < x+width; i++ {
		termbox.SetCell(i, y, '─', borderColor, termbox.ColorDefault)
		termbox.SetCell(i, y+height-1, '─', borderColor, termbox.ColorDefault)
	}
	for i := y; i < y+height; i++ {
		termbox.SetCell(x, i, '│', borderColor, termbox.ColorDefault)
		termbox.SetCell(x+width-1, i, '│', borderColor, termbox.ColorDefault)
	}

	// Box corners
	termbox.SetCell(x, y, '┌', borderColor, termbox.ColorDefault)
	termbox.SetCell(x+width-1, y, '┐', borderColor, termbox.ColorDefault)
	termbox.SetCell(x, y+height-1, '└', borderColor, termbox.ColorDefault)
	termbox.SetCell(x+width-1, y+height-1, '┘', borderColor, termbox.ColorDefault)

	// Add title
	for i, ch := range title {
		termbox.SetCell(x+i+1, y, ch, termbox.ColorDefault, termbox.ColorDefault)
	}
}

// Main display, efficient for all boxes and hangman
func (hang *HangManData) display() {
	// Main box
	drawBox(0, 0, 100, 24, termbox.ColorWhite, "main")

	// First box inside the main box
	drawBox(55, 0, 45, 15, termbox.ColorLightYellow, "Hangman")
	// HangMan in the first box
	if hang.HangmanPositions >= 0 && hang.HangmanPositions <= 9 {
		hang.DisplayHangman(55+18, 4, termbox.ColorBlue)
	}

	// Second box inside the main box
	drawBox(0, 0, 50, 8, termbox.ColorBlue, "Word...")
	drawBox(25, 0, 25, 8, termbox.ColorBlue, "Attempts")

	// Third box inside the main box
	drawBox(0, 8, 50, 8, termbox.ColorGreen, "Letter")

	// Fourth box inside the main box
	drawBox(0, 16, 50, 8, termbox.ColorLightMagenta, "Used letter/words")
}

// drawText is a function that draws text
// It takes a slice of runes (text), x and y coordinates, a color attribute, and a cursor flag.
func drawText(text []rune, x, y int, color termbox.Attribute, cursor bool) {
	// Initialize variables to track line and space offsets
	ligne := 0
	space := 0

	// Loop through the text and process each character
	for i, ch := range text {
		// Determine the line and space offsets based on the character index
		switch {
		case i > 45 && i <= 91:
			space = 46
			ligne = 1
		case i > 91 && i <= 137:
			space = 92
			ligne = 2
		case i > 137 && i <= 183:
			space = 138
			ligne = 3
		case i > 183 && i <= 229:
			space = 184
			ligne = 4
		default:
			ligne = 0
			space = 0
		}

		// Set the character at the specified position with the given color
		termbox.SetCell(x+i-space, y+ligne, ch, termbox.ColorDefault, termbox.ColorDefault)
	}

	// If the cursor flag is true, draw a cursor character at a specific position
	if cursor {
		termbox.SetCell(2+len(text)-space, 10+ligne, '_', termbox.ColorDefault, termbox.ColorDefault)
	}
}

// Display HangMan on the right position
func (hang *HangManData) DisplayHangman(x, y int, borderColor termbox.Attribute) {
	hangMan := readHang("Ressources/HangMan_Position/hangman.txt")
	for i := 0; i <= 7; i++ { // Display ligne by ligne
		runes := []rune(hangMan[hang.HangmanPositions][i])
		for index, j := range runes { // Display rune by rune
			termbox.SetCell(x+index, y+i, j, borderColor, termbox.ColorDefault)
		}
	}
}


// TermBoxGame is a function that handles the main game loop for a Hangman game using the termbox library.
// It takes the HangManData and Game structs as input parameters.
func (HangMan HangManData) TermBoxGame(game Game) {
	// Initialize the termbox library and handle errors
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	// Initialize variable
	word := "/"
	userInput := ""
	gameOver := false
	empty := "Empty or already proposed!"

	for {
		// Clear the screen and set up user interface
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		game.asciiBox(word)
		game.AsciiCounter(HangMan.Attempts)

		HangMan.display()

		// Display text
		drawText([]rune(userInput), 2, 10, termbox.ColorDefault, true)
		drawText(HangMan.Word, 2, 4, termbox.ColorDefault, false)
		drawText(HangMan.ListLetter, 2, 17, termbox.ColorDefault, false)
		for i := range HangMan.ListWord {
			drawText([]rune(HangMan.ListWord[i]), 2, 18+i, termbox.ColorDefault, false)
		}
		termbox.Flush()

		// Poll for user input events
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			if ev.Key == termbox.KeyEsc {
				return // Exit the game loop
			} else if ev.Key == termbox.KeySpace || ev.Key == termbox.KeyEnter {
				if !gameOver {
					if userInput != "" && !HangMan.UsedVerif(userInput) && userInput != empty {
						// Check if the user's input is a valid guess and update the word or game status
						if HangMan.mainMecanics(userInput) {
							word = "win"
							gameOver = true
						} else {
							word = userInput
						}
						userInput = "" // Clear user input
					} else {
						userInput = empty
					}
				} else {
					if userInput == "QUIT" {
						return
					}
				}
			} else if ev.Key == termbox.KeyDelete {
				userInput = "" // Clear user input
			} else if ev.Key == termbox.KeyBackspace || ev.Key == termbox.KeyBackspace2 {
				if userInput != "" && userInput != empty {
					userInput = userInput[:len(userInput)-1] // Remove the last character from user input
				}
			} else {
				if userInput == empty {
					userInput = ""
				}
				userInput += string(ev.Ch) // Add the character to user input
			}
		}

		// Check if the game has ended
		if HangMan.endGame() {
			if HangMan.Attempts <= 0 {
				word = "lose"
				HangMan.Word = []rune(HangMan.ToFind)
			} else {
				word = "win"
			}
			gameOver = true
		}
	}
}


// This is the hangman Ascii game
func (game HangManData) AsciiGame(data Game) {
	var inputs string
	gameOver := false
	fmt.Printf("Good Luck, you have %d attempts.\n", game.Attempts)
	data.displayAsciiText(game.Word)

	for !gameOver { // Game loop
		// Display input
		letter := input("\nChoose : ", inputs)

		// Verify input
		if letter != "" && !game.UsedVerif(letter) {
			if game.mainMecanics(letter) {
				gameOver = true
			}

			if game.LastFail {
				fmt.Printf("Not present in the word, %d attempts remaining\n", game.Attempts)
			}

			// Display AsciiText and HangMan
			data.displayAsciiText(game.Word)
			if game.HangmanPositions >= 0 {
				game.displayHangmanClassic()
			}

		} else {
			fmt.Println("Empty or already proposed!")
		}

		// Verify if it's the end of the game
		if game.endGame() {
			gameOver = true
		}
	}

	// Announcement of results
	if game.Attempts > 0 {
		fmt.Println("Congrats !")
	} else {
		fmt.Println("The word was " + game.ToFind + ". You'll do better next time!!!")
	}
}


// Displays a given ascii character in x y
func (data *Game) displayAscii(x, y, version int, borderColor termbox.Attribute) {
	var ascii [95][9]string
	switch data.letterFile { //choose the right font for ascii art
	case "shadow.txt":
		ascii = readAscii("Ressources/Ascii_Letter/shadow.txt")
	case "standard.txt":
		ascii = readAscii("Ressources/Ascii_Letter/standard.txt")
	case "thinkertoy.txt":
		ascii = readAscii("Ressources/Ascii_Letter/thinkertoy.txt")
	default:
		ascii = readAscii("Ressources/Ascii_Letter/standard.txt")
	}
	for i := 0; i <= 8; i++ { //displays the correct character
		runes := []rune(ascii[version-32][i])
		for index, j := range runes {
			termbox.SetCell(x+index, y+i, j, borderColor, termbox.ColorDefault)
		}
	}
}

// Displays the last letter entered by a user in the terminal, followed by the final result (win or lose).
func (data *Game) asciiBox(word string) {
	switch word {
	case "win": //display WIN if player win
		data.displayAscii(55+19, 15, 'I', termbox.ColorGreen)
		data.displayAscii(59, 15, 'W', termbox.ColorGreen)
		data.displayAscii(55+16+11, 15, 'N', termbox.ColorGreen)
	case "lose": //display lose if player lose
		data.displayAscii(56, 15, 'L', termbox.ColorRed)
		data.displayAscii(55+11, 15, 'O', termbox.ColorRed)
		data.displayAscii(55+16+5, 15, 'S', termbox.ColorRed)
		data.displayAscii(55+16+14, 15, 'E', termbox.ColorRed)
	default: //displays the first rune of the last input
		runes := []rune(word)
		if int(runes[0]) > 33 && int(runes[0]) < 126 {
			data.displayAscii(55+16, 15, int(runes[0]), termbox.ColorLightRed)
		} else {
			data.displayAscii(55+16, 15, '/', termbox.ColorLightRed)
		}
	}
}

// displayAsciiText displays ASCII art text using the 'letterFile' font.
// The function selects the font based on 'letterFile' and displays the ASCII art text.
func (data *Game) displayAsciiText(words []rune) {
	// Define a 2D array 'ascii' to hold ASCII art characters.
	var ascii [95][9]string

	// Select the font for the ASCII art based on the 'letterFile' field.
	switch data.letterFile {
	case "shadow.txt":
		ascii = readAscii("Ressources/Ascii_Letter/shadow.txt")
	case "standard.txt":
		ascii = readAscii("Ressources/Ascii_Letter/standard.txt")
	case "thinkertoy.txt":
		ascii = readAscii("Ressources/Ascii_Letter/thinkertoy.txt")
	default:
		ascii = readAscii("Ressources/Ascii_Letter/standard.txt")
	}

	// Loop through each line of the ASCII art (9 lines in total).
	for line := 0; line <= 8; line++ {
		// Loop through each letter (rune) in the 'words' slice.
		for _, letter := range words {
			// Print the ASCII character for the current letter on the current line.
			// The ASCII value of the letter is used to index 'ascii' array.
			fmt.Printf(ascii[letter-32][line])
		}
		// After printing a line of text, move to the next line.
		fmt.Println("")
	}
}

func (data Game) AsciiCounter(attempts int) {
	switch attempts {
	case 1:
		data.displayAscii(35, -1, '1', termbox.ColorLightGray)
	case 10:
		data.displayAscii(33, -1, '1', termbox.ColorLightGray)
		data.displayAscii(36, -1, '0', termbox.ColorLightGray)
	default:
		data.displayAscii(34, -1, attempts+'0', termbox.ColorLightGray)
	}
}


// Displays the hangman in the terminal
func (hang HangManData) displayHangmanClassic() {
	hangMan := readHang("Ressources/HangMan_Position/hangman.txt")
	fmt.Println("")
	for i := 0; i <= 7; i++ {
		fmt.Println(hangMan[hang.HangmanPositions][i])
	}
}

// displays the rune array given as a parameter in the terminal
func printRune(tab []rune) {
	for _, runes := range tab {
		fmt.Print(string(runes))
		fmt.Print(" ")
	}
	fmt.Println("")
}

// return user input
func input(s string, inputs string) string {
	fmt.Print(s)
	fmt.Scanln(&inputs)
	return inputs
}


// This is the hangman classic game
func (game HangManData) ClassicGame() {
	var inputs string
	gameOver := false
	fmt.Printf("Good Luck, you have %d attempts.\n", game.Attempts)
	printRune(game.Word)

	for !gameOver { // Game loop
		// Display word and attempts
		letter := input("\nChoose : ", inputs)

		// Verify input
		if letter != "" && !game.UsedVerif(letter) {
			if game.mainMecanics(letter) {
				gameOver = true
			}

			if game.LastFail {
				fmt.Printf("Not present in the word, %d attempts remaining\n", game.Attempts)
			}

			// Display words and HangMan
			printRune(game.Word)
			if game.HangmanPositions >= 0 {
				game.displayHangmanClassic()
			}

			// Verify if it's the end of the game
			if game.endGame() {
				gameOver = true
			}
		} else {
			fmt.Println("Empty or already proposed!")
		}
	}

	// Announcement of results
	if game.Attempts > 0 {
		fmt.Println("Congrats !")
	} else {
		fmt.Println("The word was " + game.ToFind + ". You'll do better next time!!!")
	}
}


// Handling arguments and adding values to Game structure parameters
func SortArguments() Game {
	var game Game
	arguments := os.Args[1:]
	needFile := true                    // If needFile is true, this means that the last argument given is an option requesting a file.
	for index, arg := range arguments { // Review all the arguments
		switch arg {
		case "--startWith", "-sw":
			if needFile && index != 0 {
				fmt.Println("Invalid argument")
				os.Exit(3)
			} else {
				needFile = true
			}
			game.save = true
			if len(arguments) > index+1 {
				game.saveFile = arguments[index+1] // The backup file is saved in game.SaveFile
			} else {
				fmt.Println("Invalid argument")
				os.Exit(3)
			}

		case "--classic", "-c":
			if needFile && index != 0 {
				os.Exit(3)
			} else {
				needFile = false
			}
			if game.ascii {
				fmt.Println("Two arguments not compatible")
				os.Exit(4)
			} else {
				game.classic = true
			}
		case "--ascii", "-a":
			if needFile && index != 0 {
				fmt.Println("Invalid argument")
				os.Exit(3)
			} else {
				needFile = false
			}
			if game.classic { // Cause GameAscii is not compatible with GameClassic
				fmt.Println("Two arguments not compatible")
				os.Exit(4)
			} else {
				game.ascii = true
			}
		case "--letterFile", "-lf":
			if needFile && index != 0 {
				fmt.Println("Invalid argument")
				os.Exit(3)
			} else {
				needFile = true
			}
			if game.classic { // Cause letterFile is not compatible with GameClassic
				fmt.Println("Two arguments not compatible")
				os.Exit(4)
			} else {
				game.letter = true
			}
			if len(arguments) > index+1 {
				game.letterFile = arguments[index+1] // The ascii art font file is saved in game.letterFile
			} else {
				fmt.Println("Invalid argument")
				os.Exit(3)
			}
		case "--rules", "-r":
			if index != 0 || len(arguments) != 1 {
				fmt.Println("Invalid argument")
				os.Exit(3)
			}
			Rules()
			os.Exit(0)
		case "--help", "-h":
			if index != 0 || len(arguments) != 1 {
				fmt.Println("Invalid argument")
				os.Exit(3)
			}
			Help()
			os.Exit(0)
		default:
			if needFile && index == 0 {
				game.dico = arguments[0]
			} else if !needFile {
				fmt.Println("Invalid argument")
				os.Exit(3)
			}
			needFile = false
		}
	}
	return game
}

// Using the arguments, generates HangManData's parameter values
func ExploitingArgument(game Game) {
	var data HangManData
	if game.save { // Set HangManData
		var err error
		data, err = Load("Ressources/Save/" + game.saveFile)
		if err != nil {
			fmt.Println("Error while loading the game state:", err)
			os.Exit(2)
		}
	} else {
		data.setData()
		dico := ReadTheDico(game.dico)
		data.SetWord(dico)
	}
	if !game.letter {
		game.letterFile = "standard.txt"
	} else {
		if game.letterFile != "standard.txt" && game.letterFile != "thinkertoy.txt" && game.letterFile != "shadow.txt" {
			fmt.Println("Unrecognized letterFile (i.e. letterFile will be standard.txt)\nPress enter to accept, otherwise ^C")
			var inputs string
			fmt.Scanln(&inputs)
		}
	}
	if game.classic {
		data.ClassicGame()
		os.Exit(0)
	}
	if game.ascii {
		data.AsciiGame(game)
		os.Exit(0)
	}
	data.TermBoxGame(game) // If no mode is launched, the default mode is TermboxGame
}


// This function reads the given ascii file and returns a [95][9]string containing the ascii art characters.
func readAscii(fichier string) [95][9]string {
	var ascii [95][9]string

	readFile, err := os.Open(fichier)
	if err != nil {
		fmt.Print(err)
	}

	fileScanner := bufio.NewScanner(readFile) // Creates a scanner to read the file.

	fileScanner.Split(bufio.ScanLines) // Divides the file into lines.

	i, j := 0, 0
	// Browse each line of the file.
	for fileScanner.Scan() {
		ascii[i][j] = fileScanner.Text()
		j++
		if j == 9 {
			i++
			j = 0
		}
	}

	readFile.Close()

	return ascii
}

// This function reads the given hangman file and returns a [10][8]string containing the hangman position.
func readHang(fichier string) [10][8]string {
	var hangman [10][8]string

	readFile, err := os.Open(fichier)
	if err != nil {
		fmt.Print(err)
	}

	fileScanner := bufio.NewScanner(readFile) // Creates a scanner to read the file.

	fileScanner.Split(bufio.ScanLines) // Divides the file into lines.

	i, j := 0, 0
	// Browse each line of the file.
	for fileScanner.Scan() {
		// Adds the text of the current line to the FindWord slice.
		hangman[i][j] = fileScanner.Text()
		j++
		if j == 8 {
			i++
			j = 0
		}
	}

	readFile.Close()

	return hangman
}

//########### Dictionary function ##################

// The readfile function returns an array of strings containing all the words in a dictionary
func readFile(fichier string) []string {
	var dictio []string

	readFile, err := os.Open(fichier)
	if err != nil {
		fmt.Print(err)
	}

	fileScanner := bufio.NewScanner(readFile) // Creates a scanner to read the file.

	fileScanner.Split(bufio.ScanLines) // Divides the file into lines.

	// Browse each line of the file.
	for fileScanner.Scan() {
		// Adds the text of the current line to the FindWord slice.
		dictio = append(dictio, fileScanner.Text())
	}

	readFile.Close()

	return dictio
}

// The listDictio function returns all files in the Dictinonary directory
func listDictio() []string {
	var listDico []string
	entries, err := os.ReadDir("Ressources/Dictionary/")
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		listDico = append(listDico, e.Name())
	}
	return listDico
}

// The ReadAllDico function returns an array of strings containing all the words in the various files in the dictionary folder.
func readAllDico() []string {
	listDico := listDictio()
	var dico []string
	for i := 0; i < len(listDico); i++ {
		newDico := readFile("Ressources/Dictionary/" + listDico[i])
		dico = append(dico, newDico...)
	}
	return dico
}

// This function returns an array of words, depending on the file entered as a parameter only if the file exists, otherwise it uses the other files.
func ReadTheDico(file string) []string {
	listDico := listDictio()
	if listDico == nil {
		fmt.Println("No file in Dictionary")
		os.Exit(3)
	}
	for _, j := range listDico { // Check if the requested dictionary exists
		if file == j {
			dico := readFile("Ressources/Dictionary/" + file)
			return dico
		}
	}
	fmt.Println("Unspecified or unrecognized dictionaries (i.e. words chosen at random from all dictionaries)\nPress enter to accept, otherwise ^C")
	var inputs string
	fmt.Scanln(&inputs)
	return readAllDico()
}


// Main mecanic of the game which gathers several functions, return true if the game is finished, otherwise false.
func (hang *HangManData) mainMecanics(input string) bool {
	if utf8.RuneCountInString(input) > 1 { // If it's a word
		if hang.IsThisTheWord(input) {
			hang.Word = []rune(hang.ToFind)
			hang.LastFail = false
			return true
		} else if input == "STOP" { // If the input is STOP, save the game
			err := hang.Save("Ressources/Save/save.txt")
			if err != nil {
				termbox.Close()
				fmt.Println("Game save failed :", err)
				os.Exit(2)
			}
			termbox.Close()
			fmt.Println("Game save in save.txt")
			os.Exit(0)
		} else if input == "QUIT" { // If the input is QUIT, quit the game
			termbox.Close()
			os.Exit(0)
		} else {
			hang.UsedWord(input)
			hang.Attempts -= 2
			hang.HangmanPositions += 2
			hang.LastFail = true
			if hang.HangmanPositions > 9 { // Avoid out of range
				hang.HangmanPositions = 9
				hang.Attempts = 0
			}
		}
	} else { // If it's a letter
		oneRune := []rune(input)
		hang.LetterInWord(oneRune[0])
		hang.UsedLetter(oneRune[0])
	}
	return false
}

// Check if the game is finished or not
func (game *HangManData) endGame() bool {
	if game.Attempts <= 0 { // No more attempts
		return true
	}
	for _, runes := range game.Word {
		if runes == '_' {
			return false
		}
	}
	return true // Words found
}
