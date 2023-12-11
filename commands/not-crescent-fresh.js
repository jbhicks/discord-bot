import { SlashCommandBuilder } from "@discordjs/builders";
import { playVideo } from "../lib/common-functions.js";

export const data = new SlashCommandBuilder()
    .setName("not-crescent-fresh")
    .setDescription('Good day, sir.');

export async function execute(interaction) {
    const crescentFreshUrl = 'https://youtu.be/S0gJZMx79pA';
    await playVideo(interaction, crescentFreshUrl);
    await interaction.reply('Not Crescent Fresh');
};