<script lang="ts">
    import AddQuestion from "./AddQuestion.svelte";
    import { getQuestions, questions } from "../lib/questions";
    import { activeSessison } from "../lib/session";

    $: promise = getQuestions();
</script>

{#if !$activeSessison}
    <div>keine aktive Q & A Session</div>
{:else}
    <AddQuestion />
    <div class="container">
        {#each $questions as question}
            <div>{JSON.stringify(question)}</div>
        {/each}
    </div>
{/if}

<style>
    .container {
        margin-top: 32px;
    }
</style>
