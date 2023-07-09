package khata

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/cmseguin/khata/internal/colors"
)

const (
	DEFAULT_EXIT_CODE  = 1
	DEFAULT_ERROR_CODE = -1
	DEFAULT_MESSAGE    = "error"
	DEFAULT_ERROR_TYPE = "KhataError"
)

type KhataTrace struct {
	file         string
	line         int
	functionName string
}

func (kt *KhataTrace) File() string {
	return kt.file
}

func (kt *KhataTrace) Line() int {
	return kt.line
}

func (kt *KhataTrace) FunctionName() string {
	return kt.functionName
}

type KhataExplanation struct {
	message      string
	file         string
	line         int
	functionName string
}

func (ke *KhataExplanation) Message() string {
	return ke.message
}

func (ke *KhataExplanation) File() string {
	return ke.file
}

func (ke *KhataExplanation) Line() int {
	return ke.line
}

func (ke *KhataExplanation) FunctionName() string {
	return ke.functionName
}

type KhataTemplate struct {
	message    string
	errorCode  int
	errorType  string
	exitCode   int
	properties map[string]interface{}
	parent     *KhataTemplate
}

// Create a new khata error with the template and given message
func (kt *KhataTemplate) NewWithMessage(message string) *Khata {
	return kt.Wrap(errors.New(message))
}

// Create a new khata error with the template
func (kt *KhataTemplate) New(message ...string) *Khata {
	var inputMessage string

	for i, m := range message {
		inputMessage += m
		if i < len(message)-1 {
			inputMessage += "\n"
		}
	}

	if inputMessage == "" {
		inputMessage = kt.message
	}

	return kt.Wrap(errors.New(inputMessage))
}

// Wraps an error with a Khata object while using the template
func (kt *KhataTemplate) Wrap(err error) *Khata {
	return &Khata{
		Err:              err,
		createdAt:        time.Now().UTC(),
		errorCode:        kt.errorCode,
		errorType:        kt.errorType,
		exitCode:         kt.exitCode,
		explanationStack: []KhataExplanation{},
		properties:       kt.properties,
		traceStack:       []KhataTrace{},
		template:         kt,
	}
}

func (kt *KhataTemplate) Apply(k *Khata) *Khata {
	k.errorCode = kt.Code()
	k.errorType = kt.Type()
	k.exitCode = kt.ExitCode()

	for key, value := range kt.properties {
		k.properties[key] = value
	}

	k.template = kt

	return k
}

// Allows to extend or copy the template
func (kt *KhataTemplate) Extend() *KhataTemplate {
	return &KhataTemplate{
		errorCode:  kt.errorCode,
		errorType:  kt.errorType,
		exitCode:   kt.exitCode,
		properties: kt.properties,
		message:    kt.message,
		parent:     kt,
	}
}

// Returns the message associated with the template
func (kt *KhataTemplate) Message() string {
	return kt.message
}

// Sets the message associated with the template
func (kt *KhataTemplate) SetMessage(message string) *KhataTemplate {
	kt.message = message
	return kt
}

// Returns the error code associated with the template
func (kt *KhataTemplate) Code() int {
	return kt.errorCode
}

// Sets the error code associated with the template
func (kt *KhataTemplate) SetCode(code int) *KhataTemplate {
	kt.errorCode = code
	return kt
}

// Returns the error type associated with the template
func (kt *KhataTemplate) Type() string {
	return kt.errorType
}

// Sets the error type associated with the template
func (kt *KhataTemplate) SetType(errorType string) *KhataTemplate {
	kt.errorType = errorType
	return kt
}

// Returns the exit code associated with the template
func (kt *KhataTemplate) ExitCode() int {
	return kt.exitCode
}

// Sets the exit code associated with the template
func (kt *KhataTemplate) SetExitCode(code int) *KhataTemplate {
	kt.exitCode = code
	return kt
}

// Set a property on the template
func (kt *KhataTemplate) SetProperty(key string, value interface{}) *KhataTemplate {
	kt.properties[key] = value
	return kt
}

// Returns the keys of all the properties set on the error
func (kt *KhataTemplate) PropertiesKeys() []string {
	keys := make([]string, len(kt.properties))
	for key := range kt.properties {
		keys = append(keys, key)
	}
	return keys
}

// Returns the value of a property associated with the template
func (kt *KhataTemplate) GetProperty(key string) interface{} {
	return kt.properties[key]
}

// Remove a property from the template
func (kt *KhataTemplate) RemoveProperty(key string) *KhataTemplate {
	delete(kt.properties, key)
	return kt
}

// Check if the template has the given property
func (kt *KhataTemplate) HasProperty(key string) bool {
	return kt.properties[key] != nil
}

