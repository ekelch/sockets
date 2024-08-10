<script lang="ts">
    import { onMount } from "svelte";
    import {
        DefaultMessages,
        type DirectMessage,
        type Message,
        type UserClient,
    } from "./types";
    import "./app.css";
    import MsgBox from "./components/MsgBox.svelte";
    import Dm from "./components/DM.svelte";
    let messages: Message[] = [];
    let socket: WebSocket;
    let connStatus: number = WebSocket.CLOSED;

    let username: string;
    let messageState: string;
    let oClients: Map<string, UserClient> = new Map();
    let numClients = 0;

    onMount(async () => {
        countClients();
    });

    const countClients = async () => {
        fetch("http://localhost:8080/count")
            .then((resp) => resp.json())
            .then((data) => (numClients = data));
    };

    const connect = () => {
        socket = new WebSocket("ws://localhost:8080/ws");
        connStatus = WebSocket.OPEN;

        socket.addEventListener("open", (event) => {
            sendMsg(DefaultMessages.CONNECT);
        });

        // Listen for messages
        socket.addEventListener("message", (event) => {
            const t: Message = JSON.parse(event.data);
            oClients = new Map(Object.entries(t.clients));
            numClients = oClients.size;
            oClients.delete(username);
            if (t.message && t.username) {
                messages = [t, ...messages];
            }
        });
    };

    const disconnect = () => {
        sendMsg(DefaultMessages.DISCONNECT);
        oClients = new Map();
        username = "";
        socket.close();
        connStatus = socket.CLOSED;
        messages = [];
        countClients();
    };

    const sendMsg = (message: string) => {
        socket.send(JSON.stringify({ username, message }));
    };

    const sendAndReset = (message: string) => {
        sendMsg(message);
        messageState = "";
    };

    const sendDm = (event: any) => {
        event.detail.fromUser = username;
        console.log(event.detail);
    };
</script>

<main>
    <div class="container">
        <div class="placeholder"></div>
        <div class="chat-col">
            <p>There are {numClients} users connected to the chat</p>
            {#if connStatus === WebSocket.OPEN}
                <button on:click={() => disconnect()}
                    >Disconnect from chat</button
                >
                <br /><br />

                <MsgBox
                    {messages}
                    on:sendMessage={(msg) => sendAndReset(msg.detail)}
                />
            {:else}
                <!-- svelte-ignore a11y-autofocus -->
                <input
                    autofocus
                    bind:value={username}
                    on:change={() => connect()}
                    placeholder="enter your username"
                />
                <button on:click={() => connect()}>Connect to chat</button>
            {/if}
        </div>
        <Dm clients={oClients} on:sendDM={(event) => sendDm(event)} />
    </div>
</main>

<style lang="css">
    .container {
        display: grid;
        grid-template-columns: 1fr 1fr 1fr;
        margin: auto;
    }

    .chat-col {
        width: 75%;
        margin: auto;
    }
</style>
