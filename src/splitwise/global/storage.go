package global

import (
	"math"
	"sort"
	"splitwise/model"
)

// GlobalMapStorage encapsulates all the housemates and their dues
type GlobalMapStorage struct {
	housemates   map[string]bool
	dues         map[string]map[string]int64
	simplifydues map[string]map[string]int64
}

// NewGlobalMapStorage initializes a new GlobalMapStorage with empty maps
func NewGlobalMapStorage() *GlobalMapStorage {
	return &GlobalMapStorage{
		housemates:   make(map[string]bool),
		dues:         make(map[string]map[string]int64),
		simplifydues: make(map[string]map[string]int64),
	}
}

// AddHousemate adds a new housemate to the system
func (g *GlobalMapStorage) AddHousemate(housemate string) {
	g.housemates[housemate] = true
	g.dues[housemate] = make(map[string]int64)
	g.simplifydues[housemate] = make(map[string]int64)
	for name := range g.housemates {
		if name == housemate {
			continue
		}
		g.dues[housemate][name] = 0
		g.dues[name][housemate] = 0
		g.simplifydues[housemate][name] = 0
		g.simplifydues[name][housemate] = 0
	}
}

// RemoveHousemate removes a housemate and cleans up their dues
func (g *GlobalMapStorage) RemoveHousemate(housemate string) {
	delete(g.housemates, housemate)
	g.cleanupDues(housemate)
	delete(g.dues, housemate)
	delete(g.simplifydues, housemate)
}

// cleanupDues removes all dues related to a housemate
func (g *GlobalMapStorage) cleanupDues(housemate string) {
	for name := range g.dues {
		if name != housemate {
			delete(g.dues[name], housemate)
			delete(g.simplifydues[name], housemate)
		}
	}
}

// AddOrUpdateDue adds or updates a due from one housemate to another
func (g *GlobalMapStorage) AddOrUpdateDue(from, to string, amount int64) {
	adjustDue(g.dues, from, to, amount)
	if amount == model.ZERO_DUE {
		adjustDue(g.simplifydues, from, to, amount)
	}
}

// SimplifyDebt simplifies all dues by minimizing the transactions
func (g *GlobalMapStorage) SimplifyDebt() {
	inOuts := g.calculateNetBalances()
	g.resetSimplifiedDues()
	g.minimizeTransactions(inOuts)
}

// minimizeTransactions reduces the number of transactions required to settle debts
func (g *GlobalMapStorage) minimizeTransactions(balances map[string]int64) {
	nonZeroBalances := extractNonZeroBalances(balances)

	if len(nonZeroBalances) == 0 {
		return
	}

	minAmount, maxAmount := findMinMaxBalances(nonZeroBalances)
	minHousemate := findHousemateByBalance(balances, minAmount)
	maxHousemate := findHousemateByBalance(balances, maxAmount)

	g.handleTransaction(balances, minHousemate, maxHousemate, minAmount, maxAmount)
}

// handleTransaction processes a transaction between two housemates
func (g *GlobalMapStorage) handleTransaction(balances map[string]int64, minHousemate, maxHousemate string, minAmount, maxAmount int64) {
	leftAmount := maxAmount + minAmount
	if leftAmount >= 0 {
		g.processPositiveTransaction(balances, minHousemate, maxHousemate, minAmount, leftAmount)
	} else {
		g.processNegativeTransaction(balances, minHousemate, maxHousemate, maxAmount, leftAmount)
	}
	g.minimizeTransactions(balances)
}

// processPositiveTransaction processes a positive transaction between two housemates
func (g *GlobalMapStorage) processPositiveTransaction(balances map[string]int64, minHousemate, maxHousemate string, minAmount, leftAmount int64) {
	balances[minHousemate] = 0
	balances[maxHousemate] = leftAmount
	g.simplifydues[minHousemate][maxHousemate] = int64(math.Abs(float64(minAmount)))
}

// processNegativeTransaction processes a negative transaction between two housemates
func (g *GlobalMapStorage) processNegativeTransaction(balances map[string]int64, minHousemate, maxHousemate string, maxAmount, leftAmount int64) {
	balances[minHousemate] = leftAmount
	balances[maxHousemate] = 0
	g.simplifydues[minHousemate][maxHousemate] = maxAmount
}

