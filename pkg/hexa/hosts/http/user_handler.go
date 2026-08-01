package httphost

import (
    "encoding/json"
    "net/http"

    "hexa/domain"
)

type UserHandler struct {
    useCase domain.UserUseCase
}

func NewUserHandler(uc domain.UserUseCase) *UserHandler {
    return &UserHandler{useCase: uc}
}

// Request and Response DTOs
type UserResponse struct {
    ID    string `json:"id"`
    Email string `json:"email"`
    Name  string `json:"name"`
}

type CreateUserRequest struct {
    Email string `json:"email"`
    Name  string `json:"name"`
}

type UpdateUserRequest struct {
    Email string `json:"email"`
    Name  string `json:"name"`
}

func mapUserToResponse(u domain.User) UserResponse {
    return UserResponse{
        ID:    u.ID,
        Email: u.Email,
        Name:  u.Name,
    }
}

// RegisterRoutes registers the CRUD routes on a ServeMux (Go 1.22+)
func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
    mux.HandleFunc("POST /users", h.HandleCreate)
    mux.HandleFunc("GET /users/{id}", h.HandleGet)
    mux.HandleFunc("GET /users", h.HandleList)
    mux.HandleFunc("PUT /users/{id}", h.HandleUpdate)
    mux.HandleFunc("DELETE /users/{id}", h.HandleDelete)
}

func (h *UserHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "bad request", http.StatusBadRequest)
        return
    }

    user, err := h.useCase.RegisterUser(req.Email, req.Name)
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnprocessableEntity)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(mapUserToResponse(user))
}

func (h *UserHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    user, err := h.useCase.GetUser(id)
    if err != nil {
        if err == domain.ErrUserNotFound {
            http.Error(w, "not found", http.StatusNotFound)
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(mapUserToResponse(user))
}

func (h *UserHandler) HandleList(w http.ResponseWriter, r *http.Request) {
    users, err := h.useCase.ListUsers()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    var res []UserResponse
    for _, u := range users {
        res = append(res, mapUserToResponse(u))
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(res)
}

func (h *UserHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    var req UpdateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "bad request", http.StatusBadRequest)
        return
    }

    user, err := h.useCase.UpdateUser(id, req.Email, req.Name)
    if err != nil {
        if err == domain.ErrUserNotFound {
            http.Error(w, "not found", http.StatusNotFound)
            return
        }
        http.Error(w, err.Error(), http.StatusUnprocessableEntity)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(mapUserToResponse(user))
}

func (h *UserHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    err := h.useCase.DeleteUser(id)
    if err != nil {
        if err == domain.ErrUserNotFound {
            http.Error(w, "not found", http.StatusNotFound)
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}
