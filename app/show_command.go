package app

import (
	"squirrel/data"
	"squirrel/types"
	"strconv"
)

var ()

func ShowCommand(p types.Printer, d types.Decryptor) Command {
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
			p("{gray}Show an entry by ID{/gray}\n")
			id = readId(p)
		}

		ent, err := data.LoadEntry(id)
		if err != nil {
			p("{red}Loading entry with ID {0} failed! {1}{/red}\n", id, err)
			return
		}

		decrypt(&ent, d)

		display(ent, p)
	}
}

func decrypt(ent *data.Entry, d types.Decryptor) error {
	var error error

	ent.Username, error = d(ent.Username)
	if error != nil {
		return error
	}

	ent.Address, error = d(ent.Address)
	if error != nil {
		return error
	}

	ent.Password, error = d(ent.Password)
	if error != nil {
		return error
	}

	ent.Notes, error = d(ent.Notes)
	if error != nil {
		return error
	}

	return nil
}
