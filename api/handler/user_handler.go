package handler

import (
	"encoding/json"
	"go_msg_broker/models"
	"go_msg_broker/utils"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type UserHandler struct {
	db *pgx.Conn
}

func NewUserHandler(db *pgx.Conn) *UserHandler {
	return &UserHandler{
		db: db,
	}
}

func (u *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req models.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponsWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err := u.db.Exec(r.Context(), "INSERT INTO users (name, password) VALUES ($1,$2)", req.Name, req.Password)
	if err != nil {
		utils.ResponsWithError(w, http.StatusInternalServerError, "gagal insert ke database")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, req)
}
