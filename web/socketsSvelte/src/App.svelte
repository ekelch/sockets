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
    import DmOptions from "./components/DMoptions.svelte";
    import DMclient from "./components/DMclient.svelte";
    let broadcastMsgs: BroadcastMessage[] = [];
    let socket: WebSocket;
    let connStatus: number = WebSocket.CLOSED;

    let username: string;
    let oClients: Map<string, UserClient> = new Map();
    let dmClients: Map<string, DirectMessage[]> = new Map();
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
                const dm: DirectMessage = jsonRec.rawMsg;
                const oUser = dm.toUser === username ? dm.fromUser : dm.toUser;
                console.log(dm);
                dmClients = dmClients.has(oUser)
                    ? dmClients.set(oUser, [dm, ...dmClients.get(oUser)!])
                    : dmClients.set(oUser, [dm]);
            }
        });
    };

    const disconnect = () => {
        const msg: JsonMessage = {
            type: "disconnect",
            rawMsg: { fromUser: username, message: DefaultMessages.DISCONNECT },
        };
        oClients = new Map();
        username = "";
        connStatus = socket.CLOSED;
        broadcastMsgs = [];

        socket.send(JSON.stringify(msg));
    };

    const sendBroadcast = (message: string) => {
        const msg: JsonMessage = {
            type: "broadcast",
            rawMsg: { fromUser: username, message: message },
        };
        socket.send(JSON.stringify(msg));
    };

    const sendDm = (event: any) => {
        event.detail.fromUser = username;
        const msg: JsonMessage = {
            type: "direct",
            rawMsg: event.detail,
        };
        socket.send(JSON.stringify(msg));
    };

    const openDm = (event: any) => {
        dmClients = dmClients.set(event.detail, []);
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
        <DmOptions clients={oClients} on:openDM={(event) => openDm(event)} />
    </div>
    <div class="dm-container">
        {#each dmClients as [toUser, messages]}
            <div class="dm-client">
                <span>{toUser}</span>
            </div>
            <DMclient
                fromUser={username}
                {toUser}
                messageHist={messages}
                on:sendDM={sendDm}
            />
        {/each}
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

    .dm-container {
        margin: auto 0 0;
        display: flex;
        flex-direction: row;
    }

    .dm-client {
    }
</style>
