package configuration_yaml_file

import (
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/BoRuDar/configuration/v5"
	"gopkg.in/yaml.v2"
)

const (
	YAMLFileProviderName = `YAMLFileProvider`
	YAMLFileProviderTag  = `file_yml`
)

var ErrFileMustHaveYMLExt = errors.New("file must have .yaml/.yml extension")

// NewYAMLFileProvider creates new provider which reads values from YAML files.
func NewYAMLFileProvider(fileName string) *FileProvider {
	return &FileProvider{fileName: fileName}
}

type FileProvider struct {
	fileName string
	fileData any
}

func (fp *FileProvider) Name() string {
	return YAMLFileProviderName
}

func (fp *FileProvider) Tag() string {
	return YAMLFileProviderName
}

func (fp *FileProvider) Init(_ any) error {
	file, err := os.Open(fp.fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	name := strings.ToLower(fp.fileName)

	if !strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml") {
		return ErrFileMustHaveYMLExt
	}

	return yaml.Unmarshal(b, &fp.fileData)
}

func (fp *FileProvider) Provide(field reflect.StructField, v reflect.Value) error {
	path := field.Tag.Get(YAMLFileProviderTag)
	if len(path) == 0 {
		// field doesn't have a proper tag
		return fmt.Errorf("%s: key is empty", YAMLFileProviderName)
	}

	valStr, ok := findValStrByPath(fp.fileData, strings.Split(path, "."))
	if !ok {
		return fmt.Errorf("%s: findValStrByPath returns empty value", YAMLFileProviderName)
	}

	return configuration.SetField(field, v, valStr)
}

func findValStrByPath(i any, path []string) (string, bool) {
	if len(path) == 0 {
		return "", false
	}
	firstInPath := strings.ToLower(path[0])

	currentFieldIface, ok := i.(map[interface{}]interface{}) // unmarshal from yaml
	if !ok {
		return "", false
	}

	currentFieldStr := map[string]interface{}{}
	for k, v := range currentFieldIface {
		currentFieldStr[fmt.Sprint(k)] = v
	}

	if len(path) == 1 {
		val, ok := currentFieldStr[firstInPath]
		if !ok {
			return "", false
		}

		return fmt.Sprint(val), true
	}

	return findValStrByPath(currentFieldStr[firstInPath], path[1:])
}
