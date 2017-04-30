package parser

const (
	OBJECT_CONFIG  = "object_config"
	ELEM_CONFIG    = "elem_config"
	UNKNOWN_CONFIG = "unknown_config"
)

const (
	TYPE_INT    = "int"
	TYPE_FLOAT  = "float"
	TYPE_STRING = "string"
)

const (
	CONTENT_TYPE_JSON   = "json"
	CONTENT_TYPE_HTML   = "html"
	CONTENT_TYPE_STRING = "string"
)

const (
	INPUT_TYPE = "_input_type"
)

type Selector struct {
	Xpath     string   `json:"xpath,omitempty"`
	Attr      string   `json:"attr,omitempty"`
	JsonKey   string   `json:"json_key,omitempty"`
	Contains  []string `json:"contains,omitempty"`
	Extractor string   `json:"extractor,omitempty"`
	Index     string   `json:"index,omitempty"`
}

type InputConfig struct {
	ContentType string `json:"content_type,omitempty"`
	Extractor   string `json:"extractor,omitempty"`
}

type ElemConfig struct {
	Selector     *Selector     `json:"selector,omitempty"`
	ObjectConfig *ObjectConfig `json:"object_config,omitempty"`
	Type         string        `json:"type,omitempty"`
	SavePath     string        `json:"save_path,omitemtpty"`
	SaveName     string        `json:"save_name,omitempty"`
}

type ObjectConfig struct {
	RootPath *Selector              `json:"root_path,omitemtpy"`
	Fields   map[string]interface{} `json:"fields,omitempty"`
}
