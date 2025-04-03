package processor

import (
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/receipt-processor/models"
)

type PointsCalculator interface {
	CalculatePoints(receipt models.Receipt) int
}

func CalculatePoints(receipt models.Receipt) int {
	points := 0

	points += countAlphanumeric(receipt.Retailer)

	total, _ := strconv.ParseFloat(receipt.Total, 64)
	if total == math.Floor(total) {
		points += 50
	}

	if math.Mod(total*100, 25) == 0 {
		points += 25
	}

	points += (len(receipt.Items) / 2) * 5

	for _, item := range receipt.Items {
		trimDesc := strings.TrimSpace(item.ShortDescription)
		if len(trimDesc)%3 == 0 && len(trimDesc) > 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			points += int(math.Ceil(price * 0.2))
		}
	}

	purchaseDate, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
	if purchaseDate.Day()%2 == 1 {
		points += 6
	}

	purchaseTime, _ := time.Parse("15:04", receipt.PurchaseTime)
	afternoon2pm, _ := time.Parse("15:04", "14:00")
	afternoon4pm, _ := time.Parse("15:04", "16:00")

	if purchaseTime.After(afternoon2pm) && purchaseTime.Before(afternoon4pm) {
		points += 10
	}

	return points
}

func countAlphanumeric(s string) int {
	count := 0
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			count++
		}
	}
	return count
}
