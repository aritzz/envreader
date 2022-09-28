package envreader_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/aritzz/envreader"
)

// Structures for testing purpose
type User struct {
	Id          int       `env:"ID"`
	IdEnv       int       `env:"IDENV" default:"100"`
	Name        string    `env:"NAME" default:"Pepito"`
	Age         int64     `env:"AGE" default:"37"`
	AgeInt32    int32     `env:"AGE" default:"400"`
	Size        float64   `env:"SIZE" default:"4.8"`
	Hobbies     []string  `env:"HOBBIES" default:"dancing, sleeping,gaming"`
	Counter     []int     `env:"COUNTER" default:"1, 7,2,3"`
	DataFloat   []float32 `env:"DATAFLOAT" default:"1.1,2.2,2.3,3.1"`
	DataFloat64 []float64 `env:"DATAFLOAT64" default:"1.1,2.2,2.3,3.1"`
	Enabled     bool      `env:"ENABLED" default:"true"`
	Email       string
}

var (
	testErrFormat string = "`%v` = '%v', expected '%v'"
)

func TestTagName(t *testing.T) {
	reader := envreader.EnvReader{}
	reader.Init()
	tagName := "environ"
	reader.SetTagName(tagName)

	if reader.GetTagName() != tagName {
		t.Fatalf(testErrFormat, "reader.GetTagName() ", reader.GetTagName(), tagName)
	}
}

func TestTagNameDefault(t *testing.T) {
	reader := envreader.EnvReader{}
	reader.Init()

	tagName := "environ-default"
	reader.SetTagNameDefault(tagName)

	if reader.GetTagNameDefault() != tagName {
		t.Fatalf(testErrFormat, "reader.GetTagNameDefault() ", reader.GetTagNameDefault(), tagName)
	}
}

func TestReader(t *testing.T) {

	// We will load all environment variables into this structure, replacing if needed
	var user = &User{
		Id:    1,
		Name:  "John Doe",
		Email: "john@example",
	}

	// Set IDENV to 2 as environment value
	os.Setenv("IDENV", "2")

	// This is the expected result after processing the user object
	var userToTest = User{
		Id:          1,
		IdEnv:       2,
		Name:        "Pepito",
		Age:         37,
		AgeInt32:    400,
		Size:        4.8,
		Hobbies:     []string{"dancing", "sleeping", "gaming"},
		Counter:     []int{1, 7, 2, 3},
		DataFloat:   []float32{1.1, 2.2, 2.3, 3.1},
		DataFloat64: []float64{1.1, 2.2, 2.3, 3.1},
		Enabled:     true,
		Email:       "john@example",
	}

	// Initialize reader and read into the user structure
	reader := envreader.EnvReader{}
	reader.Init()
	err := reader.Read(user)
	if err != nil {
		t.Fatalf(testErrFormat, "reader.Read(user)", err.Error(), nil)
	}

	// Compare both structures and throw an error if is not equal
	if !reflect.DeepEqual(*user, userToTest) {
		t.Fatalf(testErrFormat, "reflect.DeepEqual(*user, userToTest)", false, true)
	}
}

func TestReaderFailInt(t *testing.T) {
	var user = &User{
		Id:    1,
		Name:  "John Doe",
		Email: "john@example",
	}

	// Set INT value as a string
	os.Setenv("IDENV", "hello")

	reader := envreader.EnvReader{}
	reader.Init()
	err := reader.Read(user)
	if err == nil {
		t.Fatalf(testErrFormat, "reader.Read(user)", nil, "error")
	}
	os.Setenv("IDENV", "")
}

func TestReaderFailFloat(t *testing.T) {
	var user = &User{
		Id:    1,
		Name:  "John Doe",
		Email: "john@example",
	}

	// Set FLOAT value as a string

	os.Setenv("SIZE", "AOSFJOEPWF OAS DDSA FJIODSA FIO DASIOFJ 2!!")

	reader := envreader.EnvReader{}
	reader.Init()
	err := reader.Read(user)
	if err == nil {
		t.Fatalf(testErrFormat, "reader.Read(user)", nil, "error")
	}
	os.Setenv("SIZE", "")
}

func ExampleEnvReader() {
	// Package initialization
	reader := envreader.EnvReader{}
	reader.Init()
}

func ExampleEnvReader_Read() {
	// Package initialization
	reader := envreader.EnvReader{}
	reader.Init()

	// Define some data structure with environment values
	// and defaults
	type Config struct {
		ListenHost string `env:"LISTEN_HOST" default:"127.0.0.1"`
		ListenPort string `env:"LISTEN_PORT" default:"5000"`
		Debug      bool   `env:"ENABLE_DEBUG"`
		Numbers    []int  `env:"NUMBERS" default:"1,2,3,4"`
	}

	// Initialize this structure
	var configuration Config

	// Read values into structure
	if err := reader.Read(&configuration); err != nil {
		panic(err)
	}

	// Print values if you want
	fmt.Printf("%+v\n", configuration)
	// Output: {ListenHost:127.0.0.1 ListenPort:5000 Debug:false Numbers:[1 2 3 4]}
}
