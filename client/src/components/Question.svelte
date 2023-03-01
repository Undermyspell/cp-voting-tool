<script lang="ts">
    import { voteQuestion } from "../lib/questions";
    import type { Question } from "../models/question";

    export let question: Question;
</script>

<div class="container">
    <div>{question.Text}</div>
    <div>Votes: {question.Votes}</div>
    <button on:click={() => voteQuestion(question.Id)}>UpVote</button>
</div>

<style>
    .container {
        display: flex;
        gap: 16px;
        justify-content: flex-end;
        place-items: center;
    }
</style>
