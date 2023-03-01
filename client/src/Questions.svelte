<script lang="ts">
    import { getQuestions, questions } from "./lib/questions";
    import { activeSessison } from "./lib/session";

    $: promise = getQuestions();
</script>

{#if !$activeSessison}
    <div>keine aktive Q & A Session</div>
{:else}
    {#each $questions as question}
        <div>{JSON.stringify(question)}</div>
    {/each}
{/if}
