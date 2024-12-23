package configuration_yaml_file

import (
	"reflect"
	"testing"
	"time"

	"github.com/BoRuDar/configuration/v5"
	"gopkg.in/yaml.v2"
)

func TestYAMLFileProvider(t *testing.T) {
	type test struct {
		Timeout time.Duration `file_yml:"service.timeout"`
	}

	testObj := test{}
	expected := test{
		Timeout: time.Millisecond * 15,
	}

	fieldType := reflect.TypeOf(&testObj).Elem().Field(0)
	fieldVal := reflect.ValueOf(&testObj).Elem().Field(0)

	p := NewYAMLFileProvider("./testdata/input.yml")
	if err := p.Init(&testObj); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := p.Provide(fieldType, fieldVal); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert(t, expected, testObj)
}

func TestFindValStrByPath(t *testing.T) {
	type service struct {
		Beta int `file_yml:"service.alfa.beta"`
	}

	type testStruct struct {
		Name    string        `file_yml:"service.name"`
		Timeout time.Duration `file_yml:"service.timeout"`
		Alfa    service
	}

	var testObj any
	data, _ := yaml.Marshal(testStruct{
		Name: "test",
		Alfa: service{Beta: 42},
	})
	_ = yaml.Unmarshal(data, &testObj)

	tests := []struct {
		name         string
		input        any
		path         []string
		expectedStr  string
		expectedBool bool
	}{
		{
			name:         "empty path",
			path:         nil,
			expectedStr:  "",
			expectedBool: false,
		},
		{
			name:         "at root level | Name | json",
			input:        testObj,
			path:         []string{"Name"},
			expectedStr:  "test",
			expectedBool: true,
		},
		{
			name:         "substructures | Alfa.Beta | json",
			input:        testObj,
			path:         []string{"Alfa", "Beta"},
			expectedStr:  "42",
			expectedBool: true,
		},
		{
			name:         "not found",
			input:        testObj,
			path:         []string{"notfound"},
			expectedStr:  "",
			expectedBool: false,
		},
	}

	for _, tt := range tests {
		test := tt

		t.Run(test.name, func(t *testing.T) {
			gotStr, gotBool := findValStrByPath(test.input, test.path)
			if gotStr != test.expectedStr || gotBool != test.expectedBool {
				t.Fatalf("expected: [%q %v] but got [%q %v]", test.expectedStr, test.expectedBool, gotStr, gotBool)
			}
		})
	}
}

func TestFileProvider_Init(t *testing.T) {
	type i struct {
		Test int `file_json:"void."`
	}

	_, err := configuration.New[i](NewYAMLFileProvider("./testdata/dummy.file"))
	assert(t, "cannot init [YAMLFileProvider] provider: file must have .yaml/.yml extension", err.Error())
}

func TestConfiguration(t *testing.T) {
	type Cfg struct {
		Name string `file_yml:"service.name"`
	}

	cfg, err := configuration.New[Cfg](
		NewYAMLFileProvider("./testdata/input.yml"),
	)

	assert(t, nil, err)
	assert(t, "serv1", cfg.Name)
}
