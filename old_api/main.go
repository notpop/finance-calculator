package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LoanDetails struct {
	Month              int     `json:"month"`
	Year               float64 `json:"year"`
	Principal          float64 `json:"principal"`
	InterestRate       float64 `json:"interestRate"`
	MonthlyPayment     float64 `json:"monthlyPayment"`
	Interest           float64 `json:"interest"`
	PrincipalReduction float64 `json:"principalReduction"`
	TotalPaid          float64 `json:"totalPaid"`
}

func main() {
	r := gin.Default()
	r.GET("/loan", func(c *gin.Context) {
		principal, _ := strconv.ParseFloat(c.Query("principal"), 64)
		interestRate, _ := strconv.ParseFloat(c.Query("interestRate"), 64)
		monthlyPayment, _ := strconv.ParseFloat(c.Query("monthlyPayment"), 64)

		monthlyInterestRate := interestRate / 100 / 12

		var response []LoanDetails
		month := 1
		totalPaid := 0.0

		for principal > 0 {
			interest := principal * monthlyInterestRate
			totalDue := principal + interest
			thisPayment := monthlyPayment

			if thisPayment > totalDue {
				thisPayment = totalDue
			}

			principalReduction := thisPayment - interest
			principal -= principalReduction
			totalPaid += thisPayment

			response = append(response, LoanDetails{
				Month:              month,
				Year:               float64(month) / 12,
				Principal:          principal,
				InterestRate:       interestRate,
				MonthlyPayment:     thisPayment,
				Interest:           interest,
				PrincipalReduction: principalReduction,
				TotalPaid:          totalPaid,
			})

			month++
		}

		c.JSON(http.StatusOK, response)
	})
	r.Run(":8080")
}
