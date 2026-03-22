package feed

import "github.com/tbressel/daily-games-api/pkg"

// defaultFeeds is the list of all configured RSS feed sources.
// These are static for now and will eventually be loaded from a database.
var defaultFeeds = []pkg.Feed{
	// --- 🔥 NEXT-GEN & PRO ---
	// -- FR
	{ID: "jeux-video", Name: "Jeux Video.com", URL: "https://www.jeuxvideo.com/rss/rss-news.xml", Category: "nextgen", Lang: "fr", Description: "Le leader de l'actualité vidéoludique en France : news, tests et vidéos.", IsActive: true},
	{ID: "game-kult", Name: "Game Kult", URL: "https://www.gamekult.com/feed.xml", Category: "nextgen", Lang: "fr", Description: "Retrouvez toute l'actualité en temps réel et les tests des derniers jeux vidéo fraîchement sortis, servis par la rédaction Gamekult ! Découvrez toutes nos émissions, nos guides d'achat pour choisir le meilleur matériel ainsi que nos soluces et astuces pour profiter au maximum de vos jeux préférés.", IsActive: true},
	{ID: "xbox-gamer", Name: "Xbox Gamer", URL: "https://www.xbox-gamer.net/rss.php", Category: "nextgen", Lang: "fr", Description: "Toute l'actualité Xbox Series X|S, Xbox Game Pass, Xbox One (S/X), Xbox 360, Xbox, news, test, preview", IsActive: true},
	{ID: "canard-pc", Name: "Canard PC", URL: "https://www.canardpc.com/feed/", Category: "nextgen", Lang: "fr", Description: "Actualité et critique indépendante du jeu vidéo par une rédaction historique.", IsActive: true},
	{ID: "p-nintendo", Name: "Puissance Nintendo", URL: "https://feeds.feedburner.com/pn-majs", Category: "nextgen", Lang: "fr", Description: "Actualité des consoles Switch, 3DS, Wii U", IsActive: true},
	{ID: "origami", Name: "Origami", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UC73Te9mCdJcHxr_7hPBeGBw", Category: "nextgen", Lang: "fr", Description: "Analyses, débats et actualité du jeu vidéo par un média indépendant d'experts.", IsActive: true},
	{ID: "bistro", Name: "Le bistro du jeu video", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCRXcryyD7dzNQzd0Zkbj3ug", Category: "nextgen", Lang: "fr", Description: "Émission sur l'actualité du jeu vidéo rétro et moderne.", IsActive: true},
	{ID: "sega-mag", Name: "Sega-Mag", URL: "https://www.sega-mag.com/rss.xml", Category: "nextgen", Lang: "fr", Description: "Toute l'actualité SEGA : news, tests et dossiers par la rédaction de Sega-Mag.", IsActive: true},
	{ID: "factor-news", Name: "Factor News", URL: "https://www.factornews.com/rss.xml", Category: "nextgen", Lang: "fr", Description: "Jouer moins cher et légalement - Nintendo Switch Online, Epic Games Store - News - Factornews", IsActive: true},
	{ID: "jeux-on-line", Name: "Jeux On Line", URL: "https://www.jeuxonline.info/rss/actualites/rss.xml", Category: "nextgen", Lang: "fr", Description: "JeuxOnLine est un site d'information et communautaire sur le jeu vidéo, notamment MMO (massivement multi-joueurs) et sur les cultures populaires.", IsActive: true},
	// -- EN
	{ID: "vgc", Name: "VGC", URL: "https://www.videogameschronicle.com/feed/", Category: "nextgen", Lang: "en", Description: "Unbiased news, reporting and features from the global gaming industry.", IsActive: true},
	{ID: "gameradar", Name: "GamesRadar+", URL: "https://www.gamesradar.com/feeds.xml", Category: "nextgen", Lang: "en", Description: "Breaking news, reviews and features from the world of gaming.", IsActive: true},

	// --- 🕹️ RETROGAMING ---
	// -- FR
	{ID: "recalbox", Name: "Recalbox", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCfcqrtnHwB84YQlVN75PRfQ", Category: "retrogaming", Lang: "fr", Description: "Le portail de l'émulation francophone de référence pour les nostalgique.", IsActive: true},
	{ID: "mo5", Name: "Association MO5", URL: "https://mag.mo5.com/feed/", Category: "retrogaming", Lang: "fr", Description: "L'actualité de la préservation du patrimoine numérique par l'association MO5.", IsActive: true},
	{ID: "rom-game", Name: "Rom Game", URL: "https://www.rom-game.fr/rss/rss.rss", Category: "retrogaming", Lang: "fr", Description: "Toute l'actualité du retrogaming, du homebrew et des éditions physiques.", IsActive: true},
	{ID: "abandonware", Name: "Abandonware France", URL: "https://www.abandonware-france.org/rss/abandonware/", Category: "retrogaming", Lang: "fr", Description: "L'histoire du jeu vidéo sur PC à travers les titres du patrimoine.", IsActive: true},
	{ID: "backin", Name: "Back in Toys TV", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCy2_IMeOgTZ-NYsmGB2_qUw", Category: "retrogaming", Lang: "fr", Description: "Actualité geek, jouets vintage et jeux vidéo rétro.", IsActive: true},
	{ID: "conkerax", Name: "Conkerax", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCnyXbcCPqBOf_qXjyNF7dlg", Category: "retrogaming", Lang: "fr", Description: "Bienvenue à tous sur la chaîne d'un passionné de jeux vidéo, collectionneur, player, avec un petit (gros !) faible pour la Nintendo Gamecube.", IsActive: true},
	{ID: "oldschoolisbeautifull", Name: "Old School Is Beautifull", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCMUvAvf9UMdyq_VCog9CZVA", Category: "retrogaming", Lang: "fr", Description: "Plonge dans les histoires des consoles oubliées, des jeux cultes et des nouveaux jeux pour ces anciennes machines !", IsActive: true},
	{ID: "passionjeuxvideotv", Name: "Passion Jeux Vidéo TV", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UC9ci-PpcfhI9uN3qfqHjqlQ", Category: "retrogaming", Lang: "fr", Description: "Prêt à remonter le temps ? 🕹️ Viens (re)plonger dans l'âge d'or du jeu vidéo !", IsActive: true},
	{ID: "retrogamer", Name: "Retrogamer", URL: "https://retrogamer.cc/feed/", Category: "retrogaming", Lang: "fr", Description: "Site indépendant dédié au rétrogaming et à l'actualité tech. Pas de publicité, juste du contenu pour les passionnés", IsActive: true},
	{ID: "jdg", Name: "Le Joueur du Grenier", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UC_yP2DpIgs5Y1uWC0T03Chw", Category: "retrogaming", Lang: "fr", Description: "Test de jeux à la con !", IsActive: true},
	// -- EN
	// {ID: "vintageisthenewold", Name: "Vintage is the New Old", URL: "https://www.vintageisthenewold.com/feed/", Category: "retrogaming", Lang: "en", Description: "News about classic systems, computing history and homebrew releases.", IsActive: true},
	{ID: "reddit-retro", Name: "Reddit Retrogaming", URL: "https://www.reddit.com/r/retrogaming/.rss", Category: "retrogaming", Lang: "en", Description: "The premier Reddit community for classic gaming enthusiasts.", IsActive: true},

	// --- 💎 INDIE & DÉCOUVERTES ---
	// -- FR
	{ID: "indiemag", Name: "IndieMag", URL: "https://www.indiemag.fr/feed/rss.xml", Category: "indie", Lang: "fr", Description: "Le portail francophone spécialisé dans l'actualité du jeu indépendant.", IsActive: true},
	{ID: "at0mium", Name: "At0mium", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCI0LNmSlhS-H9mGNPWM8gzQ", Category: "indie", Lang: "fr", Description: "Chroniques et découvertes quotidiennes de pépites de la scène indépendante.", IsActive: true},
	// -- EN
	{ID: "indie-retro-news", Name: "Indie Retro News", URL: "http://www.indieretronews.com/feeds/posts/default?alt=rss", Category: "indie", Lang: "en", Description: "Focus on indie games and retro-styled modern titles across all platforms.", IsActive: true},
	{ID: "indiedb", Name: "Indie DB", URL: "https://rss.indiedb.com/articles/feed/rss.xml", Category: "indie", Lang: "en", Description: "Comprehensive database and latest news about independent video games.", IsActive: true},
	{ID: "itchio", Name: "Itch.io News", URL: "https://itch.io/feed/new.xml", Category: "indie", Lang: "en", Description: "Showcasing the best independent and experimental games on Itch.io.", IsActive: true},
	// {ID: "itchio", Name: "Itch.io News", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UC64fwl47Wrc6VJcskll7vsA", Category: "indie", Lang: "en", Description: "Showcasing the best independent and experimental games on Itch.io.", IsActive: true},

	// --- 🛠️ HOMEBREW & TECH ---
	// -- EN
	{ID: "wololo", Name: "Wololo.net", URL: "https://wololo.net/feed/", Category: "homebrew", Lang: "en", Description: "The home of the console homebrew, hacking and customization scene.", IsActive: true},
	{ID: "retro-rgb", Name: "Retro RGB", URL: "https://www.retrorgb.com/feed", Category: "homebrew", Lang: "en", Description: "High-quality video output, hardware mods and technical guides for retro consoles.", IsActive: true},
	{ID: "gbatemp", Name: "GBAtemp", URL: "https://gbatemp.net/feed/news", Category: "homebrew", Lang: "en", Description: "Independent gaming community focused on console hacking and homebrew.", IsActive: true},
	{ID: "scene-world", Name: "Scene World", URL: "https://feeds.feedburner.com/sceneworldpodcast", Category: "homebrew", Lang: "en", Description: "The digital magazine and podcast covering the international computing scene.", IsActive: true},
	// --FR

	// --- 📦 MACHINES ---
	// -- EN
	{ID: "atarilegend", Name: "Atari Legend", URL: "https://www.atarilegend.com/feed", Category: "computing", Lang: "en", Description: "Information, reviews and comments about Atari ST games, interviews of famous Atari ST game developers, contribute missing information to the database.", IsActive: true},
	// {ID: "atariage", Name: "Atariage", URL: "https://www.atariage.com/news/rss.php", Category: "computing", Lang: "en", Description: "The premier website for everything related to Atari systems and homebrew.", IsActive: true},
	// {ID: "octoate", Name: "Octoate", URL: "https://www.octoate.de/feed/", Category: "computing", Lang: "en", Description: "The Amstrad CPC resource for news, articles and technical information.", IsActive: true},
	// -- FR
	{ID: "vretrocomputing", Name: "Vretro Computing", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCG4S3PerB8tmodN-tpGQthA", Category: "computing", Lang: "fr", Description: "Utiliser de nos jours les ordinateurs populaires des années 80. Atari ST, Amiga... mais pas seulement. Modifications hardware et émulateurs. Développer pour ces machines avec les outils modernes.", IsActive: true},
	{ID: "amstrad-eu", Name: "Amstrad.eu", URL: "https://amstrad.eu/feed/", Category: "computing", Lang: "fr", Description: "Le portail communautaire francophone de référence pour l'Amstrad CPC.", IsActive: true},
	{ID: "amiga-impact", Name: "Amiga Impact", URL: "https://www.amigaimpact.org/feed/", Category: "computing", Lang: "fr", Description: "Actualités et ressources pour les utilisateurs d'AmigaOS, AROS et MorphOS", IsActive: true},
	{ID: "asmtariste", Name: "ASMtariSTe", URL: "https://www.asmtariste.fr/feed", Category: "computing", Lang: "fr", Description: "Apprendre l'assembleur 68000 sur Atari ST/STe", IsActive: true},
	// {ID: "64nops", Name: "64nops", URL: "https://64nops.wordpress.com/feed/", Category: "computing", Lang: "fr", Description: "Le blog de la scène Amstrad CPC, programmation et nouveautés.", IsActive: true},
	// {ID: "ucpm", Name: "ùCPM Blog", URL: "https://ucpmblog.ovh/index.php/feed/", Category: "computing", Lang: "fr", Description: "Blog dédié à l'actualité et au développement sur Amstrad CPC.", IsActive: true},
	// {ID: "amstariga", Name: "Amstariga", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCZnozTrLo1y-4VSdTYfC5dQ", Category: "computing", Lang: "fr", Description: "Chaîne YouTube dédiée aux machines 8 bits et à l'Amstrad CPC.", IsActive: true},
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

// GetByLang returns all active feeds matching the given language code ("fr" or "en").
func GetByLang(lang string) []pkg.Feed {
	var result []pkg.Feed
	for _, f := range defaultFeeds {
		if f.IsActive && f.Lang == lang {
			result = append(result, f)
		}
	}
	return result
}

// GetByCategoryAndLang returns all active feeds matching both category and language.
func GetByCategoryAndLang(category, lang string) []pkg.Feed {
	var result []pkg.Feed
	for _, f := range defaultFeeds {
		if f.IsActive && f.Category == category && f.Lang == lang {
			result = append(result, f)
		}
	}
	return result
}
