// Code generated by "enumer -type RSAHasher -trimprefix RSAHasher -text -json -linecomment"; DO NOT EDIT.

package crypto

import (
	"encoding/json"
	"fmt"
	"strings"
)

const _RSAHasherName = "SHA256SHA384SHA512"

var _RSAHasherIndex = [...]uint8{0, 0, 6, 12, 18}

const _RSAHasherLowerName = "sha256sha384sha512"

func (i RSAHasher) String() string {
	if i < 0 || i >= RSAHasher(len(_RSAHasherIndex)-1) {
		return fmt.Sprintf("RSAHasher(%d)", i)
	}
	return _RSAHasherName[_RSAHasherIndex[i]:_RSAHasherIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _RSAHasherNoOp() {
	var x [1]struct{}
	_ = x[RSAHasherUnspecified-(0)]
	_ = x[RSAHasherSHA256-(1)]
	_ = x[RSAHasherSHA384-(2)]
	_ = x[RSAHasherSHA512-(3)]
}

var _RSAHasherValues = []RSAHasher{RSAHasherUnspecified, RSAHasherSHA256, RSAHasherSHA384, RSAHasherSHA512}

var _RSAHasherNameToValueMap = map[string]RSAHasher{
	_RSAHasherName[0:0]:        RSAHasherUnspecified,
	_RSAHasherLowerName[0:0]:   RSAHasherUnspecified,
	_RSAHasherName[0:6]:        RSAHasherSHA256,
	_RSAHasherLowerName[0:6]:   RSAHasherSHA256,
	_RSAHasherName[6:12]:       RSAHasherSHA384,
	_RSAHasherLowerName[6:12]:  RSAHasherSHA384,
	_RSAHasherName[12:18]:      RSAHasherSHA512,
	_RSAHasherLowerName[12:18]: RSAHasherSHA512,
}

var _RSAHasherNames = []string{
	_RSAHasherName[0:0],
	_RSAHasherName[0:6],
	_RSAHasherName[6:12],
	_RSAHasherName[12:18],
}

// RSAHasherString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func RSAHasherString(s string) (RSAHasher, error) {
	if val, ok := _RSAHasherNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _RSAHasherNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to RSAHasher values", s)
}

// RSAHasherValues returns all values of the enum
func RSAHasherValues() []RSAHasher {
	return _RSAHasherValues
}

// RSAHasherStrings returns a slice of all String values of the enum
func RSAHasherStrings() []string {
	strs := make([]string, len(_RSAHasherNames))
	copy(strs, _RSAHasherNames)
	return strs
}

// IsARSAHasher returns "true" if the value is listed in the enum definition. "false" otherwise
func (i RSAHasher) IsARSAHasher() bool {
	for _, v := range _RSAHasherValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for RSAHasher
func (i RSAHasher) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for RSAHasher
func (i *RSAHasher) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("RSAHasher should be a string, got %s", data)
	}

	var err error
	*i, err = RSAHasherString(s)
	return err
}

// MarshalText implements the encoding.TextMarshaler interface for RSAHasher
func (i RSAHasher) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for RSAHasher
func (i *RSAHasher) UnmarshalText(text []byte) error {
	var err error
	*i, err = RSAHasherString(string(text))
	return err
}