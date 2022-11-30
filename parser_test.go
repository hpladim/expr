package expr

import (
	"testing"
)

func TestParse(t *testing.T) {
	t.Run("BinaryBool", func(t *testing.T) { RunParseBinaryBoolTest(t) })
	t.Run("ConcatString", func(t *testing.T) { RunParseStringTest(t) })
}

func RunParseBinaryBoolTest(t *testing.T) {
	RunExprTest(t, "true == true", true)
	RunExprTest(t, "true == false", false)
	RunExprTest(t, "false == true", false)
	RunExprTest(t, "true != true", false)
	RunExprTest(t, "false != true", true)
	RunExprTest(t, "true != false", true)
	RunExprTest(t, "true && true", true)
	RunExprTest(t, "false && true", false)
	RunExprTest(t, "true || false", true)
	RunExprTest(t, "true || true", true)
	RunExprTest(t, "false || false", false)
	//With sub expressions
	RunExprTest(t, "true && (true && true)", true)
	RunExprTest(t, "true && (false && true)", false)
	RunExprTest(t, "(false || true) && (false || true)", true)
	//Conditional statements
	RunExprTest(t, "false == true?true:false", false)
	RunExprTest(t, "1 == 2?true:false", false)
	//Conditional statement with sub expression
	RunExprTest(t, "1 == 1?(true == true?true:false):false", true)
}

func RunParseStringTest(t *testing.T) {
	//String concat test
	RunExprTest(t, "\"expr\" + \" \" + \"rules!\"", "expr rules!")
	//In Evaluates to bool
	RunExprTest(t, "\"lunch\" in [\"breakfast\", \"lunch\", \"dinner\", \"supper\"]", true)
}

// Runs tests in a separate goroutine. Enables paralell testing.
func RunExprTest(t *testing.T, testString string, expectedValue interface{}) {
	t.Run(testString, func(t *testing.T) { exprTest(t, testString, expectedValue) })
}

func exprTest(t *testing.T, testString string, expectedValue interface{}) {
	//Step 1: Create environment
	env := NewEnvironment()
	//Step 2: Get parser
	parser := env.GetParser()
	//Step 3: Parse the input
	ex, err := parser.Parse(testString)
	if err != nil {
		t.Errorf("Parse failed: %v. \n\n Error: %s\n", ex, err.Error())
		return
	}
	//Step 4: Evaluate the parsed text
	exEv, errEv := ex.Evaluate(NewEnvironment())
	if errEv != nil {
		t.Errorf("Eval failed:\n\n %v \n\n Evaluated to:\n %v \n Error: %s\n", ex.Literal(), exEv, errEv.Error())
		return
	}
	//Step 5: Check result
	if expectedValue != exEv.Value() {
		t.Errorf("Compare failed:\n\n %v \n\n Evaluated to:\n\n %v \n", ex, exEv.Value())
		return
	}
}
