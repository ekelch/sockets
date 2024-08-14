<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import type { DirectMessage, UserClient } from "../types";
    import MsgBox from "./MsgBox.svelte";

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
    <MsgBox messages={messageHist} on:sendMessage={sendDM} />
</main>

<style>
</style>
