# Environment reader

EnvReader is a simple library to ease environment variable parsing into a given structure. This library reads a given structure, parses the wanted environment value and defaults, and fills the structure with this values.

# How to use

In order to use this library, you need to define a data structure in Go, for example, we can define a simple configuration:

```
type Config struct {
    ListenHost string
    ListenPort string
    Debug      bool
    Numbers    []int
}
```

Once is defined, you need to define environment variable tags. You also can define a fallback value if this environment variable does not exist. In this example, we will define a default value for ListenHost and ListenPort, but not for Debug. Slice types are splitted with comma (,) value.

```
type Config struct {
    ListenHost string `env:"LISTEN_HOST" default:"127.0.0.1"`
    ListenPort string `env:"LISTEN_PORT" default:"5000"`
    Debug      bool   `env:"ENABLE_DEBUG"`
    Numbers    []int  `env:"NUMBERS" default:"1,2,3,4"`
}
```

Now, we will initialize the library and read all environment values into the structure:

```
var configuration Config

reader := envreader.EnvReader{}
reader.Init()
if err := reader.Read(&configuration); err != nil {
    panic(err)
}
```

If we get no error, now we can take a look on our _configuration_ variable to check all variables are loaded correctly:

```
fmt.Printf("%+v\n", configuration)
```

Expected output:
```
{ListenHost:127.0.0.1 ListenPort:5000 Debug:false Numbers:[1 2 3 4]}
```

# Supported datatypes

Currently, the following basic datatypes are supported:
- Boolean
- Integers: Int, Int64, Int32, Int16, Int8
- Floats: Float64, Float32
- Slices: []string, []int, []float32, []float64

# Examples

Check the `examples/` directory to get some examples.

# Documentation & testing

This library is documented using GoDoc, so you can access to this documentation [here](). Or you can check directly using this command:
```
$ godoc -http=:6060
```

# License
EnvReader is released under GNU General Public License. For more details, take a look at the LICENSE
