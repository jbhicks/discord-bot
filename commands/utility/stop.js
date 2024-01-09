const { SlashCommandBuilder } = require("@discordjs/builders");
const { leaveVoiceChannel } = require("../../lib/common-functions.js");

const data = new SlashCommandBuilder()
    .setName('stop')
    .setDescription('stop playing a song');

async function execute(interaction, player) {
    if (player) player.stop();
    leaveVoiceChannel(interaction.member.guild.id);
    await interaction.reply('Record scratch!');
};

module.exports = { data, execute };

