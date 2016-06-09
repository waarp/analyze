package waarp

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

//----------------------------------------------------------------------
// TransferStatus
//----------------------------------------------------------------------

// Transferstatus lists the several states a transfer can have in its lifetime.
type TransferStatus uint8

const (
	TS_UNKNOWN TransferStatus = iota
	TS_NOTUPDATED
	TS_INTERRUPTED
	TS_TOSUBMIT
	TS_INERROR
	TS_RUNNING
	TS_DONE

	// Additional statuses
	TS_CANCELED TransferStatus = 98
	TS_STOPPED  TransferStatus = 99
)

var transferStatusNames = map[TransferStatus]string{
	TS_UNKNOWN:     "unknown",
	TS_NOTUPDATED:  "not updated",
	TS_INTERRUPTED: "interrupted",
	TS_TOSUBMIT:    "to submit",
	TS_INERROR:     "error",
	TS_RUNNING:     "running",
	TS_DONE:        "done",
	TS_CANCELED:    "canceled",
	TS_STOPPED:     "stopped",
}

// Implements Stringer()
func (ts TransferStatus) String() string {
	if rv, ok := transferStatusNames[ts]; ok {
		return rv
	}
	return transferStatusNames[TS_UNKNOWN]
}

func TransferStatusByName(name string) TransferStatus {
	for item, str := range transferStatusNames {
		if str == name {
			return item
		}
	}
	return TS_UNKNOWN
}

// JSON Serialization
func (t *TransferStatus) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.String() + `"`), nil
}

// JSON Deserialization
func (t *TransferStatus) UnmarshalJSON(val []byte) error {
	val2 := strings.Trim(string(val), `"'`)
	if v, err := strconv.ParseUint(val2, 10, 8); err == nil {
		*t = TransferStatus(uint8(v))
	} else {
		*t = TransferStatusByName(val2)
	}
	return nil
}

// CSV Serialization
func (t TransferStatus) MarshalCSV() []string {
	return []string{t.String()}
}

// *-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-
//								RULE MODE
// *-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-

type RuleMode uint8

const (
	UNKNOWN_MODE RuleMode = 0
	SEND         RuleMode = iota
	RECV
	SENDMD5
	RECVMD5
	SENDTHROUGH
	RECVTHROUGH
	SENDMD5THROUGH
	RECVMD5THROUGH

	RM_NOOP RuleMode = 99
)

var ruleModeNames = map[RuleMode]string{
	UNKNOWN_MODE:   "unknown",
	SEND:           "send",
	RECV:           "receive",
	SENDMD5:        "send+md5",
	RECVMD5:        "receive+md5",
	SENDTHROUGH:    "sendthrough",
	RECVTHROUGH:    "receivethrough",
	SENDMD5THROUGH: "sendthrough+md5",
	RECVMD5THROUGH: "receivethrough+md5",
	RM_NOOP:        "noop",
}

func RuleModeByName(s string) RuleMode {
	for m, n := range ruleModeNames {
		if n == s {
			return m
		}
	}
	return UNKNOWN_MODE
}

func (r RuleMode) String() string {
	return ruleModeNames[r]
}

// SQL serialization
func (r RuleMode) Value() (driver.Value, error) {
	return r.String(), nil
}

// SQL deserialization
func (r *RuleMode) Scan(src interface{}) error {
	switch rv := src.(type) {
	case string:
		*r = RuleModeByName(rv)
	case []byte:
		*r = RuleModeByName(string(rv))
	default:
		return fmt.Errorf("Cannot convert '%v' to RuleMode", src)
	}
	return nil
}

// JSON Serialization
func (t *RuleMode) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.String() + `"`), nil
}

// JSON Deserialization
func (t *RuleMode) UnmarshalJSON(val []byte) error {
	val2 := strings.Trim(string(val), `"`)
	if n, err := strconv.ParseUint(val2, 10, 8); err == nil && n <= uint64(len(ruleModeNames)-1) {
		*t = RuleMode(uint8(n))
		return nil
	}
	*t = RuleModeByName(val2)
	return nil
}

// CSV Serialization
func (t RuleMode) MarshalCSV() []string {
	return []string{t.String()}
}
