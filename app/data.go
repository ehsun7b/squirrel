package app

import (
	"os"
	d "squirrel/data"
	l "squirrel/log"
)

func ReadState() d.State {
	if d.HasDataFile() {
		count, err := d.CountEntries()
		if err != nil {
			l.Println("{red}Error in reading data file!{/red}{0}", err)
			l.Println("Run using --fix-data-file")
			os.Exit(3)
		}

		lastId, err := d.GetLargestId()
		if err != nil {
			l.Println("{red}Error in reading data file!{/red}{0}", err)
			l.Println("Run using --fix-data-file")
			os.Exit(3)
		}

		return d.State{
			Count:  count,
			LastId: lastId,
		}
	}

	return d.State{}
}
