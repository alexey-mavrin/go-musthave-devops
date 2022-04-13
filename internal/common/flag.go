package common

import (
	"fmt"
	"time"
)

// StringFlag holds string flag
type StringFlag struct {
	Value  *string
	Option string
	Set    bool
}

// BoolFlag holds bool flags
type BoolFlag struct {
	Value  *bool
	Option string
	Set    bool
}

// TimeFlag holds time flags
type TimeFlag struct {
	Value  *time.Duration
	Option string
	Set    bool
}

func (s StringFlag) String() string {
	var a = AnyFlag{
		value:  s.Value,
		option: s.Option,
		set:    s.Set,
	}
	return a.String()
}

func (t TimeFlag) String() string {
	var a = AnyFlag{
		value:  t.Value,
		option: t.Option,
		set:    t.Set,
	}
	return a.String()
}

func (b BoolFlag) String() string {
	var a = AnyFlag{
		value:  b.Value,
		option: b.Option,
		set:    b.Set,
	}
	return a.String()
}

// AnyFlag is used to convert flag values into a string
type AnyFlag struct {
	value  interface{}
	option string
	set    bool
}

func (s AnyFlag) String() string {
	var ret = "{option: " + s.option + ", "

	if s.value != nil {
		switch v := s.value.(type) {
		case *string:
			ret += fmt.Sprintf("value: %s", *v)
		case *bool:
			ret += fmt.Sprintf("value: %t", *v)
		case *time.Duration:
			ret += fmt.Sprintf("value: %v", *v)
		}
	} else {
		ret += "(nil)"
	}

	if s.set {
		ret += ", set explicitly}"
	} else {
		ret += ", not set explicitly}"
	}
	return ret
}
