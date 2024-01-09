const { SlashCommandBuilder } = require("@discordjs/builders");
const { playVideo } = require("../../lib/common-functions.js");

const data = new SlashCommandBuilder()
    .setName("llama-school")
    .setDescription('Good day, sir.');

async function execute(interaction, client) {

    const crescentFreshUrl = 'https://www.youtube.com/watch?v=FOfa6JEIVkk';
    await playVideo(client, interaction, crescentFreshUrl);
    await interaction.reply('how do youget the llama to school?');
};

module.exports = { data, execute };
