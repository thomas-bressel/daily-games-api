package feed

import "github.com/tbressel/daily-games-api/pkg"

// defaultFeeds is the list of all configured RSS feed sources.
// These are static for now and will eventually be loaded from a database.
var defaultFeeds = []pkg.Feed{
	// --- 🔥 NEXT-GEN & PRO ---
	{ID: "jeux-video", Name: "Jeux Video.com", URL: "https://www.jeuxvideo.com/rss/rss-news.xml", Category: "nextgen", Description: "Le leader de l'actualité vidéoludique en France : news, tests et vidéos.", IsActive: true},
	{ID: "canard-pc", Name: "Canard PC", URL: "https://www.canardpc.com/feed/", Category: "nextgen", Description: "Actualité et critique indépendante du jeu vidéo par une rédaction historique.", IsActive: true},
	{ID: "origami", Name: "Origami", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCvX68V95zZpP8PZ6yU_7_pA", Category: "nextgen", Description: "Analyses, débats et actualité du jeu vidéo par un média indépendant d'experts.", IsActive: true},
	{ID: "vgc", Name: "VGC", URL: "https://www.videogameschronicle.com/feed/", Category: "nextgen", Description: "Unbiased news, reporting and features from the global gaming industry.", IsActive: true},
	{ID: "gameradar", Name: "GamesRadar+", URL: "https://www.gamesradar.com/feeds.xml", Category: "nextgen", Description: "Breaking news, reviews and features from the world of gaming.", IsActive: true},
	{ID: "bistro", Name: "Le bistro du jeu video", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCRXcryyD7dzNQzd0Zkbj3ug", Category: "nextgen", Description: "Émission sur l'actualité du jeu vidéo rétro et moderne.", IsActive: true},

	// --- 🕹️ RETROGAMING ---
	{ID: "mo5", Name: "Association MO5", URL: "https://mag.mo5.com/feed/", Category: "retrogaming", Description: "L'actualité de la préservation du patrimoine numérique par l'association MO5.", IsActive: true},
	{ID: "rom-game", Name: "Rom Game", URL: "https://www.rom-game.fr/rss/rss.rss", Category: "retrogaming", Description: "Toute l'actualité du retrogaming, du homebrew et des éditions physiques.", IsActive: true},
	{ID: "abandonware", Name: "Abandonware France", URL: "https://www.abandonware-france.org/rss/abandonware/", Category: "retrogaming", Description: "L'histoire du jeu vidéo sur PC à travers les titres du patrimoine.", IsActive: true},
	{ID: "backin", Name: "Back in Toys TV", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCy2_IMeOgTZ-NYsmGB2_qUw", Category: "retrogaming", Description: "Actualité geek, jouets vintage et jeux vidéo rétro.", IsActive: true},
	{ID: "reddit-retro", Name: "Reddit Retrogaming", URL: "https://www.reddit.com/r/retrogaming/.rss", Category: "retrogaming", Description: "The premier Reddit community for classic gaming enthusiasts.", IsActive: true},
	{ID: "vintageisthenewold", Name: "Vintage is the New Old", URL: "https://www.vintageisthenewold.com/feed/", Category: "retrogaming", Description: "News about classic systems, computing history and homebrew releases.", IsActive: true},

	// --- 💎 INDIE & DÉCOUVERTES ---
	{ID: "indiemag", Name: "IndieMag", URL: "https://www.indiemag.fr/rss", Category: "indie", Description: "Le portail francophone spécialisé dans l'actualité du jeu indépendant.", IsActive: true},
	{ID: "at0mium", Name: "At0mium", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UClS99rS86pS6vKueY_mS_7w", Category: "indie", Description: "Chroniques et découvertes quotidiennes de pépites de la scène indépendante.", IsActive: true},
	{ID: "indie-retro-news", Name: "Indie Retro News", URL: "http://www.indieretronews.com/feeds/posts/default?alt=rss", Category: "indie", Description: "Focus on indie games and retro-styled modern titles across all platforms.", IsActive: true},
	{ID: "indiedb", Name: "Indie DB", URL: "https://rss.indiedb.com/articles/feed/rss.xml", Category: "indie", Description: "Comprehensive database and latest news about independent video games.", IsActive: true},
	{ID: "itchio", Name: "Itch.io News", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UC64fwl47Wrc6VJcskll7vsA", Category: "indie", Description: "Showcasing the best independent and experimental games on Itch.io.", IsActive: true},

	// --- 🛠️ HOMEBREW & TECH ---
	{ID: "wololo", Name: "Wololo.net", URL: "https://wololo.net/feed/", Category: "homebrew", Description: "The home of the console homebrew, hacking and customization scene.", IsActive: true},
	{ID: "retro-rgb", Name: "Retro RGB", URL: "https://www.retrorgb.com/feed", Category: "homebrew", Description: "High-quality video output, hardware mods and technical guides for retro consoles.", IsActive: true},
	{ID: "gbatemp", Name: "GBAtemp", URL: "https://gbatemp.net/feed/news", Category: "homebrew", Description: "Independent gaming community focused on console hacking and homebrew.", IsActive: true},
	{ID: "scene-world", Name: "Scene World", URL: "https://feeds.feedburner.com/sceneworldpodcast", Category: "homebrew", Description: "The digital magazine and podcast covering the international computing scene.", IsActive: true},

	// --- 📦 MACHINES (NICHE) ---
	{ID: "amstrad-eu", Name: "Amstrad.eu", URL: "https://amstrad.eu/feed/", Category: "niche", Description: "Le portail communautaire francophone de référence pour l'Amstrad CPC.", IsActive: true},
	{ID: "amstariga", Name: "Amstariga", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCZnozTrLo1y-4VSdTYfC5dQ", Category: "niche", Description: "Chaîne YouTube dédiée aux machines 8 bits et à l'Amstrad CPC.", IsActive: true},
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
