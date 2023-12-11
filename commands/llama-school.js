import { SlashCommandBuilder } from "@discordjs/builders";
import { playVideo } from "../lib/common-functions.js";

export const data = new SlashCommandBuilder()
    .setName("llama-school")
    .setDescription('Good day, sir.');

export async function execute(interaction) {
    const crescentFreshUrl = 'https://www.youtube.com/watch?v=FOfa6JEIVkk';
    await playVideo(interaction, crescentFreshUrl);
    await interaction.reply('how do youget the llama to school?');
};