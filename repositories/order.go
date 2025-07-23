package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/sweatshirt2/go-analytics/models"
)

type OrderRepo struct {
	Client *redis.Client

}

func orderIdKey(id uint64) string {
	return fmt.Sprintf("order:%d", id)
}

func (r *OrderRepo) Insert(ctx context.Context, order models.Order) error {
	data, err := json.Marshal(order)

	if err != nil {
		return fmt.Errorf("failed to encode order: %w", err)
	}

	key := orderIdKey(order.OrderId)

	res := r.Client.SetNX(ctx, key, string(data), 0)
	if err := res.Err(); err != nil {
		return fmt.Errorf("failed to set order: %w", err)
	}

	return nil
}

func (r *OrderRepo) FindById(ctx context.Context, id uint64) (models.Order, error) {
	key := orderIdKey(id)

	value, err := r.Client.Get(ctx, key).Result()
	
	if errors.Is(err, redis.Nil) {
		return models.Order{}, errors.New("Order does not exist")
	}

	if err != nil {
		return models.Order{}, fmt.Errorf("Get order: %w", err)
	}

	var order models.Order
	err = json.Unmarshal([]byte(value), &order)

	if err != nil {
		return models.Order{}, fmt.Errorf("Get order: %w", err)
	}

	return order, nil
}
