package casesensitiveheaders

import (
	"context"
	"net/http"
)

type AddHeaderConfig struct {
	Name  string `json:"headerName,omitempty"`
	Value string `json:"headerValue,omitempty"`
}

func CreateAddHeaderConfig(name string, value string) *AddHeaderConfig {
	return &AddHeaderConfig{
		Name:  name,
		Value: value,
	}
}

type RemoveHeaderConfig struct {
	Names []string `json:"headerName,omitempty"`
}

func CreateRemoveHeaderConfig(names []string) *RemoveHeaderConfig {
	return &RemoveHeaderConfig{
		Names: names,
	}
}

type ModifyHeaderConfig struct {
	From             string `json:"from,omitempty"`
	To               string `json:"to,omitempty"`
	Prefix           string `json:"prefix,omitempty"`
	Suffix           string `json:"suffix,omitempty"`
	RemoveOriginal   bool   `json:"removeOriginal,omitempty"`
	OverwriteIfExist bool   `json:"overwriteIfExist,omitempty"`
}

func CreateModifyHeaderConfig(from string, to string, prefix string, suffix string, removeOriginal bool, overwriteIfExist bool) *ModifyHeaderConfig {
	return &ModifyHeaderConfig{
		From:             from,
		To:               to,
		Prefix:           prefix,
		Suffix:           suffix,
		RemoveOriginal:   removeOriginal,
		OverwriteIfExist: overwriteIfExist,
	}
}

type HeaderConfig struct {
	AddHeaders    []*AddHeaderConfig    `json:"addHeaders,omitempty"`
	RemoveHeaders []*RemoveHeaderConfig `json:"removeHeaders,omitempty"`
	ModifyHeaders []*ModifyHeaderConfig `json:"modifyHeaders,omitempty"`
}

func CreateHeaderConfig() *HeaderConfig {
	return &HeaderConfig{}
}

type Config struct {
	Headers []*HeaderConfig `json:"headers,omitempty"`
}

func CreateConfig(headers []*HeaderConfig) *Config {
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
	for _, header := range headerRewrite.config.Headers {
		rewriteHeaders(&req.Header, header)
	}
	headerRewrite.next.ServeHTTP(rw, req)
}

func rewriteHeaders(headers *http.Header, headerConfig *HeaderConfig) {
	for _, addHeaderConfig := range headerConfig.AddHeaders {
		(*headers)[addHeaderConfig.Name] = []string{addHeaderConfig.Value}
	}

	for _, removeHeaderConfig := range headerConfig.RemoveHeaders {
		for _, headerName := range removeHeaderConfig.Names {
			headers.Del(headerName)
			delete((*headers), headerName)
		}

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
					headerValue = headerValue + modifyHeaderConfig.Suffix
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
