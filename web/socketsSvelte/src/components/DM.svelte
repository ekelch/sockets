<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import MsgBox from "./MsgBox.svelte";
    import { type DirectMessage, type BroadcastMessage } from "../types";

    export let clients;
    let dmClient: string;
    let dmMessages: DirectMessage[] = [];

    const dispatch = createEventDispatcher();

    function sendDM(msgEvent: any) {
        let dm: DirectMessage = {
            fromUser: "",
            toUser: dmClient,
            message: msgEvent.detail,
        };
        dispatch("sendDM", dm);
    }

    const msgUser = (toUsername: string) => {
        dmClient = toUsername;
    };
</script>

<main>
    <div class="container">
        {#if clients?.size}
            <div>
                <p>Other online users:</p>
                {#each [...clients] as [key, value]}
                    <div class="dm-col">
                        <span>{value.username}</span>
                        <button on:click={() => msgUser(key)}>dm</button>
                    </div>
                {/each}
            </div>
        {/if}
        {#if dmClient}
            <MsgBox
                messages={dmMessages}
                on:sendMessage={(msgEvent) => sendDM(msgEvent)}
            />
            <button on:click={() => msgUser("")}>close</button>
        {/if}
    </div>
</main>

<style>
    .container {
        width: 50%;
    }
    .dm-col {
        display: flex;
        flex-direction: row;
        gap: 8px;
        vertical-align: text-bottom;
    }
</style>
