package http_entity

var defaultFetchOption *FetchOption = &FetchOption{
	useWapUserAgent: false,
	followRedirect:  false,
	transformToUtf8: false,
	render:          false,
}

func GetDefaultFetchOption() *FetchOption {
	return defaultFetchOption
}

// fetch option, all field is immutable, initialize with builder, only provide accessor
// TODO: not implement builder now
type FetchOption struct {
	useWapUserAgent bool
	followRedirect  bool
	transformToUtf8 bool
	render          bool
}

func newFetchOption(b *FetchOptionBuilder) *FetchOption {
	opt := new(FetchOption)
	opt.followRedirect = b.followRedirect
	opt.render = b.render
	opt.transformToUtf8 = b.transformToUtf8
	opt.useWapUserAgent = b.useWapUserAgent

	return opt
}

func (fo *FetchOption) GetUseWapUserAgent() bool {
	return fo.useWapUserAgent
}

func (fo *FetchOption) GetFollowRedirect() bool {
	return fo.followRedirect
}

func (fo *FetchOption) GetTransformToUtf8() bool {
	return fo.transformToUtf8
}

func (fo *FetchOption) GetRender() bool {
	return fo.render
}

type FetchOptionBuilder struct {
	useWapUserAgent bool
	followRedirect  bool
	transformToUtf8 bool
	render          bool
}

func NewFetchOptionBuilder() *FetchOptionBuilder {
	return new(FetchOptionBuilder)
}

func (b *FetchOptionBuilder) SetUseWapUserAgent(useWapUserAgent bool) *FetchOptionBuilder {
	b.useWapUserAgent = useWapUserAgent
	return b
}

func (b *FetchOptionBuilder) SetFollowRedirect(followRedirect bool) *FetchOptionBuilder {
	b.followRedirect = followRedirect
	return b
}

func (b *FetchOptionBuilder) SetTransformToUtf8(transformToUtf8 bool) *FetchOptionBuilder {
	b.transformToUtf8 = transformToUtf8
	return b
}

func (b *FetchOptionBuilder) SetRender(render bool) *FetchOptionBuilder {
	b.render = render
	return b
}

func (b *FetchOptionBuilder) Build() *FetchOption {
	return newFetchOption(b)
}
