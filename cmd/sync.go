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
package cmd

import (
	"context"
	"fmt"

	"github.com/omerzamir/airflow-vars/internal"

	"github.com/apache/airflow-client-go/airflow"
	"github.com/spf13/cobra"
)

const (
	defaultInputFilesSize = 50
)

var (
	yesFlag    = "yes"
	dryFlag    = "dry"
	unsafeFlag = "unsafe"
)

type Exists struct{}

// importCmd represents the import command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "sync your variables files with your airflow cluster",
	Long:  `sync will read the given file/directory and sync your airflow cluster with the given state.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cobra.CheckErr(fmt.Errorf("missing input file/directory, enter \".\" to run in the current directory"))
		}

		ctx, cli := initAirflowCli(cmd)

		files := make([]*internal.InputFile, 0, defaultInputFilesSize)

		for _, path := range args {
			f, err := internal.ReadInputFiles(path)
			cobra.CheckErr(err)

			files = append(files, f...)
		}

		p := make(map[string]any)

		unsafeMode, err := cmd.Flags().GetBool(unsafeFlag)
		cobra.CheckErr(err)

		for _, file := range files {
			if file == nil {
				continue
			}

			if file.Config.Prefix == "" {
				if unsafeMode {
					p = map[string]any{"": &Exists{}}

					break
				} else {
					cobra.CheckErr(fmt.Errorf(`a config file without prefix was provided.
this options is disabled by default for safety reasons.
if you want to enable empty/undefined prefixes, use the --unsafe flag`))
				}

			}

			p[file.Config.Prefix] = &Exists{}
		}

		prefixes := make([]string, 0, len(p))
		for k := range p {
			prefixes = append(prefixes, k)
		}

		relevantVariables := internal.GetAllVariables(ctx, cli, prefixes)
		zippedVariables := internal.ZipVariablesByKey(files, relevantVariables)

		hasChange := internal.PrintDiff(zippedVariables)

		dryRun, err := cmd.Flags().GetBool(dryFlag)
		cobra.CheckErr(err)

		if dryRun {
			return
		}

		if !hasChange {
			fmt.Print("No changes. Exiting... \n")
			return
		}

		approved, err := cmd.Flags().GetBool(yesFlag)
		cobra.CheckErr(err)

		if !approved && !internal.YesNoPrompt("Approve plan?", false) {
			fmt.Print("Plan not approved. Exiting... \n")
			return
		}
		fmt.Print("Approved, executing plan... \n")
		internal.ApplyChanges(ctx, cli, zippedVariables)
		fmt.Print("Done. \n")
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().BoolP(yesFlag, "y", false, "proceed without confirm changes")
	syncCmd.Flags().Bool(dryFlag, false, "dry run mode")
	syncCmd.Flags().Bool(unsafeFlag, false, "unsafe mode, configs without prefix can be provided")
}

func initAirflowCli(cmd *cobra.Command) (context.Context, *airflow.APIClient) {
	conf := airflow.NewConfiguration()

	host, err := cmd.Flags().GetString(hostFlag)
	cobra.CheckErr(err)
	conf.Host = host

	scheme, err := cmd.Flags().GetString(schemeFlag)
	cobra.CheckErr(err)
	conf.Scheme = scheme

	cli := airflow.NewAPIClient(conf)
	cred := airflow.BasicAuth{
		UserName: username,
		Password: password,
	}
	ctx := context.WithValue(context.Background(), airflow.ContextBasicAuth, cred)

	return ctx, cli
}
