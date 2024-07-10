// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package subscription

import (
	"errors"
	"strings"
	"testing"
)

func TestFindYoutubePlaylistFeed(t *testing.T) {
	scenarios := map[string]string{
		"https://www.youtube.com/playlist?list=PLOOwEPgFWm_NHcQd9aCi5JXWASHO_n5uR":            "https://www.youtube.com/feeds/videos.xml?playlist_id=PLOOwEPgFWm_NHcQd9aCi5JXWASHO_n5uR",
		"https://www.youtube.com/playlist?list=PLOOwEPgFWm_N42HlCLhqyJ0ZBWr5K1QDM":            "https://www.youtube.com/feeds/videos.xml?playlist_id=PLOOwEPgFWm_N42HlCLhqyJ0ZBWr5K1QDM",
		"https://www.youtube.com/watch?v=dQw4w9WgXcQ&list=PLOOwEPgFWm_N42HlCLhqyJ0ZBWr5K1QDM": "https://www.youtube.com/feeds/videos.xml?playlist_id=PLOOwEPgFWm_N42HlCLhqyJ0ZBWr5K1QDM",
	}

	for websiteURL, expectedFeedURL := range scenarios {
		subscriptions, localizedError := NewSubscriptionFinder(nil).FindSubscriptionsFromYouTubePlaylistPage(websiteURL)
		if localizedError != nil {
			t.Fatalf(`Parsing a correctly formatted YouTube playlist page should not return any error: %v`, localizedError)
		}

		if len(subscriptions) != 1 {
			t.Fatal(`Incorrect number of subscriptions returned`)
		}

		if subscriptions[0].URL != expectedFeedURL {
			t.Errorf(`Unexpected Feed, got %s, instead of %s`, subscriptions[0].URL, expectedFeedURL)
		}
	}
}

func TestItDoesNotConsiderPlaylistWatchPageAsVideoWatchPage(t *testing.T) {
	_, localizedError := NewSubscriptionFinder(nil).FindSubscriptionsFromYouTubeVideoPage("https://www.youtube.com/watch?v=dQw4w9WgXcQ&list=PLOOwEPgFWm_N42HlCLhqyJ0ZBWr5K1QDM")
	if localizedError != nil {
		t.Fatalf(`Should not consider a playlist watch page as a video watch page`)
	}
}

func TestYoutubeIdExtractor(t *testing.T) {
	type testResult struct {
		ID    string
		Kind  youtubeKind
		error error
	}
	urls := map[string]testResult{
		"https://www.youtube.com/watch?v=dQw4w9WgXcQ": {
			ID:    "dQw4w9WgXcQ",
			Kind:  youtubeIDKindVideo,
			error: nil,
		},
		"https://www.youtube.com/watch?v=dQw4w9WgXcQ&t=1": {
			ID:    "dQw4w9WgXcQ",
			Kind:  youtubeIDKindVideo,
			error: nil,
		},
		"https://www.youtube.com/watch?t=1&v=dQw4w9WgXcQ": {
			ID:    "dQw4w9WgXcQ",
			Kind:  youtubeIDKindVideo,
			error: nil,
		},
		"https://www.youtube.com/watch?v=dQw4w9WgXcQ&list=PLOOwEPgFWm_N42HlCLhqyJ0ZBWr5K1QDM": {
			ID:    "PLOOwEPgFWm_N42HlCLhqyJ0ZBWr5K1QDM",
			Kind:  youtubeIDKindPlaylist,
			error: nil,
		},
		"https://www.youtube.com/playlist?list=PLOOwEPgFWm_NHcQd9aCi5JXWASHO_n5uR": {
			ID:    "PLOOwEPgFWm_NHcQd9aCi5JXWASHO_n5uR",
			Kind:  youtubeIDKindPlaylist,
			error: nil,
		},
		"https://www.youtube.com/channel/UC-Qj80avWItNRjkZ41rzHyw": {
			ID:    "UC-Qj80avWItNRjkZ41rzHyw",
			Kind:  youtubeIDKindChannel,
			error: nil,
		},
		"https://www.example.com/channel/UC-Qj80avWItNRjkZ41rzHyw": {
			ID:    "",
			Kind:  "",
			error: errNotYoutubeUrl,
		},
	}

	for websiteURL, expected := range urls {
		kind, id, err := youtubeURLIDExtractor(websiteURL)
		if !errors.Is(err, expected.error) {
			t.Fatalf(`Unexpected error: %v got %v`, expected.error, err)
		}
		if id != expected.ID {
			t.Fatalf(`Unexpected ID: %v got %v`, expected.ID, id)
		}
		if kind != expected.Kind {
			t.Fatalf(`Unexpected Kind: %v got %v`, expected.Kind, kind)
		}
	}
}

