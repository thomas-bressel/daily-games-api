package rss

import (
	"html"
	"regexp"
	"strings"

	"github.com/mmcdole/gofeed"
)

// firstImgRe extracts the first <img src="..."> URL from HTML content.
var firstImgRe = regexp.MustCompile(`(?i)<img[^>]+src=["']([^"']+)["']`)

// textInnerRe extracts all content inside et_pb_text_inner divs (Divi/ET Builder).
var textInnerRe = regexp.MustCompile(`(?i)<div[^>]+class="et_pb_text_inner"[^>]*>([\s\S]*?)</div>`)

// segaTextRe extracts content from Drupal field--name-field-text divs (Sega-Mag).
var segaTextRe = regexp.MustCompile(`(?i)<div[^>]+class="[^"]*field--name-field-text[^"]*"[^>]*>([\s\S]*?)</div>`)

// htmlTagRe strips all HTML tags.
var htmlTagRe = regexp.MustCompile(`<[^>]+>`)

// cdataURLRe extracts a URL from a CDATA-wrapped string like <![CDATA[https://...]]>
var cdataURLRe = regexp.MustCompile(`(?i)<!\[CDATA\[(https?://[^\]]+)]]>`)

// Enricher receives the raw feed item and returns corrected image URL and description.
// Return empty strings to keep the already-extracted values unchanged.
type Enricher func(item *gofeed.Item) (imageURL string, description string)

// enrichers maps a feed ID to its custom enricher function.
// Add a new entry here whenever a feed requires special handling.
var enrichers = map[string]Enricher{
	"abandonware": enrichAbandonware,
	"amstrad-eu":  enrichAmstradEU,
	"sega-mag":    enrichSegaMag,
}

// enrichAmstradEU handles the Amstrad.eu WordPress/Divi feed:
//  1. Image: first <img src> found in content:encoded.
//  2. Description: first et_pb_text_inner block with real text (skips Divi widget blocks with % signs).
func enrichAmstradEU(item *gofeed.Item) (imageURL string, description string) {
	if item.Content == "" {
		return "", ""
	}

	// Extract image from first <img>
	if m := firstImgRe.FindStringSubmatch(item.Content); len(m) > 1 {
		imageURL = m[1]
	}

	// Extract first et_pb_text_inner with real narrative text
	matches := textInnerRe.FindAllStringSubmatch(item.Content, -1)
	for _, m := range matches {
		text := strings.TrimSpace(htmlTagRe.ReplaceAllString(m[1], " "))
		text = strings.Join(strings.Fields(text), " ")
		// Skip Divi widget blocks (counters, percentages, empty)
		if len(text) < 40 || strings.Contains(text, "%") {
			continue
		}
		description = text
		break
	}

	return imageURL, description
}

// enrichSegaMag handles the Sega-Mag Drupal feed:
//  1. Image: first <img src> in the description HTML.
//  2. Description: content of the field--name-field-text div (clean HTML, no double-encoding).
func enrichSegaMag(item *gofeed.Item) (imageURL string, description string) {
	raw := item.Description
	if raw == "" {
		return "", ""
	}

	if m := firstImgRe.FindStringSubmatch(raw); len(m) > 1 {
		imageURL = m[1]
	}

	if m := segaTextRe.FindStringSubmatch(raw); len(m) > 1 {
		text := strings.TrimSpace(htmlTagRe.ReplaceAllString(m[1], " "))
		description = strings.Join(strings.Fields(text), " ")
	}

	return imageURL, description
}

// enrichAbandonware handles two quirks of the Abandonware France RSS feed:
//  1. Image is in the <im:image> extension tag, not in standard enclosure/media fields.
//  2. Description is double-encoded (e.g. &amp;egrave; instead of &egrave;),
//     requiring two passes of html.UnescapeString to get clean text.
func enrichAbandonware(item *gofeed.Item) (imageURL string, description string) {
	// Extract image from <im:image> extension
	if im := item.Extensions["im"]; im != nil {
		if images, ok := im["image"]; ok && len(images) > 0 {
			raw := images[0].Value
			// Value may be wrapped in a CDATA-like string — strip it
			raw = strings.TrimSpace(raw)
			if strings.HasPrefix(raw, "http") {
				imageURL = raw
			} else if m := cdataURLRe.FindStringSubmatch(raw); len(m) > 1 {
				imageURL = m[1]
			}
		}
	}

	// Double-decode the description
	raw := item.Description
	if raw == "" {
		raw = item.Content
	}
	description = html.UnescapeString(html.UnescapeString(raw))

	return imageURL, description
}
