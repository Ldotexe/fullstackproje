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

type messageResponse struct {
	IsMine bool   `json:"is_mine"`
	Text   string `json:"text"`
}

type send struct {
	Text string `json:"text"`
}

type auth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Handler struct {
	dm datamanager.DataManager
}

func New(dm datamanager.DataManager) *Handler {
	return &Handler{dm: dm}
}

func (h *Handler) UsersHandler(w http.ResponseWriter, _ *http.Request) {
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

	body, err := json.Marshal(logins)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(body)
	return
}

func (h *Handler) MessagesHandler(w http.ResponseWriter, r *http.Request) {
	userLogin, err := GetLogin(r)
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

func (h *Handler) SendHandler(w http.ResponseWriter, r *http.Request) {
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

	h.dm.AddMessage(ctx, &datamanager.Message{FromLogin: userLogin, ToLogin: toLogin, Text: txt.Text, Ts: time.Now().Format(time.RFC3339Nano)})
	//TODO: handle err and evrth

}

func (h *Handler) AuthHandler(w http.ResponseWriter, r *http.Request) {
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

	cookie := http.Cookie{
		Name:     "ssjwt",
		Value:    t,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("cookie set!"))
}

func GetLogin(r *http.Request) (string, error) {
	cookie, err := r.Cookie("ssjwt")
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

	// do something with decoded claims
	userLogin := claims["login"]

	return userLogin.(string), nil
}
