// Package envreader is a library to ease environment variable parsing into a given structure.
// This library reads a given structure, parses the wanted environment value and defaults, and fills the structure with this values.
package envreader

import (
	"os"
	"reflect"
	"strconv"
	"strings"
)

// EnvReader Environment reader main object with tag definitions
type EnvReader struct {
	tag        string
	tagDefault string
}

// define default tag names
const (
	tagName        string = "env"
	tagNameDefault string = "default"
)

// Init Initializes environment reader with default values
func (env *EnvReader) Init() {
	(*env).tag = tagName
	(*env).tagDefault = tagNameDefault
}

// SetTagName Set custom tag name for this instance of the reader
func (env *EnvReader) SetTagName(name string) {
	(*env).tag = name
}

// SetTagNameDefault Set custom tag name for default values
func (env *EnvReader) SetTagNameDefault(name string) {
	(*env).tagDefault = name
}

// GetTagName Get current tag name
func (env *EnvReader) GetTagName() string {
	return (*env).tag
}

// GetTagNameDefault Get current default tag name
func (env *EnvReader) GetTagNameDefault() string {
	return (*env).tagDefault
}

// Read Reads a structure and sets values set in environment variables
func (env *EnvReader) Read(data interface{}) error {
	elemType, elemTypeVal := reflect.TypeOf(data).Elem(), reflect.ValueOf(data)

	// Iterate over all available fields and read the tag value
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)

		tagEnvValue := field.Tag.Get((*env).tag)
		tagEnvValueDefault := field.Tag.Get((*env).tagDefault)
		elemValAddr := elemTypeVal.Elem().Field(i)
		err := (*env).setEnvByType(&elemValAddr, field.Type, tagEnvValue, tagEnvValueDefault)
		if err != nil {
			return err
		}
	}
	return nil
}

// setEnvByType Get environment value by type and set his value or fallback if any
func (env *EnvReader) setEnvByType(field *reflect.Value, envFieldType reflect.Type, envValue, envDefaultValue string) error {

	// Get environment value
	environmentInfo := (*env).getEnvFallback(envValue, envDefaultValue)
	if len(strings.TrimSpace(environmentInfo)) == 0 {
		return nil
	}

	// Get kind and bits count
	envType := envFieldType.Kind()

	switch envType {

	// Boolean: Parse numeric or true/false string values.
	case reflect.Bool:
		reflectVal := (environmentInfo == "1") || (strings.ToLower(environmentInfo) == "true")
		field.SetBool(reflectVal)

	// Integer: Read all types of integers and parse it.
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		envBits := envFieldType.Bits()
		reflectVal, err := strconv.ParseInt(environmentInfo, 10, envBits)
		if err != nil {
			return err
		}
		field.SetInt(reflectVal)

	// Float: Read float32/64 type values.
	case reflect.Float64, reflect.Float32:
		envBits := envFieldType.Bits()
		reflectVal, err := strconv.ParseFloat(environmentInfo, envBits)
		if err != nil {
			return err
		}
		field.SetFloat(reflectVal)

	// Slices: Parse int or string slices
	case reflect.Slice:
		reflectVal := (*env).createSlice(environmentInfo, envFieldType)
		field.Set(reflectVal)

	// String: Set string type slice as-is
	case reflect.String:
		field.SetString(environmentInfo)
	}

	return nil
}

// createSlice Creates a slice from a string
func (env *EnvReader) createSlice(envInfo string, envInfoType reflect.Type) reflect.Value {
	envInfoClean := strings.Split(strings.Replace(envInfo, " ", "", -1), ",")
	envInfoCleanLen := len(envInfoClean)
	reflectedSlice := reflect.MakeSlice(envInfoType, envInfoCleanLen, envInfoCleanLen)

	// Parse values depending on slice type and store it
	for count, value := range envInfoClean {
		switch envInfoType {
		// String slice
		case reflect.TypeOf([]string{}):
			reflectedSlice.Index(count).SetString(value)

		// Int slice
		case reflect.TypeOf([]int{}):
			var intVal int
			reflectVal, err := strconv.ParseInt(value, 10, reflect.TypeOf(intVal).Bits())
			if err == nil {
				reflectedSlice.Index(count).SetInt(reflectVal)
			}

		// Float32 slice
		case reflect.TypeOf([]float32{}):
			reflectVal, err := strconv.ParseFloat(value, 32)
			if err == nil {
				reflectedSlice.Index(count).SetFloat(reflectVal)
			}

		// Float64 slice
		case reflect.TypeOf([]float64{}):
			reflectVal, err := strconv.ParseFloat(value, 64)
			if err == nil {
				reflectedSlice.Index(count).SetFloat(reflectVal)
			}
		}

	}

	return reflectedSlice
}

// getEnvFallback Get environment value and provide fallback if empty
func (env *EnvReader) getEnvFallback(envName, fallback string) string {
	envVal := os.Getenv(envName)
	if len(strings.TrimSpace(envVal)) == 0 {
		return fallback
	}

	return envVal
}
