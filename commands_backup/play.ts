import { SlashCommandBuilder } from "@discordjs/builders";
import { joinVoiceChannel, StreamType, createAudioPlayer, createAudioResource, entersState, VoiceConnection, VoiceConnectionStatus, AudioPlayerStatus } from '@discordjs/voice';
import { createDiscordJSAdapter } from "../adapter";
const { channelInfo } = require("../config.json");
const player = createAudioPlayer();
let connection = {} as VoiceConnection;

module.exports = {

  data: new SlashCommandBuilder()
    .setName("play")
    .setDescription("Plays a song!")
    .addStringOption(option =>
      option.setName('target')
        .setDescription('The song to play')
        .setRequired(true)),
  async execute(interaction: any) {
    const target = interaction.options.getString('target') ?? 'No target provided';
    this.createConnection(interaction);
    if (connection.state.status === VoiceConnectionStatus.Ready) {
      this.playSong(target);
    }
    await interaction.reply(`Playing song at ${target}`);
  },

  createConnection(interaction: any) {
    if (interaction !== null) {
      const channelObj = interaction.guild.channels.cache.get(channelInfo);
      connection = joinVoiceChannel({
        channelId: channelInfo.id,
        guildId: channelInfo.guild.id,
        adapterCreator: createDiscordJSAdapter(channelObj),
      });
    }

    try {
      entersState(connection, VoiceConnectionStatus.Ready, 30_0000);
      console.log("Connected!");
      return connection;
    } catch (error) {
      connection.destroy();
      throw error;
    }
  },

  playSong(target: string) {
    const resource = createAudioResource(target, {
      inputType: StreamType.Arbitrary,
    });
    player.play(resource);
    return entersState(player, AudioPlayerStatus.Playing, 5000);
  }

}


