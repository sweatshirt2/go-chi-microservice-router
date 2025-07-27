package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/sweatshirt2/go-analytics/models"
)

var ErrNotExist = errors.New("Order does not exist")

type OrderRepo struct {
	Client *redis.Client
}

func OrderIdKey(id uint64) string {
	return fmt.Sprintf("order:%d", id)
}

func (r *OrderRepo) Insert(ctx context.Context, order models.Order) error {
	data, err := json.Marshal(order)

	if err != nil {
		return fmt.Errorf("failed to encode order: %w", err)
	}

	key := OrderIdKey(order.OrderId)
	tx := r.Client.TxPipeline()

	res := tx.SetNX(ctx, key, string(data), 0)
	if err := res.Err(); err != nil {
		tx.Discard()
		return fmt.Errorf("failed to set order: %w", err)
	}

	if err = tx.SAdd(ctx, "orders", key).Err(); err != nil {
		tx.Discard()
		return fmt.Errorf("failed to add to orders list, %w", err)
	}

	if _, err = tx.Exec(ctx); err != nil {
		return fmt.Errorf("failed to finish the process of adding order, %w", err)
	}

	return nil
}

func (r *OrderRepo) FindById(ctx context.Context, id uint64) (models.Order, error) {
	key := OrderIdKey(id)

	value, err := r.Client.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return models.Order{}, ErrNotExist
	}

	if err != nil {
		return models.Order{}, fmt.Errorf("Get order: %w", err)
	}

	var order models.Order
	err = json.Unmarshal([]byte(value), &order)

	if err != nil {
		return models.Order{}, fmt.Errorf("Error decoding order: %w", err)
	}

	return order, nil
}

func (r *OrderRepo) Delete(ctx context.Context, id uint64) error {
	key := OrderIdKey(id)
	tx := r.Client.TxPipeline()

	// if count, err := tx.Del(ctx, key).Result(); count == 0 {
	// 	tx.Discard()
	// 	return ErrNotExist
	// } else if err != nil {
	// 	tx.Discard()
	// 	return fmt.Errorf("error deleting order: %w", err)
	// }

	if _, err := tx.Del(ctx, key).Result(); err != nil {
		tx.Discard()
		return fmt.Errorf("error deleting order: %w", err)
	}

	if err := tx.SRem(ctx, "orders", key).Err(); err != nil {
		tx.Discard()
		return fmt.Errorf("failed to remove order from the list: %w", err)
	}

	if _, err := tx.Exec(ctx); err != nil {
		tx.Discard()
		return fmt.Errorf("failed to finish the process of removing order: %w", err)
	}

	return nil
}

func (r *OrderRepo) Update(ctx context.Context, order models.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to encode order: %w", err)
	}

	key := OrderIdKey(order.OrderId)

	err = r.Client.SetXX(ctx, key, string(data), 0).Err()

	if errors.Is(err, redis.Nil) {
		return ErrNotExist
	} else if err != nil {
		return fmt.Errorf("error updating order: %w", err)
	}

	return nil
}

type FindAllPage struct {
	Size uint
	Offset uint64
}

type FindResult struct {
	Orders []models.Order
	Cursor uint64
}

func (r *OrderRepo) GetAll(ctx context.Context, page FindAllPage) (FindResult, error) {
	res := r.Client.SScan(ctx, "orders", uint64(page.Offset), "*", int64(page.Size))
	keys, cursor, err := res.Result()

	if err != nil {
		return FindResult{}, nil
	}

	if len(keys) == 0 {
		return FindResult{
			Orders: []models.Order{},
		}, nil
	}

	xs, err := r.Client.MGet(ctx, keys...).Result()

	if err != nil {
		return FindResult{}, fmt.Errorf("error getting list of orders, %w", err)
	}

	orders := make([]models.Order, len(xs))

	for i, x := range xs {
		x := x.(string)
		var order models.Order

		err := json.Unmarshal([]byte(x), &order)
		if err != nil {
			return FindResult{}, fmt.Errorf("error decoding order at index %d: %w", i, err)
		}

		orders[i] = order
	}

	return FindResult{Orders: orders, Cursor: cursor}, nil
}
