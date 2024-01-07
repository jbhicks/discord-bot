import { SlashCommandBuilder } from "@discordjs/builders";
import { playVideo } from "../../lib/common-functions.js";

export const data = new SlashCommandBuilder()
    .setName("mercy")
    .setDescription('Good day, sir.');

export async function execute(interaction) {
    const crescentFreshUrl = 'https://youtu.be/5hfYJsQAhl0';
    await playVideo(interaction, crescentFreshUrl);
    await interaction.reply('may God have mercy upon your soul');
};
