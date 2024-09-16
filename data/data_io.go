package data

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
)

var ErrEntryExists = errors.New("entry with this ID already exists")
var ErrEntryNotFound = errors.New("entry not found")

const dataFile = "data.bin"
const stateFile = "state.bin"
const passwordVerifyFile = "enc.bin"

func SaveEntry(entry Entry) error {
	// Check if an entry with the same ID already exists
	_, err := LoadEntry(entry.Id)
	if err == nil {
		// If no error, it means the entry exists, so we return an error
		return ErrEntryExists
	}

	file, err := os.OpenFile(dataFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	err = binary.Write(file, binary.LittleEndian, entry.Id)
	if err != nil {
		return err
	}

	if err := writeString(file, entry.Title); err != nil {
		return err
	}
	if err := writeString(file, entry.Username); err != nil {
		return err
	}
	if err := writeString(file, entry.Password); err != nil {
		return err
	}
	if err := writeString(file, entry.Address); err != nil {
		return err
	}
	if err := writeString(file, entry.Notes); err != nil {
		return err
	}

	return nil
}

func SavePassVerify(encrypted string) error {

	file, err := os.OpenFile(passwordVerifyFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := writeString(file, encrypted); err != nil {
		return err
	}

	return nil
}

func LoadPassVerify() (string, error) {
	file, err := os.Open(passwordVerifyFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var result = ""
	result, err = readString(file)

	if err != nil {
		return "", err
	}

	return result, nil
}

func DeleteEntry(entryID int64) error {
	sourceFile, err := os.Open(dataFile)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	tempFile, err := os.OpenFile("temp_"+dataFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	entryFound := false

	for {
		var entry Entry

		err = binary.Read(sourceFile, binary.LittleEndian, &entry.Id)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		entry.Title, err = readString(sourceFile)
		if err != nil {
			return err
		}
		entry.Username, err = readString(sourceFile)
		if err != nil {
			return err
		}
		entry.Password, err = readString(sourceFile)
		if err != nil {
			return err
		}
		entry.Address, err = readString(sourceFile)
		if err != nil {
			return err
		}
		entry.Notes, err = readString(sourceFile)
		if err != nil {
			return err
		}

		if entry.Id == entryID {
			entryFound = true
			continue
		}

		err = binary.Write(tempFile, binary.LittleEndian, entry.Id)
		if err != nil {
			return err
		}
		if err := writeString(tempFile, entry.Title); err != nil {
			return err
		}
		if err := writeString(tempFile, entry.Username); err != nil {
			return err
		}
		if err := writeString(tempFile, entry.Password); err != nil {
			return err
		}
		if err := writeString(tempFile, entry.Address); err != nil {
			return err
		}
		if err := writeString(tempFile, entry.Notes); err != nil {
			return err
		}
	}

	if !entryFound {
		return ErrEntryNotFound
	}

	err = os.Remove(dataFile)
	if err != nil {
		return err
	}

	err = os.Rename("temp_"+dataFile, dataFile)
	if err != nil {
		return err
	}

	return nil
}

// In-memory version of DeleteEntry
func DeleteEntryInMemory(id int64) error {
	// Read all entries into memory
	file, err := os.OpenFile(dataFile, os.O_RDWR, 0644)
	if err != nil {
		return err
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
			return err
		}

		entry.Title, err = readString(file)
		if err != nil {
			return err
		}
		entry.Username, err = readString(file)
		if err != nil {
			return err
		}
		entry.Password, err = readString(file)
		if err != nil {
			return err
		}
		entry.Address, err = readString(file)
		if err != nil {
			return err
		}
		entry.Notes, err = readString(file)
		if err != nil {
			return err
		}

		entries = append(entries, entry)
	}

	// Filter out the entry to be deleted
	var newEntries []Entry
	for _, entry := range entries {
		if entry.Id != id {
			newEntries = append(newEntries, entry)
		}
	}

	// Re-open the file to overwrite it
	file, err = os.OpenFile(dataFile, os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, entry := range newEntries {
		err = binary.Write(file, binary.LittleEndian, entry.Id)
		if err != nil {
			return err
		}
		if err := writeString(file, entry.Title); err != nil {
			return err
		}
		if err := writeString(file, entry.Username); err != nil {
			return err
		}
		if err := writeString(file, entry.Password); err != nil {
			return err
		}
		if err := writeString(file, entry.Address); err != nil {
			return err
		}
		if err := writeString(file, entry.Notes); err != nil {
			return err
		}
	}

	return nil
}

func UpdateEntry(entryId int64, updatedEntry Entry) error {
	sourceFile, err := os.Open(dataFile)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	tempFile, err := os.OpenFile("temp_"+dataFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	entryFound := false

	for {
		var entry Entry

		err = binary.Read(sourceFile, binary.LittleEndian, &entry.Id)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		entry.Title, err = readString(sourceFile)
		if err != nil {
			return err
		}
		entry.Username, err = readString(sourceFile)
		if err != nil {
			return err
		}
		entry.Password, err = readString(sourceFile)
		if err != nil {
			return err
		}
		entry.Address, err = readString(sourceFile)
		if err != nil {
			return err
		}
		entry.Notes, err = readString(sourceFile)
		if err != nil {
			return err
		}

		if entry.Id == entryId {
			entry = updatedEntry
			entry.Id = entryId
			entryFound = true
		}

		err = binary.Write(tempFile, binary.LittleEndian, entry.Id)
		if err != nil {
			return err
		}
		if err := writeString(tempFile, entry.Title); err != nil {
			return err
		}
		if err := writeString(tempFile, entry.Username); err != nil {
			return err
		}
		if err := writeString(tempFile, entry.Password); err != nil {
			return err
		}
		if err := writeString(tempFile, entry.Address); err != nil {
			return err
		}
		if err := writeString(tempFile, entry.Notes); err != nil {
			return err
		}
	}

	if !entryFound {
		return ErrEntryNotFound
	}

	err = os.Remove(dataFile)
	if err != nil {
		return err
	}

	err = os.Rename("temp_"+dataFile, dataFile)
	if err != nil {
		return err
	}

	return nil
}

// In-memory version of UpdateEntry
func UpdateEntryInMemory(id int64, updatedEntry Entry) error {
	// Read all entries into memory
	file, err := os.OpenFile(dataFile, os.O_RDWR, 0644)
	if err != nil {
		return err
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
			return err
		}

		entry.Title, err = readString(file)
		if err != nil {
			return err
		}
		entry.Username, err = readString(file)
		if err != nil {
			return err
		}
		entry.Password, err = readString(file)
		if err != nil {
			return err
		}
		entry.Address, err = readString(file)
		if err != nil {
			return err
		}
		entry.Notes, err = readString(file)
		if err != nil {
			return err
		}

		entries = append(entries, entry)
	}

	// Update the entry in memory
	for i, entry := range entries {
		if entry.Id == id {
			entries[i] = updatedEntry
			break
		}
	}

	// Re-open the file to overwrite it
	file, err = os.OpenFile(dataFile, os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, entry := range entries {
		err = binary.Write(file, binary.LittleEndian, entry.Id)
		if err != nil {
			return err
		}
		if err := writeString(file, entry.Title); err != nil {
			return err
		}
		if err := writeString(file, entry.Username); err != nil {
			return err
		}
		if err := writeString(file, entry.Password); err != nil {
			return err
		}
		if err := writeString(file, entry.Address); err != nil {
			return err
		}
		if err := writeString(file, entry.Notes); err != nil {
			return err
		}
	}

	return nil
}

func LoadEntry(id int64) (Entry, error) {
	file, err := os.Open(dataFile)
	if err != nil {
		return Entry{}, err
	}
	defer file.Close()

	var entry Entry

	for {
		err = binary.Read(file, binary.LittleEndian, &entry.Id)
		if err == io.EOF {
			return Entry{}, ErrEntryNotFound
		} else if err != nil {
			return Entry{}, err
		}

		if entry.Title, err = readString(file); err != nil {
			return Entry{}, err
		}
		if entry.Username, err = readString(file); err != nil {
			return Entry{}, err
		}
		if entry.Password, err = readString(file); err != nil {
			return Entry{}, err
		}
		if entry.Address, err = readString(file); err != nil {
			return Entry{}, err
		}
		if entry.Notes, err = readString(file); err != nil {
			return Entry{}, err
		}

		if entry.Id == id {
			return entry, nil
		}
	}
}

func GetLargestId() (int64, error) {
	if !HasDataFile() {
		return 0, nil
	}

	file, err := os.Open(dataFile)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var largestId int64
	for {
		var entry Entry

		err := binary.Read(file, binary.LittleEndian, &entry.Id)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				break
			} else if err == io.EOF {
				break
			}
		}

		if _, err := readString(file); err != nil {
			return 0, err
		}
		if _, err := readString(file); err != nil {
			return 0, err
		}
		if _, err := readString(file); err != nil {
			return 0, err
		}
		if _, err := readString(file); err != nil {
			return 0, err
		}
		if _, err := readString(file); err != nil {
			return 0, err
		}

		if entry.Id > largestId {
			largestId = entry.Id
		}
	}

	return largestId, nil
}

func CountEntries() (int64, error) {
	file, err := os.Open(dataFile)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var count int64 = 0
	for {
		var entry Entry

		// Attempt to read the ID
		err = binary.Read(file, binary.LittleEndian, &entry.Id)
		if err != nil {
			if err == io.EOF {
				// We reached the end of the file, break out of the loop
				break
			}

			return 0, err
		}

		// Read the strings (Title, Username, Password, Address, Notes)
		if _, err := readString(file); err != nil {
			return 0, err
		}
		if _, err := readString(file); err != nil {
			return 0, err
		}
		if _, err := readString(file); err != nil {
			return 0, err
		}
		if _, err := readString(file); err != nil {
			return 0, err
		}
		if _, err := readString(file); err != nil {
			return 0, err
		}

		// Increment the count for each successfully read entry
		count++
	}

	return count, nil
}

func Entries(o Order, limit int) ([]Entry, error) {
	// Retrieve all entries
	allEntries, err := entries()
	if err != nil {
		return nil, err
	}

	// Sort entries based on the Order
	switch o {
	case ByTitle:
		sort.SliceStable(allEntries, func(i, j int) bool {
			return allEntries[i].Title < allEntries[j].Title
		})
	case ByUsername:
		sort.SliceStable(allEntries, func(i, j int) bool {
			return allEntries[i].Username < allEntries[j].Username
		})
	default:
		return nil, fmt.Errorf("unknown order: %v", o)
	}

	// Apply the limit to the number of entries returned
	if limit > len(allEntries) {
		limit = len(allEntries)
	}

	return allEntries[:limit], nil
}

// entries reads and returns all entries from a file
func entries() ([]Entry, error) {
	file, err := os.Open(dataFile) // Replace dataFile with your file path
	if err != nil {
		return nil, err
	}
	defer file.Close()

	count, err := CountEntries()
	if err != nil {
		return nil, err
	}

	var entries = make([]Entry, 0, count)

	for {
		var entry Entry

		// Read the ID
		if err := binary.Read(file, binary.LittleEndian, &entry.Id); err != nil {
			if err == io.EOF {
				// End of file reached
				break
			}
			return nil, err
		}

		// Read the strings for Title, Username, Password, Address, and Notes
		if entry.Title, err = readString(file); err != nil {
			return nil, err
		}
		if entry.Username, err = readString(file); err != nil {
			return nil, err
		}
		if entry.Password, err = readString(file); err != nil {
			return nil, err
		}
		if entry.Address, err = readString(file); err != nil {
			return nil, err
		}
		if entry.Notes, err = readString(file); err != nil {
			return nil, err
		}

		// Append the successfully read entry to the list
		entries = append(entries, entry)
	}

	return entries, nil
}

func SaveState(state State) error {
	file, err := os.OpenFile(stateFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	err = binary.Write(file, binary.LittleEndian, state)
	if err != nil {
		return err
	}

	return nil
}

func LoadState() (State, error) {
	file, err := os.Open(stateFile)
	if err != nil {
		return State{}, err
	}
	defer file.Close()

	var state State

	err = binary.Read(file, binary.LittleEndian, &state)
	if err != nil {
		return State{}, err
	}

	return state, nil
}

func HasStateFile() bool {
	return fileExists(stateFile)
}

func HasDataFile() bool {
	return fileExists(dataFile)
}

func HasPassVerifyFile() bool {
	return fileExists(passwordVerifyFile)
}

func writeString(file *os.File, str string) error {
	strBytes := []byte(str)
	// Write the length of the string (as a varint)
	length := uint64(len(strBytes))
	if err := binary.Write(file, binary.LittleEndian, length); err != nil {
		return err
	}
	_, err := file.Write(strBytes)
	return err
}

func readString(file *os.File) (string, error) {
	var length uint64
	if err := binary.Read(file, binary.LittleEndian, &length); err != nil {
		return "", err
	}
	strBytes := make([]byte, length)
	if _, err := io.ReadFull(file, strBytes); err != nil {
		return "", err
	}
	return string(strBytes), nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
