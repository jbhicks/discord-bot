const { SlashCommandBuilder } = require("@discordjs/builders");
const { createAudioResource, createAudioPlayer } = require("@discordjs/voice");
const ytdl = require("ytdl-core");

module.exports = {
    data: new SlashCommandBuilder()
        .setName("crescent-fresh")
        .setDescription('plays the Crescent Fresh song'),

    async execute(interaction, player, url) {
        console.log(`handling crescent-fresh command with ${player} and URL ${url}`);
        const crescentFreshUrl = 'https://youtu.be/_qU_gEiSbIU';
        const streamOptions = { seek: 0, volume: 1 };
        const stream = ytdl(crescentFreshUrl, { filter: 'audioonly' });
        console.log(`creating audio resource from ${stream}\n ______________________________________________`);
        const audioResource = createAudioResource(stream);
        player.play(audioResource);
        console.log(`audio player is now status ${player.state.status}`);
        await interaction.reply('Now playing the Crescent Fresh song');
    },
};
