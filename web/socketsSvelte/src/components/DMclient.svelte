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
</script>

<main>
    <div class="chat-window">
        <div class="name-header">{toUser}</div>
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
        background-color: gray;
        height: fit-content;
        padding: 6px;
        border-radius: 8px 8px 0 0;
    }
</style>
