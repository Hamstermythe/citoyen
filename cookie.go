package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"strings"
)

func getCookie(r *http.Request) (string, string) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return "", ""
	}
	// obtenir le username et l'uuid dans le cookies crypté
	username, uuid, err := decrypt(cookie.Value)
	if err != nil {
		return "", ""
	}
	return username, uuid
}

func setCookie(w http.ResponseWriter, username, uuid string) {
	// créer un cookie crypté contenant le username et l'uuid
	encryptId, err := encrypt(username, uuid)
	if err != nil {
		return
	}
	cookie := http.Cookie{
		Name:     "session",
		Value:    encryptId, // encrypt(username, uuid),
		HttpOnly: true,
	}
	cookie.SameSite = http.SameSiteNoneMode
	cookie.Secure = true
	// ajouter le cookie dans la réponse
	http.SetCookie(w, &cookie)
}

const (
	key = "une-cle-secrete-de-16-ou-32-caracteres" // Remplacez par une vraie clé secrète
)

func encrypt(username, uuid string) (string, error) {
	plaintext := []byte(username + "|" + uuid)

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func decrypt(encryptedValue string) (string, string, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(encryptedValue)
	if err != nil {
		return "", "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", "", errors.New("Ciphertext trop court")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	// La valeur décryptée devrait être sous la forme "username|uuid"
	decryptedValue := string(ciphertext)
	parts := strings.Split(decryptedValue, "|")
	if len(parts) != 2 {
		return "", "", errors.New("Format de valeur décryptée incorrect")
	}

	return parts[0], parts[1], nil
}
