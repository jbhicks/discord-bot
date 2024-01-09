const { SlashCommandBuilder } = require("@discordjs/builders");
const { playVideo } = require("../../lib/common-functions.js");

const data = new SlashCommandBuilder()
    .setName("us-whatever")
    .setDescription('whatever man');

async function execute(interaction, client) {

    const crescentFreshUrl = 'https://youtu.be/viaTT859Yk0';
    await playVideo(client, interaction, crescentFreshUrl);
    await interaction.reply('US Whatever');
};

module.exports = { data, execute };
