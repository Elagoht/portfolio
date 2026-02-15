package handlers

type Link struct {
	Title string
	Href  string
}

type LinkGroup struct {
	Label string
	Links []Link
}

var SiteName = "Furkan Baytekin"
var SiteEmail = "furkan@baytekin.dev"

var SiteLinks = []Link{
	{Title: "GitHub", Href: "https://github.com/Elagoht"},
	{Title: "LinkedIn", Href: "https://linkedin.com/in/furkan-baytekin"},
	{Title: "Website", Href: "https://furkanbaytekin.dev"},
	{Title: "YouTube", Href: "https://youtube.com/@furkanbytekin"},
	{Title: "Itch.io", Href: "https://elagoht.itch.io"},
	{Title: "Udemy", Href: "https://www.udemy.com/user/furkan-baytekin/"},
	{Title: "X", Href: "https://x.com/furkanbytekin"},
	{Title: "Telegram", Href: "https://t.me/furkanbytekin"},
	{Title: "Reddit", Href: "https://reddit.com/u/furkanbytekin"},
	{Title: "Spotify", Href: "https://open.spotify.com/user/furkanbytekin"},
	{Title: "RSS", Href: "https://furkanbaytekin.dev/rss"},
}

var FooterGroups = []LinkGroup{
	{
		Label: "Direct Message",
		Links: []Link{
			{Title: "Telegram", Href: "https://t.me/furkanbytekin"},
			{Title: "Send Email", Href: "mailto:" + SiteEmail},
			{Title: "LinkedIn", Href: "https://linkedin.com/in/furkan-baytekin"},
		},
	},
	{
		Label: "Social Media",
		Links: []Link{
			{Title: "X", Href: "https://x.com/furkanbytekin"},
			{Title: "Reddit", Href: "https://reddit.com/u/furkanbytekin"},
			{Title: "Spotify", Href: "https://open.spotify.com/user/furkanbytekin"},
		},
	},
	{
		Label: "My Productive Hours",
		Links: []Link{
			{Title: "YouTube", Href: "https://youtube.com/@furkanbytekin"},
			{Title: "GitHub", Href: "https://github.com/Elagoht"},
			{Title: "RSS", Href: "https://furkanbaytekin.dev/rss"},
		},
	},
}

func BaseData(lang string) map[string]any {
	return map[string]any{
		"Lang":         lang,
		"Name":         SiteName,
		"Email":        SiteEmail,
		"Links":        SiteLinks,
		"FooterGroups": FooterGroups,
	}
}
