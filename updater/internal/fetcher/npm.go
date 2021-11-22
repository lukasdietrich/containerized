package fetcher

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type npm struct {
	client *http.Client
}

func (n npm) latest(name string) (*Release, error) {
	var res struct {
		Time map[string]time.Time `json:"time"`
		Tags struct {
			Latest string `latest`
		} `json:"dist-tags"`
	}

	if err := n.fetch(name, &res); err != nil {
		return nil, fmt.Errorf("npm: could not fetch releases for %q: %w", name, err)
	}

	if res.Tags.Latest == "" {
		return nil, fmt.Errorf("npm: no latest tag for %q found", name)
	}

	tag := res.Tags.Latest
	release := Release{
		Name:        name,
		Version:     tag,
		PublishedAt: res.Time[tag],
	}

	return &release, nil
}

func (n npm) fetch(name string, v interface{}) error {
	res, err := n.client.Get(n.buildURL(name))
	if err != nil {
		return err
	}

	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(v)
}

func (npm) buildURL(name string) string {
	u := url.URL{
		Scheme: "https",
		Host:   "registry.npmjs.org",
		Path:   name,
	}

	return u.String()
}
