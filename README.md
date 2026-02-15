# forage

cli tool to find and download similar music from a spotify track

## what it does

paste a spotify track url, forage finds similar songs on last.fm and downloads them as mp3s with metadata and album art

## install

you need:
- go 1.21+
- yt-dlp: `brew install yt-dlp`
- ffmpeg: `brew install ffmpeg`
```bash
git clone https://github.com/jackmayhew/forage.git
cd forage
go build
```

## setup

create a `.env` file:
```
SPOTIFY_CLIENT_ID=your_id
SPOTIFY_CLIENT_SECRET=your_secret
LASTFM_API_KEY=your_key
```

get spotify credentials: https://developer.spotify.com/dashboard
get lastfm key: https://www.last.fm/api/account/create

## usage

basic:
```bash
./forage "https://open.spotify.com/track/..."
```

options:
```bash
./forage --count 5 --output ~/Music "spotify-url"
./forage --quiet "spotify-url"
```

flags:
- `--count N` - number of similar tracks (default: 10)
- `--output DIR` - where to save files (default: ./downloads)
- `--quiet` - minimal output

## how it works

1. gets track info from spotify
2. finds similar tracks on last.fm
3. downloads from youtube as mp3
4. adds metadata and album art

files are named `Artist - Track.mp3` and automatically skips download if already downloaded

## notes

- requires free spotify and last.fm api accounts
- downloads audio from youtube (quality varies)
- first download can be slow while yt-dlp searches