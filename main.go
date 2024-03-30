package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LoanDetails struct {
	Month              int     `json:"month"`
	Year               float64 `json:"year"`
	OriginalPrincipal  float64 `json:"originalPrincipal"`
	RemainingPrincipal float64 `json:"remainingPrincipal"`
	InterestRate       float64 `json:"interestRate"`
	MonthlyPayment     float64 `json:"monthlyPayment"`
	Interest           float64 `json:"interest"`
	PrincipalReduction float64 `json:"principalReduction"`
	TotalPaid          float64 `json:"totalPaid"`
}

type LoanRequest struct {
	Principal      float64 `json:"principal"`
	InterestRate   float64 `json:"interestRate"`
	MonthlyPayment float64 `json:"monthlyPayment"`
}

type MultiLoanDetails struct {
	LoanDetails             []LoanDetails `json:"loanDetails"`
	TotalOriginalPrincipal  float64       `json:"totalOriginalPrincipal"`
	TotalRemainingPrincipal float64       `json:"totalRemainingPrincipal"`
	TotalMonthlyPayment     float64       `json:"totalMonthlyPayment"`
	TotalInterest           float64       `json:"totalInterest"`
	TotalPrincipalReduction float64       `json:"totalPrincipalReduction"`
	TotalPaid               float64       `json:"totalPaid"`
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func extractLoanDetails(loans [][]LoanDetails, month int) []LoanDetails {
	var details []LoanDetails
	for i, loan := range loans {
		if month < len(loan) {
			detail := loan[month]
			if month > 0 {
				detail.OriginalPrincipal = loans[i][month-1].RemainingPrincipal
			}
			details = append(details, detail)
		}
	}
	return details
}

func calculateMultiLoan(requests []LoanRequest) []MultiLoanDetails {
	maxMonths := 0
	loans := make([][]LoanDetails, len(requests))

	for i, request := range requests {
		monthlyInterestRate := request.InterestRate / 100 / 12
		principal := request.Principal
		totalPaid := 0.0
		month := 1

		for principal > 0 {
			interest := principal * monthlyInterestRate
			thisPayment := min(request.MonthlyPayment, principal+interest)
			principalReduction := thisPayment - interest
			totalPaid += thisPayment

			loans[i] = append(loans[i], LoanDetails{
				Month:              month,
				Year:               float64(month) / 12,
				OriginalPrincipal:  principal,
				RemainingPrincipal: principal - principalReduction,
				InterestRate:       request.InterestRate,
				MonthlyPayment:     thisPayment,
				Interest:           interest,
				PrincipalReduction: principalReduction,
				TotalPaid:          totalPaid,
			})

			principal -= principalReduction
			month++
		}

		if month-1 > maxMonths {
			maxMonths = month - 1
		}
	}

	return aggregateLoanDetails(loans, maxMonths)
}

func aggregateLoanDetails(loans [][]LoanDetails, maxMonths int) []MultiLoanDetails {
	multiLoanDetails := make([]MultiLoanDetails, maxMonths)
	totalPaidPerLoan := make([]float64, len(loans))

	for month := 0; month < maxMonths; month++ {
		var totalMonthlyPayment, totalInterest, totalPrincipalReduction, totalOriginalPrincipal, totalRemainingPrincipal, totalPaid float64

		for i, loan := range loans {
			if month < len(loan) {
				detail := loan[month]
				totalMonthlyPayment += detail.MonthlyPayment
				totalInterest += detail.Interest
				totalPrincipalReduction += detail.PrincipalReduction
				totalOriginalPrincipal += detail.OriginalPrincipal
				totalRemainingPrincipal += detail.RemainingPrincipal
				totalPaidPerLoan[i] = detail.TotalPaid
			}
			totalPaid += totalPaidPerLoan[i]
		}

		multiLoanDetails[month] = MultiLoanDetails{
			LoanDetails:             extractLoanDetails(loans, month),
			TotalOriginalPrincipal:  totalOriginalPrincipal,
			TotalRemainingPrincipal: totalRemainingPrincipal,
			TotalMonthlyPayment:     totalMonthlyPayment,
			TotalInterest:           totalInterest,
			TotalPrincipalReduction: totalPrincipalReduction,
			TotalPaid:               totalPaid,
		}
	}
	return multiLoanDetails
}

func roundToTwoDecimals(value float64) float64 {
	rounded, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return rounded
}

func formatMultiLoanDetails(details []MultiLoanDetails) []MultiLoanDetails {
	for i := range details {
		details[i].TotalOriginalPrincipal = roundToTwoDecimals(details[i].TotalOriginalPrincipal)
		details[i].TotalRemainingPrincipal = roundToTwoDecimals(details[i].TotalRemainingPrincipal)
		details[i].TotalMonthlyPayment = roundToTwoDecimals(details[i].TotalMonthlyPayment)
		details[i].TotalInterest = roundToTwoDecimals(details[i].TotalInterest)
		details[i].TotalPrincipalReduction = roundToTwoDecimals(details[i].TotalPrincipalReduction)
		details[i].TotalPaid = roundToTwoDecimals(details[i].TotalPaid)

		for j := range details[i].LoanDetails {
			details[i].LoanDetails[j].Year = roundToTwoDecimals(details[i].LoanDetails[j].Year)
			details[i].LoanDetails[j].OriginalPrincipal = roundToTwoDecimals(details[i].LoanDetails[j].OriginalPrincipal)
			details[i].LoanDetails[j].RemainingPrincipal = roundToTwoDecimals(details[i].LoanDetails[j].RemainingPrincipal)
			details[i].LoanDetails[j].MonthlyPayment = roundToTwoDecimals(details[i].LoanDetails[j].MonthlyPayment)
			details[i].LoanDetails[j].Interest = roundToTwoDecimals(details[i].LoanDetails[j].Interest)
			details[i].LoanDetails[j].PrincipalReduction = roundToTwoDecimals(details[i].LoanDetails[j].PrincipalReduction)
			details[i].LoanDetails[j].TotalPaid = roundToTwoDecimals(details[i].LoanDetails[j].TotalPaid)
		}
	}
	return details
}

func main() {
	r := gin.Default()

	r.POST("/loan", func(c *gin.Context) {
		var requests []LoanRequest
		if err := c.ShouldBindJSON(&requests); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		response := calculateMultiLoan(requests)
		formattedResponse := formatMultiLoanDetails(response)
		c.JSON(http.StatusOK, formattedResponse)
	})

	r.Run(":8080")
}
