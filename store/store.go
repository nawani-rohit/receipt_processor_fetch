package store

import (
	"errors"
	"math"
	"receipt-processor/models"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

// global variables
var (
	receipt_points = make(map[string]int)
	data_lock      sync.Mutex
)

// ValidateReceipt
func ValidateReceipt(r *models.Receipt) error {
	if r.Retailer == "" {
		return errors.New("retailer name is required")
	}
	if r.PurchaseDate == "" {
		return errors.New("purchaseDate is required")
	}
	if r.PurchaseTime == "" {
		return errors.New("purchaseTime is required")
	}
	if r.Total == "" {
		return errors.New("total is required")
	}
	if len(r.Items) == 0 {
		return errors.New("at least one item is required")
	}

	// Validate retailer using regex
	retailer_regex := regexp.MustCompile(`^[\w\s\-\&]+$`)
	if !retailer_regex.MatchString(r.Retailer) {
		return errors.New("retailer name format not correct")
	}

	// Validate total using regex
	total_regex := regexp.MustCompile(`^\d+\.\d{2}$`)
	if !total_regex.MatchString(r.Total) {
		return errors.New("total amount format not correct - use XX.XX format")
	}

	// Validate purchaseDate using the layout
	_, date_err := time.Parse("2006-01-02", r.PurchaseDate)
	if date_err != nil {
		return errors.New("date format not correct - use YYYY-MM-DD format")
	}

	// Validate purchaseTime using the layout
	_, time_err := time.Parse("15:04", r.PurchaseTime)
	if time_err != nil {
		return errors.New("time format not correct - use HH:MM format")
	}

	// Validate each item
	item_description_regex := regexp.MustCompile(`^[\w\s\-]+$`)
	for _, item := range r.Items {
		if item.ShortDescription == "" {
			return errors.New("item shortDescription is required")
		}
		if !item_description_regex.MatchString(item.ShortDescription) {
			return errors.New("item shortDescription format is not correct")
		}
		if item.Price == "" {
			return errors.New("item price is required")
		}
		if !total_regex.MatchString(item.Price) {
			return errors.New("item price format is not correct - use XX.XX format")
		}
	}

	return nil
}

// CalculatePoints
func CalculatePoints(receipt *models.Receipt) int {
	total_points := 0

	// Rule 1: points for retailer name
	for _, ch := range receipt.Retailer {
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) {
			total_points++
		}
	}

	// Rule 2: points for total amount
	total, convert_err := strconv.ParseFloat(receipt.Total, 64)
	if convert_err == nil {
		if total == math.Floor(total) {
			total_points += 50
		}
		if int(math.Round(total*100))%25 == 0 {
			total_points += 25
		}
		if total > 10.00 {
			total_points += 5
		}
	}

	// Rule 3: points for number of items
	total_points += (len(receipt.Items) / 2) * 5

	// Rule 4: points for item description length
	for _, item := range receipt.Items {
		trimmed_description := strings.TrimSpace(item.ShortDescription)
		if len(trimmed_description)%3 == 0 {
			price, price_err := strconv.ParseFloat(item.Price, 64)
			if price_err == nil {
				total_points += int(math.Ceil(price * 0.2))
			}
		}
	}

	// Rule 5: points for purchase date
	purchaseDate, err := time.Parse("2006-01-02", receipt.PurchaseDate)
	if err == nil && purchaseDate.Day()%2 != 0 {
		total_points += 6
	}

	// Rule 6: points for purchase time
	purchase_time, time_err := time.Parse("15:04", receipt.PurchaseTime)
	if time_err == nil {
		base_time := time.Date(0, 1, 1, 14, 0, 0, 0, time.UTC)

		start_time := base_time
		end_time := base_time.Add(2 * time.Hour)

		check_time := time.Date(0, 1, 1, purchase_time.Hour(), purchase_time.Minute(), 0, 0, time.UTC)
		if check_time.After(start_time) && check_time.Before(end_time) {
			total_points += 10
		}
	}

	return total_points
}

// SaveReceipt
func SaveReceipt(id string, points int) {
	data_lock.Lock()
	receipt_points[id] = points
	defer data_lock.Unlock()
}

// GetPoints
func GetPoints(id string) (int, bool) {
	data_lock.Lock()
	points, found := receipt_points[id]
	defer data_lock.Unlock()
	return points, found
}
