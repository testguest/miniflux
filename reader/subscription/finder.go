// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package subscription // import "miniflux.app/reader/subscription"

import (
	"io"
	"strings"

	"miniflux.app/errors"
	"miniflux.app/http/client"
	"miniflux.app/reader/browser"
	"miniflux.app/reader/parser"
	"miniflux.app/url"

	"github.com/PuerkitoBio/goquery"
)

var (
	errUnreadableDoc = "Unable to analyze this page: %v"
)

// FindSubscriptions downloads and try to find one or more subscriptions from an URL.
func FindSubscriptions(websiteURL, userAgent, username, password string) (Subscriptions, *errors.LocalizedError) {
	request := client.New(websiteURL)
	request.WithCredentials(username, password)
	request.WithUserAgent(userAgent)
	response, err := browser.Exec(request)
	if err != nil {
		return nil, err
	}

	body := response.BodyAsString()
	if format := parser.DetectFeedFormat(body); format != parser.FormatUnknown {
		var subscriptions Subscriptions
		subscriptions = append(subscriptions, &Subscription{
			Title: response.EffectiveURL,
			URL:   response.EffectiveURL,
			Type:  format,
		})

		return subscriptions, nil
	}

	return parseDocument(response.EffectiveURL, strings.NewReader(body))
}

func parseDocument(websiteURL string, data io.Reader) (Subscriptions, *errors.LocalizedError) {
	var subscriptions Subscriptions
	queries := map[string]string{
		"link[type='application/rss+xml']":  "rss",
		"link[type='application/atom+xml']": "atom",
		"link[type='application/json']":     "json",
	}

	doc, err := goquery.NewDocumentFromReader(data)
	if err != nil {
		return nil, errors.NewLocalizedError(errUnreadableDoc, err)
	}

	for query, kind := range queries {
		doc.Find(query).Each(func(i int, s *goquery.Selection) {
			subscription := new(Subscription)
			subscription.Type = kind

			if title, exists := s.Attr("title"); exists {
				subscription.Title = title
			} else {
				subscription.Title = "Feed"
			}

			if feedURL, exists := s.Attr("href"); exists {
				subscription.URL, _ = url.AbsoluteURL(websiteURL, feedURL)
			}

			if subscription.Title == "" {
				subscription.Title = subscription.URL
			}

			if subscription.URL != "" {
				subscriptions = append(subscriptions, subscription)
			}
		})
	}

	return subscriptions, nil
}
