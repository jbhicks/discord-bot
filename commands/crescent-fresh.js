const { SlashCommandBuilder } = require("@discordjs/builders");
const { createAudioResource } = require("@discordjs/voice");
const ytdl = require("ytdl-core");

module.exports = {
    data: new SlashCommandBuilder()
        .setName("crescent-fresh")
        .setDescription('plays the Crescent Fresh song'),
    
    async execute(interaction, player, url) {
        console.log(`handling crescent-fresh command with ${player} and URL ${url}`);
        const crescentFreshUrl = 'https://youtu.be/_qU_gEiSbIU';
        const streamOptions = { seek: 0, volume: 1 };
        const process = ytdl(crescentFreshUrl, { filter: 'audioonly' });
        const stream = process.stdout;
        console.log(`creating audio resource from ${stream}\n ______________________________________________`);
        const audioResource = createAudioResource(stream);
        player.play(audioResource);
        await interaction.reply('Now playing the Crescent Fresh song');
    },
};

const crescentFreshUrl = 'https://youtu.be/_qU_gEiSbIU';
const streamOptions = { seek: 0, volume: 1 };
const stream = ytdl(crescentFreshUrl, { filter: 'audioonly' });
console.log(`creating audio resource from ${JSON.stringify(stream, null, 2)}\n ______________________________________________`);
const audioResource = createAudioResource(stream);
console.log(`created audio resource ${JSON.stringify(audioResource, null, 2)} from stream`);
 
 
