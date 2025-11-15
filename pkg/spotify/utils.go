package spotify

import (
	"github.com/zmb3/spotify/v2"
)

func getArtists(artists []spotify.SimpleArtist) string {
	data := ""

	for i, artist := range artists {
		data += artist.Name

		if i < len(artists)-2 {
			data += ", "
		} else if i == len(artists)-2 {
			data += " & "
		}
	}

	return data
}
