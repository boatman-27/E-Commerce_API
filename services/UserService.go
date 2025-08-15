package services

import (
	"eCommerce/helpers"
	"eCommerce/models"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserService struct {
	DB *sqlx.DB
}

func NewUserService(db *sqlx.DB) *UserService {
	return &UserService{
		DB: db,
	}
}

func (us *UserService) Login(creds *models.Credentials) (*models.User, string, string, error) {
	var user models.User
	query := `
	SELECT userid, name, email, password, phone_number, verified, verification_token, role FROM users WHERE email = $1
	`
	// Get User
	err := us.DB.Get(&user, query, creds.Email)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to fetch user: %w", err)
	}
	if !user.Verified {
		return nil, "", "", fmt.Errorf("user not verified, please check your inbox for the verification link: %w", err)
	}

	// Compare passwords
	if !helpers.CheckPasswords(creds.Password, user.Password) {
		return nil, "", "", fmt.Errorf("passwords don't match")
	}

	accessToken, err := helpers.GenerateAccessToken(user.UserId.String(), user.Email, user.Role)
	if err != nil {
		return nil, "", "", fmt.Errorf("error generating access token: %w", err)
	}

	refreshToken, err := helpers.GenerateRefreshToken(user.UserId.String(), user.Email, user.Role)
	if err != nil {
		return nil, "", "", fmt.Errorf("error generating refresh Token: %w", err)
	}

	return &user, accessToken, refreshToken, nil
}

func (us *UserService) Signup(user *models.User) error {
	// checks if entered email is used
	emailAvailable, err := helpers.IsEmailAvailable(us.DB, user.Email)
	if err != nil {
		return fmt.Errorf("could not check email availability: %w", err)
	}
	if !emailAvailable {
		return fmt.Errorf("email already taken")
	}

	if !helpers.IsValidEmail(user.Email) {
		return fmt.Errorf("invalid email format")
	}

	// generate new userId format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	user.UserId = uuid.New()

	hashedPassword, err := helpers.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("could not hash password: %w", err)
	}

	user.Password = hashedPassword
	user.VerifiedToken = uuid.New().String()

	fmt.Println(user.PhoneNumber)

	// Insert into DB and return inserted user
	query := `
		INSERT INTO users (userid, name, email, password, phone_number, verification_token, role)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = us.DB.Exec(query, user.UserId, user.Name, user.Email, user.Password, user.PhoneNumber, user.VerifiedToken, user.Role)
	if err != nil {
		return fmt.Errorf("error inserting user: %w", err)
	}

	err = helpers.SendVerificationEmail(user.VerifiedToken, user.Email)
	if err != nil {
		return fmt.Errorf("could not send verification email: %w", err)
	}

	return nil
}

func (us *UserService) UpdateUser(user *models.User, userId string) (*models.User, string, string, error) {
	setClauses := []string{}
	args := []any{}
	argIndex := 1
	needNewTokens := false

	if user.Name != "" {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, user.Name)
		argIndex++
	}

	if user.Email != "" {
		if !helpers.IsValidEmail(user.Email) {
			return nil, "", "", fmt.Errorf("invalid email format")
		}

		setClauses = append(setClauses, fmt.Sprintf("email = $%d", argIndex))
		args = append(args, user.Email)
		argIndex++
		needNewTokens = true
	}

	if user.Password != "" {
		hashedPassword, err := helpers.HashPassword(user.Password)
		if err != nil {
			return nil, "", "", fmt.Errorf("failed to hash password: %w", err)
		}
		setClauses = append(setClauses, fmt.Sprintf("password = $%d", argIndex))
		args = append(args, hashedPassword)
		argIndex++
	}

	if user.PhoneNumber != "" {
		setClauses = append(setClauses, fmt.Sprintf("phone_number = $%d", argIndex))
		args = append(args, user.PhoneNumber)
		argIndex++
	}

	setClauses = append(setClauses, "updatedat = CURRENT_TIMESTAMP")

	if len(setClauses) == 1 {
		return nil, "", "", fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf(`
		UPDATE users
		SET %s
		WHERE userid = $%d
		RETURNING *
	`, strings.Join(setClauses, ", "), argIndex)

	args = append(args, userId)

	err := us.DB.Get(user, query, args...)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to update user: %w", err)
	}

	if needNewTokens {
		accessToken, err := helpers.GenerateAccessToken(user.UserId.String(), user.Email, user.Role)
		if err != nil {
			return nil, "", "", fmt.Errorf("error generating access token: %w", err)
		}

		refreshToken, err := helpers.GenerateRefreshToken(user.UserId.String(), user.Email, user.Role)
		if err != nil {
			return nil, "", "", fmt.Errorf("error generating refresh Token: %w", err)
		}

		return user, accessToken, refreshToken, nil
	}

	return user, "", "", nil
}

func (us *UserService) VerifyUser(verificationToken string) error {
	var user models.User
	query := `
	SELECT userid, name, email, password, phone_number, verified, verification_token, role FROM users WHERE verification_token = $1
	`
	// Get User
	err := us.DB.Get(&user, query, verificationToken)
	if err != nil {
		return fmt.Errorf("failed to fetch user: %w", err)
	}

	if user.Verified {
		return fmt.Errorf("user already verified, logging in")
	}

	tx, err := us.DB.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
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

	verificationQuery := `
	UPDATE users
	SET verified = $1, verification_token = $2
	WHERE verification_token = $3
	`
	_, err = tx.Exec(verificationQuery, true, "", verificationToken)
	if err != nil {
		return fmt.Errorf("error verifying user: %w", err)
	}

	cartQuery := `
	INSERT INTO carts (userid)
	VALUES ($1)`
	_, err = tx.Exec(cartQuery, user.UserId)
	if err != nil {
		return fmt.Errorf("error creating cart for new user: %w", err)
	}

	return nil
}

func (us *UserService) GetUserProfile(userId string) (*models.User, error) {
	var user models.User

	query := `SELECT * FROM users WHERE userid = $1 LIMIT 1`
	err := us.DB.Get(&user, query, userId)
	if err != nil {
		return nil, fmt.Errorf("error fetching user profile: %w", err)
	}

	return &user, nil
}
