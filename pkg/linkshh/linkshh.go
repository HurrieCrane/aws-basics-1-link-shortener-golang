package linkshh

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"linkshh/pkg/linkshh/store"
	"net/url"
	"os"
)

const (
	defaultTinyDomain = "https://link.thestudio.com"
	domainEnvKey      = "LINK_DOMAIN"
)

type URLShortener struct {
}

func ShortenUri(ctx context.Context, uri *url.URL) (*url.URL, error) {
	h := sha1.Sum([]byte(uri.String()))

	var domain string
	if domain = os.Getenv(domainEnvKey); domain == "" {
		domain = defaultTinyDomain
	}

	shortUri, err := url.Parse(fmt.Sprintf("%slink/%x", domain, h))
	if err != nil {
		return nil, err
	}

	s, err := store.NewStore()
	if err != nil {
		return nil, err
	}

	err = s.StoreUrl(ctx, fmt.Sprintf("%x", h), uri.String())
	if err != nil {
		return nil, err
	}

	return shortUri, nil
}

func ExpandHash(ctx context.Context, hash string) (*url.URL, error) {
	s, err := store.NewStore()
	if err != nil {
		return nil, err
	}

	rawURL, err := s.RetrieveUrlForHash(ctx, hash)
	if err != nil {
		return nil, err
	}

	tinyUrl, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	if tinyUrl == nil {
		return nil, errors.New("no tiny url found")
	}

	return tinyUrl, err
}
