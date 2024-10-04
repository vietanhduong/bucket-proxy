package proxy

type Option func(*Proxy)

func WithWebMode(enabled bool) Option {
	return func(p *Proxy) {
		p.webMode = enabled
	}
}

func WithIndexPage(page string) Option {
	return func(p *Proxy) {
		if page != "" {
			p.indexPage = page
		}
	}
}

func WithNotFoundPage(page string) Option {
	return func(p *Proxy) {
		if page != "" {
			p.notFoundPage = page
		}
	}
}
