package client

type FormContentType string

const (
	FORM_URLENCODED FormContentType = "application/x-www-form-urlencoded"
	JSON            FormContentType = "application/json"
	HTML            FormContentType = "text/html"
	XML             FormContentType = "text/xml"
	FORM_DATA       FormContentType = "multipart/form-data"
)
