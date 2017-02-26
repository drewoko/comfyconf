#ComfyConf [![Build Status](https://travis-ci.org/drewoko/comfyconf.svg?branch=master)](https://travis-ci.org/drewoko/comfyconf) [![Coverage Status](https://coveralls.io/repos/github/drewoko/comfyconf/badge.svg?branch=master)](https://coveralls.io/github/drewoko/comfyconf?branch=master)

Yet another configuration framework for GO, inspired by different existing frameworks and with additional 
functionality, like parsing arrays, parameter existence and structure assembly. 
This framework uses middleware conception for allowing to make complex applications with different configuration sources. 

**Note! This library is not production ready, API may change soon**.

## Features

- Simple in use, API is quite similar to standard `flags` library
- Middleware model with already existing JSON, Command Line (flags) and Environment middlewares
- Parsing strings, integers, booleans and slices from configuration sources
- Parameter existence in middleware
- Custom help printer

## Installation
For downloading latest version:
    
    go get github.com/drewoko/comfyconf

gopkg.in link will be available when API will be freezed.

## Getting Started

### Initialization

For creating ComfyConf instance need to call New function from package, defining middlewares.
```go
conf := comfyconf.New(comfyconf.NewFlags(), comfyconf.NewEnv())
```
alternatively middlewares can be defined separately by calling AddMiddleware
```go
conf.AddMiddleware(comfyconf.NewFlags())
``` 

### Performing parameter parsing

For performing parameter parsing, need to call `Parse` function after parameter declaration.

```go
conf := comfyconf.New(comfyconf.NewFlags(), comfyconf.NewEnv())
//Parameter declaration here
conf.Parse()
```

### Declaring parameters

ComfyConf package have four types for configuration variables. These are Integer, String, Boolean and Interface Slice.
For all variable types are two methods of declaration - declaration with passing existing variable or declaration with variable creation.

They have similar declaration signature (Except Existence)

    val := conf.TYPE("shortName", "fullName", defaultValue, "Description string, that can be used for printing help")

or if you want to use own variable then
    
    var type TYPE
    conf.TYPEVar("shortName", "fullName", defaultValue, &type, "Description string, that can be used for printing help")

#### Integer
```go
intParam := conf.Int("short", "FullName", 10, "Description for that parameter")
``` 
or
```go
var intParam int
conf.IntVar("short", "FullName", 10, &intParam, "Description for that parameter")
``` 

#### String
```go
strParam := conf.String("short", "FullName", "Default String", "Description for that parameter")
``` 
or
```go
var strParam string
conf.StringVar("short", "FullName", "Default String", &strParam, "Description for that parameter")
``` 

#### Boolean
```go
boolParam := conf.Bool("short", "FullName", true, "Description for that parameter")
``` 
or
```go
var boolParam bool
conf.BoolVar("short", "FullName", true, &boolParam, "Description for that parameter")
``` 

#### Slice
```go
sliceParam := conf.Slice("short", "FullName", make([]interface{}, 0), "Description for that parameter")
``` 
or
```go
var sliceParam []interface{}
conf.SliceVar("short", "FullName", make([]interface{}, 0), &sliceParam, "Description for that parameter")
``` 

#### Existence
Parameter that returns true if parameter with selected name exists in configuration source
```go
boolParam := conf.Exist("short", "FullName", "Description for that parameter")
``` 
or
```go
var boolParam bool
conf.ExistVar("short", "FullName", &boolParam, "Description for that parameter")
``` 

### Middlewares

All middlewares should implement Middleware interface, so you can make own middleware.

Order of adding middlewares is important. Latest added middleware have higher priority. That means that if same parameter
exists for Env and Flags, and Flags is defined last, then value will be taken from Flags middleware.

#### Command line flags

Command line flags middleware are Flags struct, that implements Middleware interface.
For making developers life easier we have functions that prepares struct for work.

```go
NewFlags()
```
will make Flags struct with `=` assignment, what means that `--full=Test -s=t` will be parsed.

```go
NewFlagsWithCustomAssignment(" ")
```
Same NewFlags, but with custom assignment, which makes possible to parse  `--full=val -s val` flags.

```go
NewFlagsWithCustomParser(func(arg string) (string, string) {
    return "key", "value"
})
```
Another possibility for command line flags parsing. You can provide own function, that will parameter one by one returning key-value pair. 

#### Environment

Environment middleware are Env struct, that implements Middleware interface. 

```go
NewEnv()
```
will return Environment middleware with `ENV_` prefix for environment parameters. That means, all environment parameters 
with name starting on `ENV_` will be taken in account. 

For setting custom prefix function `NewEnvWithPrefix` can be used.

```go
NewEnvWithPrefix("TEST_")
```

#### JSON

JSON middleware (struct) allows to read configuration parameters from JSON files. FullName in that environment is used as JSON path.

For example we have JSON
```json
{
  "level1": {
    "level2": "Value"
  }
}
```
then full name (path) to "Value" in json will be `level1.level2` and short name will be `level2`.

For creating JSON middleware with predefined path `NewJSON` function should be used
```go 
NewJSON("path/to/config.json")
```

Also can be provided custom file reader using `NewJSONWithCustomReader`. It receives function, that returns tuple with
byte array and error.
```go 
NewJSONWithCustomReader(func(j *JSON) ([]byte, error) {
    return []byte(`
        {}
    `), nil
})
```

`NewJSONWithCustomFileMiddleware` powerful function, that receives middlewares, that used for obtaining JSON configuration path.
 
```go
NewJSONWithCustomFileMiddleware("shortName", "fullName", "/path/to/config.json", NewFlags(), NewEnv())
```

## Other

### ToStruct

Library tries to map configuration parameters to predefined struct using `comfyname` tag by it value. It can work with referenced and 
with direct values.

```go
var testStruct struct {
    TestString string        `comfyname:"test"`
    TestPrtString *string    `comfyname:"t2"`
}

conf := comfyconf.New(comfyconf.NewFlags(), comfyconf.NewEnv())

conf.String("t", "test", "", "Basic description")
conf.String("t2", "test2", "", "Basic description")

conf.ToStruct(&testStruct)
```

### PrintHelp

Function for printing help. It receives function, that get as parameter options `map[OptionKey]*Option`. ComfyConf have
default help printer called `DefaultHelpPrinter`

```go
conf.PrintHelp(DefaultHelpPrinter)
```






