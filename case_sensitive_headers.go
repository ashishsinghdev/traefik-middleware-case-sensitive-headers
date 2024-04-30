package traefik_middleware_case_sensitive_headers

import (
	"context"
	"net/http"
)

type addHeaderConfig struct {
	Name  string `json:"headerName,omitempty"`
	Value string `json:"headerValue,omitempty"`
}

func createAddHeaderConfig(name string, value string) *addHeaderConfig {
	return &addHeaderConfig{
		Name:  name,
		Value: value,
	}
}

type removeHeaderConfig struct {
	Names []string `json:"headerName,omitempty"`
}

func createRemoveHeaderConfig(names []string) *removeHeaderConfig {
	return &removeHeaderConfig{
		Names: names,
	}
}

type modifyHeaderConfig struct {
	From             string `json:"from,omitempty"`
	To               string `json:"to,omitempty"`
	Prefix           string `json:"prefix,omitempty"`
	Suffix           string `json:"suffix,omitempty"`
	RemoveOriginal   bool   `json:"removeOriginal,omitempty"`
	OverwriteIfExist bool   `json:"overwriteIfExist,omitempty"`
}

func createModifyHeaderConfig(from string, to string, prefix string, suffix string, removeOriginal bool, overwriteIfExist bool) *modifyHeaderConfig {
	return &modifyHeaderConfig{
		From:             from,
		To:               to,
		Prefix:           prefix,
		Suffix:           suffix,
		RemoveOriginal:   removeOriginal,
		OverwriteIfExist: overwriteIfExist,
	}
}

type headerConfig struct {
	AddHeaders    []*addHeaderConfig    `json:"addHeaders,omitempty"`
	ModifyHeaders []*modifyHeaderConfig `json:"modifyHeaders,omitempty"`
	RemoveHeaders *removeHeaderConfig   `json:"removeHeaders,omitempty"`
}

func createHeaderConfig() *headerConfig {
	return &headerConfig{}
}

type Config struct {
	Headers *headerConfig `json:"headers,omitempty"`
}

func CreateConfig() *Config {
	return &Config{}
}

type ProcessHeader struct {
	next   http.Handler
	name   string
	config *Config
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &ProcessHeader{
		next: next, config: config, name: name,
	}, nil
}

func (headerRewrite *ProcessHeader) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rewriteHeaders(&req.Header, headerRewrite.config.Headers)
	headerRewrite.next.ServeHTTP(rw, req)
}

func rewriteHeaders(headers *http.Header, headerConfig *headerConfig) {
	for _, addHeaderConfig := range headerConfig.AddHeaders {
		(*headers)[addHeaderConfig.Name] = []string{addHeaderConfig.Value}
	}

	for _, headerName := range headerConfig.RemoveHeaders.Names {
		headers.Del(headerName)
		delete((*headers), headerName)
	}

	for _, modifyHeaderConfig := range headerConfig.ModifyHeaders {
		headerValues := headers.Values(modifyHeaderConfig.From)

		if modifyHeaderConfig.OverwriteIfExist {
			headers.Del(modifyHeaderConfig.To)
			delete((*headers), modifyHeaderConfig.To)
		}

		for _, headerValue := range headerValues {
			if headerValue != "" {
				if len(modifyHeaderConfig.Prefix) > 0 {
					headerValue = modifyHeaderConfig.Prefix + headerValue
				}
				if len(modifyHeaderConfig.Suffix) > 0 {
					headerValue += modifyHeaderConfig.Suffix
				}
				(*headers)[modifyHeaderConfig.To] = []string{headerValue}
			}
		}

		if modifyHeaderConfig.RemoveOriginal {
			headers.Del(modifyHeaderConfig.From)
			delete((*headers), modifyHeaderConfig.From)
		}
	}
}
