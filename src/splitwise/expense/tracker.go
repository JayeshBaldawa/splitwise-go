package expense

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"splitwise/global"
	"splitwise/model"
)

type TrackerServiceImpl struct {
	storage *global.GlobalMapStorage
}

const MinimumHousemates = 2

// Member represents a housemate with a name and their dues.
type Member struct {
	name string
	dues int64
}

// NewTrackerServiceImpl initializes a new TrackerServiceImpl with the given storage.
func NewTrackerServiceImpl(storage *global.GlobalMapStorage) *TrackerServiceImpl {
	return &TrackerServiceImpl{
		storage: storage,
	}
}

// AddExpense adds an expense to the system and updates the dues of the beneficiaries.
func (t *TrackerServiceImpl) AddExpense(amount float64, beneficiaries []string) (string, error) {
	amountPerPerson := t.calculateDues(amount, len(beneficiaries))
	payer := beneficiaries[0]

	if err := t.validateHousemateExists(payer); err != nil {
		return "", err
	}

	if err := t.updateDuesForBeneficiaries(payer, beneficiaries[1:], amountPerPerson); err != nil {
		return "", err
	}

	t.storage.SimplifyDebt()
	return string(model.SUCCESS), nil
}

func (t *TrackerServiceImpl) validateHousemateExists(housemate string) error {
	if !t.storage.CheckHousemateExists(housemate) {
		return errors.New(string(model.MEMBER_NOT_FOUND))
	}
	return nil
}

func (t *TrackerServiceImpl) updateDuesForBeneficiaries(payer string, beneficiaries []string, amountPerPerson int64) error {
	for _, beneficiary := range beneficiaries {
		if err := t.validateHousemateExists(beneficiary); err != nil {
			return err
		}
		t.storage.AddOrUpdateDue(payer, beneficiary, amountPerPerson)
	}
	return nil
}

// ShowDues returns the list of housemates with their dues in descending order.
// If two housemates have the same dues, they are sorted in ascending order of their names.
func (t *TrackerServiceImpl) ShowDues(housemate string) ([]string, error) {
	if err := t.validateHousemateExists(housemate); err != nil {
		return nil, err
	}

	dues := t.storage.GetAllDues(housemate)
	members := t.sortMembersByDues(dues)

	return t.formatResult(members), nil
}

// ShowAllDues returns the list of all housemates with their dues in descending order.
func (t *TrackerServiceImpl) formatResult(members []Member) []string {
	result := make([]string, 0, len(members))
	for _, member := range members {
		result = append(result, fmt.Sprintf("%s %d", member.name, member.dues))
	}
	return result
}

// ClearDues clears a specified amount of dues between two housemates.
func (t *TrackerServiceImpl) ClearDues(from, to string, amount int64) (string, error) {
	if !t.storage.CheckHousemateExists(from) || !t.storage.CheckHousemateExists(to) {
		return "", errors.New(string(model.MEMBER_NOT_FOUND))
	}

	dues := t.storage.GetDue(from, to)

	if amount > dues {
		return "", errors.New(string(model.INCORRECT_PAYMENT))
	}

	t.storage.ClearDues(from, to, amount)

	return fmt.Sprintf("%d", dues-amount), nil
}

// calculateDues calculates the amount each beneficiary should pay.
func (t *TrackerServiceImpl) calculateDues(amount float64, beneficiaryCount int) int64 {
	return int64(math.Round(amount / float64(beneficiaryCount)))
}

// sortMembersByDues sorts housemates by their dues in descending order and by name in ascending order for ties.
func (t *TrackerServiceImpl) sortMembersByDues(dues map[string]int64) []Member {
	members := t.mapToMembers(dues)
	t.sortMembers(members)
	return members
}

// mapToMembers converts a map of dues to a slice of Member structs.
func (t *TrackerServiceImpl) mapToMembers(dues map[string]int64) []Member {
	members := make([]Member, 0, len(dues))
	for name, due := range dues {
		members = append(members, Member{name: name, dues: due})
	}
	return members
}

// sortMembers sorts the members first by dues descending, then by name ascending.
func (t *TrackerServiceImpl) sortMembers(members []Member) {
	sort.Slice(members, func(i, j int) bool {
		if members[i].dues == members[j].dues {
			return members[i].name < members[j].name
		}
		return members[i].dues > members[j].dues
	})
}
