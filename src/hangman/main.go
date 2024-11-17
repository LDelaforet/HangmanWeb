package hangman

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type PassedVars struct {
	DiscoveredLetters  string
	TriedLetters       string
	LifeImage          string
	Letter             string
	LettresTentees     string
	MotActuel          string
	RemainingLives     string
	TotalTimer         string
	CurrentScore       string
	CurrentPlayer      string
	TriesCount         string
	LeaderboardEntries []LeaderboardEntry
	SoloPlay           bool
}

type WordListChoose struct {
	WordLists []string
}

type Stats struct {
	DoDisplay   bool
	TotalFound  int
	TotalLenght int
	TotalLives  int
	TotalTimer  int
	MissedWord  string
}

var timer time.Time
var scoreFileContent string
var currentPlayer int // Joueur actuel
var player1Lost bool  // Joueur 1 a perdu
var player2Lost bool  // Joueur 2 a perdu

// Valeurs pour le score
var totalFound int  // Nombre de mots trouvés
var totalLenght int // Longueur totale des mots trouvés
var totalLives int  // Nombre total de vies perdues
var totalTimer int  // Temps total passé à jouer

var temp *template.Template
var errTemp error

// Initialise toutes les valeurs de base
func MainProgram(filename string) {
	*FileNamePtr = filename
	VarInit()
	http.Handle("/Images/", http.StripPrefix("/Images/", http.FileServer(http.Dir("./Images"))))
	http.Handle("/src/", http.StripPrefix("/src/", http.FileServer(http.Dir("./htmlSource"))))
	MainLoop()
}

func playerSwitch(p int) int {
	if p == 1 {
		return 2
	}
	return 1
}

