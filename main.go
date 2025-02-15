package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/godbus/dbus/v5"
)

const (
	spotifyBusName     = "org.mpris.MediaPlayer2.spotify"
	spotifyObjPath     = "/org/mpris/MediaPlayer2"
	spotifyIface       = "org.freedesktop.DBus.Properties"
	metadataProp       = "org.mpris.MediaPlayer2.Player.Metadata"
	playbackStatusProp = "org.mpris.MediaPlayer2.Player.PlaybackStatus"
	outputFile         = "spotify_now_playing.txt"
	coverFile          = "spotify_cover_url.txt"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	conn, err := dbus.SessionBus()
	if err != nil {
		log.Fatalf("Failed to connect to D-Bus: %v", err)
	}
	defer conn.Close()

	signalChan := make(chan *dbus.Signal, 10)
	conn.Signal(signalChan)

	matchRule := fmt.Sprintf("type='signal',interface='%s'", spotifyIface)
	err = conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, matchRule).Err
	if err != nil {
		log.Fatalf("Failed to add D-Bus match rule: %v", err)
	}

	log.Printf("Writing song metadata to '%s'...", filepath.Join(cwd, outputFile))
	log.Printf("Writing song cover URL to '%s'...", filepath.Join(cwd, coverFile))

	log.Println("Listening for Spotify track changes...")
	var lastTrack string
	for signal := range signalChan {
		if signal.Path == dbus.ObjectPath(spotifyObjPath) && signal.Name == "org.freedesktop.DBus.Properties.PropertiesChanged" {
			playbackStatus := getSpotifyPlaybackStatus(conn)
			if playbackStatus == "Paused" {
				clearFile(outputFile)
				clearFile(coverFile)
				lastTrack = ""
				continue
			}

			metadata, coverURL := getSpotifyMetadata(conn)
			if metadata != "" && metadata != lastTrack {
				lastTrack = metadata
				writeToFile(metadata, outputFile)
				writeToFile(coverURL, coverFile)
			}
		}
	}
}

func getSpotifyMetadata(conn *dbus.Conn) (string, string) {
	obj := conn.Object(spotifyBusName, dbus.ObjectPath(spotifyObjPath))
	variant, err := obj.GetProperty(metadataProp)
	if err != nil {
		log.Printf("Failed to get Spotify metadata: %v", err)
		return "", ""
	}

	metadata, ok := variant.Value().(map[string]dbus.Variant)
	if !ok {
		log.Println("Failed to parse metadata")
		return "", ""
	}
	title, _ := metadata["xesam:title"].Value().(string)
	artist, _ := metadata["xesam:artist"].Value().([]string)
	coverURL, _ := metadata["mpris:artUrl"].Value().(string)

	if title == "" {
		return "", ""
	}
	return fmt.Sprintf("%s - %s", title, artist), coverURL
}

func getSpotifyPlaybackStatus(conn *dbus.Conn) string {
	obj := conn.Object(spotifyBusName, dbus.ObjectPath(spotifyObjPath))
	variant, err := obj.GetProperty(playbackStatusProp)
	if err != nil {
		log.Printf("Failed to get Spotify playback status: %v", err)
		return ""
	}
	status, _ := variant.Value().(string)
	return status
}

func writeToFile(content string, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Printf("Failed to write file %s: %v", filename, err)
		return
	}
	defer file.Close()
	file.WriteString(content + "\n")
	log.Printf("Updated file %s: %s", filename, content)
}

func clearFile(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Printf("Failed to clear file %s: %v", filename, err)
		return
	}
	defer file.Close()
	log.Printf("File %s cleared due to playback pause", filename)
}
