package data

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestSaveState(t *testing.T) {
	defer os.Remove("state.bin")

	state := State{
		LastId: 0,
		Count:  0,
	}

	SaveState(state)

	loadedState, err := LoadState()

	if err != nil {
		t.Errorf("Test failed %v", err)
	}

	if loadedState != state {
		t.Errorf("%v is not equal to %v", state, loadedState)
	}

}

func TestSaveEntry(t *testing.T) {
	defer os.Remove("data.bin")

	entry := Entry{
		Id:       10,
		Title:    "some title",
		Username: "some username",
		Password: "some password",
		Address:  "http://some-address.com",
		Notes:    "note bla bla bla",
	}

	SaveEntry(entry)

	loadedEntry, _ := LoadEntry(10)

	if loadedEntry != entry {
		t.Errorf("%v is not equal to %v", entry, loadedEntry)
	}
}

func TestCountEntries(t *testing.T) {
	// Defer the removal of the test file after the test completes
	defer func() {
		if err := os.Remove("data.bin"); err != nil {
			t.Errorf("Failed to remove test file: %v", err)
		}
	}()

	// Initialize entry with a test data structure
	entry := Entry{
		Id:       1,
		Title:    "some title",
		Username: "some username",
		Password: "some password",
		Address:  "http://some-address.com",
		Notes:    "note bla bla bla",
	}

	// Save the first entry
	if err := SaveEntry(entry); err != nil {
		t.Fatalf("SaveEntry failed: %v", err)
	}

	// Save the second entry with a different ID
	entry.Id = 2
	if err := SaveEntry(entry); err != nil {
		t.Fatalf("SaveEntry failed: %v", err)
	}

	// Save the third entry with a different ID
	entry.Id = 3
	if err := SaveEntry(entry); err != nil {
		t.Fatalf("SaveEntry failed: %v", err)
	}

	// Count the entries in the file
	count, err := CountEntries()
	if err != nil {
		t.Fatalf("CountEntries failed: %v", err)
	}

	// Verify if the count matches the expected value
	if count != 3 {
		t.Errorf("Expected count to be 3 but got %v", count)
	}
}

func TestDeleteEntryWithLargeFile(t *testing.T) {
	defer os.Remove("data.bin")

	// Add multiple entries to the file
	for i := int64(1); i <= 100; i++ {
		entry := Entry{
			Id:       i,
			Title:    "Title " + fmt.Sprint(i),
			Username: "user" + fmt.Sprint(i),
			Password: "pass" + fmt.Sprint(i),
			Address:  "http://address" + fmt.Sprint(i) + ".com",
			Notes:    "Notes " + fmt.Sprint(i),
		}
		err := SaveEntry(entry)
		if err != nil {
			t.Fatalf("SaveEntry failed for ID %v: %v", i, err)
		}
	}

	// Now delete an entry with Id 50
	err := DeleteEntry(50)
	if err != nil {
		t.Errorf("DeleteEntry failed: %v", err)
	}

	// Count remaining entries, should be 99
	count, err := CountEntries()
	if err != nil {
		t.Errorf("CountEntries failed: %v", err)
	}
	if count != 99 {
		t.Errorf("Expected 99 entries, but got %v", count)
	}

	// Verify that entry with ID 50 is no longer present
	_, err = LoadEntry(50)
	if err == nil {
		t.Errorf("Expected error loading entry with ID 50, but got none")
	}
}

func TestDeleteEntryInMemory(t *testing.T) {
	defer os.Remove("data.bin")

	// Add multiple entries to the file
	for i := int64(1); i <= 10; i++ {
		entry := Entry{
			Id:       i,
			Title:    "Title " + fmt.Sprint(i),
			Username: "user" + fmt.Sprint(i),
			Password: "pass" + fmt.Sprint(i),
			Address:  "http://address" + fmt.Sprint(i) + ".com",
			Notes:    "Notes " + fmt.Sprint(i),
		}
		err := SaveEntry(entry)
		if err != nil {
			t.Fatalf("SaveEntry failed for ID %v: %v", i, err)
		}
	}

	// Now delete an entry with Id 50
	err := DeleteEntryInMemory(5)
	if err != nil {
		t.Errorf("DeleteEntry failed: %v", err)
	}

	// Count remaining entries, should be 99
	count, err := CountEntries()
	if err != nil {
		t.Errorf("CountEntries failed: %v", err)
	}
	if count != 9 {
		t.Errorf("Expected 99 entries, but got %v", count)
	}

	// Verify that entry with ID 50 is no longer present
	_, err = LoadEntry(5)
	if err == nil {
		t.Errorf("Expected error loading entry with ID 50, but got none")
	}
}

