package services

import (
	"eCommerce/models"
	"fmt"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type ProductService struct {
	DB *sqlx.DB
}

func NewProductService(db *sqlx.DB) *ProductService {
	return &ProductService{
		db,
	}
}

func (ps *ProductService) GetProducts(filters map[string]string) ([]*models.Product, error) {
	var products []*models.Product
	baseQuery := `SELECT * FROM products WHERE is_active = TRUE`
	var args []any
	var conditions []string
	argsIndex := 1

	for key, value := range filters {
		switch key {
		case "category", "brand":
			conditions = append(conditions, fmt.Sprintf("%s ILIKE $%d", key, argsIndex))
			args = append(args, "%"+value+"%")
			argsIndex++
		case "min_price":
			conditions = append(conditions, fmt.Sprintf("price >= $%d", argsIndex))
			args = append(args, value)
			argsIndex++
		case "max_price":
			conditions = append(conditions, fmt.Sprintf("price <= $%d", argsIndex))
			args = append(args, value)
			argsIndex++
		case "discount":
			conditions = append(conditions, fmt.Sprintf("discount >= $%d", argsIndex))
			args = append(args, value)
			argsIndex++
		case "keyword":
			conditions = append(conditions, fmt.Sprintf("(name ILIKE $%d OR description ILIKE $%d)", argsIndex, argsIndex))
			args = append(args, "%"+value+"%")
			argsIndex++
		}
	}
	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	// sorting
	if sortBy := filters["sort_by"]; sortBy != "" {
		switch sortBy {
		case "price", "discount", "stock", "createdat":
			baseQuery += fmt.Sprintf(" ORDER BY %s", sortBy)
		}
	} else {
		baseQuery += " ORDER BY productid ASC "
	}

	// ordering
	if order := filters["order"]; order != "" {
		switch order {
		case "ASC":
			baseQuery += " ASC"
		case "DESC":
			baseQuery += " DESC"
		default:
			baseQuery += " ASC"
		}
	}

	// pagination
	if limitStr := filters["limit"]; limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			baseQuery += fmt.Sprintf(" LIMIT %d", limit)
		}
	}
	if offsetStr := filters["offset"]; offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset > 0 {
			baseQuery += fmt.Sprintf(" OFFSET %d", offset)
		}
	}

	err := ps.DB.Select(&products, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("error fetching products: %w", err)
	}

	return products, nil
}
