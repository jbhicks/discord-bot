const { SlashCommandBuilder } = require('@discordjs/builders');
const { playVideo } = require('../../lib/common-functions.js');

const data = new SlashCommandBuilder()
    .setName('crescent-fresh')
    .setDescription('Crescent fresh, baby.');

async function execute(interaction, client) {

    const crescentFreshUrl = 'https://youtu.be/_qU_gEiSbIU';
    await playVideo(client, interaction, crescentFreshUrl);
    await interaction.reply('crescent fresh baby');
}

module.exports = { data, execute };

