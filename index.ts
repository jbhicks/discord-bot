import {
	joinVoiceChannel,
	createAudioPlayer,
	entersState,
	VoiceConnectionStatus,
} from '@discordjs/voice';
import { GatewayIntentBits, ActivityOptions, ActivityType, Client, Collection, Events, VoiceBasedChannel, VoiceChannel } from 'discord.js';
import { createDiscordJSAdapter } from './adapter';
import { token, guildId, channelInfo } from './config.json';
const player = createAudioPlayer();
const client = new Client({
	intents: [GatewayIntentBits.Guilds, GatewayIntentBits.GuildMessages, GatewayIntentBits.GuildVoiceStates],
});
const fs = require('fs');
const path = require('path');
const commands = new Collection();
const commandsPath = path.join(__dirname, 'commands');
const commandFiles = fs.readdirSync(commandsPath).filter((file: string) => file.endsWith('.js'));
const testMp3Url = 'https://www.soundhelix.com/examples/mp3/SoundHelix-Song-1.mp3';


commandSetup();
void client.login(token);

function commandSetup() {
	for (const file of commandFiles) {
		const filePath = path.join(commandsPath, file);
		const command = require(filePath);
		if ('data' in command) {
			commands.set(command.data.name, command);
		} else {
			console.error(`Command ${file} does not export a command data object`);
		}
	}
}

client.once(Events.ClientReady, async () => {
	console.log('Discord.js client is ready!');
	try {
		const channel: VoiceBasedChannel = client.channels.cache.get(channelInfo.id) as VoiceChannel;
		const connection = joinVoiceChannel({
			channelId: channelInfo.id,
			guildId: guildId,
			adapterCreator: createDiscordJSAdapter(channel),
		});
		client.user.setAvatar('./bot-avatar.jpg');

		client.user.setActivity('songs');
		if (connection) {
			console.log(`Connected to ${channel.name}!`);
			connection.subscribe(player);
			await entersState(connection, VoiceConnectionStatus.Ready, 30_000);
		} else {
			console.log('Connection failed!');
		}
	} catch (error) {
		console.error(error);
	}
});

client.on(Events.InteractionCreate, async (interaction) => {
	console.log(`handling interaction type ${interaction.type}`);
	if (!interaction.isChatInputCommand()) {
		console.log(`interaction is not a command`);
		return;
	}
	console.log(`received interaction ${interaction.commandName}`);
	
	try { 
		const command: any = commands.get(interaction.commandName);
		if (command) {
			await command.execute(interaction, player, testMp3Url);
		} else {
			console.log(`command ${interaction.commandName} not found`);
		}
	} catch (error) {
		console.error(error);
	}

});
