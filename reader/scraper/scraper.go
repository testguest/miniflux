// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package scraper // import "miniflux.app/reader/scraper"

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"miniflux.app/http/client"
	"miniflux.app/logger"
	"miniflux.app/reader/readability"
	"miniflux.app/url"

	"github.com/PuerkitoBio/goquery"
)

// Fetch downloads a web page and returns relevant contents.
func Fetch(websiteURL, rules, userAgent string) (string, error) {
	clt := client.New(websiteURL)
	if userAgent != "" {
		clt.WithUserAgent(userAgent)
	}

	response, err := clt.Get()
	if err != nil {
		return "", err
	}

	if response.HasServerFailure() {
		return "", errors.New("scraper: unable to download web page")
	}

	if !isWhitelistedContentType(response.ContentType) {
		return "", fmt.Errorf("scraper: this resource is not a HTML document (%s)", response.ContentType)
	}

	if err = response.EnsureUnicodeBody(); err != nil {
		return "", err
	}

	// The entry URL could redirect somewhere else.
	websiteURL = response.EffectiveURL

	if rules == "" {
		rules = getPredefinedScraperRules(websiteURL)
	}

	var content string
	if rules != "" {
		logger.Debug(`[Scraper] Using rules %q for %q`, rules, websiteURL)
		content, err = scrapContent(response.Body, rules)
	} else {
		logger.Debug(`[Scraper] Using readability for %q`, websiteURL)
		content, err = readability.ExtractContent(response.Body)
	}

	if err != nil {
		return "", err
	}

	return content, nil
}

func scrapContent(page io.Reader, rules string) (string, error) {
	document, err := goquery.NewDocumentFromReader(page)
	if err != nil {
		return "", err
	}

	contents := ""
	document.Find(rules).Each(func(i int, s *goquery.Selection) {
		var content string

		content, _ = goquery.OuterHtml(s)
		contents += content
	})

	return contents, nil
}

func getPredefinedScraperRules(websiteURL string) string {
	urlDomain := url.Domain(websiteURL)

	for domain, rules := range predefinedRules {
		if strings.Contains(urlDomain, domain) {
			return rules
		}
	}

	return ""
}

func isWhitelistedContentType(contentType string) bool {
	contentType = strings.ToLower(contentType)
	return strings.HasPrefix(contentType, "text/html") ||
		strings.HasPrefix(contentType, "application/xhtml+xml")
}
