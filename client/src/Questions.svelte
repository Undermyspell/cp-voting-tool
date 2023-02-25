<script lang="ts">
    import { each } from "svelte/internal";
    import { getQuestions, questions } from "./lib/questions";

    $: promise = getQuestions();
</script>

{#each $questions as question}
    <div>{question.Text}</div>
{/each}
