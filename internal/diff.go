/*
Copyright © 2022 Omer Zamir <zamir98@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package internal

import (
	"encoding/json"
	"fmt"

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
		}

		if variable.New != nil {
			err := json.Unmarshal([]byte(*variable.New.Value), &newVar)
			cobra.CheckErr(err)
		}

		diff := pretty.Compare(prev, newVar)
		if diff != "" {
			fmt.Printf("Variable: %s \n", key)
			fmt.Println(diff)
			fmt.Print("\n")
		}

		if variable.New != nil && variable.Prev != nil {
			if diff != "" {
				change += 1
			}
		} else if variable.New == nil && variable.Prev != nil {
			toDelete += 1
		} else {
			add += 1
		}
	}

	fmt.Printf("Plan: Add %d, Update %d, Delete %d \n", add, change, toDelete)

	return add > 0 || change > 0 || toDelete > 0
}
