package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func authHandler(w http.ResponseWriter, r *http.Request) {
	// lire le formulaire reçu contenant les identifiants id et mdp
	r.ParseForm()
	id := r.FormValue("id")
	mdp := r.FormValue("mdp")
	// vérifier si l'identifiant et le mot de passe sont corrects
	if userData := dbGetAuth(id, mdp); userData.Username != "" {
		// si c'est le cas, afficher la page d'accueil
		var user User
		user.Username = userData.Username
		user.UUID = uuid.New().String()
		setCookie(w, user.Username, user.UUID)
		users[user.UUID] = &user
		err := accueil.ExecuteTemplate(w, "accueil.html", userData)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		// sinon, afficher un message d'erreur
	}
}

func calcHandler(w http.ResponseWriter, r *http.Request) {
	// obtenir le username et l'uuid dans le cookies crypté
	username, uuid := getCookie(r)
	// lire le formulaire reçu contenant les données de l'utilisateur
	r.ParseForm()
	tipe := r.FormValue("type")
	// la quantité est en grammes pour les aliments et les consommables
	// la quantité est en heures pour les services
	// var userData Accueil
	if tipe == "alimentaire" {
		quantite, err := strconv.ParseFloat(strings.TrimSpace(r.FormValue("quantite poids")), 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		if quantite > 1000000 {
			return
		}
		// 1 gramme de consommation = 60 secondes de production
		// quantiteToDuration := time.Duration(time.Second * time.Duration(quantite) * 60)
		quantityPrimaire := quantite * 0.4
		quantitySecondaire := quantite * 0.2
		quantityTertiaire := quantite * 0.2
		quantityQuatiaire := quantite * 0.2
		durationPrimaire := time.Duration(time.Second * time.Duration(quantityPrimaire) * 120)
		durationSecondaire := time.Duration(time.Second * time.Duration(quantitySecondaire) * 120)
		durationTertiaire := time.Duration(time.Second * time.Duration(quantityTertiaire) * 120)
		durationQuatiaire := time.Duration(time.Second * time.Duration(quantityQuatiaire) * 120)
		var queredUserData Accueil
		if user, ok := users[uuid]; ok {
			if user.Username != username {
				dbUpdateConsommation(username, int(durationPrimaire), int(durationSecondaire), int(durationTertiaire), int(durationQuatiaire))
				queredUserData = dbGetUserData(username)
			}
		}
		queredUserData.Primaire = durationToHourMinSec(durationPrimaire)
		queredUserData.Secondaire = durationToHourMinSec(durationSecondaire)
		queredUserData.Tertiaire = durationToHourMinSec(durationTertiaire)
		queredUserData.Quatiaire = durationToHourMinSec(durationQuatiaire)
		accueilTemplate(w, queredUserData)
	} else if tipe == "consommable" {
		quantite, err := strconv.ParseFloat(strings.TrimSpace(r.FormValue("quantité poids")), 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		if quantite > 1000000 {
			return
		}
		// 1 gramme de consommation = 60 secondes de production
		// quantiteToDuration := time.Duration(time.Second * time.Duration(quantite) * 60)
		quantityPrimaire := quantite * 0.5
		quantitySecondaire := quantite * 0.2
		quantityTertiaire := quantite * 0.1
		quantityQuatiaire := quantite * 0.2
		durationPrimaire := time.Duration(time.Second * time.Duration(quantityPrimaire) * 240)
		durationSecondaire := time.Duration(time.Second * time.Duration(quantitySecondaire) * 240)
		durationTertiaire := time.Duration(time.Second * time.Duration(quantityTertiaire) * 240)
		durationQuatiaire := time.Duration(time.Second * time.Duration(quantityQuatiaire) * 240)
		var queredUserData Accueil
		if user, ok := users[uuid]; ok {
			if user.Username != username {
				dbUpdateConsommation(username, int(durationPrimaire), int(durationSecondaire), int(durationTertiaire), int(durationQuatiaire))
				queredUserData = dbGetUserData(username)
			}
		}
		queredUserData.Primaire = durationToHourMinSec(durationPrimaire)
		queredUserData.Secondaire = durationToHourMinSec(durationSecondaire)
		queredUserData.Tertiaire = durationToHourMinSec(durationTertiaire)
		queredUserData.Quatiaire = durationToHourMinSec(durationQuatiaire)
		accueilTemplate(w, queredUserData)
	} else if tipe == "service" {
		quantityHour, err := strconv.Atoi(r.FormValue("quantité heures"))
		if err != nil {
			fmt.Println(err)
			return
		}
		quantityMinute, err := strconv.Atoi(r.FormValue("quantité minutes"))
		if err != nil {
			fmt.Println(err)
			return
		}
		durationHour := time.Duration(time.Hour * time.Duration(quantityHour))
		durationMinute := time.Duration(time.Minute * time.Duration(quantityMinute))
		totDuration := durationHour + durationMinute
		durationPrimaire := time.Duration(int(float64(totDuration) * 0.4))
		durationSecondaire := time.Duration(int(float64(totDuration) * 0.1))
		durationTertiaire := time.Duration(int(float64(totDuration) * 0.4))
		durationQuatiaire := time.Duration(int(float64(totDuration) * 0.1))
		var queredUserData Accueil
		if user, ok := users[uuid]; ok {
			if user.Username != username {
				dbUpdateConsommation(username, int(durationPrimaire), int(durationSecondaire), int(durationTertiaire), int(durationQuatiaire))
				queredUserData = dbGetUserData(username)
			}
		}
		queredUserData.Primaire = durationToHourMinSec(durationPrimaire)
		queredUserData.Secondaire = durationToHourMinSec(durationSecondaire)
		queredUserData.Tertiaire = durationToHourMinSec(durationTertiaire)
		accueilTemplate(w, queredUserData)
	}
	fmt.Fprint(w, "Données calculées")
}

func durationToHourMinSec(duree time.Duration) string {
	hoursPrimaire := (duree / time.Hour) //- ((durationPrimaire%time.Hour)))
	duree -= (duree / time.Hour) * time.Hour
	minutesPrimaire := (duree / time.Minute)
	duree -= (duree / time.Minute) * time.Minute
	secondsPrimaire := (duree / time.Second)
	return strconv.Itoa(int(hoursPrimaire)) + ":" + strconv.Itoa(int(minutesPrimaire)) + ":" + strconv.Itoa(int(secondsPrimaire))
}

func accueilTemplate(w http.ResponseWriter, userData Accueil) {
	err := header.ExecuteTemplate(w, "header.html", userData)
	if err != nil {
		fmt.Println("error to execute header.html")
	}
	err = accueil.ExecuteTemplate(w, "accueil.html", userData)
	if err != nil {
		fmt.Println("error to execute accueil.html")
	}
	err = script.ExecuteTemplate(w, "script.html", nil)
	if err != nil {
		fmt.Println("error to execute script.html")
	}
}
