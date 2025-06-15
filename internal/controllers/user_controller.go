package controllers

import (
	"encoding/json"
	"net/http"

	models "github.com/coworker-match-api/gen/go"
	"github.com/coworker-match-api/internal/common"
	"github.com/coworker-match-api/internal/usecases"
)

type UserController struct {
	uu usecases.IUserUsecase
}

func NewUserController(uu usecases.IUserUsecase) *UserController {
	return &UserController{uu: uu}
}

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := common.ExtractUserIdFromContext(r)
	if err != nil {
		common.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	createdUser, err := uc.uu.CreateUser(userId, &req)
	if err != nil {
		common.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	common.RespondWithJSON(w, http.StatusOK, createdUser)
}

func (uc *UserController) GetUserById(w http.ResponseWriter, r *http.Request) {
	userId, err := common.ExtractPathParam(r, "userId")
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := uc.uu.GetUserById(userId)
	if err != nil {
		common.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	common.RespondWithJSON(w, http.StatusOK, user)
}

func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userId, err := common.ExtractPathParam(r, "userId")
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req models.UpdateUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := uc.uu.UpdateUser(userId, &req)
	if err != nil {
		common.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	common.RespondWithJSON(w, http.StatusOK, user)
}

func (uc *UserController) IsUserExist(w http.ResponseWriter, r *http.Request) {
	userId, err := common.ExtractUserIdFromContext(r)
	if err != nil {
		common.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	isExist, err := uc.uu.IsUserExist(userId)
	if err != nil {
		common.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	common.RespondWithJSON(w, http.StatusOK, map[string]bool{"exists": isExist})
}
