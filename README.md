# go-app

A small toolkit and microframework to create Golang applications, including features like:

- Dynamic configuration loading through struct tags (including nested structs).
- Configuration loading from multiple sources (command line flags, environment variables, configuration file).
- Dependency injection pattern for nested structs.
- Standard logging with [uber-go/zap](https://github.com/uber-go/zap).

## Creating an application

The code bellow is the minimal structure to create a new application with zero configurations and zero dependencies.

```go
package main

import (
    "github.com/protomesh/go-app"
)

// This is the root of the dependency tree.
// All dependency are discovered from the struct attributes of
// this dependency annotated with the `config` tag.
struct *root {
    // Embed pointer to auto-instantiated injector.
    // This injector provides two methods for every embedding struct:
    // - `Dependency()` which returns the generic specified in the *app.Injector[D]
    // - `Log()` returning the app.Logger interface
    //
    // Its important to note that only the root dependency can have a struct pointer
    // like *root as its concretion. All child dependency must specify the D generic
    // for *app.Injector[D] as an interface which specifies the required dependency or `any`.
    *app.Injector[*root]

}

// Helper function
func newRoot() *root {
    return &root{}
}

// Must implement the app.Dependency interface
func (r *root) Dependency() *root {
    return r
}

var opts = &app.AppOptions{
    // Which flag set to assign, parse and read flags.
    // flag.CommandLine is the default instance set for command line applications.
    FlagSet:   flag.CommandLine,
    // Print configurations and dependency tree
    Print:     os.Getenv("PRINT_CONFIG") == "true",
}

func main() {
    
    deps := newRoot()

    myApp := app.NewApp(deps, opts)
    // Always defer the app close, it is required to flush any remaining log buffer.
    defer myApp.Close()

    // Your service startup logic here...

    // Blocks until the process receives an interruption signal
    app.WaitInterruption()

}
```

## Dynamic configuration

To add dynamic configuration loading, you need to specify tags in the struct's attributes from the root of the dependency tree. The following code shows how to define a simple dependency in the **dependency root**:

```go
struct *root {
    *app.Injector[*root]

    // Dynamic configuration must be a public attribute.
    MyString app.Config `config:"my.string,str" default:"mydefaultval" usage:"usage instructions"`
}
```

The `config` tag is used to specify the configuration key and type in the form of `[config key],[type specifier]`. All valid types are documented bellow:

| Type specifier               | Name                                                     | Valid values                                                | Default value |
| ---------------------------- | -------------------------------------------------------- | ----------------------------------------------------------- | ------------- |
| `bool`, `boolean`            | Case-insensitive boolean                                 | `t`, `true`, `y`, `yes`, `f`, `false`, `n`, `no`, `not`     | `false`       |
| `int64`, `int`, `integer`    | 64-bit integer                                           | -                                                           | `0`           |
| `float64`, `float`, `double` | 64-bit floating point                                    | -                                                           | `0`           |
| `duration`                   | Duration string                                          | [Any valid duration](https://pkg.go.dev/time#ParseDuration) | `0`           |
| `str`, `string`              | Arbitrary length string                                  | -                                                           | `""`          |
| `time`, `datetime`, `date`   | [RFC3339 string](https://www.rfc-editor.org/rfc/rfc3339) | -                                                           | `time.Time{}` |

The `default` lets you specify a default value if not provided by any configuration source. And the `usage` tag is used to display a help message for each configuration when the user calls your application with the `-h` (help) flag.
Both tags, `default` and `usage` are optional.

### Nested configuration

When you're abstracting pieces of your application you may want to keep the configuration needed for each component in the component itself. Lets say we have the following component of the application:

```go
type MyNestedComponent struct {
    NestedVal app.Config `config:"val.number,int" default:"0" usage:"usage instructions"`
}

type MyComponent struct {
    ComponentVal app.Config `config:"val.text,str" default:"mydefaultval" usage:"usage instructions"`

    // In this case, we specify the child component only with the key parameter, the type specifier
    // default tag and usage instructions are not expected here.
    //
    // Note that this attribute must be public to be managebole by the go-app internals.
    //
    // The "nested" key is preffixed in the key with . as delimiter between the strings.
    MyNestedComponent *MyNestedComponent `config:"nested"`
}
```

To make the component configuration discoverable to the **go-app** framework, you just need to add it as a child configuration of your dependency root:

```go
struct *root {
    *app.Injector[*root]

    // Dynamic configuration must be a public attribute.
    MyNestedComponent *MyNestedComponent `config:"component"`
}
```

In this case the following configuration are available:

- `component.nested.val.number`
- `component.val.text`

## Dependency injection

The dependency tree injection feature is done with reflection. The first dependency is called **root dependency**, all other dependency are **nested dependencies**.

```go
struct *root {
    *app.Injector[*root]

    DatabaseConnectionString app.Config `config:"database.connection,str" usage:"Database connection string"`

    db *sql.DB
}

func newRoot() *root {
    return &root{}
}

func (r *root) Dependency() *root {
    return r
}

func (r *root) GetDB() *sql.DB {
    return r.db
}


var opts = &app.AppOptions{
    FlagSet:   flag.CommandLine,
    Print:     os.Getenv("PRINT_CONFIG") == "true",
}

func main() {
    
    deps := newRoot()

    myApp := app.NewApp(deps, opts)
    defer myApp.Close()

    // Your service startup logic here...

    db, err := sql.Open("postgres", deps.DatabaseConnectionString.StringVal())
    if err != nil {
        myApp.Log().Panic("Error connecting")
    }

    myApp.db = db

    // Blocks until the process receives an interruption signal
    app.WaitInterruption()

}
```

### Dependency tree

But as we mentioned in the configuration, you may want to specialize the components of your application abstracting responsibilities. Suppose that we have an `PetStoreService` that needs a `*sql.DB` pointer no matter where it came from (leaving the responsibility of injecting the dependency for the instantiating callee).

```go

type PetStoreServiceDependencies interface {
    GetDB() *sql.DB
}

type PetStoreService[D PetStoreServiceDependencies] struct {
    *app.Injector[D]

    Name string `config:"name,str" default:"Animaland" usage:"The brand name of the pet store"`
}

func (p *PetStoreService[D]) DoAnything() {

    p.Log().Info("Do anything in the pet store", "logkey", "logval")

    db := p.Dependency().GetDB()

    // Do anything in the pet store...
}

struct *root {
    *app.Injector[*root]

    DatabaseConnectionString app.Config `config:"database.connection,str" usage:"Database connection string"`

    db *sql.DB

    // Its not mandatory to put the config tag, its only here for demonstration.
    PetStoreService *PetStoreService[*root] `config:"pet.store"`
}
```

And you can nest dependency, note that only the root dependency can specify a pointer to itself as the generic parameter for the nested dependencies. So if you have a nested dependency to the `PetStoreService` it would be implemented as the example bellow:

```go

type PetStoreServiceDependencies interface {
    GetEmergencyChannel() <-chan string
}

type VeterinaryService[D PetStoreServiceDependencies] struct {
    *app.Injector[D]
}

func (v *VeterinaryService[D]) WatchForEmergencies() {

    emergencyCh := v.Dependency().GetEmergencyChannel()

    for message := range emergencyCh {
        v.Log().Warn("New emergency!!!", "message", message)
    }

}

var (
    // Ensure the correct implementation of methods for dependencies
    _ PetStoreServiceInjector = &PetStoreService[*root]{}
)

type PetStoreServiceInjector interface {
    GetEmergencyChannel() <-chan string
}

type PetStoreService[D PetStoreServiceDependencies] struct {
    *app.Injector[D]

    // Must be public
    VeterinaryService *VeterinaryService[PetStoreServiceInjector]
}

```