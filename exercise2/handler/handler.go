package handler

import (
	"gopkg.in/yaml.v2"
	"net/http"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if v, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, v, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	urls, err := parseYaml(yml)
	if err != nil {
		return nil, err
	}
	return MapHandler(urls, fallback), nil
}

func parseYaml(yml []byte) (map[string]string, error) {
	var records []shortUrl
	if err := yaml.Unmarshal(yml, &records); err != nil {
		return nil, err
	}
	urlMap := make(map[string]string)
	for _, v := range records {
		urlMap[v.Path] = v.Url
	}
	return urlMap, nil
}

type shortUrl struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}
