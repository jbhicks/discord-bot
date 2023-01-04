import type { DiscordGatewayAdapterCreator, DiscordGatewayAdapterLibraryMethods } from "@discordjs/voice";
import {
    GatewayDispatchEvents,
    GatewayVoiceServerUpdateDispatchData,
    GatewayVoiceStateUpdateDispatchData,
} from "discord-api-types/v10";
import { Client, Snowflake, Guild, VoiceBasedChannel, Constants, Events, Status } from "discord.js";
const adapters = new Map<Snowflake, DiscordGatewayAdapterLibraryMethods>();
const trackedClients = new Set<Client>();
const trackedShards = new Map<number, Set<Snowflake>>();

function trackClient(client: Client) {
    if (trackedClients.has(client)) return;
    trackedClients.add(client);
    client.ws.on(GatewayDispatchEvents.VoiceServerUpdate, (data: GatewayVoiceServerUpdateDispatchData) => {
        adapters.get(data.guild_id)?.onVoiceServerUpdate(data);
    });

    client.ws.on(GatewayDispatchEvents.VoiceStateUpdate, (data: GatewayVoiceStateUpdateDispatchData) => {
        if (data.guild_id && data.session_id && data.user_id === client.user?.id) {
            adapters.get(data.guild_id)?.onVoiceStateUpdate(data);
        }
    });

    client.on(Events.ShardDisconnect, (event, shardID) => {
        const guilds = trackedShards.get(shardID);
        if (guilds) {
            for (const guildID of guilds) {
                const adapter = adapters.get(guildID);
                if (adapter) {
                    adapter.destroy();
                    adapters.delete(guildID);
                }
            }
            trackedShards.delete(shardID);
        }
    });
}

function trackGuild(guild: Guild) {
    let guilds = trackedShards.get(guild.shardId);
    if (!guilds) {
        guilds = new Set();
        trackedShards.set(guild.shardId, guilds);
    }
    guilds.add(guild.id);
}

export function createDiscordJSAdapter(channel: VoiceBasedChannel): DiscordGatewayAdapterCreator {
    console.log(`Creating adapter for channel ${channel.name} (${channel.id}) in guild ${channel.guild.name}`);
    return (methods) => {
        adapters.set(channel.guild.id, methods);
        trackClient(channel.client);
        trackGuild(channel.guild);
        return {
            sendPayload(data) {
                if (channel.guild.shard.status === Status.Ready) {
                    channel.guild.shard.send(data);
                    return true;
                }
                return false;
            },
            destroy() {
                return adapters.delete(channel.guild.id);
            },
        };
    };
}
