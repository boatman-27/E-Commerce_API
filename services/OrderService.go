package services

import (
	"database/sql"
	"eCommerce/models"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type OrderService struct {
	DB *sqlx.DB
}

func NewOrderService(db *sqlx.DB) *OrderService {
	return &OrderService{
		db,
	}
}

func (os *OrderService) checkOrderOwnership(orderId, userId string) error {
	var ownerId string
	err := os.DB.QueryRow("SELECT userid FROM orders WHERE orderid = $1", orderId).Scan(&ownerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrNotFound
		}
		return fmt.Errorf("error checking order issuer: %w", err)
	}
	if ownerId != userId {
		return models.ErrUnauthorized
	}
	return nil
}

func (os *OrderService) TrackOrder(orderId, userId string) (string, error) {
	var status string
	query := `SELECT status FROM orders WHERE orderid = $1 AND userid = $2`
	err := os.DB.Get(&status, query, orderId, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", models.ErrNotFound
		}
		return "", fmt.Errorf("error fetching order status: %w", err)
	}
	return status, nil
}

func (os *OrderService) ViewOrders(userId string) ([]*models.Order, error) {
	query := `
		SELECT
			o.orderid,
			o.total_price,
			o.orderedat,
			p.productid,
			p.name,
			p.description,
			p.brand,
			p.category,
			oi.quantity
		FROM orders o
		JOIN order_items oi ON o.orderid = oi.orderid
		JOIN products p ON p.productid = oi.productid
		WHERE o.userid = $1
		ORDER BY o.orderedat DESC
	`
	type row struct {
		OrderId     string    `db:"orderid"`
		TotalPrice  float64   `db:"total_price"`
		OrderedAt   time.Time `db:"orderedat"`
		ProductId   string    `db:"productid"`
		Name        string    `db:"name"`
		Description string    `db:"description"`
		Brand       string    `db:"brand"`
		Category    string    `db:"category"`
		Quantity    int       `db:"quantity"`
	}

	var rows []row
	if err := os.DB.Select(&rows, query, userId); err != nil {
		return nil, fmt.Errorf("error fetching orders: %w", err)
	}

	orderMap := make(map[string]*models.Order)
	for _, r := range rows {
		if _, exists := orderMap[r.OrderId]; !exists {
			orderMap[r.OrderId] = &models.Order{
				OrderId:    r.OrderId,
				TotalPrice: r.TotalPrice,
				OrderedAt:  r.OrderedAt,
				Items:      []*models.OrderItem{},
			}
		}

		orderMap[r.OrderId].Items = append(orderMap[r.OrderId].Items, &models.OrderItem{
			ProductId:   r.ProductId,
			Name:        r.Name,
			Description: r.Description,
			Brand:       r.Brand,
			Category:    r.Category,
			Quantity:    r.Quantity,
		})
	}

	orders := make([]*models.Order, 0, len(orderMap))
	for _, o := range orderMap {
		orders = append(orders, o)
	}
	return orders, nil
}
