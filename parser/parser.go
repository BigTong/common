package parser

type Parser interface {
	Parse(string, map[string]interface{}) map[string]interface{}
	ParseFromStringConfig(string, string) (map[string]interface{}, error)
	ParseById(string, string, map[string]interface{}) map[string]interface{}
}
