package comfyconf

import (
	"os"
	"regexp"
	"strings"
)

//NewEnv creates middleware for environment variables with default prefix
func NewEnv() *Env {
	return NewEnvWithPrefix("ENV_")
}

//NewEnvWithPrefix creates middleware for environment variables with custom prefix
func NewEnvWithPrefix(prefix string) *Env {
	env := &Env{
		prefix: prefix,
	}
	env.parsed = make(map[string]string)
	env.parsedSlice = make(map[string][]interface{})

	return env
}

//Env structure that implements middleware interface for environment variables
type Env struct {
	Flags
	prefix string
}

//Init initializing middleware for environment variables
func (f *Env) Init() error {

	arrExpr := regexp.MustCompile(`^(.+)(\[[\d+]?])$`)

	for _, envPair := range os.Environ() {
		pair := strings.Split(envPair, "=")

		k := pair[0]
		v := pair[1]

		if strings.HasPrefix(k, f.prefix) {
			k = strings.TrimPrefix(k, f.prefix)

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
	}

	return nil
}
