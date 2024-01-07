import { SlashCommandBuilder } from "@discordjs/builders";
import { playVideo } from "../../lib/common-functions.js";

export const data = new SlashCommandBuilder()
    .setName("crescent-fresh")
    .setDescription("Crescent fresh, baby.");

export async function execute(interaction) {
    const crescentFreshUrl = 'https://youtu.be/_qU_gEiSbIU';
    await playVideo(interaction, crescentFreshUrl);
    await interaction.reply('crescent fresh baby');
};
