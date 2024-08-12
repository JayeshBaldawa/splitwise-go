package global

import "testing"

func TestAddOrUpdateDue(t *testing.T) {
	globalStorage := NewGlobalMapStorage()

	// Add housemates
	globalStorage.AddHousemate("Andy")
	globalStorage.AddHousemate("Woody")

	// TEST CASE 1: Default dues
	if globalStorage.dues["Andy"]["Woody"] != 0 {
		t.Errorf("Expected 0, got %d", globalStorage.dues["Andy"]["Woody"])
	}

	if globalStorage.dues["Woody"]["Andy"] != 0 {
		t.Errorf("Expected 0, got %d", globalStorage.dues["Woody"]["Andy"])
	}

	// TEST CASE 2: Add dues for the first time
	globalStorage.AddOrUpdateDue("Andy", "Woody", 1000)
	globalStorage.AddOrUpdateDue("Woody", "Andy", 500)

	// Check dues
	if globalStorage.dues["Andy"]["Woody"] != 1000 {
		t.Errorf("Expected 1000, got %d", globalStorage.dues["Andy"]["Woody"])
	}

	if globalStorage.dues["Woody"]["Andy"] != 500 {
		t.Errorf("Expected 500, got %d", globalStorage.dues["Woody"]["Andy"])
	}

	// TEST CASE 3: Update dues for existing housemates
	globalStorage.AddOrUpdateDue("Andy", "Woody", 2000)
	globalStorage.AddOrUpdateDue("Woody", "Andy", 1000)

	// Check dues
	if globalStorage.dues["Andy"]["Woody"] != 3000 {
		t.Errorf("Expected 3000, got %d", globalStorage.dues["Andy"]["Woody"])
	}

	if globalStorage.dues["Woody"]["Andy"] != 1500 {
		t.Errorf("Expected 1500, got %d", globalStorage.dues["Woody"]["Andy"])
	}
}

func TestRemoveHousemate(t *testing.T) {
	globalStorage := NewGlobalMapStorage()

	// Add housemates
	globalStorage.AddHousemate("Andy")
	globalStorage.AddHousemate("Woody")

	// Remove housemate
	globalStorage.RemoveHousemate("Andy")

	// TEST CASE 1: Check if housemate is removed
	if globalStorage.CheckHousemateExists("Andy") {
		t.Errorf("Expected housemate Andy to be removed")
	}

	// TEST CASE 2: Check if dues are cleared for the removed housemate
	if len(globalStorage.GetAllDues("Andy")) != 0 {
		t.Errorf("Expected no dues for removed housemate Andy")
	}

	// TEST CASE 3: Check if dues are cleared for other housemates
	if len(globalStorage.GetAllDues("Woody")) != 0 {
		t.Errorf("Expected no dues for housemate Woody after removing Andy")
	}

	// TEST CASE 4: Check if transactions are cleared for the removed housemate
	if len(globalStorage.simplifydues["Andy"]) != 0 {
		t.Errorf("Expected no transactions for removed housemate Andy")
	}

	// TEST CASE 5: Check if transactions are cleared for other housemates
	if len(globalStorage.simplifydues["Woody"]) != 0 {
		t.Errorf("Expected no transactions for housemate Woody after removing Andy")
	}

}

