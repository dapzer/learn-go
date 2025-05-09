package proxy

import (
	"net/http"
	"net/url"
	"proxy-api/internal/config"
)

type Opts struct {
	Cfg config.Config
}

func New(opts Opts) *Service {
	return &Service{
		cfg: opts.Cfg,
	}
}

type Service struct {
	cfg config.Config
}

func buildUrl(baseUrl string, path string, query *url.Values) string {
	u, _ := url.Parse(baseUrl)
	u = u.JoinPath(path)
	if query != nil {
		u.RawQuery = query.Encode()
	}

	return u.String()
}

func (s *Service) GetResponse(path string, query url.Values) (*http.Response, error) {
	query.Add("api_key", s.cfg.TmdbApiKey)

	u := buildUrl(s.cfg.TmdbApiUrl, path, &query)

	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
