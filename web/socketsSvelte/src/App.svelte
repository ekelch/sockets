<script lang="ts">
    import { onMount } from "svelte";
    import { DefaultMessages, type Message, type UserClient } from "./types";
    import "./app.css";
    let messages: Message[] = [];
    let socket: WebSocket;
    let connStatus: number = WebSocket.CLOSED;

    let username: string;
    let messageState: string;
    let clients: Map<string, UserClient> = new Map();
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
            clients = new Map(Object.entries(t.clients));
            numClients = clients.size;
            clients.delete(username);
            if (t.message && t.username) {
                messages = [t, ...messages];
            }
        });
    };

    const disconnect = () => {
        sendMsg(DefaultMessages.DISCONNECT);
        clients = new Map();
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
</script>

<main>
    <div class="container">
        <div class="chat-col">
            <p>There are {numClients} users connected to the chat</p>
            {#if connStatus === WebSocket.OPEN}
                <button on:click={() => disconnect()}
                    >Disconnect from chat</button
                >
                <br /><br />

                <div class="msg-box">
                    {#each messages as msg}
                        <span>{msg.username}: {msg.message}</span>
                    {/each}
                </div>
                <!-- svelte-ignore a11y-autofocus -->
                <input
                    autofocus
                    bind:value={messageState}
                    on:change={() => sendAndReset(messageState)}
                />
                <button on:click={() => sendAndReset(messageState)}
                    >Send Msg</button
                >
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
        {#if clients?.size}
            <div>
                <p>Other online users:</p>
                {#each [...clients] as [key, value]}
                    <div class="dm-col">
                        <span>{value.username}</span>
                        <button on:click={() => alert(key)}>dm</button>
                    </div>
                {/each}
            </div>
        {/if}
    </div>
</main>
