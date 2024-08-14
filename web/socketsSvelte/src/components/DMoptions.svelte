<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import MsgBox from "./MsgBox.svelte";
    import { type DirectMessage } from "../types";

    export let clients;

    const dispatch = createEventDispatcher();

    const openDm = (toUsername: string) => {
        dispatch("openDM", toUsername);
    };
</script>

<main>
    <div class="container">
        {#if clients?.size}
            <div>
                <p>Other online users:</p>
                {#each [...clients] as [key, value]}
                    <div class="dm-col">
                        <span>{value.username}</span>
                        <button on:click={() => openDm(key)}>dm</button>
                    </div>
                {/each}
            </div>
        {/if}
    </div>
</main>

<style>
    .container {
        width: 50%;
    }
    .dm-col {
        display: flex;
        flex-direction: row;
        gap: 8px;
        vertical-align: text-bottom;
    }
</style>
