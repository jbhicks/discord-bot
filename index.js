const fs = require("fs");
const path = require("path");
const commandsPath = path.join(__dirname, "commands");
const commandFiles = fs.readdirSync(commandsPath).filter((file) => file.endsWith(".js"));
const { Client, Events, Collection, GatewayIntentBits } = require("discord.js");
const { token } = require("./config.json");
const client = new Client({
  intents: [GatewayIntentBits.Guilds],
});

initClient();

function initClient() {
  client.once(Events.ClientReady, () => {
    console.log(`Logged in as ${client.user.tag}!`);
  });
  client.commands = new Collection();
  initClientCommands();
  initClientInteractions();
  client.login(token);
  client.on("error", console.error);
  client.on("warn", console.warn);
}

function initClientCommands() {
  for (const file of commandFiles) {
    const filePath = path.join(commandsPath, file);
    const command = require(filePath);
    if ("data" in command && "execute" in command) {
      client.commands.set(command.data.name, command);
    } else {
      console.warn(`Command ${file} is missing data or execute`);
    }
  }
  //   console.log(`Loaded client commands: \n ${JSON.stringify(client.commands, null, 2)}`);
}

function initClientInteractions() {
  client.on(Events.InteractionCreate, (interaction) => {
    if (!interaction.isCommand()) return;
    const { commandName } = interaction;
    if (!client.commands.has(commandName)) return;
    try {
      client.commands.get(commandName).execute(interaction);
    } catch (error) {
      console.error(error);
      interaction.reply({
        content: "There was an error while executing this command!",
        ephemeral: true,
      });
    }
  });
}
