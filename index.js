import { Client, Events, GatewayIntentBits, Collection } from 'discord.js';
import fs from 'fs';
import path from 'node:path';
import pkg from 'dotenv';
const { dotenv } = pkg;
const token = process.env.TOKEN;
console.log(token);

const client = new Client({
	intents: [GatewayIntentBits.Guilds, GatewayIntentBits.GuildMessages, GatewayIntentBits.GuildVoiceStates],
});

client.commands = new Collection();
const commandsPath = path.join('./', 'commands');
const commandFiles = fs.readdirSync(commandsPath).filter(file => file.endsWith('.js'));

for (const file of commandFiles) {
	const filePath = path.join(commandsPath, file);
	const command = require(filePath);
	// Set a new item in the Collection with the key as the command name and the value as the exported module
	if ('data' in command && 'execute' in command) {
		client.commands.set(command.data.name, command);
	} else {
		console.log(`[WARNING] The command at ${filePath} is missing a required "data" or "execute" property.`);
	}
}

void client.login(token);

client.on('ready', () => {
	console.log('Discord.js client is ready!');
});

client.on(Events.InteractionCreate, interaction => {
	console.log(interaction);
});
