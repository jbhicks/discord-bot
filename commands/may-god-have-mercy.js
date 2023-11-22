import { SlashCommandBuilder } from "@discordjs/builders";
import { createAudioResource, createAudioPlayer, getVoiceConnection } from "@discordjs/voice";
import ytdl from "ytdl-core";

export const data = new SlashCommandBuilder()
    .setName("mercy")
    .setDescription('Good day, sir.');

export async function execute(interaction) {
    const crescentFreshUrl = 'https://youtu.be/5hfYJsQAhl0';
    const streamOptions = { seek: 0, volume: 1 };
    const stream = ytdl(crescentFreshUrl, { filter: 'audioonly' });
    console.log(`getting voice connection for guild ${interaction.guildId}`);
    const guildId = interaction.guildId;
    const connection = getVoiceConnection(guildId);
    console.log(`connection: ${connection}`);
    if (!connection) {
        console.log('No voice connection in this guild');
        return;
    }

    const player = createAudioPlayer();
    connection.subscribe(player);

    const audioResource = createAudioResource(stream);
    player.play(audioResource);

    player.on('stateChange', (oldState, newState) => {
        console.log(`Player state changed from ${oldState.status} to ${newState.status}`);
    });
    await interaction.reply('may God have mercy upon your soul');
};