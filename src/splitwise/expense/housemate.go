package expense

import (
	"errors"
	"splitwise/global"
	"splitwise/model"
)

type HousemateServiceImpl struct {
	storage *global.GlobalMapStorage
}

// NewHousemateServiceImpl creates a new instance of HousemateServiceImpl with the provided storage.
func NewHousemateServiceImpl(storage *global.GlobalMapStorage) *HousemateServiceImpl {
	return &HousemateServiceImpl{
		storage: storage,
	}
}

// MoveIn adds a housemate to the house and initializes dues with existing housemates.
func (h *HousemateServiceImpl) MoveIn(housemate string) (string, error) {
	if h.storage.CheckHousemateExists(housemate) {
		return "", errors.New(string(model.MEMBER_ALREADY_EXISTS))
	}
	if h.isRoomFull() {
		return "", errors.New(string(model.HOUSEFUL))
	}
	// Add the new housemate and initialize dues.
	h.storage.AddHousemate(housemate)
	return string(model.SUCCESS), nil
}

// MoveOut removes a housemate from the house and clears their dues.
func (h *HousemateServiceImpl) MoveOut(housemate string) (string, error) {
	if !h.storage.CheckHousemateExists(housemate) {
		return "", errors.New(string(model.MEMBER_NOT_FOUND))
	}
	if h.hasPendingDue(housemate) {
		return "", errors.New(string(model.FAILURE))
	}
	h.storage.RemoveHousemate(housemate)
	return string(model.SUCCESS), nil
}

// isRoomFull checks if the house has reached its maximum capacity.
func (h *HousemateServiceImpl) isRoomFull() bool {
	return h.storage.GetNumberOfHousemates() == model.MAX_HOUSEMATES
}

// hasPendingDue checks if a housemate has any pending dues.
func (h *HousemateServiceImpl) hasPendingDue(name string) bool {
	transactions := h.storage.GetTransactions()

	// Check if the housemate owes or is owed any positive amount
	return h.hasPositiveDue(transactions[name]) || h.isOwedPositiveDue(name, transactions)
}

// hasPositiveDue checks if the housemate has any positive dues in their transactions.
func (h *HousemateServiceImpl) hasPositiveDue(transactions map[string]int64) bool {
	for _, due := range transactions {
		if due > 0 {
			return true
		}
	}
	return false
}

// isOwedPositiveDue checks if any housemate owes the given housemate a positive due.
func (h *HousemateServiceImpl) isOwedPositiveDue(name string, transactions map[string]map[string]int64) bool {
	for _, dues := range transactions {
		if dues[name] > 0 {
			return true
		}
	}
	return false
}
