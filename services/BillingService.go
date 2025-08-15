package services

import (
	"database/sql"
	"eCommerce/models"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type BillingService struct {
	DB *sqlx.DB
}

func NewBillingService(db *sqlx.DB) *BillingService {
	return &BillingService{
		db,
	}
}

func (bs *BillingService) checkOwnership(userId, billingId string) error {
	var ownerId string
	err := bs.DB.QueryRow("SELECT userid FROM billing_addresses WHERE billingid = $1", billingId).Scan(&ownerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrNotFound
		}
		return fmt.Errorf("error checking billing address issuer: %w", err)
	}
	if ownerId != userId {
		return models.ErrUnauthorized
	}
	return nil
}

func (bs *BillingService) AddBillingAddress(billingAddress *models.BillingAddress, userId string) (*models.BillingAddress, error) {
	tx, err := bs.DB.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// If new address is default, unset old default
	if billingAddress.IsDefault {
		_, err := tx.Exec(`UPDATE billing_addresses SET is_default = FALSE WHERE userid = $1 AND is_default = TRUE`, userId)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to unset old default address: %w", err)
		}
	}

	query := `
	INSERT INTO billing_addresses (userid, address_line1, address_line2, city, state, postal_code, country, is_default)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING *
	`
	err = tx.Get(billingAddress, query,
		userId,
		billingAddress.AddressLine1,
		billingAddress.AddressLine2,
		billingAddress.City,
		billingAddress.State,
		billingAddress.PostalCode,
		billingAddress.Country,
		billingAddress.IsDefault,
	)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to insert billing address: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return billingAddress, nil
}

func (bs *BillingService) DeleteBillingAddress(billingId, userId string) error {
	_, err := bs.DB.Exec("DELETE FROM billing_addresses WHERE billingid = $1 AND userid = $2", billingId, userId)
	if err != nil {
		return fmt.Errorf("failed to delete billing address: %w", err)
	}

	return nil
}

func (bs *BillingService) GetBillingAddresses(userId string) ([]*models.BillingAddress, error) {
	var addresses []*models.BillingAddress
	fetchQuery := "SELECT * FROM billing_addresses WHERE userid = $1"
	err := bs.DB.Select(&addresses, fetchQuery, userId)
	if err != nil {
		return nil, fmt.Errorf("error fetching addresses: %w", err)
	}
	return addresses, nil
}

func (bs *BillingService) UpdateBillingAddress(billingAddress *models.BillingAddress, userId, billingId string) (*models.BillingAddress, error) {
	if err := bs.checkOwnership(userId, billingId); err != nil {
		return nil, err
	}

	setClauses := []string{}
	args := []any{}
	argsIndex := 1

	if billingAddress.AddressLine1 != "" {
		setClauses = append(setClauses, fmt.Sprintf("address_line1 = $%d", argsIndex))
		args = append(args, billingAddress.AddressLine1)
		argsIndex++
	}

	if billingAddress.AddressLine2 != "" {
		setClauses = append(setClauses, fmt.Sprintf("address_line2 = $%d", argsIndex))
		args = append(args, billingAddress.AddressLine2)
		argsIndex++
	}

	if billingAddress.City != "" {
		setClauses = append(setClauses, fmt.Sprintf("city = $%d", argsIndex))
		args = append(args, billingAddress.City)
		argsIndex++
	}

	if billingAddress.PostalCode != "" {
		setClauses = append(setClauses, fmt.Sprintf("postal_code = $%d", argsIndex))
		args = append(args, billingAddress.PostalCode)
		argsIndex++
	}

	if billingAddress.State != "" {
		setClauses = append(setClauses, fmt.Sprintf("state = $%d", argsIndex))
		args = append(args, billingAddress.State)
		argsIndex++
	}

	if billingAddress.Country != "" {
		setClauses = append(setClauses, fmt.Sprintf("country = $%d", argsIndex))
		args = append(args, billingAddress.Country)
		argsIndex++
	}

	setClauses = append(setClauses, "updatedat = CURRENT_TIMESTAMP")

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf(`
		UPDATE billing_addresses
		SET %s
		WHERE billingid = $%d
		RETURNING *`, strings.Join(setClauses, ", "), argsIndex)

	args = append(args, billingAddress.BillingId)
	err := bs.DB.Get(billingAddress, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update billing Address: %w", err)
	}

	return billingAddress, nil
}

func (bs *BillingService) GetDefaultBillingAddress(userId string) (*models.BillingAddress, error) {
	var address models.BillingAddress
	fetchQuery := "SELECT * FROM billing_addresses WHERE userid = $1 AND is_default = TRUE"
	err := bs.DB.Get(&address, fetchQuery, userId)
	if err != nil {
		return nil, fmt.Errorf("error fetching addresses: %w", err)
	}
	return &address, nil
}

func (bs *BillingService) ChangeDefaultBillingAddress(userId, oldDefault, newDefault string) error {
	tx, err := bs.DB.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Check old default exists
	var oldCount int
	err = tx.Get(&oldCount, `SELECT COUNT(*) FROM billing_addresses WHERE billingid = $1 AND userid = $2`, oldDefault, userId)
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
	err = tx.Get(&newCount, `SELECT COUNT(*) FROM billing_addresses WHERE billingid = $1 AND userid = $2`, newDefault, userId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to check new default address existence: %w", err)
	}
	if newCount == 0 {
		tx.Rollback()
		return fmt.Errorf("new default address not found for user")
	}

	// Unset old default
	res, err := tx.Exec(`UPDATE billing_addresses SET is_default = FALSE WHERE billingid = $1 AND userid = $2`, oldDefault, userId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error unsetting old address: %w", err)
	}
	if affectedRows, _ := res.RowsAffected(); affectedRows == 0 {
		tx.Rollback()
		return fmt.Errorf("no rows were affected when unsetting old default")
	}

	// Set new default
	res, err = tx.Exec(`UPDATE billing_addresses SET is_default = TRUE WHERE billingid = $1 AND userid = $2`, newDefault, userId)
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
