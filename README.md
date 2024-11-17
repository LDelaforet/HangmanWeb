## Hangman Web

***Implémentation en Go et en HTML/CSS du jeu du pendu***

## Sommaire

- [Fonctionnalités](#fonctionnalités)
- [Utilisation](#utilisation)
- [Routes](#routes)

## Fonctionnalités
Hangman est une implémentation simpliste en golang du jeu du pendu. Ses fonctionnalités principales sont les suivantes :

- **Liste de mots personnalisable**
  - Listes de mots exemples incluses
  - Possibilité d'importer ses propres listes de mots sans difficulté

- **Mode multi-joueur (2 joueurs max)**
  - Possibilité de faire deviner n'importe quel mot à son adversaire

- **Système de score et de tableau des scores**
  - Sauvegarde et chargement des scores selon la liste de mots utilisée
  - Nom personnalisable dans la liste des scores

- **Interface utilisateur claire et intuitive**
  - Affichage web pour un meilleur confort et une meilleure compréhension
  - Navigation entre les menus intuitive

## Utilisation
Pour lancer le jeu, il suffit d'ouvrir la console et d'exécuter la commande suivante dans le répertoire d'Hangman Web:
```shell
go run .
```
Le jeu sera ainsi accessible a la page 127.0.0.1:8080 ainsi qu'a toutes les ip depuis lesquelle votre machine est accessible depuis le port 8080.
## Routes

### Routes de vues
- `GET /` : Affiche la page d'accueil.
- `GET /WordList` : Affiche la liste des fichiers de mots disponibles.
- `GET /Leaderboard` : Affiche le tableau des scores.
- `GET /Game/multiWin` : Affiche l'écran de victoire pour le mode multijoueur.

### Routes de traitement des données
- `POST /Game/init` : Initialise une nouvelle partie.
- `POST /Game/LeaderboardRegister/submit` : Soumet un nouveau score au tableau des scores.
- `POST /WordList/submit` : Soumet un fichier de mots pour réinitialiser la liste des mots.
