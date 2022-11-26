package user

import (
	"authService/internal/transport/httpv1"
	"authService/internal/transport/inputs"
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"net/http"
)

type AuthService interface {
	SignUp(ctx context.Context, input *inputs.CreateInput) (int, error)
	SignIn(ctx context.Context, input *inputs.AuthInput) (string, error)
	IsAuthed(ctx context.Context, token string) (int, error)
}
type authHandler struct {
	as AuthService
}

func NewAuthHandler(as AuthService) httpv1.Handler {
	return &authHandler{as: as}
}

func (h *authHandler) Register(router *httprouter.Router) {
	router.POST("/auth/sign-up", h.SignUp)
	router.POST("/auth/sign-in", h.SignIn)
}

func (h *authHandler) SignUp(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var input inputs.CreateInput

	all, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Reading error")
		return
	}
	err = json.Unmarshal(all, &input)
	if err != nil {
		log.Println("Unmarshall error")
		return
	}
	UserID, err := h.as.SignUp(r.Context(), &input)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp := make(map[string]int)
	resp["UserID"] = UserID
	response, err := json.Marshal(resp)
	if err != nil {
		return
	}
	_, err = w.Write(response)
	if err != nil {
		return
	}

}

func (h *authHandler) SignIn(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var input inputs.AuthInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	token, err := h.as.SignIn(r.Context(), &input)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp := make(map[string]string)
	resp["JWT"] = token
	response, err := json.Marshal(resp)
	if err != nil {
		return
	}
	_, err = w.Write(response)
	if err != nil {
		return
	}
}
