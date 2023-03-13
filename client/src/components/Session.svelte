<script lang="ts">
    import { startSession, stopSession } from "../lib/session";
</script>

<div class="flex gap-4">
    <button type="button" class="btn variant-filled" on:click={startSession}
        >Fragerunde starten</button
    >
    <button type="button" class="btn variant-filled" on:click={stopSession}
        >Fragerunde beenden</button
    >
</div>
