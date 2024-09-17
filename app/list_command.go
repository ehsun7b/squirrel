package app

import (
	"errors"
	"fmt"
	"squirrel/data"
	"squirrel/types"
	"strconv"
)

var (
	ErrWrongArgsForCommandList = errors.New("incorrect number of arguments")
	DefaultOrder               = data.ByTitle
	DefaultLimit               = 10
)

func ListCommand(p types.Printer, d types.Decryptor) Command {
	return func(args ...string) {
		order, limit, err := determineOrderAndLimit(args...)
		if err != nil {
			p("{gray}{0}{/gray}\n", err.Error())
			p("{red}Wrong arguments{/red}\nlist command examples:{/brightWhite}\n\tlist\n\tlist title 20\n\tlist 20\n\tlist username\n\tlist username 30{/brightWhite}\n")
			return
		}

		count, _ := data.CountEntries()
		p("There are {0} entries.\n", count)

		if count > 0 {
			entries, err := data.Entries(order, limit, d)
			if err != nil {
				p("{red}Error in loading entries!{/red}: {0}\n", err)
			}

			for i, entry := range entries {
				p("{0}. {1} \tID: {2} \tUsername: {3}\n", i+1, entry.Title, entry.Id, entry.Username)
			}
		}
	}
}

func determineOrderAndLimit(args ...string) (data.Order, int, error) {
	l := len(args)

	switch l {
	case 0:
		return DefaultOrder, DefaultLimit, nil

	case 2:
		order, err := data.OrderFromString(args[0])
		if err != nil {
			return DefaultOrder, DefaultLimit, fmt.Errorf("%v; %w", ErrWrongArgsForCommandList, err)
		}

		limit, err := strconv.Atoi(args[1])
		if err != nil {
			return DefaultOrder, DefaultLimit, fmt.Errorf("%v; %w", ErrWrongArgsForCommandList, err)
		}

		return order, limit, nil

	case 1:
		order, err1 := data.OrderFromString(args[0])
		limit, err2 := strconv.Atoi(args[0])

		if err1 != nil && err2 != nil {
			return DefaultOrder, DefaultLimit, fmt.Errorf("%v; %w", ErrWrongArgsForCommandList, fmt.Errorf("%v - %v", err1, err2))
		} else if err1 != nil {
			return DefaultOrder, limit, nil
		} else {
			return order, DefaultLimit, nil
		}
	default:
		return DefaultOrder, DefaultLimit, ErrWrongArgsForCommandList
	}
}
