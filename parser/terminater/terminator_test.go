package terminater

import (
	"encoding/json"
	"testing"

	"github.com/BigTong/common/file"
)

var htmlConfig string = `{
	"_input_type":{
		"content_type":"html"
	},
	"source":{
		"selector":{
			"xpath":"head title"
		}
	},
	"cars":{
		"object_config":{
			"root_path":{
				"xpath":"#tab-content-item1 .uibox div dl"
			},
			"fields":{
				"brand":{
					"selector":{
						"xpath":"dd .h3-tit"
					}
				},
				"products":{
					"object_config":{
						"root_path":{
							"xpath":"dd ul li"
						},
						"fields":{
							"link":{
								"selector":{
									"xpath":"h4 a",
									"attr":"href"
								}
							},
							"name":{
								"selector":{
								"xpath":"h4 a"
								}
							}
						}
					}
				}
			}
		}
	}
}`

var jsonConfig string = `{
	"_input_type":{
		"content_type":"json"
	},
	"title":{
		"selector":{
			"json_key":"parse_title"
		}
	},
	"buyers":{
		"object_config":{
			"root_path":{
				"json_key":"buyers[*]"
			},
			"fields":{
				"name":{
					"selector":{
						"json_key":"name"
					}
				},
				"addr":{
					"selector":{
						"json_key":"addr"
					}
				},
				"age":{
					"selector":{
						"json_key":"age"
					},
					"type":"int"
				}
			}
		}
	}
}
`

func TestHtmlParse(t *testing.T) {
	src := file.ReadFileToString("test_data/html")
	parser := NewTerminater()

	config := make(map[string]interface{})
	err := json.Unmarshal([]byte(htmlConfig), &config)
	if err != nil {
		t.Error(err.Error())
		return
	}

	result := parser.Parse(src, config)
	data, _ := json.Marshal(result)
	t.Log(string(data))
}

func TestJsonParse(t *testing.T) {
	src := file.ReadFileToString("test_data/json")
	parser := NewTerminater()

	config := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonConfig), &config)
	if err != nil {
		t.Error(err.Error())
		return
	}

	result := parser.Parse(src, config)
	data, _ := json.Marshal(result)
	t.Log(string(data))
}
