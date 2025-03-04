package orchestrator

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/calculate", h.AddExpression).Methods("POST")
	router.HandleFunc("/api/v1/expressions", h.GetExpressions).Methods("GET")
	router.HandleFunc("/api/v1/expressions/{id}", h.GetExpressionByID).Methods("GET")
	router.HandleFunc("/internal/task", h.GetTask).Methods("GET")
	router.HandleFunc("/internal/task", h.SubmitResult).Methods("POST")
}

// Добавление нового выражения
func (h *Handler) AddExpression(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusUnprocessableEntity)
		return
	}

	id, err := h.Service.AddExpression(req.Expression)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

// Получение списка выражений
func (h *Handler) GetExpressions(w http.ResponseWriter, r *http.Request) {
	expressions := h.Service.GetExpressions()
	json.NewEncoder(w).Encode(map[string]interface{}{"expressions": expressions})
}

// Получение выражения по ID
func (h *Handler) GetExpressionByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	expression, found := h.Service.GetExpressionByID(id)
	if !found {
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"expression": expression})
}

// Агент запрашивает задачу
func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	task, found := h.Service.GetNextTask()
	if !found {
		http.Error(w, "No available tasks", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"task": task})
}

func (h *Handler) SubmitResult(w http.ResponseWriter, r *http.Request) {
    var req struct {
        ID     string  `json:"id"`
        Result float64 `json:"result,omitempty"`
        Error  string  `json:"error,omitempty"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusUnprocessableEntity)
        return
    }

    if req.Error != "" {
        if err := h.Service.SubmitTaskError(req.ID, req.Error); err != nil {
            http.Error(w, "Failed to submit error", http.StatusInternalServerError)
            return
        }
    } else {
        if err := h.Service.SubmitTaskResult(req.ID, req.Result); err != nil {
            http.Error(w, "Failed to submit result", http.StatusInternalServerError)
            return
        }
    }

    w.WriteHeader(http.StatusOK)
}