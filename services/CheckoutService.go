package services

import (
	"eCommerce/helpers"
	"eCommerce/models"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type CheckoutService struct {
	DB *sqlx.DB
	cs CartService
	ss ShippingService
}

func NewCheckoutService(db *sqlx.DB, cs CartService, ss ShippingService) *CheckoutService {
	return &CheckoutService{
		db,
		cs,
		ss,
	}
}

func (chs *CheckoutService) OrderSummary(cartId, userId string) (*models.Summary, error) {
	cartSummary, err := chs.cs.ViewCart(cartId, userId)
	if err != nil {
		return nil, err
	}

	defaultShippingAddress, err := chs.ss.GetDefaultShippingAddress(userId)
	if err != nil {
		return nil, err
	}

	result := &models.Summary{ReturnedCart: cartSummary, ShippingAddress: defaultShippingAddress}

	return result, nil
}

func (chs *CheckoutService) ConfirmPurchase(cartId, userId, userEmail string) (*models.Summary, error) {
	tx, err := chs.DB.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	summary, err := chs.OrderSummary(cartId, userId)
	if err != nil {
		return nil, err
	}

	cart := summary.ReturnedCart

	if len(cart.Items) == 0 {
		return nil, fmt.Errorf("your cart is empty, please add some items")
	}

	var newOrderId string
	createOrderQuery := `INSERT INTO orders (cartid, userid, status, total_price)
	VALUES ($1, $2, $3, $4)
	RETURNING orderid
	`
	err = tx.Get(&newOrderId, createOrderQuery, cartId, userId, "processing", cart.Total)
	if err != nil {
		return nil, fmt.Errorf("error creating a new order: %w", err)
	}

	// Collect all product IDs from cart items
	productIDs := make([]string, 0, len(cart.Items))
	for _, item := range cart.Items {
		productIDs = append(productIDs, item.ProductId)
	}

	// Lock all products in one query
	// sqlx.In: it takes your query with a placeholder ? and a slice of values, and expands that slice into the proper number of placeholders.
	// So, before rebinding, the query string would look like:
	// SELECT productid, stock FROM products WHERE productid IN (?, ?, ?, ...) FOR UPDATE, with as many ? placeholders as there are product IDs in your productIDs slice.
	// After Rebind, it converts those ? placeholders to whatever the database driver expects $1, $2, $3, ...
	query, args, err := sqlx.In(`SELECT productid, stock FROM products WHERE productid IN (?) FOR UPDATE`, productIDs)
	if err != nil {
		return nil, fmt.Errorf("error preparing stock lock query: %w", err)
	}
	query = tx.Rebind(query)

	rows, err := tx.Queryx(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error locking product stocks: %w", err)
	}
	defer rows.Close()

	// Map productid -> current stock
	stockMap := make(map[string]int)
	for rows.Next() {
		var pid string
		var stock int
		if err := rows.Scan(&pid, &stock); err != nil {
			return nil, fmt.Errorf("error scanning locked stock row: %w", err)
		}
		stockMap[pid] = stock
	}

	// Check stock availability
	for _, item := range cart.Items {
		currentStock, exists := stockMap[item.ProductId]
		if !exists {
			return nil, fmt.Errorf("product %s no longer exists", item.Name)
		}
		if currentStock < item.Quantity {
			return nil, fmt.Errorf("not enough stock for %s", item.Name)
		}
	}

	// Update stock for each product
	for _, item := range cart.Items {
		updateQuery := `UPDATE products SET stock = stock - $1 WHERE productid = $2 AND stock >= $1`
		res, err := tx.Exec(updateQuery, item.Quantity, item.ProductId)
		if err != nil {
			return nil, fmt.Errorf("error updating stock for %s: %w", item.Name, err)
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("error checking update result for %s: %w", item.Name, err)
		}
		if rowsAffected == 0 {
			return nil, fmt.Errorf("failed to update stock for %s, possibly insufficient stock", item.Name)
		}
	}
	// insert into order_items
	for _, item := range cart.Items {
		insertQuery := `INSERT INTO order_items (orderid, productid, quantity) VALUES($1, $2, $3)`
		_, err := tx.Exec(insertQuery, newOrderId, item.ProductId, item.Quantity)
		if err != nil {
			return nil, fmt.Errorf("error inserting order items: %w", err)
		}
	}

	err = chs.cs.DeleteCart(cartId, userId)
	if err != nil {
		return nil, fmt.Errorf("error emptying cart: %w", err)
	}

	err = helpers.SendConfirmationEmail(newOrderId, userEmail, cart.Total, time.Now())
	if err != nil {
		return nil, fmt.Errorf("error sending confirmation email: %w", err)
	}

	return summary, nil
}
