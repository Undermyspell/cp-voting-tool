<script lang="ts">
    export let showModal;

    let dialog;

    $: if (dialog && showModal) dialog.showModal();
    $: if (dialog && !showModal) dialog.close();
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<dialog
    class="bg-surface-100-800-token rounded-container-token"
    bind:this={dialog}
    on:close={() => (showModal = false)}
>
    <div on:click|stopPropagation>
        <slot name="header" />
        <slot />
        <hr />
        <!-- svelte-ignore a11y-autofocus -->
        <div class="flex justify-end gap-4 pt-4">
            <button
                type="button"
                class="btn btn-sm variant-filled"
                on:click={() => dialog.close()}>Abbrechen</button
            >
            <slot name="action" />
        </div>
    </div>
</dialog>

<style>
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
