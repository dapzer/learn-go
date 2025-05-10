package proxy

import (
	"github.com/h2non/bimg"
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

func buildUrl(baseUrl string, path []string, query *url.Values) string {
	u, _ := url.Parse(baseUrl)
	u = u.JoinPath(path...)
	if query != nil {
		u.RawQuery = query.Encode()
	}

	return u.String()
}

func (s *Service) GetResponse(path string, query url.Values) (*http.Response, error) {
	query.Add("api_key", s.cfg.TmdbApiKey)

	u := buildUrl(s.cfg.TmdbApiUrl, []string{path}, &query)

	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s Service) GetImage(path string) (*http.Response, error) {
	u := buildUrl(s.cfg.TmdbImageApiUrl, []string{"w500", path}, nil)
	resp, err := http.Get(u)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s Service) ProccessImage(sourceImage []byte, size *int) ([]byte, error) {
	processedImage := sourceImage
	if size != nil && *size != 0 {
		resizedImage, err := bimg.NewImage(processedImage).Process(bimg.Options{Width: *size})
		if err != nil {
			return nil, err
		}
		processedImage = resizedImage
	}

	convertedImage, err := bimg.NewImage(processedImage).Convert(bimg.WEBP)
	if err != nil {
		return nil, err
	}
	processedImage = convertedImage

	return processedImage, nil
}
