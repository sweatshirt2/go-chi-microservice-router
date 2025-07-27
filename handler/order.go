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
	"github.com/sweatshirt2/go-analytics/models"
	repository "github.com/sweatshirt2/go-analytics/repositories"
)


const decimal = 10
const bitSize = 64

type OrderController struct {
	Repo *repository.OrderRepo
}

func (o *OrderController) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CustomerId uuid.UUID `json:"customer_id"`
		Items []models.Item	`json:"items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		println("error parsing json...")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()

	order := models.Order{
		OrderId: rand.Uint64(),
		CreatedAt: &now,
		CustomerId: body.CustomerId,
		Items: body.Items,
	}

	err := o.Repo.Insert(r.Context(), order)

	if err != nil {
		fmt.Println("failed to save order")
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}

	res, err := json.Marshal(order)

	if err != nil {
		fmt.Println("failed to marshall order")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(res)
	w.WriteHeader(http.StatusCreated)
}

func (o *OrderController) GetAll(w http.ResponseWriter, r *http.Request) {
	cursorStr := r.URL.Query().Get("cursor")

	if cursorStr == "" {
		cursorStr = "0"
	}

	cursor, err := strconv.ParseUint(cursorStr, decimal, bitSize)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	const size = 50
	res, err := o.Repo.GetAll(r.Context(), repository.FindAllPage{
		Offset: cursor,
		Size: size,
	})

	if err != nil {
		fmt.Println("failed to find orders, %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var response struct {
		Orders[]models.Order `json:"orders"`
		Next uint64
	}
	response.Orders = res.Orders
	response.Next = res.Cursor

	data, err := json.Marshal(response)

	if err != nil {
		fmt.Println("failed to marshall orders: %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func (o *OrderController) GetById(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	orderId, err := strconv.ParseUint(idParam, decimal, bitSize)

	if err != nil {
		fmt.Println("failed to parse id: %w", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	order, err := o.Repo.FindById(r.Context(), orderId)

	if err != nil {
		fmt.Println("failed to get order: %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(order)

	if err != nil {
		fmt.Println("failed to marshal order: %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func (o *OrderController) Update(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		println("error parsing json...")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	idParam := chi.URLParam(r, "id")
	orderId, err := strconv.ParseUint(idParam, decimal, bitSize)

	if err != nil {
		println("error parsing id")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	order, err := o.Repo.FindById(r.Context(), orderId)

	if errors.Is(err, repository.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
	} else if err != nil {
		fmt.Println("error getting order: %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	now := time.Now().UTC()

	switch body.Status {
	case "shipped":
		if order.ShippedAt != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("order already shipped"))
			return
		}
		order.ShippedAt = &now
	case "completed":
		if order.ShippedAt != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("order already shipped"))
			return
		}
		if order.CompletedAt != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("order already completed"))
			return
		}
		order.CompletedAt = &now
	}

	err = o.Repo.Update(r.Context(), order)

	if err != nil {
		fmt.Println("error updating order: %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}

func (o *OrderController) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	orderId, err := strconv.ParseUint(idParam, decimal, bitSize)

	if err != nil {
		println("error parsing id")
		w.WriteHeader(http.StatusUnprocessableEntity)
	}

	if err := o.Repo.Delete(r.Context(), orderId); errors.Is(err, repository.ErrNotExist) {
		fmt.Println("error deleting order: %w", err)
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("error deleting order: %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}

// ? Context for next refactor

// type ControllerParams struct {
// 	w http.ResponseWriter
// 	r *http.Request
// }

// func wrapControllers(w http.ResponseWriter, r *http.Request) {
// 	o := OrderController{}
// 	o.Create(ControllerParams{w: w, r: r})
// }

// func (o *OrderController) Create(args ControllerParams) {
// 	println("creating order...")
// }