// Returns true if the template's parent is the same as the given template
func (t *KhataTemplate) IsInstanceOf(kt *KhataTemplate) bool {
	return t.parent == kt
}

// Returns true if the template is the same as the given's template parent
func (t *KhataTemplate) IsParentOf(kt *KhataTemplate) bool {
	return kt.parent == t
}

// Returns true if the templates's parent or any of its parents is the same as the given template
func (t *KhataTemplate) IsRelatedTo(kt *KhataTemplate) bool {
	templateToCheck := t.parent
	for {
		if templateToCheck == nil {
			return false
		}
		if templateToCheck == kt {
			return true
		}
		templateToCheck = templateToCheck.parent
	}
}

type Khata struct {
	errorCode        int
	errorType        string
	exitCode         int
	createdAt        time.Time
	Err              error
	properties       map[string]interface{}
	traceStack       []KhataTrace
	explanationStack []KhataExplanation
	template         *KhataTemplate
}

// Expose the error so it behaves like a normal error
func (k *Khata) Error() string {
	return k.Err.Error()
}

// Allows you to change the initial error (This should be used with caution)
func (k *Khata) SetError(err error) *Khata {
	k.Err = err
	return k
}

// Check if the error is the same as the given error
func (k *Khata) Is(err error) bool {
	if err == nil || k.Err == nil {
		return false
	}
	return k.Err == err
}

// Check if the error is any of the given errors
func (k *Khata) IsAny(errs ...error) bool {
	for _, err := range errs {
		if k.Is(err) {
			return true
		}
	}
	return false
}

// Returns true if the error's template is the same as the given template
func (k *Khata) IsInstanceOf(kt *KhataTemplate) bool {
	return k.template == kt
}

// Returns true if the error's template or any of its parents is the same as the given template
func (k *Khata) IsRelatedTo(kt *KhataTemplate) bool {
	templateToCheck := k.template
	for {
		if templateToCheck == nil {
			return false
		}
		if templateToCheck == kt {
			return true
		}
		templateToCheck = templateToCheck.parent
	}
}

// Returns the trace stack as an array of objects. Will be computed at the time of calling.
func (k *Khata) Trace() []KhataTrace {
	trace := collectTrace()
	k.traceStack = trace

	return k.traceStack
}

// Returns the code of the error. If not set, defaults to -1
func (k *Khata) Code() int {
	return k.errorCode
}

// Set the error code. If not set, defaults to -1
func (k *Khata) SetCode(code int) *Khata {
	k.errorCode = code
	return k
}

// Check if the error is the same as the given code
func (k *Khata) IsCode(code int) bool {
	return k.errorCode == code
}

// Check if the error is any of the given codes
func (k *Khata) IsAnyCode(codes ...int) bool {
	for _, code := range codes {
		if k.IsCode(code) {
			return true
		}
	}
	return false
}

// The type of the error. If not set, defaults to "KhataError"
func (k *Khata) Type() string {
	return k.errorType
}

// Allows you to change the type of the error
func (k *Khata) SetType(errorType string) *Khata {
	k.errorType = errorType
	return k
}

// Check if the error is the same as the given type
func (k *Khata) IsType(errorType string) bool {
	return k.errorType == errorType
}

// Check if the error is any of the given types
func (k *Khata) IsAnyType(errorTypes ...string) bool {
	for _, errorType := range errorTypes {
		if k.IsType(errorType) {
			return true
		}
	}
	return false
}

// Returns the keys of all the properties set on the error
func (k *Khata) PropertiesKeys() []string {
	keys := make([]string, len(k.properties))
	for key := range k.properties {
		keys = append(keys, key)
	}
	return keys
}

// Returns the value for a given property key. If the property is not set, returns nil
func (k *Khata) GetProperty(key string) interface{} {
	return k.properties[key]
}

// Check if a property is set within the khata error
func (k *Khata) HasProperty(key string) bool {
	return k.properties[key] != nil
}

// Set a property within the khata error
func (k *Khata) SetProperty(key string, value interface{}) *Khata {
	k.properties[key] = value
	return k
}

// Remove a property from the khata error
func (k *Khata) RemoveProperty(key string) *Khata {
	delete(k.properties, key)
	return k
}

// Returns the exit code from the khata error. If not set, defaults to 1
func (k *Khata) ExitCode() int {
	return k.exitCode
}

// Set the exit code of the program. If not set, defaults to 1
func (k *Khata) SetExitCode(code int) *Khata {
	k.exitCode = code
	return k
}

// Check if the error exit code is the same as the given exit code
func (k *Khata) IsExitCode(code int) bool {
	return k.exitCode == code
}

