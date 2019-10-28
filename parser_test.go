package expr

import (
	"testing"
)

func TestParse(t *testing.T) {
	t.Run("BinaryBool", func(t *testing.T) { RunParseBinaryBoolTest(t) })
	t.Run("ConcatString", func(t *testing.T) { RunParseStringTest(t) })
}

func RunParseBinaryBoolTest(t *testing.T) {
	RunExprTest("true == true", true, t)
	RunExprTest("true == false", false, t)
	RunExprTest("false == true", false, t)
	RunExprTest("true != true", false, t)
	RunExprTest("false != true", true, t)
	RunExprTest("true != false", true, t)
	RunExprTest("true && true", true, t)
	RunExprTest("false && true", false, t)
	RunExprTest("true || false", true, t)
	RunExprTest("true || true", true, t)
	RunExprTest("false || false", false, t)
	//With sub expressions
	RunExprTest("true && (true && true)", true, t)
	RunExprTest("true && (false && true)", false, t)
	RunExprTest("(false || true) && (false || true)", true, t)
}

func RunParseStringTest(t *testing.T) {
	RunExprTest("\"expr\" + \" \" + \"rules!\"", "expr rules!", t)
	//TODO: This failes. Write lower level tests on lexer
	//In Evaluates to bool
	RunExprTest("\"lunch\" in [\"breakfast\", \"lunch\", \"dinner\", \"supper\"]", true, t)
}

func RunStringInListTest(t *testing.T) {
	RunExprTest("\"expr\" + \" \" + \"rules!\"", "expr rules!", t)
}

func RunExprTest(testString string, expectedValue interface{}, t *testing.T) {
	t.Run(testString, func(t *testing.T) { exprTest(testString, expectedValue, t) })
}

func exprTest(testString string, expectedValue interface{}, t *testing.T) {
	//Step 1: Create environment
	env := NewEnvironment()
	//Step 2: optionally extend environment
	//
	//Step 3: Get parser
	parser := env.GetParser()

	//Step 4: Parse the input
	ex, err := parser.Parse(testString)
	if err != nil {
		t.Errorf("Parse failed: %v. \n\n Error: %s\n", ex, err.Error())
		return

	}
	//Step 5: Evaluate the parsed text
	exEv, errEv := ex.Evaluate(NewEnvironment())
	if errEv != nil {
		t.Errorf("Eval failed:\n\n %v \n\n Evaluated to:\n\n %v \n\n Error: %s\n", ex, exEv, errEv.Error())
		return

	}
	//Step 6: Check result
	if expectedValue != exEv.Value() {
		t.Errorf("Compare failed:\n\n %v \n\n Evaluated to:\n\n %v \n", ex, exEv.Value())
		return
	}
}
