package clash

import (
	"errors"
	"os"
	"unsafe"

	"subscribe2clash/internal/subscribe"
	"subscribe2clash/internal/xbase64"
)

type SourceType int

const (
	Url SourceType = iota
	File
	Text
)

func Config(sourceType SourceType, source string) (string, error) {
	var contents []string
	var err error

	switch sourceType {
	case Url:
		contents, err = subscribe.GetRawProxiesFromLinks(source)
		if err != nil {
			return "", err
		}

	case File:
		file, err := os.ReadFile(source)
		if err != nil {
			return "", err
		}

		contents = append(contents, subscribe.ParseRawProxies(unsafe.String(&file[0], len(file))))
	default:
		return "", errors.New("unknown source type")
	}

	proxies := subscribe.ParseProxy(contents)
	config, err := subscribe.GenerateClashConfig(proxies)
	if err != nil {
		return "", err
	}

	return xbase64.UnicodeEmojiDecode(string(config)), nil
}

func Nodes(urls []string) (string, error) {
	proxies := subscribe.ParseProxy(urls)
	config, err := subscribe.GenerateClashConfig(proxies)
	if err != nil {
		return "", err
	}

	return xbase64.UnicodeEmojiDecode(string(config)), nil
}
