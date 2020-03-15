// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package rewrite // import "miniflux.app/reader/rewrite"

import (
	"strings"

	"miniflux.app/logger"
	"miniflux.app/url"
)

// Rewriter modify item contents with a set of rewriting rules.
func Rewriter(entryURL, entryContent, customRewriteRules string) string {
	rulesList := getPredefinedRewriteRules(entryURL)
	if customRewriteRules != "" {
		rulesList = customRewriteRules
	}

	rules := strings.Split(rulesList, ",")
	rules = append(rules, "add_pdf_download_link")

	logger.Debug(`[Rewrite] Applying rules %v for %q`, rules, entryURL)

	for _, rule := range rules {
		switch strings.TrimSpace(rule) {
		case "add_image_title":
			entryContent = addImageTitle(entryURL, entryContent)
		case "add_mailto_subject":
			entryContent = addMailtoSubject(entryURL, entryContent)
		case "add_dynamic_image":
			entryContent = addDynamicImage(entryURL, entryContent)
		case "add_youtube_video":
			entryContent = addYoutubeVideo(entryURL, entryContent)
		case "add_pdf_download_link":
			entryContent = addPDFLink(entryURL, entryContent)
		case "nl2br":
			entryContent = replaceLineFeeds(entryContent)
		case "convert_text_link", "convert_text_links":
			entryContent = replaceTextLinks(entryContent)
		}
	}

	return entryContent
}

func getPredefinedRewriteRules(entryURL string) string {
	urlDomain := url.Domain(entryURL)

	for domain, rules := range predefinedRules {
		if strings.Contains(urlDomain, domain) {
			return rules
		}
	}

	return ""
}
