package app

import (
	"fmt"
	"squirrel/data"
	"squirrel/types"
	"strconv"
)

var ()

func DeleteCommand(p types.Printer) Command {
	return func(args ...string) {
		var id int64
		if len(args) > 0 {
			passedId, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				p("{red}Bad ID! {0}{/red}\n", err)
				return
			}

			id = passedId
		} else {
			p("{gray}Delete an entry by ID{/gray}\n")
			id = readId(p)
		}

		ent, deleted, err := delete(id, p)
		if err != nil {
			p("{red}Loading or deleting entry with ID {0} failed! {1}{/red}\n", id, err)
			return
		}

		if deleted {
			p("{green}Entry '{0}' deleted. {/red}\n", ent.Title)
		}
	}
}

func delete(id int64, p types.Printer) (data.Entry, bool, error) {
	ent, err := data.LoadEntry(id)
	if err != nil {
		return data.Entry{}, false, err
	}

	if GetYesNoInput(p, fmt.Sprintf("Delete entry '%v'", ent.Title)) {
		err := data.DeleteEntryInMemory(id)
		if err != nil {
			return data.Entry{}, false, err
		}

		// deleted
		return ent, true, nil
	}

	return ent, false, nil
}

func readId(p types.Printer) int64 {
	var id int64
	ReadInput("ID", "", true, p, &id)
	return id
}
