### developer notes

when developing packages, it is best to "encapsulate" as much logic inside that package as possible. so that if one were to run the package from any surrounding environment given the same parameters, it would give the same outputs. this makes the application's implicit dependency graph cleaner and makes it easy to write tests via dependency injection.

however in the case that the package requires a specific environment to produce a given result, it is best to keep everything it requires in a separate package dedicated to that environment. (then the one writing tests can add a mechanism to control that environment)

for example, when utilizing the database in multiple files, many packages require their respective structs to be created in the database. instead of putting the details of initialization into a single place (like `main.go`) it is better for each package to provide a method to initialization. that cleans up the entrypoint file which shouldn't need to care about what in the package internally needs to be initialized.

a *semi-formal* interface of this will be defined here (as golang does not support package interfaces)

```golang
type Package interface {
    Initialize() error
    Destroy() error
}
```

you are not forced to export the two methods, however the naming and `returns` should be kept consistent
