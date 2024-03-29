const { SlashCommandBuilder } = require("@discordjs/builders");
const { createAudioResource } = require("@discordjs/voice");
const { playVideo } = require("../../lib/common-functions.js");

const data = new SlashCommandBuilder()
    .setName('play')
    .setDescription('play a song')
    .addStringOption(option => option.setName('url').setDescription('url of song to play'));

async function execute(interaction, player) {
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
            await playVideo(interaction, url);
        }
        await interaction.reply('Now playing ' + url);
    }
};

module.exports = { data, execute };

