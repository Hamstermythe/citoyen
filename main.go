package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Accueil struct {
	Username    string
	Information string
	// valeur de retour du calcul pour un produit
	Primaire   string // time.Time
	Secondaire string // time.Time
	Tertiaire  string // time.Time
	Quatiaire  string // time.Time
	// valeur de retour pour le total de consommation d'un utilisateur
	TotPrimaire   string // time.Time
	TotSecondaire string // time.Time
	TotTertiaire  string // time.Time
	TotQuatiaire  string // time.Time
}

type User struct {
	Username string
	UUID     string
}

var (
	pass    = "Ht(9@3!t3hip7hvghfFEZZ&é.,ipigbWà(56789);?NLODGRTQKOJXm/"
	user    = "root"
	db, _   = sql.Open("mysql", user+":"+pass+"@tcp(127.0.0.1:8888)/")
	accueil = template.Must(template.New("accueil").ParseFiles("accueil.html"))
)

var users = make(map[string]*User)

func MakeAuthServer(addr string, port string, sizer int) *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/Auth", authHandler)
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
	listener := &http.Server{
		Addr:         addr + ":" + port,
		Handler:      router,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	return listener
}

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
	if user, ok := users[uuid]; ok {
		if user.Username != username {
			return
		}
	}
	// lire le formulaire reçu contenant les données de l'utilisateur
	r.ParseForm()
	tipe := r.FormValue("type")
	// la quantité est en grammes pour les aliments et les consommables
	// la quantité est en heures pour les services
	// var userData Accueil
	if tipe == "alimentaire" {
		quantite, err := strconv.ParseFloat(strings.TrimSpace(r.FormValue("quantite")), 64)
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
		quantitySecondaire := quantite * 0.3
		quantityTertiaire := quantite * 0.2
		durationPrimaire := time.Duration(time.Second * time.Duration(quantityPrimaire) * 60)
		durationSecondaire := time.Duration(time.Second * time.Duration(quantitySecondaire) * 60)
		durationTertiaire := time.Duration(time.Second * time.Duration(quantityTertiaire) * 60)
		queredUserData := dbGetUserData(username)
		queredUserData.Primaire = durationToHourMinSec(durationPrimaire)
		queredUserData.Secondaire = durationToHourMinSec(durationSecondaire)
		queredUserData.Tertiaire = durationToHourMinSec(durationTertiaire)
		err = accueil.ExecuteTemplate(w, "accueil.html", queredUserData)
		if err != nil {
			fmt.Println("error to execute accueil.html")
		}
	} else if tipe == "consommable" {
		quantite, err := strconv.ParseFloat(strings.TrimSpace(r.FormValue("quantite")), 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		if quantite > 1000000 {
			return
		}
		// 1 gramme de consommation = 60 secondes de production
		// quantiteToDuration := time.Duration(time.Second * time.Duration(quantite) * 60)
		quantityPrimaire := quantite * 0.7
		quantitySecondaire := quantite * 0.2
		quantityTertiaire := quantite * 0.1
		durationPrimaire := time.Duration(time.Second * time.Duration(quantityPrimaire) * 60)
		durationSecondaire := time.Duration(time.Second * time.Duration(quantitySecondaire) * 60)
		durationTertiaire := time.Duration(time.Second * time.Duration(quantityTertiaire) * 60)
		/*
			userData.Primaire = durationToHourMinSec(durationPrimaire) // strconv.Itoa(int(hoursPrimaire)) + ":" + strconv.Itoa(int(minutesPrimaire)) + ":" + strconv.Itoa(int(secondsPrimaire))
			userData.Secondaire = durationToHourMinSec(durationSecondaire)
			userData.Tertiaire = durationToHourMinSec(durationTertiaire)
		*/
		queredUserData := dbGetUserData(username)
		queredUserData.Primaire = durationToHourMinSec(durationPrimaire)
		queredUserData.Secondaire = durationToHourMinSec(durationSecondaire)
		queredUserData.Tertiaire = durationToHourMinSec(durationTertiaire)
		err = accueil.ExecuteTemplate(w, "accueil.html", queredUserData)
		if err != nil {
			fmt.Println("error to execute accueil.html")
		}
	} else if tipe == "service" {
		// calculer les données
		// afficher les données calculées
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

func main() {
	// Create a new server
	server := MakeAuthServer("localhost", "8080", 0)
	// Start the server
	go server.ListenAndServeTLS("server.crt", "server.key")
	// Wait for input
	fmt.Println("Press any key to stop the server")
	var input string
	fmt.Scanln(&input)
	// Stop the server
	server.Close()
}
