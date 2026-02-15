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
	{Title: "YouTube", Href: "https://youtube.com/@furkanbytekin"},
	{Title: "X", Href: "https://x.com/furkanbytekin"},
	{Title: "Telegram", Href: "https://t.me/furkanbytekin"},
	{Title: "Reddit", Href: "https://reddit.com/u/furkanbytekin"},
	{Title: "Spotify", Href: "https://open.spotify.com/user/furkanbytekin"},
	{Title: "Udemy", Href: "https://www.udemy.com/user/furkan-baytekin/"},
	{Title: "Itch.io", Href: "https://elagoht.itch.io"},
}

var FooterGroups = []LinkGroup{
	{
		Label: "Social",
		Links: []Link{
			{Title: "GitHub", Href: "https://github.com/Elagoht"},
			{Title: "LinkedIn", Href: "https://linkedin.com/in/furkan-baytekin"},
			{Title: "X", Href: "https://x.com/furkanbytekin"},
			{Title: "Reddit", Href: "https://reddit.com/u/furkanbytekin"},
		},
	},
	{
		Label: "Content",
		Links: []Link{
			{Title: "YouTube", Href: "https://youtube.com/@furkanbytekin"},
			{Title: "Udemy", Href: "https://www.udemy.com/user/furkan-baytekin/"},
			{Title: "Itch.io", Href: "https://elagoht.itch.io"},
		},
	},
	{
		Label: "Other",
		Links: []Link{
			{Title: "Telegram", Href: "https://t.me/furkanbytekin"},
			{Title: "Spotify", Href: "https://open.spotify.com/user/furkanbytekin"},
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
