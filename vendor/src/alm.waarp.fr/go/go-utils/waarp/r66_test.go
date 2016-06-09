package waarp

import (
	"testing"
)

func TestRuleMode(t *testing.T) {
	if RECV.String() != "receive" {
		t.Error("Name of rule mode cannot be found")
	}

	if RuleModeByName("receive") != RECV {
		t.Error("Cannot get the right mode by name")
	}

	if RuleModeByName("does not exist") != UNKNOWN_MODE {
		t.Error("RuleModeByName does accepts unknown modes")
	}

	v, err := SEND.Value()
	if err != nil {
		t.Fatalf("Cannot get Value() of a rulemode: %s", err.Error())
	}
	rv, ok := v.(string)
	if !ok {
		t.Fatalf("Expected a string representation from Value(). Got %T (%#v)", v, v)
	}
	if rv != "send" {
		t.Fatalf("Value() returned a wrong string representation. expected 'send', got '%s'", rv)
	}
	{
		t.Log("Testing SQL scanning")
		var r RuleMode
		srcList := map[RuleMode]interface{}{
			SENDTHROUGH: "sendthrough",
			RECVTHROUGH: []byte("receivethrough"),
		}
		for expected, src := range srcList {
			err := r.Scan(src)
			if err != nil {
				t.Errorf("Cannot Scan() '%v': %s", src, err.Error())
			}
			if r != expected {
				t.Errorf("Scan returned the wrong value. expected %v, got %v", expected, r)
			}
		}

		var badSrc interface{} = 3
		err = r.Scan(badSrc)
		if err == nil {
			t.Errorf("Scan did not returned an error for bad values")
		}
	}

	{
		t.Log("Testing json unmarshalling")
		var r RuleMode
		srcList := map[string]RuleMode{
			"1":        SEND,
			"receive":  RECV,
			"send+md5": SENDMD5,
			"foobar":   UNKNOWN_MODE,
			"1048":     UNKNOWN_MODE,
			"42":       UNKNOWN_MODE,
			"8":        RECVMD5THROUGH,
			`"receivethrough+md5"`: RECVMD5THROUGH,
			`"foobar"`:             UNKNOWN_MODE,
			`"42"`:                 UNKNOWN_MODE,
		}
		for src, expected := range srcList {
			err := r.UnmarshalJSON([]byte(src))
			if err != nil {
				t.Errorf("Cannot Scan() '%v': %s", src, err.Error())
			}
			if r != expected {
				t.Errorf("Scan returned the wrong value. expected %v, got %v", expected, r)
			}
		}
	}
}
