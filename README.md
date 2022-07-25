# airflow vars

airflow vars is a cli intends to help you throughout your airflow deployment process.

It'll help you manage your airflow variables with yaml files and deploy all of your variables directly to airflow.

## Installation

There are 2 options in order to install the cli.

 1. Install with go.
 2. Install from binary release.

### Install with go
Install the airflow vars with the command `go install github.com/omerzamir/airflow-vars@latest`. 
Go will automatically install it in your `$GOPATH/bin` directory which should be in your $PATH. 

### Install from binary release
1. Download your [desired version](https://github.com/omerzamir/airflow-vars/releases)
2. Unpack it (tar -zxvf airflow-vars_0.0.1_Linux_arm64.tar.gz)
3. Find the `airflow-vars` binary in the unpacked directory, and move it to its desired destination (`mv airflow-vars_0.0.1_Linux_arm64/airflow-vars /usr/local/bin/airflow-vars`)

## Usage
At any time, you can view usage instructions by entering `airflow-vars --help`.

### airflow-vars sync
sync will read the given file/directory and sync your airflow cluster with the given state.

The given input path should represent the state of all of the prefixes defined in the files.

create a file `example.yaml`:
```
config:
  prefix: example
vars:
  a:
    test1: "test value 1"
    test_arr: ["test1", "test2"]

```

This file defines only 1 variable, with the prefix `example`.
Running the sync command the `airflow-vars` will manage all of the airflow variables starts with `example` and will reflect the state from the yaml onto airflow.
e.g. let's assume the `example_a` variable does not exist in our airflow cluster, `airflow-vars` will create a json variable named `example_a` with json value of: 
```
{
  "test1": "test value 1",
  "test_arr": ["test1", "test2"]
}
```
If any variables with prefix of example exists in the cluster, `airflow-vars` will delete them.


## License

Airflow-Vars is released under the Apache 2.0 license. See [LICENSE](https://github.com/omerzamir/airflow-vars/blob/master/LICENSE)
