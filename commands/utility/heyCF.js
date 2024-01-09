const { SlashCommandBuilder } = require("@discordjs/builders");

const data = new SlashCommandBuilder()
    .setName('heycrescentfresh')
    .setDescription('Get Crescent Fresh!');

async function execute(interaction, client) {

    console.log(`replying to ${interaction}`);
    await interaction.reply('Crescent Fresh!');
};

module.exports = { data, execute };

