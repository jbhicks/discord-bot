import { SlashCommandBuilder } from "@discordjs/builders";

export const data = new SlashCommandBuilder()
    .setName('stop')
    .setDescription('stop playing a song');

export async function execute(interaction, player) {
    player.stop();
    await interaction.reply('Record scratch!');
};