// Check if the error exit code is any of the given exit codes
func (k *Khata) IsAnyExitCode(codes ...int) bool {
	for _, code := range codes {
		if k.IsExitCode(code) {
			return true
		}
	}
	return false
}

// Returns the explanations for the error
func (k *Khata) Explanations() []KhataExplanation {
	return k.explanationStack
}

// Add an explanation to the error
func (k *Khata) Explain(explanation string) *Khata {
	lastTrace := collectCallerTrace()

	k.explanationStack = append(k.explanationStack, KhataExplanation{
		message:      explanation,
		file:         lastTrace.file,
		line:         lastTrace.line,
		functionName: lastTrace.functionName,
	})

	return k
}

// Explainf is a wrapper around Explain that accepts a format string
func (k *Khata) Explainf(format string, args ...interface{}) *Khata {
	lastTrace := collectCallerTrace()

	k.explanationStack = append(k.explanationStack, KhataExplanation{
		message:      fmt.Sprintf(format, args...),
		file:         lastTrace.file,
		line:         lastTrace.line,
		functionName: lastTrace.functionName,
	})

	return k
}

// Check if the error is fatal. Fatal errors are those that should stop the program.
func (k *Khata) IsFatal() bool {
	return k.exitCode != -1
}

// Print the error in a console friendly way
func (k *Khata) Debug() *Khata {
	handledAt := time.Now().UTC()
	diff := handledAt.Sub(k.createdAt)

	// Print error
	p := fmt.Sprintf(
		"\n%s%s%s",
		colors.BoldRed,
		k.Err,
		colors.Reset,
	)
	println(p)

	// Print explanations
	explanations := k.Explanations()

	println(fmt.Sprintf("\n=== %sExplanations%s", colors.BoldYellow, colors.Reset))

	for _, explanation := range explanations {
		file := tryTrimmingPath(explanation.file)
		funcName := tryTrimmingFunc(explanation.functionName)
		p := fmt.Sprintf(
			"  %s%s%s:%s%d%s (%s%s%s)\n  └── %s%s%s",
			colors.UnderlineGray,
			file,
			colors.Reset,
			colors.Green,
			explanation.line,
			colors.Reset,
			colors.Cyan,
			funcName,
			colors.Reset,
			colors.BoldWhite,
			explanation.message,
			colors.Reset,
		)
		fmt.Println(p)
	}

	// Print trace
	trace := k.Trace()

	println(fmt.Sprintf("\n=== %sTrace%s", colors.BoldYellow, colors.Reset))

	for _, trace := range trace {
		file := tryTrimmingPath(trace.file)
		funcName := tryTrimmingFunc(trace.functionName)
		p := fmt.Sprintf(
			"  %s%s%s:%s%d%s (%s%s%s)",
			colors.UnderlineGray,
			file,
			colors.Reset,
			colors.Green,
			trace.line,
			colors.Reset,
			colors.Cyan,
			funcName,
			colors.Reset,
		)
		fmt.Println(p)
	}

	println(fmt.Sprintf("\n=== %sDetails%s", colors.BoldYellow, colors.Reset))

	println(fmt.Sprintf("  %sError Type%s: %s%s%s", colors.BoldWhite, colors.Reset, colors.Cyan, k.errorType, colors.Reset))
	println(fmt.Sprintf("  %sError Code%s: %s%d%s", colors.BoldWhite, colors.Reset, colors.Cyan, k.errorCode, colors.Reset))
	println(fmt.Sprintf("  %sExit Code%s: %s%d%s", colors.BoldWhite, colors.Reset, colors.Cyan, k.exitCode, colors.Reset))
	println(fmt.Sprintf("  %sError At%s: %s%s%s", colors.BoldWhite, colors.Reset, colors.Cyan, k.createdAt.Format("2006/01/02 15:04:05 0.000ms"), colors.Reset))
	println(fmt.Sprintf("  %sHandled At%s: %s%s%s", colors.BoldWhite, colors.Reset, colors.Cyan, handledAt.Format("2006/01/02 15:04:05 0.000ms"), colors.Reset))
	println(fmt.Sprintf("  %sEnlapse Time%s: %s%.3fs%s", colors.BoldWhite, colors.Reset, colors.Cyan, (float64(diff.Milliseconds()) / 1000), colors.Reset))

	if len(k.properties) == 0 {
		fmt.Println()
		return k
	}

	println(fmt.Sprintf("\n=== %sProperties%s", colors.BoldYellow, colors.Reset))

	longestKey := 0

	for key := range k.properties {
		if len(key) > longestKey {
			longestKey = len(key)
		}
	}

	for key, value := range k.properties {
		spaces := ""

		for i := 0; i < longestKey-len(key); i++ {
			spaces += " "
		}

		p := fmt.Sprintf(
			"  %s%s%s%s -> %s%v%s",
			colors.BoldWhite,
			key,
			colors.Reset,
			spaces,
			colors.Cyan,
			value,
			colors.Reset,
		)
		fmt.Println(p)
	}

	fmt.Println()

	return k
}

