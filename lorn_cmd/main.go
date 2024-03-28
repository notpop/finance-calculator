package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("元金を入力してください (例: 3000000): ")
	principalStr, _ := reader.ReadString('\n')
	principal, _ := strconv.ParseFloat(strings.TrimSpace(principalStr), 64)

	fmt.Print("年利率を入力してください (例: 10): ")
	interestRateStr, _ := reader.ReadString('\n')
	interestRate, _ := strconv.ParseFloat(strings.TrimSpace(interestRateStr), 64)
	monthlyInterestRate := interestRate / 100 / 12

	fmt.Print("月額支払いを入力してください (例: 100000): ")
	paymentStr, _ := reader.ReadString('\n')
	monthlyPayment, _ := strconv.ParseFloat(strings.TrimSpace(paymentStr), 64)

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

		fmt.Printf("%dヶ月目: 元金 %.2f 利率 %.2f%% 月額 %.2f 利息 %.2f 返済額 %.2f 累計返済額 %.2f\n", month, principal, interestRate, thisPayment, interest, principalReduction, totalPaid)
		month++
	}
}
