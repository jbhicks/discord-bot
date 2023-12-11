import { joinVoiceChannel, getVoiceConnection, createAudioPlayer, createAudioResource, AudioPlayerStatus } from '@discordjs/voice';
import { client } from "../bot.js";
import ytdl from 'ytdl-core'; // Import the 'ytdl-core' package

export async function playVideo(interaction, url) {
    joinChannel();
    // const streamOptions = { seek: 0, volume: 1 };
    const stream = ytdl(url, { filter: 'audioonly' });
    console.log(`getting voice connection for guild ${interaction.guildId}`);
    const guildId = interaction.guildId;
    const connection = getVoiceConnection(guildId);
    console.log(`connection: ${connection}`);
    if (!connection) {
        console.log('No voice connection in this guild');
        return;
    }
    const player = createAudioPlayer();
    connection.subscribe(player);

    console.log(`creating audio resource`)
    const audioResource = createAudioResource(stream);
    player.play(audioResource);

    player.on('stateChange', (oldState, newState) => {
        console.log(`Player state changed from ${oldState.status} to ${newState.status}`);
        if (newState.status === AudioPlayerStatus.Idle) {
            leaveVoiceChannel(guildId);
            console.log('Finished playing and left the voice channel');
        }

    });
}

function joinChannel() {
    const voiceChannel = client.channels.cache.get('939640186604781661');
    console.log(`voiceChannel: ${JSON.stringify(voiceChannel)}`);
    if (voiceChannel && voiceChannel.type === 2) {
        const connection = joinVoiceChannel({
            channelId: voiceChannel.id,
            guildId: voiceChannel.guild.id,
            adapterCreator: voiceChannel.guild.voiceAdapterCreator,
        });
        console.log(`Joined voice channel ${voiceChannel.name}`);
    }
}

export function leaveVoiceChannel(guildId) {
    const connection = getVoiceConnection(guildId);
    if (connection) {
        connection.destroy();
    }
}