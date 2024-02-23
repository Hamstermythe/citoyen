package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func dbCreatTable() {
	// requête SQL pour créer la table users qui comprend username, password, identityFile(carte  identité, passeport), primaire, secondaire, tertiaire, quatiaire
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS users (username TEXT, password TEXT, identityFile BLOB, primaire INTEGER, secondaire INTEGER, tertiaire INTEGER, quatiaire INTEGER)")
	if err != nil {
		fmt.Println(err)
	}
}

func dbInsertUser(id string, mdp string) {
	// requête SQL pour insérer un nouvel utilisateur dans la base de données
	_, err := db.Exec("INSERT INTO users (username, password, primaire, secondaire, tertiaire, quatiaire) VALUES (?, ?, ?, ?, ?, ?)", id, mdp, 0, 0, 0, 0)
	if err != nil {
		fmt.Println(err)
	}
}

func dbGetAuth(id string, mdp string) Accueil {
	// requête SQL pour vérifier si l'identifiant et le mot de passe sont corrects
	rows, err := db.Query("SELECT username, primaire, secondaire, tertiaire, quatiaire FROM users WHERE username = ? AND password = ?", id, mdp)
	if err != nil {
		fmt.Println(err)
		return Accueil{}
	}
	rows.Close()
	// lire les données de l'utilisateur
	userData := Accueil{}
	for rows.Next() {
		err = rows.Scan(&userData.Username, &userData.TotPrimaire, &userData.TotSecondaire, &userData.TotTertiaire, &userData.TotQuatiaire)
		if err != nil {
			fmt.Println(err)
			return Accueil{}
		}
	}
	return userData
}

func dbGetUserData(id string) Accueil {
	// requête SQL pour obtenir les données de consommation de l'utilisateur
	rows, err := db.Query("SELECT username, primaire, secondaire, tertiaire, quatiaire FROM users WHERE username = ?", id)
	if err != nil {
		fmt.Println(err)
		return Accueil{}
	}
	rows.Close()
	// lire les données de l'utilisateur
	userData := Accueil{}
	for rows.Next() {
		err = rows.Scan(&userData.Username, &userData.TotPrimaire, &userData.TotSecondaire, &userData.TotTertiaire, &userData.TotQuatiaire)
		if err != nil {
			fmt.Println(err)
			return Accueil{}
		}
	}
	return userData
}

func dbGetPassword(id string) string {
	// requête SQL pour obtenir le mot de passe de l'utilisateur
	rows, err := db.Query("SELECT password FROM users WHERE username = ?", id)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	rows.Close()
	// lire les données de l'utilisateur
	var password string
	for rows.Next() {
		err = rows.Scan(&password)
		if err != nil {
			fmt.Println(err)
			return ""
		}
	}
	return password
}

func dbUpdatePassword(id string) {
	// génerer un nouveau mot de passe
	password := genRandomPassword()
	// requête SQL pour mettre à jour le mot de passe de l'utilisateur
	_, err := db.Exec("UPDATE users SET password = ? WHERE username = ?", password, id)
	if err != nil {
		fmt.Println(err)
	}
}

func genRandomPassword() string {
	const passwordLength = 16
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}|;:'\",.<>/?`~"
	randomPassword := make([]byte, passwordLength)
	charsetLength := big.NewInt(int64(len(charset)))
	for i := 0; i < passwordLength; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			// Gérer l'erreur, par exemple, en renvoyant une chaîne vide ou une valeur par défaut
			return ""
		}
		randomPassword[i] = charset[randomIndex.Int64()]
	}
	return string(randomPassword)
}

func dbUpdateConsommation(id string, primaire string, secondaire string, tertiaire string, quatiaire string) {
	// requête SQL pour mettre à jour les données de consommation de l'utilisateur
	_, err := db.Exec("UPDATE users SET primaire = ?, secondaire = ?, tertiaire = ?, quatiaire = ? WHERE username = ?", primaire, secondaire, tertiaire, quatiaire, id)
	if err != nil {
		fmt.Println(err)
	}
}

func dbUpdateIdentity(id string, identityFile []byte) {
	// requête SQL pour mettre à jour le fichier d'identité de l'utilisateur
	_, err := db.Exec("UPDATE users SET identityFile = ? WHERE username = ?", identityFile, id)
	if err != nil {
		fmt.Println(err)
	}
}
