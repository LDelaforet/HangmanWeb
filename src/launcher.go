package main

import (
	"fmt"
	hangman "hangman"
	"os"
	"strings"
)

func main() {
	fichier := "wordLists\\words.txt"
	args := os.Args[1:]
	if len(args) == 1 {
		fichier = "wordLists\\" + args[0]
	}

	helpMessage := `
A simple CLI only implementation of the famous hangman game in golang.

Usage:
  hangman <file>
  hangman (-h | --help)

Options:
  -h, --help     Show this screen.

Arguments:
  <file>        The file with the word list for the game.
  `

	// J'essaie au maximum de respecter docopt car j'estime que les normes sont une bonne chose
	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			// Soit on met help soit on met un fichier
			fmt.Print(helpMessage)
			return
		}
		if strings.HasPrefix(arg, "-") {
			fmt.Println("Fatal: Invalid option " + arg)
			return
		}
	}
	if !hangman.FileExists(fichier) {
		fmt.Println("Fatal: The file " + fichier + " does not exist.")
		return
	}
	hangman.ClearScreen()
	hangman.MainProgram(fichier)
}
