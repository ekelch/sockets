<script lang="ts">
    import { onMount } from "svelte";
    import {
        DefaultMessages,
        type BroadcastMessage,
        type DirectMessage,
        type JsonMessage,
        type UserClient,
    } from "./types";
    import "./app.css";
    import MsgBox from "./components/MsgBox.svelte";
    import Dm from "./components/DM.svelte";
    let broadcastMsgs: BroadcastMessage[] = [];
    let socket: WebSocket;
    let connStatus: number = WebSocket.CLOSED;

    let username: string;
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
            sendBroadcast(DefaultMessages.CONNECT);
        });

        // Listen for messages
        socket.addEventListener("message", (event) => {
            const jsonRec: JsonMessage = JSON.parse(event.data);
            if (jsonRec.type === "broadcast") {
                const bMsg: BroadcastMessage = jsonRec.rawMsg;
                oClients = new Map(Object.entries(bMsg.clients));
                numClients = oClients.size;
                oClients.delete(username);
                if (bMsg.message && bMsg.fromUser) {
                    broadcastMsgs = [bMsg, ...broadcastMsgs];
                }
            } else if (jsonRec.type === "direct") {
                console.log(jsonRec.rawMsg);
            }
        });
    };

    const disconnect = () => {
        sendBroadcast(DefaultMessages.DISCONNECT);
        oClients = new Map();
        username = "";
        connStatus = socket.CLOSED;
        broadcastMsgs = [];
        console.log("closing websocket connection:");
        console.log(socket);
        socket.close(1000);
        countClients();
    };

    const sendBroadcast = (message: string) => {
        const msg: JsonMessage = {
            type: "broadcast",
            rawMsg: { fromUser: username, message: message },
        };
        console.log(JSON.stringify(msg));
        socket.send(JSON.stringify(msg));
    };

    const sendDm = (event: any) => {
        event.detail.fromUser = username;
        const msg: JsonMessage = {
            type: "direct",
            rawMsg: JSON.stringify(event.detail),
        };
        console.log(msg);
        socket.send(JSON.stringify(msg));
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
                    messages={broadcastMsgs}
                    on:sendMessage={(msg) => sendBroadcast(msg.detail)}
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
