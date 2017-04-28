package filter

type Filter interface {
	NeedFilter(string) bool
	SetFilter(string)
	Reset()
	Dump()
}
