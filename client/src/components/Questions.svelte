<script lang="ts">
    import AddQuestion from "./AddQuestion.svelte";
    import Question from "./Question.svelte";
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
            <Question {question} />
        {/each}
    </div>
{/if}

<style>
    .container {
        margin-top: 32px;
        display: flex;
        flex-direction: column;
        gap: 16px;
    }
</style>
