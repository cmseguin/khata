package khata_test

import (
	"errors"
	"os"
	"testing"

	"github.com/cmseguin/khata"
)

func TestNewKhata(t *testing.T) {
	errorMessage := "This is an error message"
	k := khata.New(errorMessage)

	if k.Error() != errorMessage {
		t.Error("New() did not set the error message")
	}
}

func TestWrapKhata(t *testing.T) {
	error := errors.New("This is an error message")
	k := khata.Wrap(error)

	if k.Error() != error.Error() {
		t.Error("Wrap() did not set the error message")
	}
}

func subFnExplain(err *khata.Khata) *khata.Khata {
	err.Explain("This is an explanation for subFnExplain")
	return err
}
func TestKhataExplain(t *testing.T) {
	k := khata.New("This is an error message")

	k = subFnExplain(k)

	if len(k.Explanations()) != 1 {
		t.Error("Explain() did not set the explanation")
		return
	}

	if k.Explanations()[0].FunctionName != "github.com/cmseguin/khata_test.subFnExplain" {
		t.Error("Explain() did not set the function name")
		return
	}

	if k.Explanations()[0].Message != "This is an explanation for subFnExplain" {
		t.Error("Explain() did not set the explanation message")
		return
	}
}

func TestKhataExplainf(t *testing.T) {
	k := khata.New("This is an error message")

	k.Explainf("%s is an explanation", "This")

	if len(k.Explanations()) != 1 {
		t.Error("Explainf() did not set the explanation")
		return
	}
	println(k.Explanations()[0].FunctionName)
	if k.Explanations()[0].FunctionName != "github.com/cmseguin/khata_test.TestKhataExplainf" {
		t.Error("Explainf() did not set the function name")
		return
	}

	if k.Explanations()[0].Message != "This is an explanation" {
		t.Error("Explainf() did not set the explanation message")
		return
	}
}

func TestKhataTemplate(t *testing.T) {
	template := khata.NewTemplate().SetType("MyNewTemplate")

	if template.Type() != "MyNewTemplate" {
		t.Error("NewTemplate() did not set the template name")
		return
	}

	template.SetType("MyNewTemplate2")

	if template.Type() != "MyNewTemplate2" {
		t.Error("SetType() did not set the template name")
		return
	}

	if template.ExitCode() != 1 {
		t.Error("Default exit code was not set to 1")
		return
	}

	template.SetExitCode(2)

	if template.ExitCode() != 2 {
		t.Error("SetExitCode() did not set the exit code")
		return
	}

	if template.Code() != -1 {
		t.Error("Default code was not set to -1")
		return
	}

	template.SetCode(1)

	if template.Code() != 1 {
		t.Error("SetCode() did not set the code")
		return
	}

	if template.HasProperty("test") {
		t.Error("HasProperty() returned true for a property that does not exist")
		return
	}

	template.SetProperty("test", "testValue")

	if !template.HasProperty("test") {
		t.Error("HasProperty() returned false for a property that exists")
		return
	}

	if template.GetProperty("test") != "testValue" {
		t.Error("Property() did not return the correct value")
		return
	}

	template.RemoveProperty("test")

	if template.HasProperty("test") {
		t.Error("HasProperty() returned true for a property that was removed")
		return
	}

	if template.GetProperty("test") != nil {
		t.Error("Property() did not return nil for a property that was removed")
		return
	}

	template.SetProperty("test2", "testValue2")

	template2 := template.Extend()

	if template2.Type() != template.Type() {
		t.Error("Extend() did not set the template name")
		return
	}

	if template2.ExitCode() != template.ExitCode() {
		t.Error("Extend() did not set the exit code")
		return
	}

	if template2.Code() != template.Code() {
		t.Error("Extend() did not set the code")
		return
	}

	if template2.HasProperty("test") {
		t.Error("Extend() did not remove the property")
		return
	}

	if template2.GetProperty("test") != nil {
		t.Error("Extend() did not return nil for a property that was removed")
		return
	}

	if template2.GetProperty("test2") != "testValue2" {
		t.Error("Extend() did not return the correct value for a property that was not removed")
		return
	}
}

func TestKhataErrorFromTemplate(t *testing.T) {
	template := khata.NewTemplate().SetType("MyNewTemplate")
	template2 := khata.NewTemplate().SetType(template.Type())

	template.SetCode(404)
	template.SetExitCode(2)
	template.SetProperty("test", "testValue")
	template.SetProperty("test2", "testValue2")
	template.SetType("MyType")

	template.RemoveProperty("test2")

	k := template.New()

	if k.Error() != "error" {
		t.Error("Error() did not set the default error message")
		return
	}

	if k.ExitCode() != 2 {
		t.Error("Error() did not set the exit code")
		return
	}

	if k.Code() != 404 {
		t.Error("Error() did not set the code")
		return
	}

	if k.Type() != "MyType" {
		t.Error("Error() did not set the type")
		return
	}

	if !k.HasProperty("test") {
		t.Error("Error() did not set the property")
		return
	}

	if k.GetProperty("test") != "testValue" {
		t.Error("Error() did not set the property value")
		return
	}

	if k.IsInstanceOf(template2) {
		t.Error("IsInstanceOf() did not return true for the incorrect template")
		return
	}

	if !k.IsInstanceOf(template) {
		t.Error("IsInstanceOf() did not return false for the correct template")
		return
	}
}

func TestKhataErrorChaining(t *testing.T) {
	k := khata.
		New("This is an error message").
		Explain("This is an explanation").
		SetCode(404).
		SetExitCode(2).
		SetType("MyType")

	if k.Error() != "This is an error message" {
		t.Error("Error() did not set the error message")
		return
	}

	if k.ExitCode() != 2 {
		t.Error("Error() did not set the exit code")
		return
	}

	if k.Code() != 404 {
		t.Error("Error() did not set the code")
		return
	}

	if k.Type() != "MyType" {
		t.Error("Error() did not set the type")
		return
	}

	if k.Explanations()[0].Message != "This is an explanation" {
		t.Error("Explain() did not set the explanation message")
		return
	}
}

func TestDebugOutput(t *testing.T) {
	os.Setenv("KHATA_FUNC_TRUNC_PREFIX", "github.com/cmseguin/")
	os.Setenv("KHATA_PATH_TRUNC_PREFIX", "/Users/cmseguin/dev/git/khata/")

	k := khata.
		New("Not Found").
		Explain("This is an explanation of not found").
		Explain("This is an other explanation of not found").
		SetCode(404).
		SetExitCode(1).
		SetType("HTTP").
		SetProperty("test", "testValue").
		SetProperty("test2", "testValue2")

	k.Debug()
}
