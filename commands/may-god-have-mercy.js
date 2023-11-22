import { SlashCommandBuilder } from "@discordjs/builders";
import { createAudioResource, createAudioPlayer } from "@discordjs/voice";
import ytdl from "ytdl-core";

export const data = new SlashCommandBuilder()
    .setName("mercy")
    .setDescription('Good day, sir.');

export async function execute(interaction) {
    const crescentFreshUrl = 'https://youtu.be/5hfYJsQAhl0';
    const streamOptions = { seek: 0, volume: 1 };
    const stream = ytdl(crescentFreshUrl, { filter: 'audioonly' });
    console.log(`creating audio resource from ${stream}\n ______________________________________________`);
    const audioResource = createAudioResource(stream);
    console.dir(audioResource, { depth: null });

    const player = createAudioPlayer();
    player.play(audioResource);

    console.log(`audio player is now status ${player.state.status}`);
    await interaction.reply('may God have mercy upon your soul');
};