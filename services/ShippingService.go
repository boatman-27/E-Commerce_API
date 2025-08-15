package services

import (
	"database/sql"
	"eCommerce/models"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type ShippingService struct {
	DB *sqlx.DB
}

func NewShippingService(db *sqlx.DB) *ShippingService {
	return &ShippingService{
		db,
	}
}

func (ss *ShippingService) checkOwnership(userId, shippinggId string) error {
	var ownerId string
	err := ss.DB.QueryRow("SELECT userid FROM shipping_addresses WHERE shippingid = $1", shippinggId).Scan(&ownerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrNotFound
		}
		return fmt.Errorf("error checking shipping address issuer: %w", err)
	}
	if ownerId != userId {
		return models.ErrUnauthorized
	}
	return nil
}

func (ss *ShippingService) AddShippingAddress(shippingAddress *models.ShippingAddress, userId string) (*models.ShippingAddress, error) {
	tx, err := ss.DB.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// If new address is default, unset old default
	if shippingAddress.IsDefault {
		_, err := tx.Exec(`UPDATE shipping_addresses SET is_default = FALSE WHERE userid = $1 AND is_default = TRUE`, userId)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to unset old default address: %w", err)
		}
	}

	query := `
	INSERT INTO shipping_addresses (userid, address_line1, address_line2, city, state, postal_code, country, is_default)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING *
	`
	err = tx.Get(shippingAddress, query,
		userId,
		shippingAddress.AddressLine1,
		shippingAddress.AddressLine2,
		shippingAddress.City,
		shippingAddress.State,
		shippingAddress.PostalCode,
		shippingAddress.Country,
		shippingAddress.IsDefault,
	)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to insert shipping address: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return shippingAddress, nil
}

func (ss *ShippingService) DeleteShippingAddress(shippingId, userId string) error {
	_, err := ss.DB.Exec("DELETE FROM shipping_addresses WHERE shippingid = $1 AND userid = $2", shippingId, userId)
	if err != nil {
		return fmt.Errorf("failed to delete shipping address: %w", err)
	}

	return nil
}

func (ss *ShippingService) GetShippingAddresses(userId string) ([]*models.ShippingAddress, error) {
	var addresses []*models.ShippingAddress
	fetchQuery := "SELECT * FROM shipping_addresses WHERE userid = $1"
	err := ss.DB.Select(&addresses, fetchQuery, userId)
	if err != nil {
		return nil, fmt.Errorf("error fetching addresses: %w", err)
	}
	return addresses, nil
}

func (ss *ShippingService) UpdateShippingAddress(shippingAddress *models.ShippingAddress, userId, shippingId string) (*models.ShippingAddress, error) {
	if err := ss.checkOwnership(userId, shippingId); err != nil {
		return nil, err
	}

	setClauses := []string{}
	args := []any{}
	argsIndex := 1

	if shippingAddress.AddressLine1 != "" {
		setClauses = append(setClauses, fmt.Sprintf("address_line1 = $%d", argsIndex))
		args = append(args, shippingAddress.AddressLine1)
		argsIndex++
	}

	if shippingAddress.AddressLine2 != "" {
		setClauses = append(setClauses, fmt.Sprintf("address_line2 = $%d", argsIndex))
		args = append(args, shippingAddress.AddressLine2)
		argsIndex++
	}

	if shippingAddress.City != "" {
		setClauses = append(setClauses, fmt.Sprintf("city = $%d", argsIndex))
		args = append(args, shippingAddress.City)
		argsIndex++
	}

	if shippingAddress.PostalCode != "" {
		setClauses = append(setClauses, fmt.Sprintf("postal_code = $%d", argsIndex))
		args = append(args, shippingAddress.PostalCode)
		argsIndex++
	}

	if shippingAddress.State != "" {
		setClauses = append(setClauses, fmt.Sprintf("state = $%d", argsIndex))
		args = append(args, shippingAddress.State)
		argsIndex++
	}

	if shippingAddress.Country != "" {
		setClauses = append(setClauses, fmt.Sprintf("country = $%d", argsIndex))
		args = append(args, shippingAddress.Country)
		argsIndex++
	}

	setClauses = append(setClauses, "updatedat = CURRENT_TIMESTAMP")

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf(`
		UPDATE shipping_addresses
		SET %s
		WHERE shippingid = $%d
		RETURNING *`, strings.Join(setClauses, ", "), argsIndex)

	args = append(args, shippingAddress.ShippingId)
	err := ss.DB.Get(shippingAddress, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update shipping Address: %w", err)
	}

	return shippingAddress, nil
}

func (ss *ShippingService) GetDefaultShippingAddress(userId string) (*models.ShippingAddress, error) {
	var address models.ShippingAddress
	fetchQuery := "SELECT * FROM shipping_addresses WHERE userid = $1 AND is_default = TRUE"
	err := ss.DB.Get(&address, fetchQuery, userId)
	if err != nil {
		return nil, fmt.Errorf("error fetching addresses: %w", err)
	}
	return &address, nil
}

func (ss *ShippingService) ChangeDefaultShippingAddress(userId, oldDefault, newDefault string) error {
	tx, err := ss.DB.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Check old default exists
	var oldCount int
	err = tx.Get(&oldCount, `SELECT COUNT(*) FROM shipping_addresses WHERE shippingid = $1 AND userid = $2`, oldDefault, userId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to check old default address existence: %w", err)
	}
	if oldCount == 0 {
		tx.Rollback()
		return fmt.Errorf("old default address not found for user")
	}

	// Check new default exists
	var newCount int
	err = tx.Get(&newCount, `SELECT COUNT(*) FROM shipping_addresses WHERE shippingid = $1 AND userid = $2`, newDefault, userId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to check new default address existence: %w", err)
	}
	if newCount == 0 {
		tx.Rollback()
		return fmt.Errorf("new default address not found for user")
	}

	// Unset old default
	res, err := tx.Exec(`UPDATE shipping_addresses SET is_default = FALSE WHERE shippingid = $1 AND userid = $2`, oldDefault, userId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error unsetting old address: %w", err)
	}
	if affectedRows, _ := res.RowsAffected(); affectedRows == 0 {
		tx.Rollback()
		return fmt.Errorf("no rows were affected when unsetting old default")
	}

	// Set new default
	res, err = tx.Exec(`UPDATE shipping_addresses SET is_default = TRUE WHERE shippingid = $1 AND userid = $2`, newDefault, userId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error setting new default address: %w", err)
	}
	if affectedRows, _ := res.RowsAffected(); affectedRows == 0 {
		tx.Rollback()
		return fmt.Errorf("no rows were affected when setting new default")
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
