# khata
**ḵaṭā - ﺧَﻄَﺄ**

*Khata is a simple and easy-to-use error management library for Go. It helps you write cleaner and more readable code by providing a more expressive way to handle errors. With Khata, you can easily add context to your errors, wrap them with additional information, and even customize the error messages.* 

*But that's not all. Khata also helps you with debugging by adding a lot more capabilities to the default errors. You can extract the stack trace, print the error in a more readable format, and even log the error to a file or a remote service. All of this makes it easier to diagnose and fix issues in your code.*

*In short, Khata is a must-have library for any Go developer who wants to write better code and improve their debugging experience.*

## Installation

To install Khata, simply run the following command in your terminal:

```bash
go get github.com/cmseguin/khata
```

## Usage

### Creating an Error

To create an error, you can use the `khata.New` function. It takes a message as an argument and returns a reference to a new khata error object. The message can be any string value, but it's recommended to use a constant or a variable so that you can easily change it later if needed.

```go
err := khata.New("something went wrong")
```

### Wrapping an error

To wrap an error, you can use the `khata.Wrap` function. It takes an error object as arguments and returns a reference to a new khata error object.

```go
err := errors.New("something went wrong")
// ...
k := khata.Wrap(err)
```

### Adding context to an error

To add context to an error, multiple methods are available. You can use most of them directly on the error object. The following methods are available:

- `SetCode(code int)`: Sets the error code.
- `SetExitCode(code int)`: Sets the exit code.
- `SetError(err error)`: Sets the error wraped in the khata object.
- `SetType(type string)`: Sets the type of the error.
- `SetProperty(key string, value interface{})`: Sets a custom property on the error object.
- `RemoveProperty(key string)`: Removes a custom property from the error object.
- `Explain(message string)`: Adds an explanation to the error.
- `Explainf(format string, args ...interface{})`: Adds an explanation to the error using a format string.

Those methods can be chained together as they all return a reference to the error object.

```go
err := khata.New("something went wrong").
    SetCode(500).
    SetExitCode(1).
    SetError(err).
    SetType("internal").
    SetProperty("foo", "bar")
```

### Reading the context of the error

To read the context of an error, multiple methods are also available. You can use most of them directly on the error object. However these methods cannot be chained because they do not return the reference to the khata error. The following methods are available:

- `Code() int`: Returns the error code.
- `ExitCode() int`: Returns the exit code.
- `Error() string`: Returns the wrapped error's message.
- `Type() string`: Returns the type of the error.
- `PropertiesKeys() []string`: Returns the keys of the custom properties.
- `GetProperty(key string) interface{}`: Returns the value of a custom property.
- `HasProperty(key string) bool`: Returns whether a custom property exists.
- `Explanations() []KhataExplanation`: Returns the explanations of the error.
- `Trace() []KhataTrace`: Returns the stack trace of the error.

```go
code := khata.Code(err)
```

### Utility methods on the error object

The khata error object also provides some utility methods. The following methods are available:

- `Is(err error) bool`: Returns whether the error is the same as the one passed as an argument.
- `IsAny(errs ...error) bool`: Returns whether the error is the same as one of the errors passed as arguments.
- `IsType(type string) bool`: Returns whether the error has the same type as the one passed as an argument.
- `IsAnyType(types ...string) bool`: Returns whether the error has the same type as one of the types passed as arguments.
- `IsCode(code int) bool`: Returns whether the error has the same code as the one passed as an argument.
- `IsAnyCode(codes ...int) bool`: Returns whether the error has the same code as one of the codes passed as arguments.
- `IsExitCode(code int) bool`: Returns whether the error has the same exit code as the one passed as an argument.
- `IsAnyExitCode(codes ...int) bool`: Returns whether the error has the same exit code as one of the codes passed as arguments.
- `IsTemplate(template *Template) bool`: Returns whether the error has the same template as the one passed as an argument.
- `IsAnyTemplate(templates ...*Template) bool`: Returns whether the error has the same template as one of the templates passed as arguments.
- `IsFatal() bool`: Returns whether the error is fatal. Fatal errors are errors that have an exit code other than -1.

```go
if khata.Is(err) {
  // do something if khata wraps the err
}

if khata.IsType() {
  // do something if khata has the same type as the one passed as an argument
}

// etc...
```

### Printing the error

To print the error, you can use the `khata.Debug` function. This function will output a lot of information about the error, including the message, the code, the type, the explanations, the stack trace, and the custom properties. It's very useful for debugging purposes.

