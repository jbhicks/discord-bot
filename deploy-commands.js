import { REST, Routes } from "discord.js";
import dotenv from 'dotenv';
import fs from 'fs';
import { fileURLToPath } from 'url';
import { dirname } from 'path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

dotenv.config({ path: __dirname + '/.env' });

// rest of your code

const token = process.env.TOKEN;
const guildId = process.env.GUILD_ID;
const clientId = process.env.CLIENT_ID;

const commands = [];
const commandFiles = fs.readdirSync("./commands").filter((file) => file.endsWith(".ts") || file.endsWith(".js"));

const loadCommands = async () => {
  for (const file of commandFiles) {
    const command = await import(`./commands/${file}`);
    commands.push(command.data.toJSON());
  }
}

const rest = new REST({ version: "10" }).setToken(token);

(async () => {
  try {
    await loadCommands();
    console.log(`Started refreshing ${commands.length} application (/) commands`);
    console.log(`using token ${token} and guildId ${guildId} and clientId ${clientId}`);
    const data = await rest.put(Routes.applicationGuildCommands(clientId, guildId), {
      body: commands,
    });

    console.log(`Successfully reloaded ${data.length} application (/) commands.`);
  } catch (error) {
    console.error(error);
  }
})();