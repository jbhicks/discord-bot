const { SlashCommandBuilder } = require("@discordjs/builders");
const { createAudioResource } = require("@discordjs/voice");

module.exports = {
    data: new SlashCommandBuilder()
        .setName('play')
        .setDescription('plays a song')
        .addStringOption(option => option.setName('url').setDescription('url of song to play')),
    async execute(interaction, player, url) {
        console.log(`replying to ${interaction} with url ${url}`);
        const resource = createAudioResource(url);
        player.play(resource);
        await interaction.reply('Now playing ' + url);
    },
};


