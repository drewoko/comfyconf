package comfyconf

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"strings"
)

//NewJSON returns pointer to instance of JSON configuration middleware
//with default file from where JSON will be read
func NewJSON(file string) *JSON {
	return &JSON{
		flags:  false,
		path:   file,
		reader: DefaultJSONReader,
	}
}

//NewJSONWithCustomReader returns pointer to instance of JSON configuration middleware
//with custom file reader
func NewJSONWithCustomReader(reader func(j *JSON) ([]byte, error)) *JSON {
	return &JSON{
		flags:  false,
		reader: reader,
	}
}

//NewJSONWithCustomFileMiddleware returns pointer to instance of JSON configuration middleware.
//Getting JSON configuration location performed by provided middlewares using short and long name of flag
func NewJSONWithCustomFileMiddleware(shortName string, fullName string, defaultFile string, middlewareList ...Middleware) *JSON {

	var file string

	for _, middleware := range middlewareList {
		err := middleware.Init()

		if err != nil {
			continue
		}

		parsedFile, isOk := middleware.ParseString(shortName, fullName)

		if isOk {
			file = parsedFile
		}
	}

	if len(file) == 0 {
		file = defaultFile
	}

	return &JSON{
		flags:  true,
		path:   file,
		reader: DefaultJSONReader,
	}
}

//JSON structure implemenents Middleware instance for parsing JSON configuration
type JSON struct {
	flags  bool
	path   string
	reader func(j *JSON) ([]byte, error)

	parsed     map[string]interface{}
	shortIndex map[string]string
}

//DefaultJSONReader default JSON file reader
func DefaultJSONReader(j *JSON) ([]byte, error) {
	return ioutil.ReadFile(j.path)
}

//Init initializing middleware for JSON configuration
func (j *JSON) Init() error {

	contentBytes, err := j.reader(j)

	if err != nil {
		return err
	}

	tmpParsed := make(map[string]interface{})

	err = json.Unmarshal(contentBytes, &tmpParsed)

	if err != nil {
		return err
	}

	j.parsed = j.parse(tmpParsed)

	j.prepareIndex()

	return nil
}

func (j *JSON) prepareIndex() {

	j.shortIndex = make(map[string]string)

	for k := range j.parsed {
		s := strings.Split(k, ".")

		if len(s) == 1 {
			j.shortIndex[s[0]] = k
			continue
		}

		j.shortIndex[s[len(s)-1]] = k
	}
}

func (j *JSON) parse(data map[string]interface{}) map[string]interface{} {

	tmp := make(map[string]interface{})

	parseOne(tmp, "", data)

	return tmp
}

func parseOne(tmp map[string]interface{}, key string, data map[string]interface{}) {

	for k, v := range data {

		if len(key) != 0 {
			k = key + "." + k
		}

		t := reflect.TypeOf(v).Kind()

		if t == reflect.Slice {
			v1, isOk := v.([]interface{})

			if !isOk {
				continue
			}

			tmp[k] = v1

			continue
		}

		if t == reflect.Map {
			v1, isOk := v.(map[string]interface{})

			if !isOk {
				continue
			}

			parseOne(tmp, k, v1)

			continue
		}

		tmp[k] = v
	}
}

func (j *JSON) get(shortName string, fullName string) (interface{}, bool) {

	var v interface{}

	v = j.parsed[fullName]

	if v != nil {
		return v, true
	}

	v = j.parsed[j.shortIndex[shortName]]

	if v != nil {
		return v, true
	}

	return "", false
}

//ParseInt tries to get int from JSON configuration
func (j *JSON) ParseInt(shortName string, fullName string) (int, bool) {

	v, isOk := j.get(shortName, fullName)

	if !isOk {
		return 0, false
	}

	switch v.(type) {
	case float32:
		v1, isOk := v.(float32)
		return int(v1), isOk
	case float64:
		v1, isOk := v.(float64)
		return int(v1), isOk
	default:
		v1, isOk := v.(int)
		return v1, isOk
	}
}

//ParseString tries to get string from JSON configuration
func (j *JSON) ParseString(shortName string, fullName string) (string, bool) {
	v, isOk := j.get(shortName, fullName)

	if !isOk {
		return "", false
	}

	v1, isOk := v.(string)
	return v1, isOk
}

//ParseBool tries to get bool from JSON configuration
func (j *JSON) ParseBool(shortName string, fullName string) (bool, bool) {
	v, isOk := j.get(shortName, fullName)

	if !isOk {
		return false, false
	}

	v1, isOk := v.(bool)
	return v1, isOk
}

//ParseExistence tries to check that element exists in JSON configuration
func (j *JSON) ParseExistence(shortName string, fullName string) (bool, bool) {

	if j.parsed[fullName] != nil {
		return true, true
	}

	_, isOk := j.shortIndex[shortName]

	return isOk, true
}

//ParseSlice tries to get slice from JSON configuration
func (j *JSON) ParseSlice(shortName string, fullName string) ([]interface{}, bool) {

	v, isOk := j.get(shortName, fullName)

	if !isOk {
		return make([]interface{}, 0), false
	}

	v1, isOk := v.([]interface{})

	return v1, isOk
}
