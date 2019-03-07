package cmdutil

import (
	"fmt"

	"github.com/spf13/pflag"
)

type StringEnumValue struct {
	Choices []string

	value string
}

func (v *StringEnumValue) Set(arg string) error {
	var ok bool
	for _, choice := range v.Choices {
		if choice == arg {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("%q is not valid (choices: %q)", arg, v.Choices)
	}
	v.value = arg
	return nil
}

func (v *StringEnumValue) Type() string {
	return "CHOICE"
}

func (v *StringEnumValue) String() string {
	if v.value != "" {
		return v.value
	} else if len(v.Choices) > 0 {
		return v.Choices[0]
	} else {
		return ""
	}
}

var _ pflag.Value = (*StringEnumValue)(nil)
