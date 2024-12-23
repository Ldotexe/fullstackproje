package handler

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"sort"
	"ss/internal/service/datamanager"
	"strings"
	"time"
)

var jwtSecretKey = []byte("very-secret-key")
var jwtCookieName = "ssjwt"

type Handler interface {
	UsersHandler(w http.ResponseWriter, _ *http.Request)
	MessagesHandler(w http.ResponseWriter, r *http.Request)
	SendHandler(w http.ResponseWriter, r *http.Request)
	AuthHandler(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	dm datamanager.DataManager
}

func New(dm datamanager.DataManager) Handler {
	return &handler{dm: dm}
}

func (h *handler) UsersHandler(w http.ResponseWriter, r *http.Request) {
    userLogin, err := getLogin(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ctx := context.Background()

	logins, err := h.dm.GetUsers(ctx)
	if err != nil {
		if errors.Is(err, datamanager.ErrObjectNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

    withoutUser := make([]string, 0)
    for _, login := range logins {
       if login != userLogin {
          withoutUser = append(withoutUser, login)
       }
    }
	body, err := json.Marshal(withoutUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	_, err = w.Write(body)
	if err != nil {
		log.Println(err)
	}
	return
}

func (h *handler) MessagesHandler(w http.ResponseWriter, r *http.Request) {
	userLogin, err := getLogin(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	toLogin := strings.TrimPrefix(r.URL.Path, "/messages/")

	ctx := context.Background()

	messages, err := h.dm.GetMessages(ctx, &datamanager.Message{FromLogin: userLogin, ToLogin: toLogin})
	if err != nil {
		if errors.Is(err, datamanager.ErrObjectNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sort.Slice(messages, func(i, j int) bool {
		t1, _ := time.Parse(time.RFC3339Nano, messages[i].Ts)
		t2, _ := time.Parse(time.RFC3339Nano, messages[j].Ts)
		return t1.Before(t2)
	})
	ans := make([]messageResponse, 0, len(messages))
	for _, msg := range messages {
		ans = append(ans, messageResponse{Text: msg.Text, IsMine: msg.FromLogin == userLogin})
	}

	body, err := json.Marshal(ans)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(body)
}

func (h *handler) SendHandler(w http.ResponseWriter, r *http.Request) {
	userLogin, err := getLogin(r)
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

	h.dm.AddMessage(ctx, &datamanager.Message{FromLogin: userLogin, ToLogin: toLogin, Text: txt.Text, Ts: time.Now().Format(time.RFC3339Nano)})
	//TODO: handle err and evrth

}

func (h *handler) AuthHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var txt auth
	decoder.Decode(&txt)
	if txt.Login == "" || txt.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hasher := md5.New()
	hasher.Write([]byte(txt.Password))
	hash := hex.EncodeToString(hasher.Sum(nil))

	user, err := h.dm.GetUser(ctx, &datamanager.User{Login: txt.Login})
	if err != nil {
		err = h.dm.AddUser(ctx, &datamanager.User{Login: txt.Login, Password: hash})
		if err != nil {
			log.Println(err)
		}
	} else {
		if user.Password != hash {
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

	cookie := createCookie(jwtCookieName, t)
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"status\": \"success\"}"))
}

func getLogin(r *http.Request) (string, error) {
	cookie, err := r.Cookie(jwtCookieName)
	if err != nil || cookie == nil {
		return "", err
	}
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})
	if err != nil {
		return "", err
	}

	userLogin := claims["login"]

	return userLogin.(string), nil
}

func createCookie(name, value string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
}