func MainLoop() {
	temp, errTemp = template.ParseGlob("htmlSource/*.html")
	if errTemp != nil {
		fmt.Printf("Error: %v\n", errTemp)
		return
	}

	passedVariables := PassedVars{}

	soloPlay := true // On joue seul si il est sur true

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "mainMenu", passedVariables)
	})

	http.HandleFunc("/1Player", func(w http.ResponseWriter, r *http.Request) {
		soloPlay = true
		totalFound = 0
		totalLenght = 0
		totalLives = 0
		totalTimer = 0
		player1Lost = false
		player2Lost = false
		currentPlayer = 1
		http.Redirect(w, r, "/Game/init", http.StatusSeeOther)
	})

	http.HandleFunc("/2Player", func(w http.ResponseWriter, r *http.Request) {
		soloPlay = false
		totalFound = 0
		totalLenght = 0
		totalLives = 0
		totalTimer = 0
		player1Lost = false
		player2Lost = false
		currentPlayer = 1
		http.Redirect(w, r, "/Game/init", http.StatusSeeOther)
	})

	http.HandleFunc("/Game/init", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if r.FormValue("inputWord") != "" {
				*CurrentWordPtr = []rune(r.FormValue("inputWord"))
				*FoundLettersPtr = make([]rune, len(*CurrentWordPtr))
				currentPlayer = playerSwitch(currentPlayer)
			} else {
				*CurrentWordPtr = nil
			}
		}
		if !soloPlay && *CurrentWordPtr == nil {
			passedVariables.CurrentPlayer = strconv.Itoa(currentPlayer)
			passedVariables.SoloPlay = false
			temp.ExecuteTemplate(w, "Inputword", passedVariables)
			return
		} else if *CurrentWordPtr == nil {
			passedVariables.SoloPlay = true
			ChooseWord()
		}
		*RemainingLivesPtr = 9
		*TriedLettersPtr = make([]rune, 0)
		*TriesPtr = 1

		RevealLetter(2) // Révèle le nombre donné en argument de lettres, si l'argument est 99: révèle la première et la dernière lettre
		timer = StartTimer()

		var triedLettersStr []string
		for _, r := range TriedLetters {
			triedLettersStr = append(triedLettersStr, string(r))
		}

		passedVariables.TriedLetters = strings.Join(triedLettersStr, ", ") // je sais pas pourquoi y'a tried et tentées, le code est trop vieux
		passedVariables.LifeImage = "lifeCounter_" + strconv.Itoa(RemainingLives) + ".png"
		passedVariables.LettresTentees = strings.Join(triedLettersStr, ", ")
		passedVariables.Letter = " "
		passedVariables.MotActuel = string(CurrentWord)
		passedVariables.RemainingLives = strconv.Itoa(RemainingLives)
		passedVariables.TotalTimer = strconv.Itoa(totalTimer)
		passedVariables.CurrentScore = strconv.Itoa(Score)
		passedVariables.CurrentPlayer = strconv.Itoa(currentPlayer)
		passedVariables.TriesCount = strconv.Itoa(*TriesPtr)

		http.Redirect(w, r, "/Game", http.StatusSeeOther)
	})

	http.HandleFunc("/Game", func(w http.ResponseWriter, r *http.Request) {
		passedVariables.LifeImage = "lifeCounter_" + strconv.Itoa(9-RemainingLives) + ".png"
		passedVariables.DiscoveredLetters = ""

		for _, letter := range *FoundLettersPtr {
			if letter == 0 {
				passedVariables.DiscoveredLetters += " _"
			} else {
				passedVariables.DiscoveredLetters += " " + string(letter)
			}
		}
		temp.ExecuteTemplate(w, "gameMenu", passedVariables)
	})

	http.HandleFunc("/Game/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			input := []rune(r.FormValue("letterTB"))
			if len(input) == 0 {
				http.Redirect(w, r, "/Game", http.StatusSeeOther)
			} else if len(input) == 1 {
				if !checkLetter(input[0]) {
					if !checkIfTried(input[0]) {
						*RemainingLivesPtr -= 1
					}
				}
				*TriedLettersPtr = append(*TriedLettersPtr, input[0])
				passedVariables.TriedLetters = string(TriedLetters)
			} else {
				if !checkWholeWord(input) {
					*RemainingLivesPtr -= 2
				}
			}
			*TriesPtr++
			if checkWord(FoundLetters) {
				http.Redirect(w, r, "/Game/Win", http.StatusSeeOther)
			} else if RemainingLives <= 0 {
				http.Redirect(w, r, "/Game/GameOver", http.StatusSeeOther)
			} else {
				http.Redirect(w, r, "/Game", http.StatusSeeOther)
			}
		}
	})

	http.HandleFunc("/Game/GameOver", func(w http.ResponseWriter, r *http.Request) {
		doDisplay := true
		if totalFound == 0 || totalLenght == 0 || totalLives == 0 || totalTimer == 0 {
			doDisplay = false
		}
		if !soloPlay {
			http.Redirect(w, r, "/Game/multiWin", http.StatusSeeOther)
		}
		stats := Stats{doDisplay, totalFound, totalLenght, totalLives, totalTimer, string(CurrentWord)}
		*CurrentWordPtr = nil
		temp.ExecuteTemplate(w, "loseScreen", stats)
	})

	http.HandleFunc("/Game/Win", func(w http.ResponseWriter, r *http.Request) {
		totalFound++
		totalLenght += len(CurrentWord)
		totalLives += RemainingLives
		totalTimer += StopTimer(timer)

		Score = ScoreCalc(totalFound, totalLenght, totalLives, totalTimer)

		passedVariables.MotActuel = string(CurrentWord)
		passedVariables.RemainingLives = strconv.Itoa(RemainingLives)
		passedVariables.TriesCount = strconv.Itoa(*TriesPtr)
		passedVariables.TotalTimer = strconv.Itoa(int(time.Since(timer).Seconds()))
		passedVariables.CurrentScore = strconv.Itoa(Score)

		*CurrentWordPtr = nil

		temp.ExecuteTemplate(w, "winScreen", passedVariables)
	})

	http.HandleFunc("/Game/LeaderboardRegister/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if r.FormValue("nameTB") != "" {
				AddToLeaderboard(r.FormValue("nameTB"), ScoreCalc(totalFound, totalLenght, totalLives, totalTimer))
				http.Redirect(w, r, "/Leaderboard", http.StatusSeeOther)
			} else {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
		}
	})

	http.HandleFunc("/Leaderboard", func(w http.ResponseWriter, r *http.Request) {
		passedVariables.LeaderboardEntries = readLeaderBoard()
		temp.ExecuteTemplate(w, "leaderBoard", passedVariables)
	})

	http.HandleFunc("/WordList", func(w http.ResponseWriter, r *http.Request) {
		files := ListFiles("wordLists")
		wordList := WordListChoose{WordLists: files}
		temp.ExecuteTemplate(w, "wordListChoose", wordList)
	})
	http.HandleFunc("/WordList/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if r.FormValue("wordList") != "" {
				filename := "wordLists\\" + r.FormValue("wordList")
				reInitWordList(filename)
			} else {
				http.Redirect(w, r, "/WordList", http.StatusSeeOther)
			}
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
	http.HandleFunc("/Game/multiWin", func(w http.ResponseWriter, r *http.Request) {
		currentPlayer = playerSwitch(currentPlayer)
		passedVariables.CurrentPlayer = strconv.Itoa(currentPlayer)
		temp.ExecuteTemplate(w, "multiWinScreen", passedVariables)
	})

	fmt.Println("Serveur lancé sur le port 8080")
	http.ListenAndServe("0.0.0.0:8080", nil)
}

func ListFiles(directory string) []string {
	dir, _ := os.Open(directory)
	defer dir.Close()

	fileInfos, _ := dir.Readdir(0)

	var result []string
	for _, fileInfo := range fileInfos {
		result = append(result, fileInfo.Name())
	}
	return result
}
