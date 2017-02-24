package comfyconf

//Middleware interface for different configuration parsing. Can be used for external configuration parsers
type Middleware interface {
	//Initialization for Middleware
	Init() error

	//ParseInt tries to get int from Middleware by flag name and returns int and fetching status
	ParseInt(shortName string, fullName string) (int, bool)
	//ParseString tries to get string from Middleware by flag name and returns string and fetching status
	ParseString(shortName string, fullName string) (string, bool)
	//ParseBool tries to get bool from Middleware by flag name and returns bool and fetching status
	ParseBool(shortName string, fullName string) (bool, bool)
	//ParseExistence tries to check that flag exists
	ParseExistence(shortName string, fullName string) (bool, bool)
	//ParseSlice tries to get slice from Middleware by flag name and returns slice and fetching status
	ParseSlice(shortName string, fullName string) ([]interface{}, bool)
}
