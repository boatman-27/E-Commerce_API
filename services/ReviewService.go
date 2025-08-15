package services

import (
	"database/sql"
	"eCommerce/models"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type ReviewService struct {
	DB *sqlx.DB
}

func NewReviewService(db *sqlx.DB) *ReviewService {
	return &ReviewService{
		db,
	}
}

func (rs *ReviewService) getVendorId(productId string) (string, error) {
	var vendorId string
	query := `SELECT vendorid FROM products WHERE productid = $1`
	err := rs.DB.Get(&vendorId, query, productId)
	if err != nil {
		return "", fmt.Errorf("error fetching vendorId: %w", err)
	}
	return vendorId, nil
}

func (rs *ReviewService) checkOwnership(reviewId, userId string) error {
	var ownerId string
	err := rs.DB.QueryRow("SELECT userid FROM reviews WHERE reviewid = $1", reviewId).Scan(&ownerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrNotFound
		}
		return fmt.Errorf("error checking review ownership: %w", err)
	}
	if ownerId != userId {
		return models.ErrUnauthorized
	}
	return nil
}

func (rs *ReviewService) SubmitReview(userId, productId string, reviewData *models.ReviewData) error {
	vendorId, err := rs.getVendorId(productId)
	if err != nil {
		return err
	}

	if userId == vendorId {
		return fmt.Errorf("cant review your own product :)")
	}

	query := `
	INSERT INTO reviews (userid, vendorid, productid, message, stars)
	VALUES ($1, $2, $3, $4, $5)
	`

	if reviewData.Stars < 1 || reviewData.Stars > 5 {
		return fmt.Errorf("stars only range from 1 to 5 (inclusive)")
	}

	_, err = rs.DB.Exec(query, userId, vendorId, productId, reviewData.ReviewMessage, reviewData.Stars)
	if err != nil {
		return fmt.Errorf("error submitting review: %w", err)
	}

	return nil
}

func (rs *ReviewService) DeleteReview(userId, reviewId string) error {
	query := `DELETE FROM reviews WHERE reviewid = $1 AND userid = $2`
	rows, err := rs.DB.Exec(query, reviewId, userId)
	if err != nil {
		return fmt.Errorf("error deleting review: %w", err)
	}

	affectedRows, err := rows.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking deleted rows: %w", err)
	}
	if affectedRows == 0 {
		return fmt.Errorf("review doesn't exist")
	}
	return nil
}

func (rs *ReviewService) EditReview(reviewId, userId string, reviewData *models.ReviewData) error {
	if err := rs.checkOwnership(reviewId, userId); err != nil {
		return err
	}

	setClauses := []string{}
	args := []any{}
	argsIndex := 1

	if reviewData.ReviewMessage != "" {
		setClauses = append(setClauses, fmt.Sprintf("message = $%d", argsIndex))
		args = append(args, reviewData.ReviewMessage)
		argsIndex++
	}

	if reviewData.Stars != 0 {
		if reviewData.Stars < 1 || reviewData.Stars > 5 {
			return fmt.Errorf("stars only range from 1 to 5 (inclusive)")
		}
		setClauses = append(setClauses, fmt.Sprintf("stars = $%d", argsIndex))
		args = append(args, reviewData.Stars)
		argsIndex++
	}
	setClauses = append(setClauses, "updatedat = CURRENT_TIMESTAMP")

	if len(setClauses) == 0 {
		return fmt.Errorf("no fields to update")
	}

	updateQuery := fmt.Sprintf(`
	UPDATE reviews 
	SET %s
	WHERE reviewid = $%d`, strings.Join(setClauses, ", "), argsIndex)
	args = append(args, reviewId)

	_, err := rs.DB.Exec(updateQuery, args...)
	if err != nil {
		return fmt.Errorf("error update review: %w", err)
	}

	return nil
}

func (rs *ReviewService) GetReviews(userId string) ([]*models.Review, error) {
	query := `SELECT * FROM reviews WHERE userid= $1`
	var reviews []*models.Review

	err := rs.DB.Select(&reviews, query, userId)
	if err != nil {
		return nil, fmt.Errorf("error fetching reviews for current user: %w", err)
	}

	return reviews, nil
}

func (rs *ReviewService) GetVendorReviews(vendorId string) ([]*models.Review, error) {
	query := `SELECT * FROM reviews WHERE vendorid= $1`
	var reviews []*models.Review

	err := rs.DB.Select(&reviews, query, vendorId)
	if err != nil {
		return nil, fmt.Errorf("error fetching reviews for current vendor: %w", err)
	}

	return reviews, nil
}
