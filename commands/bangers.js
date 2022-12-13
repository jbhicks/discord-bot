const { SlashCommandBuilder } = require("@discordjs/builders");
const https = require("https");
const url = "https://us-central1-my-web-site-27a15.cloudfunctions.net/getSCStream?offset=0&limit=100&length=300000";

module.exports = {
  data: new SlashCommandBuilder().setName("bangers").setDescription("Retrieves a list of mixes from the Bangers playlist"),
  async execute(interaction) {
    
    await interaction.reply("Pong!");
  },
};

const trackList = getTrackList();
console.log(`trackList: ${trackList}`);

function getTrackList() {
  // send http get request to url to get track list
  https.get(url, (res) => {
    let data = "";
    res.on("data", (chunk) => {
      data += chunk;
    });
    res.on("end", () => {
      return data;
    });
    return data;
  });
}