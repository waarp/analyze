package waarp

import (
	"strconv"
)

type Step int

const S_WM_UNKNOWN Step = -1
const (
	S_NOTASK Step = iota
	S_PRETASK
	S_TRANSFERTASK
	S_POSTTASK
	S_ALLDONETASK
	S_ERRORTASK
)

var stepNames = map[Step]string{
	S_WM_UNKNOWN:   "S_WM_UNKNOWN",
	S_NOTASK:       "NOTASK",
	S_PRETASK:      "PRETASK",
	S_TRANSFERTASK: "TRANSFERTASK",
	S_POSTTASK:     "POSTTASK",
	S_ALLDONETASK:  "ALLDONETASK",
	S_ERRORTASK:    "ERRORTASK",
}

func (s Step) String() string {
	return stepNames[s]
}

// JSON Serialization
func (s *Step) MarshalJSON() ([]byte, error) {
	return []byte(`"` + s.String() + `"`), nil
}

// JSON Deserialization
func (s *Step) UnmarshalJSON(val []byte) error {
	value, err := strconv.ParseInt(string(val), 10, 64)
	if err != nil {
		*s = S_WM_UNKNOWN
		return nil
	}
	if _, ok := stepNames[Step(value)]; ok {
		*s = Step(value)
	} else {
		*s = S_WM_UNKNOWN
	}
	return nil
}
