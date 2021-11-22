package fetcher

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"
)

type pypi struct {
	client *http.Client
}

func (p pypi) latest(name string) (*Release, error) {
	var res struct {
		XMLName xml.Name `xml:"rss"`
		Items   []struct {
			Title   string `xml:"title"`
			PubDate string `xml:"pubDate"`
		} `xml:"channel>item"`
	}

	if err := p.fetch(name, &res); err != nil {
		return nil, fmt.Errorf("pypi: could not fetch releases for %q: %w", name, err)
	}

	if len(res.Items) == 0 {
		return nil, fmt.Errorf("pypi: no latest tag for %q found", name)
	}

	item := res.Items[0]
	publishedAt, err := time.Parse(time.RFC1123, item.PubDate)
	if err != nil {
		return nil, fmt.Errorf("pypi: could not parse pubDate %q: %w", item.PubDate, err)
	}

	release := Release{
		Name:        name,
		Version:     item.Title,
		PublishedAt: publishedAt,
	}

	return &release, nil
}

func (p pypi) fetch(name string, v interface{}) error {
	res, err := p.client.Get(p.buildURL(name))
	if err != nil {
		return err
	}

	defer res.Body.Close()
	return xml.NewDecoder(res.Body).Decode(v)
}

func (pypi) buildURL(name string) string {
	u := url.URL{
		Scheme: "https",
		Host:   "pypi.org",
		Path:   path.Join("/rss/project", name, "releases.xml"),
	}

	return u.String()
}
