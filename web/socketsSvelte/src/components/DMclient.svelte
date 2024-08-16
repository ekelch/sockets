<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import type { DirectMessage, UserClient } from "../types";
    import DirectMsgBox from "./DirectMsgBox.svelte";

    export let fromUser: string;
    export let toUser: string;
    export let messageHist: DirectMessage[];
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
</script>

<main>
    <div class="chat-window">
        <div class="name-header">
            <span class="header-text">{toUser}</span>
            <button on:click={close} tabindex="0" class="header-x"
                >&#x2715;</button
            >
        </div>
        <DirectMsgBox messages={messageHist} on:sendMessage={sendDM} />
    </div>
</main>

<style>
    .chat-window {
        background-color: #c9deff;
        border-radius: 8px;
        display: flex;
        flex-direction: column;
    }

    .name-header {
        display: flex;
        background-color: gray;
        height: fit-content;
        padding: 6px;
        border-radius: 8px 8px 0 0;
    }

    .header-text {
        width: fit-content;
    }

    .header-x {
        width: fit-content;
        margin: 0 0 0 auto;
        cursor: pointer;
    }
</style>
