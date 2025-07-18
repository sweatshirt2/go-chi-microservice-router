package handler

import "net/http"


type OrderController struct {}

func (o *OrderController) Create(w http.ResponseWriter, r *http.Request) {
	println("creating order...")
}

func (o *OrderController) GetAll(w http.ResponseWriter, r *http.Request) {
	println("getting all orders...")
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
