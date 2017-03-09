package comfyconf

//OptionKey key pair for indicating flag configuration
type OptionKey struct {
	shortName string
	fullName  string
}

//GetShort get short option definition name
func (ok *OptionKey) GetShort() string {
	return ok.shortName
}

//GetFull get full option definition name
func (ok *OptionKey) GetFull() string {
	return ok.fullName
}

//Option structure that used for holding information about flags
type Option struct {
	defaultValue interface{}
	variable     interface{}
	optionType   OptionType
	description  string
}

//GetDescription returns option description
func (o *Option) GetDescription() string {
	return o.description
}

//GetOptionType returns option type
func (o *Option) GetOptionType() OptionType {
	return o.optionType
}

//GetDefaultValue returns option default value
func (o *Option) GetDefaultValue() interface{} {
	return o.defaultValue
}

//Put binds value to variable
func (o *Option) Put(value interface{}) {
	if o.isOptionType(value) {
		switch value.(type) {
		case string:
			*o.variable.(*string) = value.(string)
		case int:
			*o.variable.(*int) = value.(int)
		case bool:
			*o.variable.(*bool) = value.(bool)
		case []interface{}:
			*o.variable.(*[]interface{}) = value.([]interface{})
		}
	}
}

func (o *Option) isOptionType(value interface{}) bool {
	switch value.(type) {
	case string:
		return o.optionType == stringType
	case int:
		return o.optionType == intType
	case bool:
		return o.optionType == boolType || o.optionType == existenceType
	case []interface{}:
		return o.optionType == sliceType
	}

	return false
}

//OptionType type of variable
type OptionType int

const (
	intType OptionType = iota + 1
	stringType
	boolType
	existenceType
	sliceType
)
