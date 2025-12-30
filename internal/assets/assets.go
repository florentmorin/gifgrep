package assets

import _ "embed"

//go:embed giphy-32.png
var giphyIcon32PNG []byte

func GiphyIcon32PNG() []byte {
	return giphyIcon32PNG
}

