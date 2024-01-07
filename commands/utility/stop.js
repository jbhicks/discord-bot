import { SlashCommandBuilder } from "@discordjs/builders";
import { leaveVoiceChannel } from "../../lib/common-functions.js";

export const data = new SlashCommandBuilder()
    .setName('stop')
    .setDescription('stop playing a song');

export async function execute(interaction, player) {
    if (player) player.stop();
    leaveVoiceChannel(interaction.member.guild.id);
    await interaction.reply('Record scratch!');
};
