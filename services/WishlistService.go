package services

import (
	"database/sql"
	"eCommerce/models"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type WishlistService struct {
	DB *sqlx.DB
}

func NewWishlistService(db *sqlx.DB) *WishlistService {
	return &WishlistService{
		db,
	}
}

func (ws *WishlistService) checkOwnership(userId, wishlistId string) error {
	var ownerId string
	err := ws.DB.QueryRow("SELECT userid FROM wishlists WHERE wishlistid = $1", wishlistId).Scan(&ownerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrNotFound
		}
		return fmt.Errorf("error checking wishlist ownership: %w", err)
	}
	if ownerId != userId {
		return models.ErrUnauthorized
	}
	return nil
}

func (ws *WishlistService) CreateWishlist(userId string, info *models.WishlistInfo) error {
	insertQuery := `
	INSERT INTO wishlists (userid, name, description)
	VALUES ($1, $2, $3)
	`
	_, err := ws.DB.Exec(insertQuery, userId, info.Name, info.Description)
	if err != nil {
		return fmt.Errorf("error creating a new wishlist: %w", err)
	}

	return nil
}

func (ws *WishlistService) GetWishlistItems(userId, wishlistId string) (*models.ReturnedWishlist, error) {
	var items []*models.ItemData
	query := `SELECT
	p.productid,
	p.name,
	p.description,
	p.sku,
	p.price,
	p.discount,
	p.brand,
	p.category
	FROM products p 
	JOIN wishlist_items wi ON wi.productid = p.productid
	JOIN wishlists w on w.wishlistid = wi.wishlistid
	WHERE w.userid = $1 AND wi.wishlistid = $2
	`

	err := ws.DB.Select(&items, query, userId, wishlistId)
	if err != nil {
		return nil, fmt.Errorf("error fetching items: %w", err)
	}

	wishlist := &models.ReturnedWishlist{
		ID:     wishlistId,
		UserId: userId,
		Items:  items,
	}

	return wishlist, nil
}

func (ws *WishlistService) DeleteWishlist(userId, wishlistId string) error {
	deleteQuery := `DELETE FROM wishlists WHERE userid = $1 AND wishlistid = $2`
	_, err := ws.DB.Exec(deleteQuery, userId, wishlistId)
	if err != nil {
		return fmt.Errorf("error deleting wishlist: %w", err)
	}

	return nil
}

func (ws *WishlistService) EditWishlist(userId, wishlistId string, info *models.WishlistInfo) error {
	if err := ws.checkOwnership(userId, wishlistId); err != nil {
		return err
	}
	setClauses := []string{}
	args := []any{}
	argsIndex := 1

	if info.Name != "" {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argsIndex))
		args = append(args, info.Name)
		argsIndex++
	}

	if info.Description != "" {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argsIndex))
		args = append(args, info.Description)
		argsIndex++
	}

	setClauses = append(setClauses, "updatedat = CURRENT_TIMESTAMP")

	if len(setClauses) == 0 {
		return fmt.Errorf("no fields to update")
	}

	updateQuery := fmt.Sprintf(`
	UPDATE wishlists
	SET %s
	WHERE wishlistid = $%d
	`, strings.Join(setClauses, ", "), argsIndex)

	args = append(args, wishlistId)

	_, err := ws.DB.Exec(updateQuery, args...)
	if err != nil {
		return fmt.Errorf("error update wishlist: %w", err)
	}

	return nil
}

func (ws *WishlistService) DeleteWishlistItem(userId, wishlistId, wishlistItemId string) error {
	if err := ws.checkOwnership(userId, wishlistId); err != nil {
		return err
	}

	deleteQuery := `
	DELETE FROM wishlist_items
	WHERE wishlistitemid = $1 AND wishlistid = $2`
	rows, err := ws.DB.Exec(deleteQuery, wishlistItemId, wishlistId)
	if err != nil {
		return fmt.Errorf("error deleting item from wishlist: %w", err)
	}

	affectedRows, err := rows.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking deleted rows: %w", err)
	}

	if affectedRows == 0 {
		return fmt.Errorf("item not in wishlist")
	}

	return nil
}

func (ws *WishlistService) MoveToCart(userId, wishlistId string) error {
	tx, err := ws.DB.Begin()
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}

	_, err = tx.Exec(`
		INSERT INTO cart_items (cartid, productid, quantity)
		SELECT c.cartid, w.productid, 1
		FROM carts c
		JOIN wishlist_items w ON w.wishlistid = $2
		WHERE c.userid = $1
		ON CONFLICT (cartid, productid)
		DO UPDATE SET quantity = cart_items.quantity + 1
	`, userId, wishlistId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error inserting into cart_items: %w", err)
	}

	_, err = tx.Exec(`DELETE FROM wishlist_items WHERE wishlistid = $1`, wishlistId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting from wishlist_items: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

func (ws *WishlistService) AddToWishlist(userId, wishlistId, productId string) (*models.ReturnedWishlist, error) {
	if err := ws.checkOwnership(userId, wishlistId); err != nil {
		return nil, err
	}

	tx, err := ws.DB.Beginx() // use Beginx for sqlx transaction support
	if err != nil {
		return nil, fmt.Errorf("starting transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	insertQuery := `INSERT INTO wishlist_items (wishlistid, productid) VALUES ($1, $2)`
	_, err = tx.Exec(insertQuery, wishlistId, productId)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error adding item to wishlist: %w", err)
	}

	items := []*models.ItemData{}
	selectQuery := `
	SELECT
		p.productid,
		p.name,
		p.description,
		p.sku,
		p.price,
		p.discount,
		p.brand,
		p.category
	FROM products p
	JOIN wishlist_items wi ON wi.productid = p.productid
	JOIN wishlists w ON w.wishlistid = wi.wishlistid
	WHERE w.userid = $1 AND wi.wishlistid = $2
	`
	err = tx.Select(&items, selectQuery, userId, wishlistId)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error fetching items: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	wishlist := &models.ReturnedWishlist{
		ID:     wishlistId,
		UserId: userId,
		Items:  items,
	}

	return wishlist, nil
}

func (ws *WishlistService) MoveItemToCart(userId, wishlistId, wishlistItemId string) error {
	tx, err := ws.DB.Begin()
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}

	_, err = tx.Exec(`
		INSERT INTO cart_items (cartid, productid, quantity)
		SELECT c.cartid, w.productid, 1
		FROM carts c
		JOIN wishlist_items w ON w.wishlistid = $2
		WHERE c.userid = $1 AND w.wishlistid = $2
		ON CONFLICT (cartid, productid)
		DO UPDATE SET quantity = cart_items.quantity + 1
	`, userId, wishlistId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error inserting into cart_items: %w", err)
	}

	_, err = tx.Exec(`DELETE FROM wishlist_items WHERE wishlistitemid = $1 AND wishlistid = $2`, wishlistItemId, wishlistId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting from wishlist_items: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}
