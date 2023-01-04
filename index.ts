const fs = require("fs");
const path = require("path");
const commandsPath = path.join(__dirname, "commands");
const commandFiles = fs.readdirSync(commandsPath).filter((file: any) => file.endsWith(".js"));
const { Client, Events, Collection, GatewayIntentBits, VoiceBasedChannel } = require("discord.js");
const { token, app_id, channel } = require("./config.json");
const { joinVoiceChannel,
  createAudioPlayer,
  createAudioResource,
  enterState,
  StreamType,
  AudioPlayerStatus,
  VoiceConnectionStatus } = require('@discordjs/voice');
const client = new Client({
  intents: [GatewayIntentBits.Guilds, GatewayIntentBits.Guilds, GatewayIntentBits.GuildVoiceStates],
});
import { entersState } from "@discordjs/voice";
import { createDiscordJSAdapter } from "./adapter";
const player = createAudioPlayer();

initClient();

function initClient() {
  client.once(Events.ClientReady, () => {
    console.log(`Logged in as ${client.user.tag}!`);
    try {
      connectToChannel();
    } catch (error) {
      console.error(error);
    }
  });
  client.commands = new Collection();
  initClientCommands();
  initClientInteractions();
  client.login(token);
  client.on("error", console.error);
  client.on("warn", console.warn);
}

function playSong() {
  const resource = createAudioResource('https://www.soundhelix.com/examples/mp3/SoundHelix-Song-1.mp3', {
    inputType: StreamType.Arbitrary,
  });
  player.play(resource);
  return enterState(player, AudioPlayerStatus.Playing, 5000);
}

async function connectToChannel() {
  const channelObj = client.channels.cache.get('939640186604781661');
  console.log(`created channel: ${channelObj}`);
  const connection = joinVoiceChannel({
    channelId: channel.id,
    guildId: channel.guild.id,
    adapterCreator: createDiscordJSAdapter(channelObj),
  });

  try {
    await entersState(connection, VoiceConnectionStatus.Ready, 30_000);
    connection.subscribe(player);
    // await playSong();
  } catch (error) {
    connection.destroy();
    throw error;
  }
}

function initClientCommands() {
  for (const file of commandFiles) {
    const filePath = path.join(commandsPath, file);
    const command = require(filePath);
    if ("data" in command && "execute" in command) {
      client.commands.set(command.data.name, command);
    } else {
      console.warn(`Command ${file} is missing data or execute`);
    }
  }
  // console.log(`Loaded client commands: \n ${JSON.stringify(client.commands, null, 2)}`);
}

function initClientInteractions() {
  client.on(Events.InteractionCreate, (interaction: { isCommand?: any; reply?: any; commandName?: any; }) => {
    if (!interaction.isCommand()) return;
    const { commandName } = interaction;
    if (!client.commands.has(commandName)) return;
    try {
      client.commands.get(commandName).execute(interaction);
    } catch (error) {
      console.error(error);
      interaction.reply({
        content: "There was an error while executing this command!",
        ephemeral: true,
      });
    }
  });
}