func TestFindYoutubeChannelFeed(t *testing.T) {
	scenarios := map[string]string{
		"https://www.youtube.com/channel/UC-Qj80avWItNRjkZ41rzHyw": "https://www.youtube.com/feeds/videos.xml?channel_id=UC-Qj80avWItNRjkZ41rzHyw",
	}

	for websiteURL, expectedFeedURL := range scenarios {
		subscriptions, localizedError := NewSubscriptionFinder(nil).FindSubscriptionsFromYouTubeChannelPage(websiteURL)
		if localizedError != nil {
			t.Fatalf(`Parsing a correctly formatted YouTube channel page should not return any error: %v`, localizedError)
		}

		if len(subscriptions) != 1 {
			t.Fatal(`Incorrect number of subscriptions returned`)
		}

		if subscriptions[0].URL != expectedFeedURL {
			t.Errorf(`Unexpected Feed, got %s, instead of %s`, subscriptions[0].URL, expectedFeedURL)
		}
	}
}

func TestParseWebPageWithRssFeed(t *testing.T) {
	htmlPage := `
	<!doctype html>
	<html>
		<head>
			<link href="http://example.org/rss" rel="alternate" type="application/rss+xml" title="Some Title">
		</head>
		<body>
		</body>
	</html>`

	subscriptions, err := NewSubscriptionFinder(nil).FindSubscriptionsFromWebPage("http://example.org/", "text/html", strings.NewReader(htmlPage))
	if err != nil {
		t.Fatalf(`Parsing a correctly formatted HTML page should not return any error: %v`, err)
	}

	if len(subscriptions) != 1 {
		t.Fatal(`Incorrect number of subscriptions returned`)
	}

	if subscriptions[0].Title != "Some Title" {
		t.Errorf(`Incorrect subscription title: %q`, subscriptions[0].Title)
	}

	if subscriptions[0].URL != "http://example.org/rss" {
		t.Errorf(`Incorrect subscription URL: %q`, subscriptions[0].URL)
	}

	if subscriptions[0].Type != "rss" {
		t.Errorf(`Incorrect subscription type: %q`, subscriptions[0].Type)
	}
}

func TestParseWebPageWithAtomFeed(t *testing.T) {
	htmlPage := `
	<!doctype html>
	<html>
		<head>
			<link href="http://example.org/atom.xml" rel="alternate" type="application/atom+xml" title="Some Title">
		</head>
		<body>
		</body>
	</html>`

	subscriptions, err := NewSubscriptionFinder(nil).FindSubscriptionsFromWebPage("http://example.org/", "text/html", strings.NewReader(htmlPage))
	if err != nil {
		t.Fatalf(`Parsing a correctly formatted HTML page should not return any error: %v`, err)
	}

	if len(subscriptions) != 1 {
		t.Fatal(`Incorrect number of subscriptions returned`)
	}

	if subscriptions[0].Title != "Some Title" {
		t.Errorf(`Incorrect subscription title: %q`, subscriptions[0].Title)
	}

	if subscriptions[0].URL != "http://example.org/atom.xml" {
		t.Errorf(`Incorrect subscription URL: %q`, subscriptions[0].URL)
	}

	if subscriptions[0].Type != "atom" {
		t.Errorf(`Incorrect subscription type: %q`, subscriptions[0].Type)
	}
}

