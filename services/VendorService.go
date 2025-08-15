package services

import (
	"database/sql"
	"eCommerce/models"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type VendorService struct {
	DB *sqlx.DB
}

func NewVendorService(db *sqlx.DB) *VendorService {
	return &VendorService{
		db,
	}
}

// VENDOR ROUTES

func (vs *VendorService) AddProduct(product *models.Product, vendorId string) (*models.Product, error) {
	query := `
	INSERT INTO products (
	  vendorid, name, description, sku,
		price, discount, stock, is_active, brand, category
	) VALUES (
		$1, $2, $3, $4, 
	  $5, $6, $7, $8, $9, $10
	)
	RETURNING *
	`
	var inserted models.Product

	err := vs.DB.Get(&inserted, query,
		vendorId,
		product.Name,
		product.Description,
		product.SKU,
		product.Price,
		product.Discount,
		product.Stock,
		product.IsActive,
		product.Brand,
		product.Category,
	)
	if err != nil {
		return nil, fmt.Errorf("error adding new product: %w", err)
	}

	return &inserted, nil
}

func (vs *VendorService) GetProductByID(productId, vendorId string) (*models.Product, error) {
	var product models.Product
	query := `SELECT * FROM products WHERE productid = $1 AND vendorid = $2`

	err := vs.DB.Get(&product, query, productId, vendorId)
	if err != nil {
		return nil, fmt.Errorf("error fetching product: %w", err)
	}

	return &product, nil
}

func (vs *VendorService) GetVendorProducts(vendorId string) ([]*models.Product, error) {
	var products []*models.Product
	query := `SELECT * FROM products WHERE vendorid = $1`

	err := vs.DB.Select(&products, query, vendorId)
	if err != nil {
		return nil, fmt.Errorf("error fetching products: %w", err)
	}

	return products, nil
}

func (vs *VendorService) DeleteProduct(productId, vendorId string) error {
	var ownerId string
	err := vs.DB.QueryRow(`SELECT vendorid FROM products WHERE productid = $1`, productId).Scan(&ownerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("product not found")
		}
		return fmt.Errorf("error checking product ownership: %w", err)
	}

	if ownerId != vendorId {
		return fmt.Errorf("unauthorized: you do not own this product")
	}

	_, err = vs.DB.Exec(`DELETE FROM products WHERE productid = $1`, productId)
	if err != nil {
		return fmt.Errorf("error deleting product: %w", err)
	}

	return nil
}

func (vs *VendorService) UpdateProduct(product *models.Product, vendorId string) (*models.Product, error) {
	var ownerId string
	err := vs.DB.QueryRow(`SELECT vendorid FROM products WHERE productid = $1`, product.ProductId).Scan(&ownerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("error checking product ownership: %w", err)
	}

	if ownerId != vendorId {
		return nil, fmt.Errorf("unauthorized: you do not own this product")
	}

	setClauses := []string{}
	args := []any{}
	argsIndex := 1

	if product.Name != "" {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argsIndex))
		args = append(args, product.Name)
		argsIndex++
	}

	if product.Description != "" {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argsIndex))
		args = append(args, product.Description)
		argsIndex++
	}

	if product.SKU != "" {
		setClauses = append(setClauses, fmt.Sprintf("sku = $%d", argsIndex))
		args = append(args, product.SKU)
		argsIndex++
	}

	if product.Price != nil {
		setClauses = append(setClauses, fmt.Sprintf("price = $%d", argsIndex))
		args = append(args, *product.Price)
		argsIndex++
	}

	if product.Discount != nil {
		setClauses = append(setClauses, fmt.Sprintf("discount = $%d", argsIndex))
		args = append(args, *product.Discount)
		argsIndex++
	}

	if product.Stock != nil {
		setClauses = append(setClauses, fmt.Sprintf("stock = $%d", argsIndex))
		args = append(args, *product.Stock)
		argsIndex++
	}

	if product.IsActive != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_active = $%d", argsIndex))
		args = append(args, *product.IsActive)
		argsIndex++
	}
	if product.Brand != "" {
		setClauses = append(setClauses, fmt.Sprintf("brand = $%d", argsIndex))
		args = append(args, product.Brand)
		argsIndex++
	}

	if product.Category != "" {
		setClauses = append(setClauses, fmt.Sprintf("category = $%d", argsIndex))
		args = append(args, product.Category)
		argsIndex++
	}

	setClauses = append(setClauses, "updatedat = CURRENT_TIMESTAMP")

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf(`
	UPDATE products
	SET %s
	WHERE productid = $%d
	RETURNING *
	`, strings.Join(setClauses, ", "), argsIndex)

	args = append(args, product.ProductId)
	err = vs.DB.Get(product, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error updating product: %w", err)
	}

	return product, nil
}
