import { Client, Events, GatewayIntentBits, Collection } from 'discord.js';
import fs from 'fs';
import path from 'node:path';
import { config } from 'dotenv';
import { fileURLToPath } from 'url';
import { dirname, resolve } from 'path';
import { joinVoiceChannel } from '@discordjs/voice';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);
config({ path: resolve(__dirname, '.env') });

const token = process.env.TOKEN;
console.log(`Token: ${token}`);

const client = new Client({
	intents: [GatewayIntentBits.Guilds, GatewayIntentBits.GuildMessages, GatewayIntentBits.GuildVoiceStates],
});

client.commands = new Collection();
const commandsPath = path.resolve('./commands');
const commandFiles = fs.readdirSync(commandsPath).filter(file => file.endsWith('.js'));

for (const file of commandFiles) {
	const filePath = path.join(commandsPath, file);
	import(filePath).then((command) => {
		// Set a new item in the Collection with the key as the command name and the value as the exported module
		if ('data' in command && 'execute' in command) {
			client.commands.set(command.data.name, command);
		} else {
			console.log(`[WARNING] The command at ${filePath} is missing a required "data" or "execute" property.`);
		}
	}).catch((error) => {
		console.error(`[ERROR] Failed to import command at ${filePath}:`, error);
	});
}

void client.login(token);

client.on('ready', () => {
	console.log('Discord.js client is ready!');
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
});

client.on(Events.InteractionCreate, async interaction => {
	if (!interaction.isCommand()) return;
	const command = client.commands.get(interaction.commandName);
	console.log(`command received: ${JSON.stringify(command)}`);
	if (!command) return;
	try {
		await command.execute(interaction);
	} catch (error) {
		console.error(error);
		await interaction.reply({ content: 'There was an error while executing this command!', ephemeral: true });
	}
});
