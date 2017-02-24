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
	switch value.(type) {
	case string:
		o.putString(value)
	case int:
		o.putInt(value)
	case bool:
		o.putBool(value)
	case []interface{}:
		o.putSlice(value)
	}
}

func (o *Option) putString(value interface{}) {
	if o.optionType != stringType {
		return
	}
	z, ok := o.variable.(*string)
	if !ok {
		return
	}

	*z = value.(string)
}

func (o *Option) putInt(value interface{}) {
	if o.optionType != intType {
		return
	}
	z, ok := o.variable.(*int)
	if !ok {
		return
	}
	*z = value.(int)
}

func (o *Option) putBool(value interface{}) {
	if o.optionType != boolType && o.optionType != existenceType {
		return
	}
	z, ok := o.variable.(*bool)
	if !ok {
		return
	}

	*z = value.(bool)
}

func (o *Option) putSlice(value interface{}) {
	if o.optionType != sliceType {
		return
	}
	z, ok := o.variable.(*[]interface{})
	if !ok {
		return
	}

	*z = value.([]interface{})
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
