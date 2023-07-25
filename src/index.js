const fs = require('fs');
const path = require('path');


function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

(async () => {
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
  const coreProps = fs.readFileSync(path.join(process.env.PROGRAMDATA, "SteelSeries", "SteelSeries Engine 3", "coreProps.json"), "utf8");
  let url = JSON.parse(coreProps).address;
  while (true) {
    let response = await fetch('http://localhost:10769/currentPlayingSong');
    let nowPlaying = await response.json();
    if (nowPlaying == null) {
    } else {
      if (nowPlaying.info == null || nowPlaying.info.status == "loading") {
      } else {
        event.game_event.data.frame.title = nowPlaying.info.name;
        event.game_event.data.frame.artist = nowPlaying.info.artistName;
        event.game_event.data.frame.album = nowPlaying.info.albumName;
        event.game_event.data.frame.state = nowPlaying.info.status;
        event.game_event.data.frame.duration = nowPlaying.info.durationInMillis;
        event.game_event.data.frame.time = nowPlaying.info.currentPlaybackTime * 1000;
        await sleep(25);
        let res = await fetch(`http://${url}/game_event`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(event.game_event),
        });
        const result = await res.json();
      }
    }
  }
})();
