package userhandler

import (
	"encoding/json"
	"net/http"
	"strconv"

	usermodel "assignment8/models/user"
)

type UserServiceInterface interface {
	CreateUser(user *usermodel.User) error
	GetUser(id int) (*usermodel.User, error)
}

type UserHandler struct {
	userService UserServiceInterface
}

func NewUserHandler(service UserServiceInterface) *UserHandler {
	return &UserHandler{
		userService: service,
	}
}

func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user usermodel.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid Body", http.StatusBadRequest)
		return
	}

	err = h.userService.CreateUser(&user)
	if err != nil {
		http.Error(w, "Invalid Body", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Unable to Find User", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUser(id)
	if err != nil {
		http.Error(w, "Unable to Find User", http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(&user)
	if err != nil {
		http.Error(w, "Unable to Find User", http.StatusInternalServerError)
		return
	}
}
