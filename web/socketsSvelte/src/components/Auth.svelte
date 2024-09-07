<script lang="ts">
    import { createEventDispatcher } from "svelte";
    const dispatch = createEventDispatcher();

    let username = "";
    const connect = () => {
        dispatch("connect", username);
    };
</script>

<main class="connect">
    <div>
        <!-- svelte-ignore a11y-autofocus -->
        <input
            autofocus
            bind:value={username}
            on:change={connect}
            placeholder="enter your username"
        />
        <button on:click={connect}>Connect to chat</button>
    </div>
</main>

<style lang="CSS">
    .connect {
        margin: auto auto auto 0;
    }
</style>
