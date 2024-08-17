<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import type { DirectMessage, UserClient } from "../types";
    import DirectMsgBox from "./DirectMsgBox.svelte";

    export let fromUser: string;
    export let toUser: string;
    export let messageHist: DirectMessage[];
    let open: boolean = true;
    const dispatch = createEventDispatcher();

    function sendDM(msgEvent: any) {
        let dm: DirectMessage = {
            fromUser: fromUser,
            toUser: toUser,
            message: msgEvent.detail,
        };
        dispatch("sendDM", dm);
    }

    const close = () => {
        dispatch("removeClient", toUser);
    };

    const minChat = () => {
        open = false;
    };

    const showChat = () => {
        open = true;
    };
</script>

<main class="main-container">
    {#if open}
        <div class="chat-window">
            <div class="chat-header">
                <span class="header-text">{toUser}</span>
                <div class="header-btn-grp">
                    <button on:click={minChat} class="header-btn">-</button>
                    <button on:click={close} class="header-btn">&#x2715;</button
                    >
                </div>
            </div>
            <DirectMsgBox
                messages={messageHist}
                on:sendMessage={sendDM}
                username={fromUser}
            />
        </div>
    {:else}
        <div
            id="closed-chat"
            class="closed"
            role="button"
            on:keypress={(event) => {
                showChat();
            }}
            tabindex="0"
            on:click={showChat}
        >
            <div class="chat-header">
                <span class="header-text">{toUser}</span>
                <div class="header-btn-grp">
                    <button on:click={close} class="header-btn">&#x2715;</button
                    >
                </div>
            </div>
        </div>
    {/if}
</main>

<style>
    .main-container {
        margin-top: auto;
    }

    .chat-window {
        background-color: #c9deff;
        border-radius: 8px;
        display: flex;
        flex-direction: column;
    }

    .chat-header {
        display: flex;
        background-color: gray;
        height: fit-content;
        padding: 6px;
        border-radius: 8px 8px 0 0;
    }

    .closed {
        border-radius: 8px;
        overflow: hidden;
        cursor: pointer;

        width: 240px;
        height: 32px;
    }

    .header-text {
        width: fit-content;
    }

    .header-btn-grp {
        display: flex;
        margin-left: auto;
        gap: 4px;
    }

    .header-btn {
        width: 23px;
        margin: 0 0 0 auto;
        cursor: pointer;
    }
</style>
