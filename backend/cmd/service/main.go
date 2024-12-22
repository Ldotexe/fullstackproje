package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"log"
	"net/http"
	"sort"
	"ss/internal/service/db"
	"strings"
	"time"
)

var jwtSecretKey = []byte("very-secret-key")

type message struct {
	FromLogin string `db:"from_login"`
	ToLogin   string `db:"to_login"`
	Txt       string `db:"txt"`
	Ts        string `db:"ts"`
}

type messageResponse struct {
	isMine bool
	text   string
}

type send struct {
	Text string `json:"text"`
}

type auth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func main() {
	ctx := context.Background()
	database, err := db.NewDB(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer database.CloseDB(ctx)

	h := handler{db: database}

	r := mux.NewRouter()
	r.HandleFunc("/users", h.UsersHandler).Methods("GET")
	r.HandleFunc("/messages/{login}", h.MessagesHandler).Methods("GET")
	r.HandleFunc("/message/send/{login}", h.SendHandler).Methods("POST")
	r.HandleFunc("/auth", h.AuthHandler).Methods("POST")
	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

type handler struct {
	db *db.Database
}

func (h *handler) UsersHandler(w http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()
	var logins []string
	h.db.Select(ctx, &logins, "SELECT login FROM users")
	w.Write([]byte(strings.Join(logins, ",")))
}

func (h *handler) MessagesHandler(w http.ResponseWriter, r *http.Request) {
	userLogin, err := GetLogin(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ctx := context.Background()
	var messages []message
	toLogin := strings.TrimPrefix(r.URL.Path, "/messages/")
	h.db.Select(ctx, &messages, "SELECT * FROM messages WHERE (from_login=$1 AND to_login=$2) OR (from_login=$2 AND to_login=$1)", toLogin, userLogin)
	sort.Slice(messages, func(i, j int) bool {
		t1, _ := time.Parse(time.RFC3339Nano, messages[i].Ts)
		t2, _ := time.Parse(time.RFC3339Nano, messages[j].Ts)
		return t1.Before(t2)
	})
	ans := make([]messageResponse, 0, len(messages))
	for _, msg := range messages {
		ans = append(ans, messageResponse{text: msg.Txt, isMine: (msg.FromLogin == userLogin)})
	}
	w.Write([]byte(messages[0].Txt))
}

func (h *handler) SendHandler(w http.ResponseWriter, r *http.Request) {
	userLogin, err := GetLogin(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ctx := context.Background()
	toLogin := strings.TrimPrefix(r.URL.Path, "/message/send/")

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var txt send
	decoder.Decode(&txt)

	h.db.Exec(
		ctx, `INSERT INTO messages(from_login, to_login, txt, ts) VALUES($1,$2,$3, $4)`, userLogin, toLogin, txt.Text, time.Now().Format(time.RFC3339Nano),
	)

}

func (h *handler) AuthHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var txt auth
	decoder.Decode(&txt)

	hasher := md5.New()
	hasher.Write([]byte(txt.Password))
	hash := hex.EncodeToString(hasher.Sum(nil))

	var password string
	err := h.db.Get(ctx, &password, "SELECT password_hash FROM users WHERE login=$1", txt.Login)
	if err != nil {
		err = h.CreateUser(ctx, txt.Login, hash)
		if err != nil {
			log.Println(err)
		}
	} else {
		if password != hash {
			w.WriteHeader(http.StatusTeapot)
			return
		}
	}

	payload := jwt.MapClaims{
		"login": txt.Login,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString(jwtSecretKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	cookie := GenerateSecureCookie("ssjwt", t)

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

func (h *handler) CreateUser(ctx context.Context, login, pass string) error {
	_, err := h.db.Exec(
		ctx, `INSERT INTO users(login,password_hash) VALUES($1,$2)`, login, pass,
	)
	return err
}

var hashKey = []byte(securecookie.GenerateRandomKey(64))
var blockKey = []byte(securecookie.GenerateRandomKey(32))
var s = securecookie.New(hashKey, blockKey)

func GenerateSecureCookie(name, value string) *http.Cookie {
	data := map[string]string{name: value}
	cookie := &http.Cookie{}
	if encoded, err := s.Encode(name, data); err == nil {
		cookie.Name = name
		cookie.Value = encoded
		cookie.Secure = true
		cookie.HttpOnly = true
	} else {
		panic(err)
	}
	return cookie
}

func GetLogin(r *http.Request) (string, error) {
	cookie, err := r.Cookie("ssjwt")
	if err != nil {

	}
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})
	if err != nil {
		return "", err
	}

	// do something with decoded claims
	userLogin := claims["login"]

	return userLogin.(string), nil
}
