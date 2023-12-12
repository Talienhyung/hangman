# HangMan

Hangman is a classic word guessing game where players have to guess a word or phrase by suggesting letters. Guess correctly to prevent the character from being hanged! A fun game to test your vocabulary and deduction skills."

***
## Table of contents

- [Installation](#installation)
- [Utilisation](#utilisation)
- [Using the Hangman game](#using-the-hangman-game)
- [Licence](#license)

***
## Installation

Before that, you'll need to download golang
https://go.dev/doc/install

```bash
$ git clone https://ytrack.learn.ynov.com/git/rsoleane/Hangman.git
$ go build -o hangman
```

***
## Utilisation
The game is launched with several arguments, such as the dictionary or the type of game you want to launch. For more information on the arguments and rules, use the following commands in your terminal :

```bash
$ ./hangman --help # See all command
$ ./hangman --rules # See the rule of game
```

***
## Using the Hangman game

The Hangman game offers three different game types: ASCII, Classic and ThemeBox, as well as the possibility of retrieving a game save. Follow the instructions below to play:

### Game Types

#### 1. ASCII
- Launch the game with ./hangman --ascii
- The game will offer you a word to guess, displayed as ASCII art.
- Suggest letters to guess the word and avoid hanging the character.

#### 2. Classic
- Launch the game with ./hangman --classic
- You'll be presented with a word to guess in a classic game environment.
- Suggest letters to find the hidden word and save the character.

#### 3. TermBox
- Launch the game with ./hangman
- Play in a theme-specific environment where the game adapts its visuals to the chosen theme.
- Guess the word by proposing letters while exploring unique thematic universes.

### Recovering a Savegame
If you've previously saved a game, you can recover it :
- Launch the game with ./hangman --startWith [file]

### Changing the ascii art font
You can choose the ascii art font :
- ./hangman --letterFile thinkertoy.txt
- ./hangman --letterFile shadow.txt
- ./hangman --letterFile standard.txt

***
## License

This game is protected by copyright and is under a proprietary license. Any use, distribution, or modification of this game without prior authorization from Rivier Soleane is strictly prohibited. For licensing or permission inquiries, please contact us at soleane.rivier@ynov.com.

All rights reserved © 2023