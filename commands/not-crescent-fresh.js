import { SlashCommandBuilder } from "@discordjs/builders";
import { createAudioResource } from "@discordjs/voice";
import ytdl from "ytdl-core";

export const data = new SlashCommandBuilder()
    .setName("not-crescent-fresh")
    .setDescription('plays the Not Crescent Fresh song');

export async function execute(interaction, player, url) {
    console.log(`handling Not Crescent Fresh command with ${player} and URL ${url}`);
    const crescentFreshUrl = 'https://youtu.be/S0gJZMx79pA';
    const streamOptions = { seek: 0, volume: 1 };
    const stream = ytdl(crescentFreshUrl, { filter: 'audioonly' });
    console.log(`creating audio resource from ${stream}\n ______________________________________________`);
    const audioResource = createAudioResource(stream);
    player.play(audioResource);
    console.log(`audio player is now status ${player.state.status}`);
    await interaction.reply('Now playing the Not Crescent Fresh song');
};