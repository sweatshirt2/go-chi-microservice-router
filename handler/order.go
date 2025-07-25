package handler

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/sweatshirt2/go-analytics/models"
	repository "github.com/sweatshirt2/go-analytics/repositories"
)


type OrderController struct {
	Repo *repository.OrderRepo
}

func (o *OrderController) Create(w http.ResponseWriter, r *http.Request) {
	println("creating order...")
	var body struct {
		CustomerId uuid.UUID `json:"customer_id"`
		Items []models.Item	`json:"items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		print("error parsing json...")
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
	println("getting all orders...")
	cursorStr := r.URL.Query().Get("cursor")

	if cursorStr == "" {
		cursorStr = "0"
	}

	const decimal = 10
	const bitSize = 64

	cursor, err := strconv.ParseUint(cursorStr, decimal, bitSize)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	const size = 50

	print(cursor)

	w.WriteHeader(http.StatusOK)

	var body struct {
		Name string `json:"name"`
	}

	body.Name = "abebe"

	val, err := json.Marshal(body)

	if err != nil {
		fmt.Println("failed to marshall orders")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(val)
}

func (o *OrderController) GetById(w http.ResponseWriter, r *http.Request) {
	println("getting order by id...")
}

func (o *OrderController) Update(w http.ResponseWriter, r *http.Request) {
	println("updating order...")
}

func (o *OrderController) Delete(w http.ResponseWriter, r *http.Request) {
	println("deleting order...")
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
