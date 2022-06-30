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
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

const (
	userNameFlag       = "username"
	promptPasswordFlag = "prompt-password"
	passwordFlag       = "password"
	hostFlag           = "host"
	schemeFlag         = "scheme"
	envPrefix          = "airflow_vars"
)

var (
	cfgFile  string
	username string
	password string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "airflow-vars",
	Short: "gitops tool for airflow variable deplyment",
	Long: `
airflow vars is a cli intends to help you throughout your airflow deployment process.
It'll help you manage your airflow variables with yaml files and deploy all of your variables directly to airflow.
	`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		promptPassword, err := cmd.Flags().GetBool(promptPasswordFlag)
		cobra.CheckErr(err)

		if promptPassword {
			fmt.Print("Please insert your password: \n")
			bytesPassword, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				cobra.CheckErr(fmt.Errorf("error getting password from prompt: %s", err))
			}

			password = string(bytesPassword)
		}

		if password == "" {
			cobra.CheckErr(fmt.Errorf("password is not set, please use --%s or --%s", passwordFlag, promptPasswordFlag))
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	rootCmd.Version = version
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.airflow-vars.yaml)")
	rootCmd.PersistentFlags().String(hostFlag, "localhost:8080", "your airflow host")
	rootCmd.PersistentFlags().String(schemeFlag, "https", "your airflow scheme")

	rootCmd.PersistentFlags().StringVar(&username, userNameFlag, "", "Username to authenticate airflow")
	cobra.CheckErr(rootCmd.MarkPersistentFlagRequired(userNameFlag))

	rootCmd.PersistentFlags().StringVar(&password, passwordFlag, "", "Password to authenticate airflow")
	rootCmd.PersistentFlags().Bool(promptPasswordFlag, false, "Interactive prompt for authentication password")
	rootCmd.MarkFlagsMutuallyExclusive(passwordFlag, promptPasswordFlag)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".airflow-vars" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".airflow-vars")
	}
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()

	bindFlags()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func bindFlags() {
	rootCmd.Flags().VisitAll(func(f *pflag.Flag) {
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			err := viper.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
			cobra.CheckErr(err)
		}

		if !f.Changed && viper.IsSet(f.Name) {
			val := viper.Get(f.Name)
			err := rootCmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
			cobra.CheckErr(err)
		}
	})
}
