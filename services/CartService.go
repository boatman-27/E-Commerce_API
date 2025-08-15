package services

import (
	"database/sql"
	"eCommerce/models"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type CartService struct {
	DB *sqlx.DB
}

func NewCartService(db *sqlx.DB) *CartService {
	return &CartService{
		db,
	}
}

func (cs *CartService) checkOwnership(cartId, userId string) error {
	var ownerId string
	err := cs.DB.QueryRow("SELECT userid FROM carts WHERE cartid = $1", cartId).Scan(&ownerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrNotFound
		}
		return fmt.Errorf("error checking cart ownership: %w", err)
	}
	if ownerId != userId {
		return models.ErrUnauthorized
	}
	return nil
}

func (cs *CartService) getAvailableStock(productId string) (int, error) {
	var availableStock int
	err := cs.DB.Get(&availableStock, `SELECT stock FROM products WHERE productid = $1`, productId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrNotFound
		}
		return 0, fmt.Errorf("error fetching stock for product: %w", err)
	}
	return availableStock, nil
}

func (cs *CartService) getExistingQty(cartId, productId string) (int, error) {
	var existingQty int
	err := cs.DB.Get(&existingQty,
		`SELECT quantity FROM cart_items WHERE cartid = $1 AND productid = $2`,
		cartId, productId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("error fetching existing quantity: %w", err)
	}
	return existingQty, nil
}

func (cs *CartService) getTotal(db sqlx.Ext, cartId, userId string) (float64, error) {
	var total float64
	err := sqlx.Get(db, &total, `
		SELECT COALESCE(SUM(
			(p.price - COALESCE(p.discount, 0)) * ci.quantity
		), 0) as total
		FROM cart_items ci
		JOIN products p ON p.productid = ci.productid
		JOIN carts c ON c.cartid = ci.cartid
		WHERE ci.cartid = $1 AND c.userid = $2
	`, cartId, userId)
	if err != nil {
		return 0, fmt.Errorf("error calculating total: %w", err)
	}
	return total, nil
}

func (cs *CartService) DeleteCart(cartId, userId string) error {
	_, err := cs.DB.Exec("DELETE FROM cart_items WHERE cartid = $1 AND userid = $2", cartId, userId)
	if err != nil {
		return fmt.Errorf("failed to delete cart: %w", err)
	}
	return nil
}

func (cs *CartService) ViewCart(cartId, userId string) (*models.ReturnedCart, error) {
	if err := cs.checkOwnership(cartId, userId); err != nil {
		return nil, err
	}

	var cartItems []*models.ItemData
	fetchQuery := `
	SELECT
    p.productid,
    p.name,
    p.description,
    p.sku,
    p.price,
    p.discount,
    p.brand,
    p.category,
    ci.quantity
		FROM products p
		JOIN cart_items ci ON p.productid = ci.productid
		WHERE ci.cartid = $1`
	err := cs.DB.Select(&cartItems, fetchQuery, cartId)
	if err != nil {
		return nil, fmt.Errorf("error fetching items data: %w", err)
	}

	total, err := cs.GetTotal(cartId, userId)
	if err != nil {
		return nil, fmt.Errorf("error calculating total: %w", err)
	}

	cart := &models.ReturnedCart{
		ID:     cartId,
		UserId: userId,
		Items:  cartItems,
		Total:  total,
	}

	return cart, nil
}

func (cs *CartService) GetTotal(cartId, userId string) (float64, error) {
	return cs.getTotal(cs.DB, cartId, userId)
}

