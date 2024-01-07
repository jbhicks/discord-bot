import { SlashCommandBuilder } from "@discordjs/builders";
import { playVideo } from "../../lib/common-functions.js";

export const data = new SlashCommandBuilder()
    .setName("us-whatever")
    .setDescription('whatever man');

export async function execute(interaction) {
    const crescentFreshUrl = 'https://youtu.be/viaTT859Yk0';
    await playVideo(interaction, crescentFreshUrl);
    await interaction.reply('US Whatever');
};
