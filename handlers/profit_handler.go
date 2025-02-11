package handlers

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

func integrateGaussian(Pc, sigma, lowerBound, upperBound float64) float64 {
	n := 1000
	h := (upperBound - lowerBound) / float64(n)
	result := 0.0

	for i := 0; i <= n; i++ {
		x := lowerBound + float64(i)*h
		weight := 1.0
		if i%2 == 1 {
			weight = 4
		} else if i > 0 && i < n {
			weight = 2
		}
		result += weight * (1 / (sigma * math.Sqrt(2*math.Pi))) * math.Exp(-math.Pow(x-Pc, 2)/(2*math.Pow(sigma, 2)))
	}
	return result * h / 3
}

func calculateProfit(Pc, sigma, B float64) float64 {
	deltaP := 0.05 * Pc
	lowerBound := Pc - deltaP
	upperBound := Pc + deltaP

	// Частка енергії без небалансів
	energyWithoutImbalance := integrateGaussian(Pc, sigma, lowerBound, upperBound)

	W1 := Pc * 24 * energyWithoutImbalance
	P1 := W1 * B // Прибуток від енергії без небалансів

	W2 := Pc * 24 * (1 - energyWithoutImbalance)
	penalty := W2 * B // Штраф за небаланси

	return P1 - penalty
}

func ProfitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		Pc, _ := strconv.ParseFloat(r.FormValue("Pc"), 64)
		oldSigma, _ := strconv.ParseFloat(r.FormValue("oldSigma"), 64)
		newSigma, _ := strconv.ParseFloat(r.FormValue("newSigma"), 64)
		B, _ := strconv.ParseFloat(r.FormValue("B"), 64)

		if Pc > 0 && oldSigma > 0 && newSigma > 0 && B > 0 {
			oldProfit := calculateProfit(Pc, oldSigma, B)
			newProfit := calculateProfit(Pc, newSigma, B)
			profitIncrease := newProfit - oldProfit

			tmpl, _ := template.ParseFiles("templates/profit.html")
			tmpl.Execute(w, map[string]string{
				"oldProfit":      fmt.Sprintf("%.2f", oldProfit),
				"newProfit":      fmt.Sprintf("%.2f", newProfit),
				"profitIncrease": fmt.Sprintf("%.2f", profitIncrease),
			})
			return
		}
	}

	// Відображення HTML-сторінки
	tmpl, _ := template.ParseFiles("templates/profitCalc.html")
	tmpl.Execute(w, nil)
}
