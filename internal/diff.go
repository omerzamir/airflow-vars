package internal

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kylelemons/godebug/pretty"
	"github.com/spf13/cobra"
)

func PrintDiff(variables map[string]*VersionedVariable) bool {
	add, change, toDelete := 0, 0, 0

	for key, variable := range variables {
		var prev any = nil
		var newVar any = nil

		if variable.Prev != nil {
			err := json.Unmarshal([]byte(*variable.Prev.Value), &prev)
			cobra.CheckErr(err)

			if variable.New != nil {
				if strings.Compare(*variable.Prev.Value, *variable.New.Value) != 0 {
					change += 1
				}
			} else {
				toDelete += 1
			}
		}

		if variable.New != nil {
			err := json.Unmarshal([]byte(*variable.New.Value), &newVar)
			cobra.CheckErr(err)

			if variable.Prev == nil {
				add += 1
			}
		}

		diff := pretty.Compare(prev, newVar)
		if diff != "" {
			fmt.Printf("Variable: %s \n", key)
			fmt.Println(diff)
			fmt.Print("\n")
		}
	}

	fmt.Printf("Plan: Add %d, Update %d, Delete %d \n", add, change, toDelete)

	return add > 0 || change > 0 || toDelete > 0
}
