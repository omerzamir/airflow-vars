package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/apache/airflow-client-go/airflow"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	JYaml "sigs.k8s.io/yaml"
)

const (
	defaultInputFilesSize = 10
)

type InputFile struct {
	Config struct {
		Prefix string
	}
	Vars map[string]any
}

type VersionedVariable struct {
	Prev *airflow.Variable
	New  *airflow.Variable
}

func ReadInputFiles(inPath string) ([]*InputFile, error) {
	if inPath == "" {
		inPath = "."
	}

	data := make([]*InputFile, 0, defaultInputFilesSize)

	err := filepath.Walk(inPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				cobra.CheckErr(err)
			}

			if info.IsDir() {
				return nil
			}

			if strings.HasSuffix(info.Name(), "yml") || strings.HasSuffix(info.Name(), "yaml") {
				f, err := readFile(path)
				if err != nil {
					cobra.CheckErr(err)
					return err
				}
				vars := make(map[string]any)
				for k, v := range f.Vars {
					vars[fmt.Sprintf("%s_%s", f.Config.Prefix, k)] = v
				}

				f.Vars = vars

				data = append(data, f)
			}

			return nil
		})

	if err != nil {
		cobra.CheckErr(err)
		return nil, err
	}

	return data, nil
}

func readFile(path string) (*InputFile, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var inputFile = &InputFile{}
	err = yaml.Unmarshal(f, inputFile)

	if err != nil {
		return nil, err
	}

	return inputFile, nil
}

func ZipVariablesByKey(files []*InputFile, variables []*airflow.Variable) map[string]*VersionedVariable {
	versionedVars := make(map[string]*VersionedVariable)

	for _, file := range files {
		for k, v := range file.Vars {
			y, err := JYaml.Marshal(v)
			cobra.CheckErr(err)

			av := &VersionedVariable{
				New: airflow.NewVariable(),
			}

			val, err := JYaml.YAMLToJSON(y)
			cobra.CheckErr(err)

			av.New.SetKey(k)
			av.New.SetValue(string(val))

			versionedVars[k] = av
		}
	}

	for _, v := range variables {
		if _, ok := versionedVars[*v.Key]; ok {
			versionedVars[*v.Key].Prev = v
		} else {
			versionedVars[*v.Key] = &VersionedVariable{
				Prev: v,
			}
		}
	}

	return versionedVars
}
