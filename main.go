/*
Copyright © 2022 Omer Zamir zamir98@gmail.com

*/
package main

import "github.com/omerzamir/airflow-vars-sync/cmd"

var version = "development"

func main() {
	cmd.Execute(version)
}
