package parser

const (
	KObjectConfig  = "kObject"
	KElemConfig    = "kElem"
	KUnknownConfig = "kUnknown"
)

const (
	KInt    = "int"
	KFloat  = "float"
	KString = "string"
)

const (
	KContentJson   = "json"
	KContentHtml   = "html"
	KContentString = "string"
)

const (
	KInputType = "_input_type"
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
