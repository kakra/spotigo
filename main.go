package main

import (
	"fmt"
	"log"
	"os"
	"github.com/godbus/dbus/v5"
)

const (
	spotifyBusName  = "org.mpris.MediaPlayer2.spotify"
	spotifyObjPath  = "/org/mpris/MediaPlayer2"
	spotifyIface    = "org.freedesktop.DBus.Properties"
	metadataProp    = "org.mpris.MediaPlayer2.Player.Metadata"
	outputFile      = "spotify_now_playing.txt"
)

func main() {
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

	log.Println("Listening for Spotify track changes...")
	var lastTrack string
	for signal := range signalChan {
		if signal.Path == dbus.ObjectPath(spotifyObjPath) && signal.Name == "org.freedesktop.DBus.Properties.PropertiesChanged" {
			metadata := getSpotifyMetadata(conn)
			if metadata != "" && metadata != lastTrack {
				lastTrack = metadata
				writeToFile(metadata)
			}
		}
	}
}

func getSpotifyMetadata(conn *dbus.Conn) string {
	obj := conn.Object(spotifyBusName, dbus.ObjectPath(spotifyObjPath))
	variant, err := obj.GetProperty(metadataProp)
	if err != nil {
		log.Printf("Failed to get Spotify metadata: %v", err)
		return ""
	}

	metadata, ok := variant.Value().(map[string]dbus.Variant)
	if !ok {
		log.Println("Failed to parse metadata")
		return ""
	}
	title, _ := metadata["xesam:title"].Value().(string)
	artist, _ := metadata["xesam:artist"].Value().([]string)

	if title == "" {
		return ""
	}
	return fmt.Sprintf("%s - %s", title, artist)
}

func writeToFile(content string) {
	file, err := os.Create(outputFile)
	if err != nil {
		log.Printf("Failed to write file: %v", err)
		return
	}
	defer file.Close()
	file.WriteString(content + "\n")
	log.Printf("Updated file: %s", content)
}
