package comfyconf

import (
	"os"
	"regexp"
	"strconv"
	"strings"
)

//NewFlags creates new flags middleware with default assignment
func NewFlags() *Flags {
	return NewFlagsWithCustomAssignment("=")
}

//NewFlagsWithCustomAssignment creates new flags middleware with custom assignment
func NewFlagsWithCustomAssignment(assignment string) *Flags {
	return NewFlagsWithCustomParser(func(arg string) (string, string) {
		return DefaultFlagsParser(arg, assignment)
	})
}

//NewFlagsWithCustomParser creates new flags middleware with custom parser
func NewFlagsWithCustomParser(parser func(arg string) (string, string)) *Flags {
	return &Flags{
		os.Args,
		parser,
		make(map[string]string),
		make(map[string][]interface{}),
	}
}

//DefaultFlagsParser default parser for flags middleware
func DefaultFlagsParser(arg string, assignment string) (string, string) {
	if strings.Contains(arg, "-") {
		var prefix string
		if strings.Contains(arg, "--") {
			prefix = "--"
		} else {
			prefix = "-"
		}

		kv := strings.SplitN(strings.TrimPrefix(arg, prefix), assignment, 2)

		if len(kv) == 0 {
			return "", ""
		}

		if len(kv) == 1 {
			return kv[0], "true"
		}

		return kv[0], kv[1]
	}

	return "", ""
}

//Flags struct that implements Middleware interface for parsing program arguments
type Flags struct {
	args   []string
	parser func(arg string) (string, string)

	parsed      map[string]string
	parsedSlice map[string][]interface{}
}

//Init initializing middleware for program arguments
func (f *Flags) Init() error {

	if f.parsed == nil {
		f.parsed = make(map[string]string)
	}

	arrExpr := regexp.MustCompile(`^(.+)(\[[\d+]?])$`)

	for _, arg := range f.args {
		k, v := f.parser(arg)
		if k == "" {
			continue
		}

		if arrExpr.MatchString(k) {
			match := arrExpr.FindStringSubmatch(k)
			if len(match) != 3 {
				continue
			}
			k = match[1]

			vSlice, isExist := f.parsedSlice[k]

			if isExist {
				f.parsedSlice[k] = append(vSlice, v)
			} else {
				f.parsedSlice[k] = append(make([]interface{}, 0), v)
			}

			continue
		}

		f.parsed[k] = v
	}

	return nil
}

func (f *Flags) get(shortName string, fullName string) (string, bool) {

	if len(f.parsed[fullName]) == 0 && len(f.parsed[shortName]) == 0 {
		return "", false
	}

	if len(f.parsed[fullName]) != 0 {
		v1, isOk := f.parsed[fullName]
		return v1, isOk
	}

	if len(f.parsed[shortName]) != 0 {
		v1, isOk := f.parsed[shortName]
		return v1, isOk
	}

	return "", false
}

//ParseInt tries to get int from flags middleware
func (f *Flags) ParseInt(shortName string, fullName string) (int, bool) {

	v, isOk := f.get(shortName, fullName)
	if !isOk {
		return 0, false
	}

	vi, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}

	return vi, true
}

//ParseString tries to get string from flags middleware
func (f *Flags) ParseString(shortName string, fullName string) (string, bool) {
	return f.get(shortName, fullName)
}

//ParseBool tries to get bool from flags middleware
func (f *Flags) ParseBool(shortName string, fullName string) (bool, bool) {
	v, isOk := f.get(shortName, fullName)

	if !isOk {
		return false, false
	}

	v1, err := strconv.ParseBool(v)
	if err != nil {
		return false, false
	}
	return v1, true
}

//ParseExistence tries to check that flag exists
func (f *Flags) ParseExistence(shortName string, fullName string) (bool, bool) {
	_, isOk := f.ParseBool(shortName, fullName)

	if isOk {
		return true, true
	}

	return false, true
}

func (f *Flags) getSlice(name string) ([]interface{}, bool) {

	v, isOk := f.parsedSlice[name]

	if isOk {
		return v, true
	}

	return nil, false
}

//ParseSlice tries to get slice from flags middleware
func (f *Flags) ParseSlice(shortName string, fullName string) ([]interface{}, bool) {

	v, isOk := f.getSlice(shortName)
	if isOk {
		return v, true
	}

	v, isOk = f.getSlice(fullName)
	if isOk {
		return v, true
	}

	return v, false
}