func TestClearDues(t *testing.T) {
	globalStorage := NewGlobalMapStorage()

	// Add housemates
	globalStorage.AddHousemate("Andy")
	globalStorage.AddHousemate("Woody")

	// Add dues for housemates
	globalStorage.AddOrUpdateDue("Andy", "Woody", 1000)
	globalStorage.AddOrUpdateDue("Woody", "Andy", 500)

	// Simplify debt before clearing dues
	globalStorage.SimplifyDebt()

	// TEST CASE 1: Clear dues between housemates
	globalStorage.ClearDues("Andy", "Woody", 500)

	// Check if dues are cleared for Andy (Woody -> Andy = 500)
	if data := globalStorage.GetNonShuffledDue("Andy", "Woody"); data != 500 {
		t.Errorf("Expected dues between Andy and Woody to be 500, got %d", data)
	}

	// Check if due are still present for Woody (Woody -> Andy = 500)
	if globalStorage.simplifydues["Woody"]["Andy"] != 500 {
		t.Errorf("Expected dues between Woody and Andy to be 500, got %d", globalStorage.simplifydues["Woody"]["Andy"])
	}

	// TEST CASE 2: Clear dues between housemates
	globalStorage.ClearDues("Woody", "Andy", 500)

	// Check if dues are cleared for Woody (Woody -> Andy = 0)
	if data := globalStorage.GetNonShuffledDue("Woody", "Andy"); data != 0 {
		t.Errorf("Expected dues between Woody and Andy to be 0, got %d", data)
	}

	// Check if dues are cleared for Woody (Woody -> Andy = 0)
	if data := globalStorage.GetNonShuffledDue("Woody", "Andy"); data != 0 {
		t.Errorf("Expected dues between Woody and Andy to be 0, got %d", data)
	}
}

func TestSimplifyDebt(t *testing.T) {
	globalStorage := NewGlobalMapStorage()

	// Add housemates
	globalStorage.AddHousemate("Andy")
	globalStorage.AddHousemate("Woody")
	globalStorage.AddHousemate("Buzz")

	// Add dues
	globalStorage.AddOrUpdateDue("Andy", "Woody", 1000)
	globalStorage.AddOrUpdateDue("Woody", "Andy", 500)
	globalStorage.AddOrUpdateDue("Buzz", "Andy", 2000)
	globalStorage.AddOrUpdateDue("Buzz", "Woody", 1500)

	// Simplify debt
	globalStorage.SimplifyDebt()

	// TEST CASE 1: Check if dues are simplified
	if globalStorage.GetDue("Andy", "Buzz") != 1500 {
		t.Errorf("Expected dues between Andy and Buzz to be 1500, got %d", globalStorage.simplifydues["Andy"]["Buzz"])
	}

	if globalStorage.GetDue("Woody", "Buzz") != 2000 {
		t.Errorf("Expected dues between Woody and Andy to be 2000, got %d", globalStorage.simplifydues["Woody"]["Andy"])
	}

	if globalStorage.GetDue("Buzz", "Woody") != 0 {
		t.Errorf("Expected dues between Buzz and Woody to be 0, got %d", globalStorage.simplifydues["Buzz"]["Woody"])
	}

	if globalStorage.GetDue("Andy", "Woody") != 0 {
		t.Errorf("Expected dues between Andy and Woody to be 0, got %d", globalStorage.simplifydues["Andy"]["Woody"])
	}

	if globalStorage.GetDue("Buzz", "Andy") != 0 {
		t.Errorf("Expected dues between Buzz and Andy to be 0, got %d", globalStorage.simplifydues["Buzz"]["Andy"])
	}

	// Add more dues
	globalStorage.AddOrUpdateDue("Andy", "Buzz", 1500)

	// Simplify debt
	globalStorage.SimplifyDebt()

	// TEST CASE 2: Check if dues are simplified
	if globalStorage.GetDue("Andy", "Buzz") != 0 {
		t.Errorf("Expected dues between Buzz and Woody to be 0, got %d", globalStorage.simplifydues["Buzz"]["Woody"])
	}

	if globalStorage.GetDue("Buzz", "Andy") != 0 {
		t.Errorf("Expected dues between Buzz and Andy to be 0, got %d", globalStorage.simplifydues["Buzz"]["Andy"])
	}

	if globalStorage.GetDue("Woody", "Buzz") != 2000 {
		t.Errorf("Expected dues between Woody and Buzz to be 2000, got %d", globalStorage.simplifydues["Woody"]["Buzz"])
	}

	// TEST CASE 3: Remove housemate
	globalStorage.RemoveHousemate("Andy")

	// check if any dues are present for Andy
	if len(globalStorage.GetAllDues("Andy")) != 0 {
		t.Errorf("Expected no dues for removed housemate Andy")
	}

}
