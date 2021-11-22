package fetcher

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/time/rate"
)

var fetcherMap map[string]fetcher

func init() {
	client := &http.Client{
		Timeout: time.Second * 30,
		Transport: &loggingRoundTripper{
			roundTripper: http.DefaultTransport,
			ratelimit:    rate.NewLimiter(rate.Every(time.Second*10), 1),
		},
	}

	fetcherMap = map[string]fetcher{
		"npm":  npm{client},
		"pypi": pypi{client},
	}
}

type loggingRoundTripper struct {
	roundTripper http.RoundTripper
	ratelimit    *rate.Limiter
}

func (l *loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := l.ratelimit.Wait(req.Context()); err != nil {
		return nil, err
	}

	log.Printf("%s %q", req.Method, req.URL)
	return l.roundTripper.RoundTrip(req)
}

type Release struct {
	Name        string
	Version     string
	PublishedAt time.Time
}

type fetcher interface {
	latest(name string) (*Release, error)
}

func Latest(origin string) (*Release, error) {
	spec, err := url.Parse(origin)
	if err != nil {
		return nil, fmt.Errorf("could not parse origin %q: %w", origin, err)
	}

	f, err := fetcherBySpec(spec)
	if err != nil {
		return nil, err
	}

	release, err := f.latest(spec.Opaque)
	if release != nil {
		release.PublishedAt = release.PublishedAt.UTC()
	}

	return release, err
}

func fetcherBySpec(spec *url.URL) (fetcher, error) {
	f, ok := fetcherMap[spec.Scheme]
	if !ok {
		return nil, fmt.Errorf("no fetcher for %q", spec)
	}

	return f, nil
}
