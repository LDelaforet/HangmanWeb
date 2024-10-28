package hangman

import (
	"fmt"
	"net/http"
	"strconv"
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
		if !soloPlay {
			temp.ExecuteTemplate(w, "Inputword", passedVariables)
		}
		ChooseWord()
		*RemainingLivesPtr = 9
		*TriedLettersPtr = make([]rune, 0)
		*TriesPtr = 1

		RevealLetter(2) // Révèle le nombre donné en argument de lettres, si l'argument est 99: révèle la première et la dernière lettre
		if soloPlay {
			timer = StartTimer()
		} else {
			ClearScreen()
			fmt.Println(ToCenter("Le joueur " + strconv.Itoa(currentPlayer) + " doit trouver le mot"))
			fmt.Print(ToCenter("Appuyez sur entrée pour continuer"))
		}

		passedVariables.TriedLetters = string(TriedLetters)
		passedVariables.LifeImage = "lifeCounter_" + strconv.Itoa(RemainingLives) + ".png"
		passedVariables.LettresTentees = string(TriedLetters)
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
		fmt.Print("Lettres trouvées: ")
		for _, letter := range *FoundLettersPtr {
			if letter == 0 {
				fmt.Print("_ ")
			} else {
				fmt.Print(string(letter) + " ")
			}
		}

		fmt.Print("\nMot entier: ")
		for _, letter := range *CurrentWordPtr {
			if letter == 0 {
				fmt.Print("_ ")
			} else {
				fmt.Print(string(letter) + " ")
			}
		}
		fmt.Println("\n")

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
			fmt.Println(r.FormValue("letterTB"))
			input := []rune(r.FormValue("letterTB"))
			fmt.Println("Input: ", input)
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
				fmt.Println(ToCenter("Vous avez gagné !"))
				http.Redirect(w, r, "/Game/Win", http.StatusSeeOther)
			} else if RemainingLives <= 0 {
				fmt.Println(ToCenter("Vous avez perdu !"))
				http.Redirect(w, r, "/Game/GameOver", http.StatusSeeOther)
			} else {
				http.Redirect(w, r, "/Game", http.StatusSeeOther)
			}
		}
	})

	http.HandleFunc("/Game/GameOver", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "gameOverScreen", passedVariables)
	})

	http.HandleFunc("/Game/Win", func(w http.ResponseWriter, r *http.Request) {
		totalFound++
		totalLenght += len(CurrentWord)
		totalLives += RemainingLives
		totalTimer += StopTimer(timer)

		Score = ScoreCalc(totalFound, totalLenght, totalLives, totalTimer)
		fmt.Println(ToCenter("Vous avez trouvé le mot !"))
		fmt.Println(ToCenter("Le mot était: " + string(CurrentWord)))
		fmt.Println(ToCenter("Il vous restait " + strconv.Itoa(RemainingLives) + " vies."))
		fmt.Println(ToCenter("Vous avez trouvé le mot en " + strconv.Itoa(*TriesPtr) + " essais."))
		fmt.Println(ToCenter("Temps écoulé: " + strconv.Itoa(int(time.Since(timer).Seconds()))))
		fmt.Println(ToCenter("Score actuel: " + strconv.Itoa(Score)))

		passedVariables.MotActuel = string(CurrentWord)
		passedVariables.RemainingLives = strconv.Itoa(RemainingLives)
		passedVariables.TriesCount = strconv.Itoa(*TriesPtr)
		passedVariables.TotalTimer = strconv.Itoa(int(time.Since(timer).Seconds()))
		passedVariables.CurrentScore = strconv.Itoa(Score)
		temp.ExecuteTemplate(w, "winScreen", passedVariables)
	})

	http.HandleFunc("/Game/LeaderboardRegister", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "leaderboardRegisterScreen", passedVariables)
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
		fmt.Println(readLeaderBoard())
		temp.ExecuteTemplate(w, "leaderBoard", passedVariables)
	})

	http.ListenAndServe("0.0.0.0:8080", nil)
}
