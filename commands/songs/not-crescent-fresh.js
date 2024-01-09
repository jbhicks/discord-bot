const { SlashCommandBuilder } = require("@discordjs/builders");
const { playVideo } = require("../../lib/common-functions.js");

const data = new SlashCommandBuilder()
    .setName("not-crescent-fresh")
    .setDescription('Good day, sir.');

async function execute(interaction, client) {

    const crescentFreshUrl = 'https://youtu.be/S0gJZMx79pA';
    await playVideo(client, interaction, crescentFreshUrl);
    await interaction.reply('Not Crescent Fresh');
};

module.exports = { data, execute };

