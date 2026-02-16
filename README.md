# 🍄 forage

cli tool to find and download similar music from a spotify track

<img src="demo.gif" width="600">

## about

paste a spotify track url, and forage finds similar songs on last.fm and downloads them as mp3s with metadata and album art

## why

i wanted a way to discover and download music without relying on streaming services. forage lets you build a local library based on what you already like

## install

**using pre-built binaries (recommended):**

1. download the latest release for your platform from [releases](https://github.com/jackmayhew/forage/releases)
   - **macos:** use `darwin-arm64` for M1/M2/M3 or `darwin-amd64` for Intel
   - **linux:** use `linux-amd64`
   - **windows:** use `windows-amd64.exe`

2. install the binary:

**macos/linux:**
```bash
# navigate to your downloads folder
cd ~/Downloads

# make the binary executable
chmod +x forage-darwin-arm64  # or forage-linux-amd64

# move to system path
sudo mv forage-darwin-arm64 /usr/local/bin/forage

# macos only: remove quarantine flag
sudo xattr -d com.apple.quarantine /usr/local/bin/forage
```

**windows:**
- move `forage-windows-amd64.exe` to a folder in your PATH
- or run it directly from the download location

3. run setup to create config template:
```bash
forage --setup
```
4. edit the config file with your API credentials:
```bash
open -e ~/.config/forage/config.yaml  # macos
nano ~/.config/forage/config.yaml     # linux
notepad %USERPROFILE%\.config\forage\config.yaml  # windows
```
   - get spotify credentials: https://developer.spotify.com/dashboard
   - get lastfm key: https://www.last.fm/api/account/create

**building from source:**

you need:
- go 1.21+
- yt-dlp: `brew install yt-dlp`
- ffmpeg: `brew install ffmpeg`
```bash
git clone https://github.com/jackmayhew/forage.git
cd forage
go build
```

for development, create a `.env` file:
```
SPOTIFY_CLIENT_ID=your_id
SPOTIFY_CLIENT_SECRET=your_secret
LASTFM_API_KEY=your_key
```

## usage

basic:
```bash
forage "https://open.spotify.com/track/2Ud3deeqLAG988pfW0Kwcl?si=e1f747637ed241b6"
```

flags:
- `--setup` - create config file template
- `--count N` - number of similar tracks (max: 50, default: 10)
- `--output DIR` - where to save files (default: ./downloads)
- `--quiet` - minimal output

**note:** 
- when building from source: use `./forage` 
- for development: use `go run .`

## how it works

1. gets track info from spotify
2. finds similar tracks on last.fm
3. downloads from youtube as mp3
4. adds metadata and album art from spotify

files are named `Artist - Track.mp3` and automatically skips download if already downloaded

## notes

- requires free spotify and last.fm api keys
- downloads audio from youtube (quality varies)
- downloads can be slow (yt-dlp searches youtube for each track)

## roadmap

**features:**
- similar artist search
- spotify playlist support
- interactive track selection
- genre filtering
- exclude artist flag

**technical:**
- concurrent downloads with worker pool
- retry logic for failed downloads
- custom metadata template options
- improved audio source matching (spotdl integration)

## disclaimer

this tool is for educational purposes only. please support artists by purchasing their music or using official streaming services. the developers of forage are not responsible for any misuse of this tool