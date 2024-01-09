const { SlashCommandBuilder } = require("@discordjs/builders");

const data = new SlashCommandBuilder()
    .setName('ping')
    .setDescription('Replies with Pong!');

async function execute(interaction, client) {

    console.log(`replying to ${interaction}`);
    await interaction.reply('Pong!');
};

module.exports = { data, execute };