// extractNonZeroBalances extracts non-zero balances from a map
func extractNonZeroBalances(balances map[string]int64) []int64 {
	var nonZeroBalances []int64
	for _, balance := range balances {
		if balance != model.ZERO_DUE {
			nonZeroBalances = append(nonZeroBalances, balance)
		}
	}
	return nonZeroBalances
}

// findMinMaxBalances finds the minimum and maximum balances from a list
func findMinMaxBalances(balances []int64) (int64, int64) {
	sort.Slice(balances, func(i, j int) bool {
		return balances[i] < balances[j]
	})
	return balances[0], balances[len(balances)-1]
}

// resetSimplifiedDues resets all simplified dues to zero
func (g *GlobalMapStorage) resetSimplifiedDues() {
	for _, transaction := range g.simplifydues {
		for housemate := range transaction {
			transaction[housemate] = model.ZERO_DUE
		}
	}
}

// findHousemateByBalance finds a housemate by their balance value
func findHousemateByBalance(balances map[string]int64, value int64) string {
	for name, balance := range balances {
		if balance == value {
			return name
		}
	}
	return ""
}

// calculateNetBalances calculates the net balances for all housemates
func (g *GlobalMapStorage) calculateNetBalances() map[string]int64 {
	netBalances := make(map[string]int64)
	for housemate := range g.housemates {
		netBalances[housemate] = model.ZERO_DUE
		netBalances[housemate] += g.GetInAmount(housemate) - g.GetOutAmount(housemate)
	}
	return netBalances
}

func (g *GlobalMapStorage) GetInAmount(housemate string) int64 {
	var amount int64
	for _, due := range g.dues[housemate] {
		amount += due
	}
	return amount
}

func (g *GlobalMapStorage) GetOutAmount(housemate string) int64 {
	var amount int64
	for _, due := range g.dues {
		amount += due[housemate]
	}
	return amount
}

// ClearDues clears a due between two housemates
func (g *GlobalMapStorage) ClearDues(from, to string, amount int64) {
	adjustDue(g.dues, from, to, -amount)
	adjustDue(g.simplifydues, from, to, -amount)
}

// GetNumberOfHousemates returns the number of housemates
func (g *GlobalMapStorage) GetNumberOfHousemates() int {
	return len(g.housemates)
}

// GetAllDues returns a copy of all dues for a given housemate
func (g *GlobalMapStorage) GetAllDues(housemate string) map[string]int64 {
	copy := make(map[string]int64)
	for k, v := range g.simplifydues[housemate] {
		copy[k] = v
	}
	return copy
}

// GetTransactions returns all simplified dues
func (g *GlobalMapStorage) GetTransactions() map[string]map[string]int64 {
	return g.simplifydues
}

// GetDue returns the due between two housemates
func (g *GlobalMapStorage) GetDue(from, to string) int64 {
	return g.simplifydues[from][to]
}

// GetNonShuffledDue returns the due between two housemates in the original map
func (g *GlobalMapStorage) GetNonShuffledDue(from, to string) int64 {
	return g.dues[from][to]
}

// CheckHousemateExists checks if a housemate exists
func (g *GlobalMapStorage) CheckHousemateExists(housemate string) bool {
	return g.housemates[housemate]
}

// GetHousemateNames returns a list of all housemate names
func (g *GlobalMapStorage) GetHousemateNames() []string {
	var housemates []string
	for housemate := range g.housemates {
		housemates = append(housemates, housemate)
	}
	return housemates
}

// adjustDue is a helper function to adjust dues between housemates
func adjustDue(mapData map[string]map[string]int64, from, to string, amount int64) {
	newAmount := mapData[from][to] + amount
	if newAmount < 0 {
		return
	}
	mapData[from][to] = newAmount
}

// Reset resets the storage to its initial state
func (g *GlobalMapStorage) Reset() {
	g.housemates = make(map[string]bool)
	g.dues = make(map[string]map[string]int64)
	g.simplifydues = make(map[string]map[string]int64)
}
