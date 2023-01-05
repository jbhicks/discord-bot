const { SlashCommandBuilder } = require("@discordjs/builders");

module.exports = {
    data: new SlashCommandBuilder()
        .setName('stop')
        .setDescription('stop playing a song'),
    async execute(interaction, player) {
        player.stop();
        await interaction.reply('Record scratch!');
    },
};
