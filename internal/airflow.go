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
	"context"
	"math"
	"strings"

	"github.com/apache/airflow-client-go/airflow"
	"github.com/spf13/cobra"
)

const (
	fetchChunkSize      = 1000
	defaultVariableSize = 100
)

func GetAllVariables(ctx context.Context, cli *airflow.APIClient, prefixes []string) []*airflow.Variable {
	data := make([]*airflow.Variable, 0, defaultVariableSize)
	scanned, offset, total := int32(0), int32(0), int32(1)

	for scanned < total {
		vars, _, err := cli.VariableApi.GetVariables(ctx).Limit(fetchChunkSize).Offset(offset).Execute()
		cobra.CheckErr(err)

		total = vars.GetTotalEntries()
		toScan := vars.GetVariables()

		for _, v := range toScan {
			for _, prefix := range prefixes {
				if strings.HasPrefix(*v.Key, prefix) {
					variable, _, err := cli.VariableApi.GetVariable(ctx, *v.Key).Execute()
					cobra.CheckErr(err)

					data = append(data, &variable)
					break
				}
			}
		}
		scanned += int32(len(toScan))

		if int64(offset)+fetchChunkSize >= math.MaxInt32 {
			return data
		}

		offset += int32(len(*vars.Variables))
	}

	return data
}

func ApplyChanges(ctx context.Context, cli *airflow.APIClient, variables map[string]*VersionedVariable) {
	for _, variable := range variables {
		if variable.Prev != nil && variable.New != nil {
			if strings.Compare(*variable.Prev.Value, *variable.New.Value) != 0 {
				_, _, err := cli.VariableApi.PatchVariable(ctx, *variable.New.Key).Variable(*variable.New).Execute()
				cobra.CheckErr(err)
			}
		} else if variable.New != nil {
			_, _, err := cli.VariableApi.PostVariables(ctx).Variable(*variable.New).Execute()
			cobra.CheckErr(err)
		} else {
			_, err := cli.VariableApi.DeleteVariable(ctx, *variable.Prev.Key).Execute()
			cobra.CheckErr(err)
		}
	}
}
