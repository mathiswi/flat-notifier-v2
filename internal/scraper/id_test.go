package scraper

import "testing"

func TestGetIDFromFlatLink(t *testing.T) {
	tests := []struct {
		name string
		link string
		want string
	}{
		{
			name: "full URL path",
			link: "/s-anzeige/erstbezug-schoene-helle-50-qm-dg-whg-balkon-ol-dietrichsfeld/1943252790-203-3112",
			want: "1943252790",
		},
		{
			name: "full URL with domain",
			link: "https://www.kleinanzeigen.de/s-anzeige/some-title/1955972992-203-3112",
			want: "1955972992",
		},
		{
			name: "different ID format",
			link: "/s-anzeige/title/3305804053-203-3124",
			want: "3305804053",
		},
		{
			name: "empty string",
			link: "",
			want: "",
		},
		{
			name: "trailing slash",
			link: "/s-anzeige/title/1234567890-203-3112/",
			want: "1234567890",
		},
		{
			name: "no dashes in last segment",
			link: "/s-anzeige/title/9876543210",
			want: "9876543210",
		},
		{
			name: "non-numeric ID",
			link: "/s-anzeige/title/abc-203-3112",
			want: "",
		},
		{
			name: "only slashes",
			link: "///",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetIDFromFlatLink(tt.link)
			if got != tt.want {
				t.Errorf("GetIDFromFlatLink(%q) = %q, want %q", tt.link, got, tt.want)
			}
		})
	}
}