// BUG  returns old total
func (cs *CartService) AddToCart(cartId, userId string, cartItem *models.CartItem) (*models.ReturnedCart, error) {
	if err := cs.checkOwnership(cartId, userId); err != nil {
		return nil, err
	}

	tx, err := cs.DB.Beginx()
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

	availableStock, err := cs.getAvailableStock(cartItem.ProductId)
	if err != nil {
		return nil, err
	}

	existingQty, err := cs.getExistingQty(cartId, cartItem.ProductId)
	if err != nil {
		return nil, err
	}

	if cartItem.Quantity+existingQty > availableStock {
		return nil, fmt.Errorf("only %d items available in stock", availableStock-existingQty)
	}

	// res, err := tx.Exec(`UPDATE products SET stock = stock - $1 WHERE productid = $2 AND stock >= $1`, cartItem.Quantity, cartItem.ProductId)
	// if err != nil {
	// 	return nil, fmt.Errorf("error updating stock: %w", err)
	// }
	// if rows, _ := res.RowsAffected(); rows == 0 {
	// 	return nil, fmt.Errorf("not enough stock available right now")
	// }

	_, err = tx.Exec(`
		INSERT INTO cart_items (cartid, productid, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (cartid, productid)
		DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity
	`, cartId, cartItem.ProductId, cartItem.Quantity)
	if err != nil {
		return nil, fmt.Errorf("error adding new item: %w", err)
	}

	var cartItems []*models.ItemData
	err = tx.Select(&cartItems, `
		SELECT
			p.productid,
			p.name,
			p.description,
			p.sku,
			p.price,
			p.discount,
			p.brand,
			p.category,
			ci.quantity
		FROM products p
		JOIN cart_items ci ON p.productid = ci.productid
		WHERE ci.cartid = $1
	`, cartId)
	if err != nil {
		return nil, fmt.Errorf("error fetching items data: %w", err)
	}

	total, err := cs.getTotal(tx, cartId, userId)
	if err != nil {
		return nil, fmt.Errorf("error calculating total: %w", err)
	}

	cart := &models.ReturnedCart{
		ID:     cartId,
		UserId: userId,
		Items:  cartItems,
		Total:  total,
	}

	return cart, nil
}

func (cs *CartService) EditCartItem(cartItem *models.CartItem, userId string) (*models.CartItem, error) {
	if err := cs.checkOwnership(cartItem.CartId, userId); err != nil {
		return nil, err
	}

	availableStock, err := cs.getAvailableStock(cartItem.ProductId)
	if err != nil {
		return nil, err
	}

	if cartItem.Quantity > availableStock {
		return nil, fmt.Errorf("only %d items available in stock", availableStock)
	}

	var newItem models.CartItem
	updateQuery := `
		UPDATE cart_items
		SET quantity = $1
		WHERE productid = $2 AND cartid = $3 AND cartitemid = $4
		RETURNING *
	`
	err = cs.DB.Get(&newItem, updateQuery, cartItem.Quantity, cartItem.ProductId, cartItem.CartId, cartItem.CartItemId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("item does not exist in cart")
		}
		return nil, fmt.Errorf("error updating item: %w", err)
	}

	return &newItem, nil
}

func (cs *CartService) DeleteCartItem(cartItemId, userId, cartId string) error {
	if err := cs.checkOwnership(cartId, userId); err != nil {
		return err
	}

	deleteQuery := `
		DELETE FROM cart_items
		WHERE cartitemid = $1 AND cartid = $2
	`
	rows, err := cs.DB.Exec(deleteQuery, cartItemId, cartId)
	if err != nil {
		return fmt.Errorf("error deleting cart item: %w", err)
	}

	affectedRows, err := rows.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking deleted rows: %w", err)
	}
	if affectedRows == 0 {
		return fmt.Errorf("cart item doesn't exist")
	}

	return nil
}

func (cs *CartService) CreateCart(userId string) error {
	var cartCount int
	checkQuery := `SELECT COUNT(cartid) FROM carts WHERE userid = $1`
	err := cs.DB.Get(&cartCount, checkQuery, userId)
	if err != nil {
		return fmt.Errorf("error fetching possible cart info: %w", err)
	}

	if cartCount != 0 {
		return fmt.Errorf("you can only have one cart, but you can make multiple whishlists")
	}

	createQuery := `
	INSERT INTO carts (userid)
	VALUES ($1)
	`
	_, err = cs.DB.Exec(createQuery, userId)
	if err != nil {
		return fmt.Errorf("error creating cart: %w", err)
	}

	return nil
}
