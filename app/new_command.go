package app

import (
	"fmt"
	"squirrel/data"
	"squirrel/types"
)

func NewCommand(p types.Printer, e types.Encryptor) Command {
	return func(args ...string) {
		var ne data.Entry
		p("{gray}New entry (all fields except title will be encrypted){/gray}\n")

		var pass, veryfy string

		ReadInput("Title", "", true, p, &ne.Title)

		for {
			ReadSecret("Password", "optional", false, p, &pass)

			if pass == "" {
				break
			}

			ReadSecret("Verify password", "optional", false, p, &veryfy)

			if pass != veryfy {
				p("{brightWhite}Password did not match! Try again.{/brightWhite}\n")
			} else {
				break
			}
		}

		ne.Password = pass
		ReadInput("Username", "optional", false, p, &ne.Username)
		ReadInput("Address", "optional", false, p, &ne.Address)
		ReadInput("Notes", "optional", false, p, &ne.Notes)

		err := encryptEntry(&ne, e, p)
		if err != nil {
			p("{red}Encrypting the new entry failed!{/red}\n", err)
			return
		}

		id, err := data.GetLargestId()
		if err != nil {
			p("{red}Getting last ID failed!{/red}\n", err)
		}

		ne.Id = id + 1

		err = data.SaveEntry(ne)
		if err != nil {
			p("{red}Saving the new entry failed!{/red}\n", err)
		} else {
			p("{green}Entry '{0}' was saved successfully. ID: {1}{/green}\n", ne.Title, ne.Id)
		}

	}
}

func descriptionOfField(desc string) string {
	if len(desc) > 0 {
		return fmt.Sprintf("(%v)", desc)
	} else {
		return ""
	}
}

func encryptEntry(ent *data.Entry, encrypt types.Encryptor, print types.Printer) error {
	var err error

	ent.Username, err = encrypt(ent.Username)
	if err != nil {
		print("{red}Error in encrypting username{/red} {0}", err)
		return err
	}

	ent.Password, err = encrypt(ent.Password)
	if err != nil {
		print("{red}Error in encrypting password{/red} {0}", err)
		return err
	}

	ent.Address, err = encrypt(ent.Address)
	if err != nil {
		print("{red}Error in encrypting Address{/red} {0}", err)
		return err
	}

	ent.Notes, err = encrypt(ent.Notes)
	if err != nil {
		print("{red}Error in encrypting notes{/red} {0}", err)
		return err
	}

	return nil
}
