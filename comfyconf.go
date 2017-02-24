package comfyconf

import (
	"fmt"
	"reflect"
	"bytes"
)

const tagName string = "comfyname"

//New returns pointer to new Conf instance with provided middleware
func New(middleware ...Middleware) *Conf {
	return &Conf{
		options:    make(map[OptionKey]*Option),
		middleware: middleware,
	}
}

//Conf struct that holds information and provides API for parsing configuration from different middleware.
//If you familiar with flags core library, then this API may look for you similar
type Conf struct {
	options    map[OptionKey]*Option
	middleware []Middleware
}

//Parse initializes middlewares and populates all arguments with parsed data
func (c *Conf) Parse() (err error) {
	err = c.prepare()

	if err != nil {
		return
	}

	for optKey, opt := range c.options {
		for _, m := range c.middleware {
			var isOk bool
			var r interface{}

			switch opt.optionType {
			case stringType:
				r, isOk = m.ParseString(optKey.shortName, optKey.fullName)
			case boolType:
				r, isOk = m.ParseBool(optKey.shortName, optKey.fullName)
			case intType:
				r, isOk = m.ParseInt(optKey.shortName, optKey.fullName)
			case existenceType:
				r, isOk = m.ParseExistence(optKey.shortName, optKey.fullName)
			case sliceType:
				r, isOk = m.ParseSlice(optKey.shortName, optKey.fullName)
			}

			if !isOk {
				continue
			}

			opt.Put(r)
		}
	}

	return nil
}

//ToStruct populates parsed data to pointed struct by comfyname
func (c *Conf) ToStruct(structure interface{}) {

	rt := reflect.TypeOf(structure)
	rv := reflect.ValueOf(structure)

	if rt.Kind() != reflect.Ptr {
		return
	}

	v := rt.Elem()

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)

		if f.Type.Kind() == reflect.Struct {

			v1 := rv.Elem().Field(i)

			if v1.CanAddr() {
				c.ToStruct(v1.Addr().Interface())
			}

			continue
		}

		name, isOk := f.Tag.Lookup(tagName)

		if !isOk {
			continue
		}

		for optKey, opt := range c.options {
			if c.isCorrectOpt(optKey, name) {
				kind := f.Type.Kind()
				field := rv.Elem().Field(i)

				if c.isValKindAllowed(opt.optionType, kind) {
					field.Set(c.getReflectValueOfVarInterface(opt.variable))
				} else if c.isCorrectTypePointer(kind, field, opt.variable) {
					field.Set(reflect.ValueOf(opt.variable))
				}

				break
			}
		}
	}
}

func (c *Conf) isCorrectOpt(optKey OptionKey, name string) bool {
	return optKey.fullName == name || optKey.shortName == name
}

func (c *Conf) getReflectValueOfVarInterface(variable interface{}) reflect.Value {
	return reflect.ValueOf(reflect.ValueOf(variable).Elem().Interface())
}

func (c *Conf) isCorrectTypePointer(kind reflect.Kind, field reflect.Value, variable interface{}) bool {
	return kind == reflect.Ptr && field.Type().Elem().Kind() == reflect.ValueOf(variable).Elem().Kind()
}

func (c *Conf) isValKindAllowed(optType OptionType, kind reflect.Kind) bool {
	return optType == sliceType && kind == reflect.Slice || (optType == boolType || optType == existenceType) && kind == reflect.Bool ||
		optType == intType && kind == reflect.Int || optType == stringType && kind == reflect.String
}

func (c *Conf) prepare() error {
	for _, middleware := range c.middleware {
		err := middleware.Init()

		if err != nil {
			return err
		}
	}

	return nil
}

//AddMiddleware appends additional middleware
func (c *Conf) AddMiddleware(middleware ...Middleware) {
	c.middleware = append(c.middleware, middleware...)
}

