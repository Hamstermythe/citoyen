package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

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
	pass  = "Ht(9@3!t3hip7hvghfFEZZ&é.,ipigbWà(56789);?NLODGRTQKOJXm/"
	user  = "root"
	db, _ = sql.Open("mysql", user+":"+pass+"@tcp(127.0.0.1:3306)/")
)

var (
	header  = template.Must(template.New("header").ParseFiles("header.html"))
	accueil = template.Must(template.New("accueil").ParseFiles("accueil.html"))
	script  = template.Must(template.New("script").ParseFiles("script.html"))
)

var users = make(map[string]*User)

func MakeAuthServer(addr string, port string, sizer int) *http.Server {
	router := mux.NewRouter()
	router.Handle("/image/", http.StripPrefix("/image/", http.FileServer(http.Dir("image/"))))
	router.HandleFunc("/Auth", authHandler)
	router.HandleFunc("/Calc", calcHandler)
	// router.HandleFunc("/SaveCalc", saveCalcHandler)
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

func main() {
	// Create a new server
	server := MakeAuthServer("localhost", "8080", 0)
	// Start the server
	go server.ListenAndServeTLS("cert/server.crt", "cert/server.key")
	// Wait for input
	fmt.Println("Press any key to stop the server")
	var input string
	fmt.Scanln(&input)
	// Stop the server
	server.Close()
}
