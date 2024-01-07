import { SlashCommandBuilder } from "@discordjs/builders";
import axios from 'axios';
// require('dotenv').config({ path: __dirname + '../.env' });
console.log(process.env.TOKEN);
const token = process.env.TOKEN;
const guildId = process.env.GUILD_ID;
const channelId = process.env.CHANNEL_ID;

export const data = new SlashCommandBuilder()
  .setName('bangers')
  .setDescription('Search for a mix to play');

export async function execute(interaction) {
  const tracks = await getResponses();
  await interaction.reply(`Found ${tracks.length} mixes. Choose wisely.`);
  for (const track of tracks) {
    await interaction.followUp(track);
  }
};

async function getResponses() {
  const tracks = await getSCStream();
  let responses = [];
  if (tracks) {
    for (const track of tracks) {
      responses.push(`Check out this ${track.track.genre} mix called ${track.track.title} by ${track.user.username} at ${track.track.permalink_url}`);
    }
    console.log(responses[responses.length - 1]);
  }
  return responses;
}

async function getSCStream() {
  const offset = 0;
  const limit = 750;
  const length = 3000000;
  const url = `https://api-v2.soundcloud.com/stream?offset=${offset}&sc_a_id=${sc_a_id}&limit=${limit}&promoted_playlist=true&client_id=${client_id}&app_version=1660231961&app_locale=en`;
  const headers = {
    "Accept": "application/json, text/javascript, */*; q=0.01",
    "Accept-Encoding": "gzip, deflate, br",
    "Accept-Language": "en-US,en;q=0.9",
    "Authorization": sc_oauth_token,
    "Connection": "keep-alive",
    "Host": "api-v2.soundcloud.com",
    "Origin": "https://soundcloud.com",
    "Referer": "https://soundcloud.com/",
    "Sec-Fetch-Dest": "empty",
    "Sec-Fetch-Mode": "cors",
    "Sec-Fetch-Site": "same-site",
    "User-Agent":
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.64 Safari/537.36 Edg/101.0.1210.53",
    "sec-ch-ua": "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"101\", \"Microsoft Edge\";v=\"101\"",
    "sec-ch-ua-mobile": "?0",
    "sec-ch-ua-platform": "\"Windows\"",
  };

  console.log(`attempting to get stream from ${url} with sc_a_id ${sc_a_id} and client_id ${client_id} with oauth token ${sc_oauth_token}`);

  const response = await get(url, { headers })
  console.log(`retrieved ${response.data.collection.length} tracks`);
  // filter out tracks that are of type playlist
  let filteredTracks = response.data.collection.filter((track) => {
    return !track.type.includes("playlist");
  });
  console.log(`num post filtering playlists: ${filteredTracks.length}`);
  // filter out tracks that are longer than the specified length
  filteredTracks = filteredTracks.filter((track) => {
    return track.track.duration > length;
  });
  console.log(`num tracks after filtering: ${filteredTracks.length}`);
  return filteredTracks;
}