//IntVar defines selected flag fullname, shortname, default value and description, binds provided integer pointer to flag.
func (c *Conf) IntVar(shortName string, fullName string, defaultValue int, variable *int, description string) {
	*variable = defaultValue
	c.createOption(shortName, fullName, defaultValue, variable, intType, description)
}

//Int defines selected flag fullname, shortname, default value and description,
//creates and returns pointer to integer variable and binds that integer to flag.
func (c *Conf) Int(shortName string, fullName string, defaultValue int, description string) *int {
	variable := new(int)
	c.IntVar(shortName, fullName, defaultValue, variable, description)
	return variable
}

//StringVar defines selected flag fullname, shortname, default value and description, binds provided string pointer to flag.
func (c *Conf) StringVar(shortName string, fullName string, defaultValue string, variable *string, description string) {
	*variable = defaultValue
	c.createOption(shortName, fullName, defaultValue, variable, stringType, description)
}

//String defines selected flag fullname, shortname, default value and description,
//creates and returns pointer to string variable and binds that string to flag.
func (c *Conf) String(shortName string, fullName string, defaultValue string, description string) *string {
	variable := new(string)
	c.StringVar(shortName, fullName, defaultValue, variable, description)
	return variable
}

//BoolVar defines selected flag fullname, shortname, default value and description, binds provided bool pointer to flag.
func (c *Conf) BoolVar(shortName string, fullName string, defaultValue bool, variable *bool, description string) {
	*variable = defaultValue
	c.createOption(shortName, fullName, defaultValue, variable, boolType, description)
}

//Bool defines selected flag fullname, shortname, default value and description,
//creates and returns pointer to bool variable and binds that bool to flag.
func (c *Conf) Bool(shortName string, fullName string, defaultValue bool, description string) *bool {
	variable := new(bool)
	c.BoolVar(shortName, fullName, defaultValue, variable, description)
	return variable
}

//ExistVar similar to BoolVar, but without default value.
func (c *Conf) ExistVar(shortName string, fullName string, variable *bool, description string) {
	*variable = false
	c.createOption(shortName, fullName, false, variable, existenceType, description)
}

//Exist similar to Bool, but without default value.
func (c *Conf) Exist(shortName string, fullName string, description string) *bool {
	variable := new(bool)
	c.ExistVar(shortName, fullName, variable, description)
	return variable
}

//SliceVar defines selected flag fullname, shortname, default value and description, binds provided slice pointer to flag.
func (c *Conf) SliceVar(shortName string, fullName string, defaultValue []interface{}, variable *[]interface{}, description string) {
	*variable = defaultValue
	c.createOption(shortName, fullName, defaultValue, variable, sliceType, description)
}

//Slice defines selected flag fullname, shortname, default value and description,
//creates and returns pointer to slice variable and binds that slice to flag.
func (c *Conf) Slice(shortName string, fullName string, defaultValue []interface{}, description string) *[]interface{} {
	variable := new([]interface{})
	c.SliceVar(shortName, fullName, defaultValue, variable, description)
	return variable
}

//PrintHelp created for executing function that will instruction
func (c *Conf) PrintHelp(printer func(options map[OptionKey]*Option)) {
	printer(c.options)
}

func (c *Conf) createOption(shortName string, fullName string, defaultValue interface{}, variable interface{}, optionType OptionType, description string) {
	c.options[OptionKey{
		shortName,
		fullName,
	}] = &Option{
		defaultValue,
		variable,
		optionType,
		description,
	}
}

//DefaultHelpPrinter function for printing help information. Should be passed to Conf.PrintHelp
func DefaultHelpPrinter(options map[OptionKey]*Option) {
	var buffer bytes.Buffer

	buffer.WriteString("Usage: app [options] \n")
	buffer.WriteString("  Options: \n")

	for def, opt := range options {
		buffer.WriteString("    -" + def.GetShort() + ", --" + def.GetFull() + "   "+ opt.description + "\n")
	}

	fmt.Println(buffer.String())
}
