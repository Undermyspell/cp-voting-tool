<script lang="ts">
    export let showModal;

    let dialog;

    $: if (dialog && showModal) dialog.showModal();
    $: if (dialog && !showModal) dialog.close();
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<dialog bind:this={dialog} on:close={() => (showModal = false)}>
    <div on:click|stopPropagation>
        <slot name="header" />
        <slot />
        <hr />
        <!-- svelte-ignore a11y-autofocus -->
        <div class="footer">
            <button on:click={() => dialog.close()}>Abbrechen</button>
            <slot name="action" />
        </div>
    </div>
</dialog>

<style>
    dialog {
        max-width: fit-content;
        border-radius: 0.2em;
        border: none;
        padding: 0;
    }
    dialog::backdrop {
        background: rgba(0, 0, 0, 0.3);
    }
    dialog > div {
        padding: 1em;
    }
    dialog .footer {
        display: flex;
        gap: 32px;
        justify-content: end;
    }
    dialog[open] {
        animation: zoom 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
    }
    @keyframes zoom {
        from {
            transform: scale(0.95);
        }
        to {
            transform: scale(1);
        }
    }
    dialog[open]::backdrop {
        animation: fade 0.2s ease-out;
    }
    @keyframes fade {
        from {
            opacity: 0;
        }
        to {
            opacity: 1;
        }
    }
    button {
        justify-self: end;
    }
</style>
