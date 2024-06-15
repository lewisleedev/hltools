package main

import (
	"fmt"
	"math"
)

func calculateStraightlineDepreciation(initValue int, scrapValue int, period int) ([]int, error) {
	if period <= 0 {
		return nil, fmt.Errorf("period must be greater than 0")
	}
	if initValue <= scrapValue {
		return nil, fmt.Errorf("initial value must be greater than scrap value")
	}

	amountDepreciation := initValue - scrapValue
	singleDepreciation := math.Floor(float64(amountDepreciation) / float64(period))
	depreciations := make([]int, period)

	for i := 0; i < period-1; i++ {
		depreciations[i] = int(singleDepreciation)
	}
	// Adjust the last depreciation to match the scrap value
	depreciations[period-1] = initValue - scrapValue - int(singleDepreciation)*(period-1)

	return depreciations, nil

}

func calculateDoubleDecliningBalance(initValue int, scrapValue int, period int) ([]int, error) {
	if period <= 0 {
		return nil, fmt.Errorf("period must be greater than 0")
	}
	if initValue <= scrapValue {
		return nil, fmt.Errorf("initial value must be greater than scrap value")
	}

	depreciations := make([]int, period)
	currentValue := float64(initValue)

	for i := 0; i < period; i++ {
		depreciation := 2.0 / float64(period) * currentValue
		if currentValue-depreciation < float64(scrapValue) {
			depreciation = currentValue - float64(scrapValue)
		}
		depreciations[i] = int(depreciation)
		currentValue -= depreciation
	}

	return depreciations, nil
}

func calculateSumOfTheYearsDigits(initValue int, scrapValue int, period int) ([]int, error) {
	if period <= 0 {
		return nil, fmt.Errorf("period must be greater than 0")
	}
	if initValue <= scrapValue {
		return nil, fmt.Errorf("initial value must be greater than scrap value")
	}

	amountDepreciation := initValue - scrapValue
	depreciations := make([]int, period)
	sumYears := period * (period + 1) / 2

	for i := 0; i < period; i++ {
		depreciations[i] = (period - i) * amountDepreciation / sumYears
	}
	return depreciations, nil
}
