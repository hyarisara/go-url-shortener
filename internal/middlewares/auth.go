package middlewares

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

//////////////////////
// SESSION
//////////////////////

func SetSession(w http.ResponseWriter, username string) {
	http.SetCookie(w, &http.Cookie{
		Name:  "session",
		Value: username,
		Path:  "/",
	})
}

func GetSession(r *http.Request) string {
	cookie, err := r.Cookie("session")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func ClearSession(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}

//////////////////////
// PASSWORD
//////////////////////

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