func TestParseWebPageWithJSONFeed(t *testing.T) {
	htmlPage := `
	<!doctype html>
	<html>
		<head>
			<link href="http://example.org/feed.json" rel="alternate" type="application/feed+json" title="Some Title">
		</head>
		<body>
		</body>
	</html>`

	subscriptions, err := NewSubscriptionFinder(nil).FindSubscriptionsFromWebPage("http://example.org/", "text/html", strings.NewReader(htmlPage))
	if err != nil {
		t.Fatalf(`Parsing a correctly formatted HTML page should not return any error: %v`, err)
	}

	if len(subscriptions) != 1 {
		t.Fatal(`Incorrect number of subscriptions returned`)
	}

	if subscriptions[0].Title != "Some Title" {
		t.Errorf(`Incorrect subscription title: %q`, subscriptions[0].Title)
	}

	if subscriptions[0].URL != "http://example.org/feed.json" {
		t.Errorf(`Incorrect subscription URL: %q`, subscriptions[0].URL)
	}

	if subscriptions[0].Type != "json" {
		t.Errorf(`Incorrect subscription type: %q`, subscriptions[0].Type)
	}
}

func TestParseWebPageWithOldJSONFeedMimeType(t *testing.T) {
	htmlPage := `
	<!doctype html>
	<html>
		<head>
			<link href="http://example.org/feed.json" rel="alternate" type="application/json" title="Some Title">
		</head>
		<body>
		</body>
	</html>`

	subscriptions, err := NewSubscriptionFinder(nil).FindSubscriptionsFromWebPage("http://example.org/", "text/html", strings.NewReader(htmlPage))
	if err != nil {
		t.Fatalf(`Parsing a correctly formatted HTML page should not return any error: %v`, err)
	}

	if len(subscriptions) != 1 {
		t.Fatal(`Incorrect number of subscriptions returned`)
	}

	if subscriptions[0].Title != "Some Title" {
		t.Errorf(`Incorrect subscription title: %q`, subscriptions[0].Title)
	}

	if subscriptions[0].URL != "http://example.org/feed.json" {
		t.Errorf(`Incorrect subscription URL: %q`, subscriptions[0].URL)
	}

	if subscriptions[0].Type != "json" {
		t.Errorf(`Incorrect subscription type: %q`, subscriptions[0].Type)
	}
}

func TestParseWebPageWithRelativeFeedURL(t *testing.T) {
	htmlPage := `
	<!doctype html>
	<html>
		<head>
			<link href="/feed.json" rel="alternate" type="application/feed+json" title="Some Title">
		</head>
		<body>
		</body>
	</html>`

	subscriptions, err := NewSubscriptionFinder(nil).FindSubscriptionsFromWebPage("http://example.org/", "text/html", strings.NewReader(htmlPage))
	if err != nil {
		t.Fatalf(`Parsing a correctly formatted HTML page should not return any error: %v`, err)
	}

	if len(subscriptions) != 1 {
		t.Fatal(`Incorrect number of subscriptions returned`)
	}

	if subscriptions[0].Title != "Some Title" {
		t.Errorf(`Incorrect subscription title: %q`, subscriptions[0].Title)
	}

	if subscriptions[0].URL != "http://example.org/feed.json" {
		t.Errorf(`Incorrect subscription URL: %q`, subscriptions[0].URL)
	}

	if subscriptions[0].Type != "json" {
		t.Errorf(`Incorrect subscription type: %q`, subscriptions[0].Type)
	}
}

