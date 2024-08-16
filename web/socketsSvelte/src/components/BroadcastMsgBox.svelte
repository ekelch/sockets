<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import {
        DefaultMessages,
        type BroadcastMessage,
        type DirectMessage,
    } from "../types";

    export let messages: (BroadcastMessage | DirectMessage)[];
    let message: string;

    const dispatch = createEventDispatcher();
    function sendMessage(msg: string) {
        if (msg) dispatch("sendMessage", msg);
        message = "";
    }
</script>

<main class="msg-box">
    <div class="messages">
        {#each messages as msg}
            <span
                class:con-msg={msg.message === DefaultMessages.CONNECT}
                class:disc-msg={msg.message === DefaultMessages.DISCONNECT}
                >{msg.fromUser}: {msg.message}</span
            >
        {/each}
    </div>
    <div class="msg-input">
        <!-- svelte-ignore a11y-autofocus -->
        <textarea
            autofocus
            bind:value={message}
            on:keydown={(event) => {
                if (event.key === "Enter") {
                    event.preventDefault();
                    sendMessage(message);
                }
            }}
        />
        <button on:click={() => sendMessage(message)}>Send Msg</button>
    </div>
</main>

<style lang="css">
    .msg-box {
        height: 50vh;
        display: flex;
        flex-direction: column;
        color: black;
        background-color: #c9deff;
        border-radius: 8px;
        margin: 16px 0;
        padding: 16px;
        overflow-y: scroll;
    }

    .con-msg {
        color: green;
    }

    .disc-msg {
        color: red;
    }

    .messages {
        flex: 1;
        margin-bottom: 8px;
        overflow-y: auto;
        display: flex;
        flex-direction: column-reverse;
    }

    .messages span {
        margin: 4px 0;
        word-wrap: break-word;
    }

    .msg-input {
        height: 48px;
        display: flex;
        flex-direction: row;
    }

    .msg-input textarea {
        flex: 1;
        resize: none;
    }
</style>
