const { SlashCommandBuilder } = require("@discordjs/builders");
const { createAudioResource } = require("@discordjs/voice");
const ytdl = require("ytdl-core");

module.exports = {
    data: new SlashCommandBuilder()
        .setName('play')
        .setDescription('play a song')
        .addStringOption(option => option.setName('url').setDescription('url of song to play')),
    async execute(interaction, player) {
        const url = interaction.options.getString('url');
        if (!url) {
            await interaction.reply('did not receive a url');
        } else {
            console.log(`playing song at ${url}`);
            if (url.includes('.mp3')) {
                console.log(`setting up mp3 at ${url}`);
                const resource = createAudioResource(url);
                console.log(`resource: ${resource}`);
                player.play(resource);
            } else {
                console.log(`setting up youtube video at ${url}`);
                const streamOptions = { seek: 0, volume: 1 };
                const stream = ytdl(url, { filter: 'audioonly' });
                const resource = createAudioResource(stream, streamOptions);
                console.log(`resource: ${resource}`);
                player.play(resource);
            }
            await interaction.reply('Now playing ' + url);
        }
    },
};
