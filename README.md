# Spotigo - Spotify Now Playing Logger

Spotigo is a lightweight tool that logs the currently playing Spotify track using D-Bus and writes the information into
text files. These files can be used in OBS to display the song title and cover art in real-time.

## Installation

### Prerequisites

- A Linux system with systemd and D-Bus
- Go installed (`go` command available in PATH)

### Steps

```sh
make install
```

This will:

1. Compile the Spotigo binary.
2. Install it into `~/.local/bin/`.
3. Install the systemd user service.
4. Reload the systemd user instance.

To enable and run the service, run:
```sh
systemctl --user enable --now spotigo.service
```

To uninstall, run:

```sh
make uninstall
```

## Integration with OBS

1. Add a **Text (GDI+)** source in OBS.
2. Set the file path to: `/run/user/<USERID>/spotigo/spotify_now_playing.txt`
3. (Optional) The cover art URL will be available in a separate file
   - Integration into OBS is currently TBD. Feel free to submit a PR.
   - This file contains the URL: `/run/user/<USERID>/spotigo/spotify_cover_url.txt`

Spotigo will keep these files updated with the current song title and cover art.

## License

GPLv3
