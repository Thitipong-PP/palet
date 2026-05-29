package plugin

type Plugin struct {
	Name		string		`yaml:"name"		json:"name"`
	Description	string		`yaml:"description"	json:"description"`
	Icon		string 		`yaml:"icon"		json:"icon"`
	Commands	[]Command	`yaml:"commands"	json:"commands"`
}

type Command struct {
	Name		string		`yaml:"name"		json:"name"`
	Description	string		`yaml:"description"	json:"description"`
	Template	string		`yaml:"template"	json:"template"`
	Args		[]Arg		`yaml:"args"		json:"args"`
}

type Arg struct {
	Name		string		`yaml:"name"		json:"name"`
	Description	string		`yaml:"description"	json:"description"`
	Type		string		`yaml:"type"		json:"type"`
	Required	bool		`yaml:"required"	json:"required"`
	Default		string		`yaml:"default"		json:"default"`
	Choices		[]string	`yaml:"choices"		json:"choices"`
	Flag		string		`yaml:"flag"		json:"flag"`
}

type ArgType string


const (
	ArgTypeString	 ArgType = "string"
	ArgTypeBool		 ArgType = "bool"
	ArgTypeEnum		 ArgType = "enum"
	ArgTypeFile 	 ArgType = "file"
	ArgTypeDir 		 ArgType = "dir"
)