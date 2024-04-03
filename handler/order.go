package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/Saywa94/go_api/model"
	"github.com/Saywa94/go_api/repository/order"
)

type Order struct {
	Repo *order.RedisRepo
}

// Creates and order
func (h *Order) Create(w http.ResponseWriter, r *http.Request) {

	var body struct {
		CustomerID uuid.UUID        `json:"customer_id"`
		LineItems  []model.LineItem `json:"line_items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		fmt.Println("Failed to decode body: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()

	order := model.Order{
		OrderID:    rand.Uint64(),
		CustomerID: body.CustomerID,
		LineItems:  body.LineItems,
		CreatedAt:  &now,
	}

	err := h.Repo.Insert(r.Context(), order)
	if err != nil {
		fmt.Println("Failed to create order: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(order)
	if err != nil {
		fmt.Println("Failed to encode order: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(res)
	w.WriteHeader(http.StatusCreated)
	return

}

// List all orders
func (h *Order) List(w http.ResponseWriter, r *http.Request) {
	cursorStr := r.URL.Query().Get("cursor")
	if cursorStr == "" {
		cursorStr = "0"
	}
	const base = 10
	const bitSize = 64
	cursor, err := strconv.ParseUint(cursorStr, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	const size = 50
	res, err := h.Repo.FindAll(r.Context(), order.FindAllPage{Size: size, Offset: cursor})
	if err != nil {
		fmt.Println("Failed to finda all: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var response struct {
		Items []model.Order `json:"items"`
		Next  uint64        `json:"next,omitempty"`
	}
	response.Items = res.Orders
	response.Next = res.Cursor

	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Failed to Marshall, ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
	w.WriteHeader(http.StatusOK)
	return

}

// Get an order by ID
func (h *Order) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	orderId, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		fmt.Println("Failed to parse id: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	o, err := h.Repo.FindByID(r.Context(), orderId)
	if errors.Is(err, order.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Encoder already writes json encoded string to stream
	if err := json.NewEncoder(w).Encode(o); err != nil {
		fmt.Println("Failed to encode order: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

// Update an order by ID
func (h *Order) UpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update an order by ID")
}

// Delete and order by ID
func (h *Order) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete and order by ID")
}
