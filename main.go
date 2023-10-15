package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Event struct {
	game_event GameEvent
}

type GameEvent struct {
	Level string `json:"level"`
	Game  string `json:"game"`
	Event string `json:"event"`
	Data  Data   `json:"data"`
}
type Frame struct {
	Album    string  `json:"album"`
	Artist   string  `json:"artist"`
	Duration int32     `json:"duration"`
	ImageURL string  `json:"imageUrl"`
	State    string  `json:"state"`
	Time     float32 `json:"time"`
	Title    string  `json:"title"`
	URL      string  `json:"url"`
}
type Data struct {
	Frame Frame `json:"frame"`
	Value int   `json:"value"`
}

type CiderResponse struct {
	Info struct {
		AlbumName  string `json:"albumName"`
		ArtistName string `json:"artistName"`
		Artwork    struct {
			BgColor    string `json:"bgColor"`
			Height     int    `json:"height"`
			TextColor1 string `json:"textColor1"`
			TextColor2 string `json:"textColor2"`
			TextColor3 string `json:"textColor3"`
			URL        string `json:"url"`
			Width      int    `json:"width"`
		} `json:"artwork"`
		AudioLocale             string   `json:"audioLocale"`
		AudioTraits             []string `json:"audioTraits"`
		ComposerName            string   `json:"composerName"`
		CurrentPlaybackProgress int      `json:"currentPlaybackProgress"`
		CurrentPlaybackTime     float32  `json:"currentPlaybackTime"`
		DiscNumber              int32      `json:"discNumber"`
		DurationInMillis        int32      `json:"durationInMillis"`
		EndTime                 int64    `json:"endTime"`
		ExtendedAssetUrls       struct {
			EnhancedHls      string `json:"enhancedHls"`
			Lightweight      string `json:"lightweight"`
			LightweightPlus  string `json:"lightweightPlus"`
			Plus             string `json:"plus"`
			SuperLightweight string `json:"superLightweight"`
		} `json:"extendedAssetUrls"`
		GenreNames                []string `json:"genreNames"`
		HasCredits                bool     `json:"hasCredits"`
		HasLyrics                 bool     `json:"hasLyrics"`
		HasTimeSyncedLyrics       bool     `json:"hasTimeSyncedLyrics"`
		IsAppleDigitalMaster      bool     `json:"isAppleDigitalMaster"`
		IsMasteredForItunes       bool     `json:"isMasteredForItunes"`
		IsPlaying                 bool     `json:"isPlaying"`
		IsVocalAttenuationAllowed bool     `json:"isVocalAttenuationAllowed"`
		Isrc                      string   `json:"isrc"`
		Kind                      string   `json:"kind"`
		Name                      string   `json:"name"`
		PlayParams                struct {
			ID   string `json:"id"`
			Kind string `json:"kind"`
		} `json:"playParams"`
		Previews []struct {
			URL string `json:"url"`
		} `json:"previews"`
		ReleaseDate   time.Time `json:"releaseDate"`
		RemainingTime float64   `json:"remainingTime"`
		SongID        string    `json:"songId"`
		StartTime     float64   `json:"startTime"`
		Status        string    `json:"status"`
		TrackNumber   int       `json:"trackNumber"`
		URL           struct {
			AppleMusic string `json:"appleMusic"`
			Cider      string `json:"cider"`
			SongLink   string `json:"songLink"`
		} `json:"url"`
	} `json:"info"`
}

type Props struct {
	Address            string `json:"address"`
	EncryptedAddress   string `json:"encryptedAddress"`
	GgEncryptedAddress string `json:"ggEncryptedAddress"`
}

func main() {
	//Setup "Game"
	url := loadProps()
	//Setup initial event
	event := Event {
		game_event: GameEvent {
			Level: "info",
			Game: "TIDAL",
			Event: "MEDIA_PLAYBACK",
			Data: Data {
				Frame: Frame {
					Album: "ALBUM",
					Artist: "ARTIST",
					Duration: 0,
					ImageURL: "IMAGE_URL",
					State: "PAUSED",
					Time: 0,
					Title: "",
					URL: "URL",
				},
				Value: 0,
			},
		},
	}
	for {
		nowPlaying := comRpc("currentPlayingSong")
		event = updateEvent(event, nowPlaying)
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(event)
		req, _ := http.NewRequest("POST", url, buf);
		client := &http.Client{}
		res, e := client.Do(req)
		if e != nil {
			panic(e)
		}
		res.Body.Close()
		fmt.Println("Response Status:", res.Status)
	}
}

func updateEvent(event Event, data CiderResponse) Event {
	frame := Frame {
		Title: data.Info.Name,
		Artist: data.Info.ArtistName,
		Album: data.Info.AlbumName,
		State: data.Info.Status,
		Duration: data.Info.DurationInMillis,
		Time: data.Info.CurrentPlaybackTime * 1000,
	}
	event.game_event.Data.Frame = frame
	return event
}

func loadProps() string {
	path := os.Getenv("PROGRAMDATA") + "\\SteelSeries\\SteelSeries Engine 3\\coreProps.json";
	data, err := os.ReadFile(path);
	if err != nil {
		panic(err)
	}
	props := Props{};
	json.Unmarshal([]byte(data), &props)
	return props.Address + "/game_event";
}

func comRpc( request string) CiderResponse {
	data, err := http.Get("http://localhost:10769/" + request);
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(data.Body)
	if err != nil {
		panic(err)
	}
	dataString := CiderResponse{};
	json.Unmarshal([]byte(body), &dataString)
	return dataString;
}