func TestUpdateEntry(t *testing.T) {
	// Clean up test files after the test
	defer os.Remove("data.bin")
	defer os.Remove("temp_data.bin")

	// Helper function to add an entry to the file
	addEntry := func(entry Entry) error {
		return SaveEntry(entry)
	}

	// Add some entries
	entry1 := Entry{
		Id:       1,
		Title:    "Old Title 1",
		Username: "Old User 1",
		Password: "Old Password 1",
		Address:  "http://old-address1.com",
		Notes:    "Old notes 1",
	}
	entry2 := Entry{
		Id:       2,
		Title:    "Old Title 2",
		Username: "Old User 2",
		Password: "Old Password 2",
		Address:  "http://old-address2.com",
		Notes:    "Old notes 2",
	}
	entry3 := Entry{
		Id:       3,
		Title:    "Old Title 3",
		Username: "Old User 3",
		Password: "Old Password 3",
		Address:  "http://old-address3.com",
		Notes:    "Old notes 3",
	}

	if err := addEntry(entry1); err != nil {
		t.Fatalf("Failed to add entry 1: %v", err)
	}
	if err := addEntry(entry2); err != nil {
		t.Fatalf("Failed to add entry 2: %v", err)
	}
	if err := addEntry(entry3); err != nil {
		t.Fatalf("Failed to add entry 3: %v", err)
	}

	// Update entry 2
	updatedEntry := Entry{
		Id:       2,
		Title:    "Updated Title 2",
		Username: "Updated User 2",
		Password: "Updated Password 2",
		Address:  "http://updated-address2.com",
		Notes:    "Updated notes 2",
	}

	if err := UpdateEntry(2, updatedEntry); err != nil {
		t.Fatalf("Failed to update entry 2: %v", err)
	}

	// Verify the entry was updated
	entries, err := readAllEntries()
	if err != nil {
		t.Fatalf("Failed to read all entries: %v", err)
	}

	// Check that the number of entries is correct
	if len(entries) != 3 {
		t.Errorf("Expected 3 entries, but got %d", len(entries))
	}

	// Check that entry 2 was updated correctly
	for _, entry := range entries {
		if entry.Id == 2 {
			if entry.Title != updatedEntry.Title ||
				entry.Username != updatedEntry.Username ||
				entry.Password != updatedEntry.Password ||
				entry.Address != updatedEntry.Address ||
				entry.Notes != updatedEntry.Notes {
				t.Errorf("Entry 2 was not updated correctly: %+v", entry)
			}
			return
		}
	}

	t.Errorf("Entry with ID 2 was not found after update")
}

func TestUpdateEntryInMemory(t *testing.T) {
	// Clean up test files after the test
	defer os.Remove("data.bin")

	// Helper function to add an entry to the file
	addEntry := func(entry Entry) error {
		return SaveEntry(entry)
	}

	// Add some entries
	entry1 := Entry{
		Id:       1,
		Title:    "Old Title 1",
		Username: "Old User 1",
		Password: "Old Password 1",
		Address:  "http://old-address1.com",
		Notes:    "Old notes 1",
	}
	entry2 := Entry{
		Id:       2,
		Title:    "Old Title 2",
		Username: "Old User 2",
		Password: "Old Password 2",
		Address:  "http://old-address2.com",
		Notes:    "Old notes 2",
	}
	entry3 := Entry{
		Id:       3,
		Title:    "Old Title 3",
		Username: "Old User 3",
		Password: "Old Password 3",
		Address:  "http://old-address3.com",
		Notes:    "Old notes 3",
	}

	if err := addEntry(entry1); err != nil {
		t.Fatalf("Failed to add entry 1: %v", err)
	}
	if err := addEntry(entry2); err != nil {
		t.Fatalf("Failed to add entry 2: %v", err)
	}
	if err := addEntry(entry3); err != nil {
		t.Fatalf("Failed to add entry 3: %v", err)
	}

	// Update entry 2
	updatedEntry := Entry{
		Id:       2,
		Title:    "Updated Title 2",
		Username: "Updated User 2",
		Password: "Updated Password 2",
		Address:  "http://updated-address2.com",
		Notes:    "Updated notes 2",
	}

	if err := UpdateEntryInMemory(2, updatedEntry); err != nil {
		t.Fatalf("Failed to update entry 2: %v", err)
	}

	// Verify the entry was updated
	entries, err := readAllEntries()
	if err != nil {
		t.Fatalf("Failed to read all entries: %v", err)
	}

	// Check that the number of entries is correct
	if len(entries) != 3 {
		t.Errorf("Expected 3 entries, but got %d", len(entries))
	}

	// Check that entry 2 was updated correctly
	for _, entry := range entries {
		if entry.Id == 2 {
			if entry.Title != updatedEntry.Title ||
				entry.Username != updatedEntry.Username ||
				entry.Password != updatedEntry.Password ||
				entry.Address != updatedEntry.Address ||
				entry.Notes != updatedEntry.Notes {
				t.Errorf("Entry 2 was not updated correctly: %+v", entry)
			}
			return
		}
	}

	t.Errorf("Entry with ID 2 was not found after update")
}

// Helper function to read all entries from the file
func readAllEntries() ([]Entry, error) {
	file, err := os.Open(dataFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []Entry

	for {
		var entry Entry

		err = binary.Read(file, binary.LittleEndian, &entry.Id)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		entry.Title, err = readString(file)
		if err != nil {
			return nil, err
		}
		entry.Username, err = readString(file)
		if err != nil {
			return nil, err
		}
		entry.Password, err = readString(file)
		if err != nil {
			return nil, err
		}
		entry.Address, err = readString(file)
		if err != nil {
			return nil, err
		}
		entry.Notes, err = readString(file)
		if err != nil {
			return nil, err
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