func TestParseWebPageWithEmptyTitle(t *testing.T) {
	htmlPage := `
	<!doctype html>
	<html>
		<head>
			<link href="/feed.json" rel="alternate" type="application/feed+json">
		</head>
		<body>
		</body>
	</html>`

	subscriptions, err := NewSubscriptionFinder(nil).FindSubscriptionsFromWebPage("http://example.org/", "text/html", strings.NewReader(htmlPage))
	if err != nil {
		t.Fatalf(`Parsing a correctly formatted HTML page should not return any error: %v`, err)
	}

	if len(subscriptions) != 1 {
		t.Fatal(`Incorrect number of subscriptions returned`)
	}

	if subscriptions[0].Title != "http://example.org/feed.json" {
		t.Errorf(`Incorrect subscription title: %q`, subscriptions[0].Title)
	}

	if subscriptions[0].URL != "http://example.org/feed.json" {
		t.Errorf(`Incorrect subscription URL: %q`, subscriptions[0].URL)
	}

	if subscriptions[0].Type != "json" {
		t.Errorf(`Incorrect subscription type: %q`, subscriptions[0].Type)
	}
}

func TestParseWebPageWithMultipleFeeds(t *testing.T) {
	htmlPage := `
	<!doctype html>
	<html>
		<head>
			<link href="http://example.org/atom.xml" rel="alternate" type="application/atom+xml" title="Atom Feed">
			<link href="http://example.org/feed.json" rel="alternate" type="application/feed+json" title="JSON Feed">
		</head>
		<body>
		</body>
	</html>`

	subscriptions, err := NewSubscriptionFinder(nil).FindSubscriptionsFromWebPage("http://example.org/", "text/html", strings.NewReader(htmlPage))
	if err != nil {
		t.Fatalf(`Parsing a correctly formatted HTML page should not return any error: %v`, err)
	}

	if len(subscriptions) != 2 {
		t.Fatal(`Incorrect number of subscriptions returned`)
	}
}

func TestParseWebPageWithDuplicatedFeeds(t *testing.T) {
	htmlPage := `
	<!doctype html>
	<html>
		<head>
			<link href="http://example.org/feed.xml" rel="alternate" type="application/rss+xml" title="Feed A">
			<link href="http://example.org/feed.xml" rel="alternate" type="application/rss+xml" title="Feed B">
		</head>
		<body>
		</body>
	</html>`

	subscriptions, err := NewSubscriptionFinder(nil).FindSubscriptionsFromWebPage("http://example.org/", "text/html", strings.NewReader(htmlPage))
	if err != nil {
		t.Fatalf(`Parsing a correctly formatted HTML page should not return any error: %v`, err)
	}

	if len(subscriptions) != 1 {
		t.Fatal(`Incorrect number of subscriptions returned`)
	}

	if subscriptions[0].Title != "Feed A" {
		t.Errorf(`Incorrect subscription title: %q`, subscriptions[0].Title)
	}

	if subscriptions[0].URL != "http://example.org/feed.xml" {
		t.Errorf(`Incorrect subscription URL: %q`, subscriptions[0].URL)
	}

	if subscriptions[0].Type != "rss" {
		t.Errorf(`Incorrect subscription type: %q`, subscriptions[0].Type)
	}
}

func TestParseWebPageWithEmptyFeedURL(t *testing.T) {
	htmlPage := `
	<!doctype html>
	<html>
		<head>
			<link href rel="alternate" type="application/feed+json" title="Some Title">
		</head>
		<body>
		</body>
	</html>`

	subscriptions, err := NewSubscriptionFinder(nil).FindSubscriptionsFromWebPage("http://example.org/", "text/html", strings.NewReader(htmlPage))
	if err != nil {
		t.Fatalf(`Parsing a correctly formatted HTML page should not return any error: %v`, err)
	}

	if len(subscriptions) != 0 {
		t.Fatal(`Incorrect number of subscriptions returned`)
	}
}

func TestParseWebPageWithNoHref(t *testing.T) {
	htmlPage := `
	<!doctype html>
	<html>
		<head>
			<link rel="alternate" type="application/feed+json" title="Some Title">
		</head>
		<body>
		</body>
	</html>`

	subscriptions, err := NewSubscriptionFinder(nil).FindSubscriptionsFromWebPage("http://example.org/", "text/html", strings.NewReader(htmlPage))
	if err != nil {
		t.Fatalf(`Parsing a correctly formatted HTML page should not return any error: %v`, err)
	}

	if len(subscriptions) != 0 {
		t.Fatal(`Incorrect number of subscriptions returned`)
	}
}