// Returns a JSON string representation of the error.
// This is useful to log or store the error.
// The handledAt and trace will be generated at the time of calling this method.
func (k *Khata) ToJSON() string {
	trace := k.Trace()
	explanations := k.Explanations()

	traceMap := make([]map[string]interface{}, len(trace))
	explanationsMap := make([]map[string]interface{}, len(explanations))

	for i, t := range trace {
		traceMap[i] = map[string]interface{}{
			"file":         t.file,
			"line":         t.line,
			"functionName": t.functionName,
		}
	}

	for i, e := range explanations {
		explanationsMap[i] = map[string]interface{}{
			"file":         e.file,
			"line":         e.line,
			"functionName": e.functionName,
			"message":      e.message,
		}
	}

	jsonStr, err := json.Marshal(map[string]interface{}{
		"trace":        traceMap,
		"explanations": explanationsMap,
		"error":        k.Err.Error(),
		"errorType":    k.errorType,
		"errorCode":    k.errorCode,
		"exitCode":     k.exitCode,
		"createdAt":    k.createdAt.Format("2006-01-02T15:04:05.000Z-0700"),
		"properties":   k.properties,
		"handledAt":    time.Now().UTC().Format("2006-01-02T15:04:05.000Z-0700"),
	})

	if err != nil {
		return ""
	}

	return string(jsonStr)
}

// Create the default khata error type
func Wrap(err error) *Khata {
	return &Khata{
		Err:              err,
		createdAt:        time.Now().UTC(),
		errorCode:        DEFAULT_ERROR_CODE,
		errorType:        DEFAULT_ERROR_TYPE,
		exitCode:         DEFAULT_EXIT_CODE,
		explanationStack: []KhataExplanation{},
		properties:       map[string]interface{}{},
		traceStack:       []KhataTrace{},
		template:         nil,
	}
}

func New(message string) *Khata {
	return Wrap(errors.New(message))
}

// KhataTemplate

// Create a new KhataTemplate with the given error type as an optional argument. Expects a string.
func NewTemplate() *KhataTemplate {
	return &KhataTemplate{
		errorCode:  DEFAULT_ERROR_CODE,
		errorType:  DEFAULT_ERROR_TYPE,
		exitCode:   DEFAULT_EXIT_CODE,
		properties: map[string]interface{}{},
		message:    DEFAULT_MESSAGE,
	}
}

// The default error handler for Khata errors.
// It will print the debugging information. Will exit the program if the error is fatal.
func HandleKhata(khataError Khata) {
	khataError.Debug()

	if khataError.IsFatal() {
		os.Exit(khataError.exitCode)
	}
}

// Private Functions

func tryTrimmingFunc(funcName string) string {
	prefix := os.Getenv("KHATA_FUNC_TRUNC_PREFIX")

	if prefix == "" {
		return funcName
	}

	if strings.HasPrefix(funcName, prefix) {
		return funcName[len(prefix):]
	}

	return funcName
}

func tryTrimmingPath(filePath string) string {
	prefix := os.Getenv("KHATA_PATH_TRUNC_PREFIX")

	if prefix == "" {
		return filePath
	}

	if strings.HasPrefix(filePath, prefix) {
		return filePath[len(prefix):]
	}

	return filePath
}

func collectCallerTrace() KhataTrace {
	var pc [128]uintptr
	depth := runtime.Callers(3, pc[:])
	frames := runtime.CallersFrames(pc[:depth])
	frame, _ := frames.Next()

	return KhataTrace{
		file:         frame.File,
		line:         frame.Line,
		functionName: frame.Function,
	}
}

func collectTrace() []KhataTrace {
	const packagePrefix = "github.com/cmseguin/khata."

	var pc [128]uintptr
	depth := runtime.Callers(1, pc[:])
	frames := runtime.CallersFrames(pc[:depth])

	var trace []KhataTrace
	for {
		frame, more := frames.Next()

		if !more {
			break
		}

		if len(frame.Function) > len(packagePrefix) && packagePrefix == frame.Function[:len(packagePrefix)] {
			continue
		}

		trace = append(trace, KhataTrace{
			file:         frame.File,
			line:         frame.Line,
			functionName: frame.Function,
		})
	}
	return trace
}
