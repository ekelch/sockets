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
    import DmOptions from "./components/DMoptions.svelte";
    import DMclient from "./components/DMclient.svelte";
    import BroadcastMsgBox from "./components/BroadcastMsgBox.svelte";
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
        dmClients = new Map();
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

    const closeDm = (event: any) => {
        dmClients.delete(event.detail);
        dmClients = dmClients;
    };
</script>

<main>
    <div class="main-container">
        {#if connStatus === WebSocket.OPEN}
            <div class="header">
                <div>Username: {username}</div>
            </div>
        {/if}
        <div class="grid-container">
            <div class="disconnect-col">
                {#if connStatus === WebSocket.OPEN}
                    <p>There are {numClients} users connected to the chat</p>
                    <button on:click={() => disconnect()}
                        >Disconnect from chat</button
                    >
                {/if}
            </div>
            <div class="chat-col">
                {#if connStatus === WebSocket.OPEN}
                    <div class="msg-box-container">
                        <BroadcastMsgBox
                            messages={broadcastMsgs}
                            {username}
                            on:sendMessage={(msg) => sendBroadcast(msg.detail)}
                        />
                    </div>
                {:else}
                    <div class="connect">
                        <!-- svelte-ignore a11y-autofocus -->
                        <input
                            autofocus
                            bind:value={username}
                            on:change={connect}
                            placeholder="enter your username"
                        />
                        <button on:click={connect}>Connect to chat</button>
                    </div>
                {/if}
            </div>
            <div class="dm-options-col">
                <DmOptions clients={oClients} on:openDM={openDm} />
            </div>
        </div>
    </div>

    <div class="dm-container-sticky">
        {#each dmClients as [toUser, messages]}
            <DMclient
                fromUser={username}
                {toUser}
                messageHist={messages}
                on:sendDM={sendDm}
                on:removeClient={closeDm}
            />
        {/each}
    </div>
</main>

<style lang="css">
    .main-container {
        display: flex;
        height: 100vh;
        flex-direction: column;
    }

    .header {
        background-color: lightseagreen;
        height: 32px;
    }

    .header div {
        width: fit-content;
        margin: 6px 16px 0 auto;
    }

    .grid-container {
        flex: 1;
        margin-top: 12px;
        display: flex;
        justify-content: space-between;
    }

    .disconnect-col {
        flex: 1;
        margin: 0 16px;
    }

    .chat-col {
        width: 50%;
        height: 40vh;
        margin: 0 auto auto;
        display: flex;
        flex-direction: column;
    }
    .chat-col button {
        width: fit-content;
        padding: 2px 24px;
    }

    .msg-box-container {
        display: flex;
        flex-direction: column;
    }

    .connect {
        margin: auto auto auto 0;
    }

    .dm-options-col {
        flex: 1;
        margin: 0 16px;
    }

    .dm-container-sticky {
        position: sticky;
        width: fit-content;
        gap: 12px;
        bottom: 12px;
        left: 12px;
        display: flex;
        flex-direction: row;
    }
</style>
