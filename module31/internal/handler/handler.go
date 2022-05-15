package handler

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"skillbox/internal/entity"
	"skillbox/internal/usecase"
)

type Handler struct {
	useCase usecase.Usecase
}

func NewHandler(useCase usecase.Usecase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (c *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &entity.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		log.Println(err)
		buildResponse(w, http.StatusInternalServerError, nil)
		return
	}
	id, err := c.useCase.CreateUser(user)
	if err != nil {
		log.Println(err)
		buildResponse(w, http.StatusInternalServerError, nil)
		return
	}
	result := map[string]string{"id": id}
	response, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
		buildResponse(w, http.StatusInternalServerError, nil)
		return
	}
	buildResponse(w, http.StatusCreated, response)
}

func (c *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := &entity.Id{}
	err := json.NewDecoder(r.Body).Decode(&userId)
	if err != nil {
		log.Println(err)
		buildResponse(w, http.StatusInternalServerError, nil)
		return
	}
	name, err := c.useCase.DeleteUser(userId.TargetId)
	if err != nil {
		log.Println(err)
		buildResponse(w, http.StatusInternalServerError, nil)
		return
	}

	response, err := json.Marshal("User " + name + " was deleted")
	if err != nil {
		log.Println(err)
		buildResponse(w, http.StatusInternalServerError, nil)
		return
	}
	buildResponse(w, http.StatusOK, response)
}

func (c *Handler) GetUsers(w http.ResponseWriter, _ *http.Request) {
	user := &entity.User{}
	result := c.useCase.GetUsers(user)
	response, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
		buildResponse(w, http.StatusInternalServerError, nil)
		return
	}
	buildResponse(w, http.StatusCreated, response)

}

func (c *Handler) UpgradeUser(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "id")
	id := params
	upgradeUser := &entity.UserUpgrade{}
	err := json.NewDecoder(r.Body).Decode(&upgradeUser)
	if err != nil {
		log.Println(err)
		buildResponse(w, http.StatusInternalServerError, nil)
		return
	}
	response, err := json.Marshal("Age update successful")
	err = c.useCase.UpdateUser(id, upgradeUser.NewAge)
	if err != nil {
		log.Println(err)
		buildResponse(w, http.StatusInternalServerError, nil)
		return
	}
	buildResponse(w, http.StatusOK, response)

}

func (c *Handler) MakeFriends(w http.ResponseWriter, r *http.Request) {
	user := &entity.Id{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		buildResponse(w, http.StatusInternalServerError, nil)
		return
	}
	name1, name2, err := c.useCase.MakeFriends(user.TargetId, user.SourceId)
	if err != nil {
		log.Println(err)
		buildResponse(w, http.StatusInternalServerError, nil)
		return
	}

	response, err := json.Marshal(name1 + " and " + name2 + " now friends")
	if err != nil {
		log.Println(err)
		buildResponse(w, http.StatusInternalServerError, nil)
		return
	}
	buildResponse(w, http.StatusOK, response)
}

func (c *Handler) GetFriends(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "id")
	id := params
	friends, err := c.useCase.GetFriends(id)
	if err != nil {
		log.Println(err)
		buildResponse(w, http.StatusInternalServerError, nil)
		return
	}
	response, err := json.Marshal(friends)

	buildResponse(w, http.StatusOK, response)

}

func buildResponse(w http.ResponseWriter, statusCode int, body []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_, err := w.Write(body)
	if err != nil {
		return
	}
}

func Build(router *chi.Mux, useCase usecase.Usecase) {
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	controller := NewHandler(useCase)
	router.Post("/create", controller.CreateUser)
	router.Delete("/user", controller.DeleteUser)
	router.Get("/users", controller.GetUsers)
	router.Put("/{id}", controller.UpgradeUser)
	router.Post("/make_friends", controller.MakeFriends)
	router.Get("/friends/{id}", controller.GetFriends)
}
