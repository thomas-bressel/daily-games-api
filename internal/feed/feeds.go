package feed

import "github.com/tbressel/daily-games-api/pkg"

// defaultFeeds is the list of all configured RSS feed sources.
// These are static for now and will eventually be loaded from a database.
var defaultFeeds = []pkg.Feed{
	{ID: "indie-retro-news", Name: "Indie Retro News", URL: "http://www.indieretronews.com/feeds/posts/default?alt=rss", Category: "nextgen", Description: "Best gaming website for Indie and Retro Gaming News", IsActive: true},
	{ID: "reddit", Name: "Reddit", URL: "https://www.reddit.com/r/retrogaming/.rss", Category: "nextgen", Description: "Reddit retrogaming community", IsActive: true},
	{ID: "mo5", Name: "MO5", URL: "https://mag.mo5.com/feed/", Category: "nextgen", Description: "MO5 gaming magazine", IsActive: true},
	{ID: "ucpm", Name: "ùCPM Blog", URL: "https://ucpmblog.ovh/index.php/feed/", Category: "nextgen", Description: "ùCPM blog", IsActive: true},
	{ID: "64nops", Name: "64nops", URL: "https://64nops.wordpress.com/feed/", Category: "nextgen", Description: "64nops blog", IsActive: true},
	{ID: "bistro", Name: "Le bistro du jeu video", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCRXcryyD7dzNQzd0Zkbj3ug", Category: "retrogaming", Description: "YouTube retrogaming channel", IsActive: true},
	{ID: "amstariga", Name: "Amstariga", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCZnozTrLo1y-4VSdTYfC5dQ", Category: "retrogaming", Description: "Amstariga YouTube channel", IsActive: true},
	{ID: "backin", Name: "Back in Toys TV", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCy2_IMeOgTZ-NYsmGB2_qUw", Category: "retrogaming", Description: "Back in Toys TV YouTube channel", IsActive: true},
	{ID: "amstrad-eu", Name: "Amstrad.eu", URL: "https://amstrad.eu/feed/", Category: "amstrad-cpc", Description: "Amstrad community and news", IsActive: true},
	{ID: "octoate", Name: "Octoate.de", URL: "https://www.octoate.de/feed/", Category: "amstrad-cpc", Description: "Amstrad CPC tricks and news", IsActive: true},
	{ID: "atariage", Name: "Atariage", URL: "https://www.atariage.com/news/rss.php", Category: "atari-st", Description: "Atari ST news", IsActive: true},
	{ID: "itchio", Name: "Itch.io", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UC64fwl47Wrc6VJcskll7vsA", Category: "homebrew", Description: "Itch.io YouTube channel", IsActive: true},
	{ID: "retro-rgb", Name: "Retro RGB", URL: "https://www.retrorgb.com/feed", Category: "homebrew", Description: "Retro RGB news", IsActive: true},
	{ID: "indiedb", Name: "Indie DB", URL: "https://rss.indiedb.com/articles/feed/rss.xml", Category: "homebrew", Description: "Indie game news", IsActive: true},
	{ID: "gameradar", Name: "Game Radar", URL: "https://www.gamesradar.com/feeds.xml", Category: "retrogaming", Description: "Game Radar news", IsActive: true},
	{ID: "vintageisthenewold", Name: "Vintage is the New Old", URL: "https://www.vintageisthenewold.com/feed/", Category: "retrogaming", Description: "Vintage gaming news", IsActive: true},
	{ID: "scene-world", Name: "Scene World", URL: "https://feeds.feedburner.com/sceneworldpodcast", Category: "retrogaming", Description: "Scene World podcast", IsActive: true},
	{ID: "jeux-video", Name: "Jeux Video.com", URL: "https://www.jeuxvideo.com/rss/rss-news.xml", Category: "nextgen", Description: "French gaming news", IsActive: true},
	{ID: "abandonware", Name: "Abandonware France", URL: "https://www.abandonware-france.org/rss/abandonware/", Category: "retrogaming", Description: "Abandonware news", IsActive: true},
	{ID: "rom-game", Name: "Rom Game", URL: "https://www.rom-game.fr/rss/rss.rss", Category: "retrogaming", Description: "Rom Game news", IsActive: true},
}

// GetAll returns all configured feeds regardless of their active status.
func GetAll() []pkg.Feed {
	return defaultFeeds
}

// GetActive returns only feeds where IsActive is true.
func GetActive() []pkg.Feed {
	var active []pkg.Feed
	for _, f := range defaultFeeds {
		if f.IsActive {
			active = append(active, f)
		}
	}
	return active
}

// GetByID returns a single feed by its ID.
// The second return value is false if no feed was found.
func GetByID(id string) (pkg.Feed, bool) {
	for _, f := range defaultFeeds {
		if f.ID == id {
			return f, true
		}
	}
	return pkg.Feed{}, false
}

// GetByCategory returns all active feeds belonging to the given category.
func GetByCategory(category string) []pkg.Feed {
	var result []pkg.Feed
	for _, f := range defaultFeeds {
		if f.IsActive && f.Category == category {
			result = append(result, f)
		}
	}
	return result
}
