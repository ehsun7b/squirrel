package app

import (
	"squirrel/data"
	"squirrel/types"
	"strconv"
)

var ()

func EditCommand(p types.Printer, e types.Encryptor, d types.Decryptor) Command {
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
			p("{gray}Edit an entry by ID{/gray}\n")
			id = readId(p)
		}

		ent, err := data.LoadEntry(id)
		if err != nil {
			p("{red}Loading entry with ID {0} failed! {1}{/red}\n", id, err)
			return
		}

		decrypt(&ent, d)

		display(ent, p)

		var newTitle, newPassword, newAddress, newUsername, newNotes string

		if GetYesNoInput(p, "Update title") {
			ReadInput("New title", "", true, p, &newTitle)
		} else {
			newTitle = ent.Title
		}

		if GetYesNoInput(p, "Update password") {
			var confirmNewPassword string

			for {
				ReadSecret("New password", "", true, p, &newPassword)
				ReadSecret("Verify password", "", true, p, &confirmNewPassword)

				if newPassword != confirmNewPassword {
					p("{red}Password does not match!{/red}\n")
				} else {
					break
				}
			}

		} else {
			newPassword = ent.Password
		}

		if GetYesNoInput(p, "Update username") {
			ReadInput("New title", "", false, p, &newUsername)
		} else {
			newUsername = ent.Username
		}

		if GetYesNoInput(p, "Update address") {
			ReadInput("New address", "", false, p, &newAddress)
		} else {
			newAddress = ent.Address
		}

		if GetYesNoInput(p, "Update notes") {
			ReadInput("New notes", "", false, p, &newNotes)
		} else {
			newNotes = ent.Notes
		}

		ent.Title = newTitle
		ent.Username = newUsername
		ent.Address = newAddress
		ent.Notes = newNotes
		ent.Password = newPassword

		p("{magenta}Will update to:{/magenta}\n")
		display(ent, p)

		if GetYesNoInput(p, "{magenta}Correct{/magenta}") {
			err := encryptEntry(&ent, e, p)
			if err != nil {
				p("{red}Encrypting entity failed!{/red} {0}", err)
				return
			}

			err = data.UpdateEntry(ent.Id, ent)
			if err != nil {
				p("{red}Updating entity failed!{/red} {0}", err)
				return
			}

			p("{green}Updated.{/green}\n")
		} else {
			p("Update canceled!\n")
		}
	}
}
