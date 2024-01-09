const { SlashCommandBuilder } = require("@discordjs/builders");
const { playVideo } = require("../../lib/common-functions.js");

const data = new SlashCommandBuilder()
    .setName("mercy")
    .setDescription('Good day, sir.');

async function execute(interaction, client) {

    const crescentFreshUrl = 'https://youtu.be/5hfYJsQAhl0';
    await playVideo(client, interaction, crescentFreshUrl);
    await interaction.reply('may God have mercy upon your soul');
};

module.exports = { data, execute };

