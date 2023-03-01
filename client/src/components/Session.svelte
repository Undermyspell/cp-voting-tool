<script lang="ts">
    import { startSession, stopSession } from "../lib/session";
</script>

<div>
    <button on:click={startSession}>Fragerunde starten</button>
    <button on:click={stopSession}>Fragerunde beenden</button>
</div>
