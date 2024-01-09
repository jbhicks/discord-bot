const { REST, Routes } = require("discord.js");
const dotenv = require('dotenv');
const fs = require('fs');
const tokenUrl = require('url');
const path = require('path');

dotenv.config({ path: `${__dirname}/.env` });

const token = process.env.TOKEN;
const guildId = process.env.GUILD_ID;
const clientId = process.env.CLIENT_ID;
const commands = [];

console.log(`[INFO] Deploying commands to guild ${guildId} with client ${clientId}`);

// Grab all the command folders from the commands directory you created earlier
const foldersPath = path.join(__dirname, 'commands');

const commandFolders = fs.readdirSync(foldersPath);

for (const folder of commandFolders) {
	// Grab all the command files from the commands directory you created earlier
	const commandsPath = path.join(foldersPath, folder);
	const commandFiles = fs.readdirSync(commandsPath).filter(file => file.endsWith('.js'));
	// Grab the SlashCommandBuilder#toJSON() output of each command's data for deployment
	for (const file of commandFiles) {
		const filePath = path.join(commandsPath, file);
		const command = require(filePath);
		if ('data' in command && 'execute' in command) {
			commands.push(command.data.toJSON());
		} else {
			console.log(`[WARNING] The command at ${filePath} is missing a required "data" or "execute" property.`);
		}
	}

}
const rest = new REST({ version: '9' }).setToken(token);

(async () => {
	try {
		console.log('Started refreshing application (/) commands.');


		await rest.put(
			Routes.applicationCommands(clientId, guildId),
			{ body: commands },

		);

		console.log('Successfully reloaded application (/) commands.');

	} catch (error) {
		console.error(error);
	}
})();