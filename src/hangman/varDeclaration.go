package hangman

import (
	"strings"
)

// Détermine si le mode debug est activé
var DebugMode bool = false
var DebugModePtr = &DebugMode

// Contient la liste des mots possibles
var WordList []string
var WordListPtr = &WordList

// Nombre de vies restantes au joueur
var RemainingLives int
var RemainingLivesPtr = &RemainingLives

// Mot actuel en slices de runes
var CurrentWord []rune
var CurrentWordPtr = &CurrentWord

// Lettres déjà trouvées dans le mot
var FoundLetters []rune
var FoundLettersPtr = &FoundLetters

// Lettres déjà essayées
var TriedLetters []rune
var TriedLettersPtr = &TriedLetters

// Nombre d'essais
var Tries int
var TriesPtr = &Tries

// Score actuel
var Score int
var ScorePtr = &Score

// Nom du fichier contenant les mots
var FileName string
var FileNamePtr = &FileName

// Leaderboard filename
var LeaderboardFileName string
var LeaderBoardFileNamePtr = &LeaderboardFileName

// Status actuel du jeu
var GameStatus string
var GameStatusPtr = &GameStatus

type LeaderboardEntry struct {
	Name  string
	Score int
}

// Initialise toute les déclarations de variables
func VarInit() {
	wordListInit()
	*RemainingLivesPtr = 9
	*LeaderBoardFileNamePtr = "leaderboards\\" + strings.Split(FileName, "wordLists\\")[1]
}

func reInitWordList(filename string) {
	filename = strings.ReplaceAll(filename, " ", "")
	*FileNamePtr = filename
	wordListInit()
	*LeaderBoardFileNamePtr = "leaderboards\\" + strings.Split(FileName, "wordLists\\")[1]
}

// Lis le fichier containant les mots et les ajoute à la liste
func wordListInit() {
	wordListFile, err := ReadFile(*FileNamePtr)
	if err != nil {
		panic(err)
	}
	*WordListPtr = SplitAndFormatLines(wordListFile)
}

func WebServerInit() {
	VarInit()
}
