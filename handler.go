package urlshort

import (
	"gopkg.in/yaml.v2"
	"net/http"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if path, ok := pathsToUrls[request.URL.Path]; ok {
			http.Redirect(writer, request, path, http.StatusFound)
		}
		fallback.ServeHTTP(writer, request)
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
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.

type Yml struct {
	path string `yaml:"path"`
	url string `yaml:"url"`
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var ymlResult []Yml

	err := yaml.Unmarshal(yml, &ymlResult)
	if err != nil {
		return nil, err
	}

	var pathsToUrls = make(map[string]string)
	for _, value := range ymlResult {
		pathsToUrls[value.path] = value.url
	}

	return MapHandler(pathsToUrls, fallback), nil
}
