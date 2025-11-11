package voice

import (
	"errors"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

var connections = make(map[string]*discordgo.VoiceConnection)

type Song struct {
	URL   string
	Title string
}

var queues = make(map[string][]Song)
var playing = make(map[string]bool)
var stopChans = make(map[string]chan bool)
var loopFlags = make(map[string]bool)

func JoinVoiceChannel(s *discordgo.Session, guildID, channelID string) error {
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}
	connections[guildID] = vc
	slog.Info("Joined voice channel", "guild", guildID, "channel", channelID)
	return nil
}

func LeaveVoiceChannel(guildID string) {
	StopPlaying(guildID)
	if vc, ok := connections[guildID]; ok {
		vc.Disconnect()
		delete(connections, guildID)
		slog.Info("Left voice channel", "guild", guildID)
	}
}

func AddToQueue(guildID, url, title string) {
	queues[guildID] = append(queues[guildID], Song{URL: url, Title: title})
}

func RemoveFromQueue(guildID string, index int) error {
	q := queues[guildID]
	if index < 0 || index >= len(q) {
		return errors.New("invalid index")
	}
	queues[guildID] = append(q[:index], q[index+1:]...)
	return nil
}

func ViewQueue(guildID string) []Song {
	return queues[guildID]
}

func ShuffleQueue(guildID string) {
	q := queues[guildID]
	rand.Shuffle(len(q), func(i, j int) { q[i], q[j] = q[j], q[i] })
	queues[guildID] = q
}

func PlayAudio(guildID, filename string) {
	vc, ok := connections[guildID]
	if !ok {
		slog.Error("No voice connection for guild", "guild", guildID)
		return
	}
	stopChan := make(chan bool)
	stopChans[guildID] = stopChan
	dgvoice.PlayAudioFile(vc, filename, stopChan)
	delete(stopChans, guildID)
}

func Skip(guildID string) {
	if stopChan, ok := stopChans[guildID]; ok {
		stopChan <- true
	}
}

func IsPlaying(guildID string) bool {
	return playing[guildID]
}

func StartPlaying(guildID string) {
	if playing[guildID] {
		return
	}
	playing[guildID] = true
	go playQueue(guildID)
}

func playQueue(guildID string) {
	for playing[guildID] && len(queues[guildID]) > 0 {
		song := queues[guildID][0]
		queues[guildID] = queues[guildID][1:]
		if strings.HasSuffix(song.URL, ".mp3") {
			resp, err := http.Get(song.URL)
			if err != nil {
				slog.Error("Failed to get audio", "error", err)
				continue
			}
			defer resp.Body.Close()
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				slog.Error("Failed to read response", "error", err)
				continue
			}
			tempFile, err := os.CreateTemp("", "*.mp3")
			if err != nil {
				slog.Error("Failed to create temp file", "error", err)
				continue
			}
			defer os.Remove(tempFile.Name())
			_, err = tempFile.Write(data)
			if err != nil {
				slog.Error("Failed to write temp file", "error", err)
				continue
			}
			tempFile.Close()
			PlayAudio(guildID, tempFile.Name())
		}
		if loopFlags[guildID] {
			queues[guildID] = append([]Song{song}, queues[guildID]...)
		}
	}
	playing[guildID] = false
}

func SetLoop(guildID string, loop bool) {
	loopFlags[guildID] = loop
}

func IsLooping(guildID string) bool {
	return loopFlags[guildID]
}

func StopPlaying(guildID string) {
	if stopChan, ok := stopChans[guildID]; ok {
		stopChan <- true
	}
	playing[guildID] = false
}
