const axios = require('axios');
const fs = require('fs');
const path = require('path');
const ConsoleWindow = require("node-hide-console-window");


function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

(async () => {
  ConsoleWindow.hideConsole();
  const event = {
    game_event: {
      level: "",
      game: "TIDAL",
      event: "MEDIA_PLAYBACK",
      data: {
        frame: {
          album: "Stick Season (We'll All Be Here Forever)",
          artist: "Noah Kahan",
          duration: 205,
          imageUrl: "https://resources.tidal.com/images/95c12523/7df7/4b46/85e3/08e4ec52dc54/1280x1280.jpg",
          state: "playing",
          time: 93,
          title: "",
          url: "https://tidal.com/browse/track/297959226"
        },
        value: 45
      }
    } 
  }
  const coreProps = await fs.readFileSync(path.join(process.env.PROGRAMDATA, "SteelSeries", "SteelSeries Engine 3", "coreProps.json"), "utf8");
  let url = JSON.parse(coreProps).address;
  while (true) {
    let nowPlaying = await axios.get('http://localhost:10769/currentPlayingSong')
    if (nowPlaying == null) {
    } else {
      if (nowPlaying.data.info == null || nowPlaying.data.info.status == "loading") {
      } else {
        event.game_event.data.frame.title = nowPlaying.data.info.name;
        event.game_event.data.frame.artist = nowPlaying.data.info.artistName;
        event.game_event.data.frame.album = nowPlaying.data.info.albumName;
        event.game_event.data.frame.state = nowPlaying.data.info.status;
        event.game_event.data.frame.duration = nowPlaying.data.info.durationInMillis;
        event.game_event.data.frame.time = nowPlaying.data.info.currentPlaybackTime * 1000;
        await sleep(25);
        await axios.post(`http://${url}/game_event`, event.game_event);
      }
    }
  }
})();