```go
khata.Debug(err)
```

You can expect to see something like this:

```
Not Found

=== Explanations
  file_name.go:273 (package_name.MyFunctionName)
  └── This is an explanation of not found
  file_name.go:274 (ackage_name.MyFunctionName)
  └── This is an other explanation of not found

=== Trace
  file_name.go:281 (package_name.MyFunctionName)
  /some/path/that/exec/go/main.go:1576 (main.main)

=== Details
  Error Type: HTTP
  Error Code: 404
  Exit Code: 1
  Error At: 2023/07/02 04:27:35 0.050ms
  Handled At: 2023/07/02 04:27:35 0.100ms
  Enlapse Time: 0.050s

=== Properties
  test  -> testValue
  test2 -> testValue2
```

### Generate a json representation of the error

To generate a json representation of the error, you can use the `khata.ToJSON` function. This function will return a string containing the json representation of the error. It's very useful for logging purposes.

```go
jsonStr := khata.ToJSON()
```

### Using templates

Khata also provides a way to create error templates. Those can be very powerful when you need to create multiple errors with the same context. To create a template, you can use the `khata.NewTemplate` function. It returns a reference to the newly created template object. From the template object, you can use the following methods to generate errors:

- `New() *Khata`: Generates a new khata error from the template.
- `NewWithMessage(message string) *Khata`: Generates a new khata error from the template with given message.
- `Wrap(err error) *Khata`: Wraps an error with the template.

```go
httpError := khata.NewTemplate().
    SetMessage("something went wrong")
    SetExitCode(-1).
    SetType("HTTP")

InternalServerError := httpErrorTemplate.Extend().
    SetCode(500).
    SetMessage("internal server error")

// ...
NotFoundServerError := httpErrorTemplate.Extend().
    SetCode(404).
    SetMessage("not found")
// ...
NotFoundServerError.New()
```

### Setting context on the template

To set context on a template, multiple methods are available. You can use most of them directly on the template object. The following methods are available:

- `SetMessage(message string) *KhataTemplate`: Sets the default message of the template. Defaults to "error".
- `SetCode(code int) *KhataTemplate`: Sets the error code.
- `SetExitCode(code int) *KhataTemplate`: Sets the exit code.
- `SetType(type string) *KhataTemplate`: Sets the type of the error.
- `SetProperty(key string, value interface{}) *KhataTemplate`: Sets a custom property on the error object.
- `RemoveProperty(key string) *KhataTemplate`: Removes a custom property from the error object.

### Accessing the context on the template

To access the context on a template, multiple methods are available. You can use most of them directly on the template object. The following methods are available:

- `Message() string`: Returns the default message of the template. Defaults to "error".
- `Code() int`: Returns the error code of the template.
- `ExitCode() int`: Returns the exit code of the template.
- `Type() string`: Returns the type of the template.
- `PropertiesKeys() []string`: Returns the keys of the custom properties of the template.
- `HasProperty(key string) bool`: Returns whether the template has a custom property with the given key.
- `GetProperty(key string) interface{}`: Returns the value of the custom property with the given key.

### Modifying an existing khata error with a template

You can also modify an existing khata error with a template. To do so, you can use the `FillWithTemplate` or `OverwriteWithTemplate` method on the khata error. The `FillWithTemplate` method will only fill the properties & fields that were not altered on the error, while the `OverwriteWithTemplate` method will overwrite all the properties & fields of the error with the ones from the template.

```go

HttpInternalError := khata.NewTemplate().
    SetCode(500).
    SetMessage("something went wrong")
    SetExitCode(-1).
    SetType("HTTP")

someUnknownKhataError = khata.New("random error")

// ...

someUnknownKhataError.FillWithTemplate(HttpInternalError)
// OR
someUnknownKhataError.OverwriteWithTemplate(HttpInternalError)
```

### Truncating the package or the file paths

You might find that your errors are too verbose, and that the package and file paths are too long. Often you don't really need to see the full path of your files when debugging. In that case, you can set the following environment variables to truncate the package and file paths:

- `KHATA_FUNC_TRUNC_PREFIX` (default: `""`): If set, the package name will be truncated by removing the prefix from the beginning.
- `KHATA_FUNC_TRUNC_PREFIX` (default: `""`): If set, the file paths will be truncated by removing the prefix from the beginning of the path.

## Contributing

Contributions are welcome! Feel free to open an issue or a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.