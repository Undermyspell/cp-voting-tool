<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import {
        answerQuestion,
        deleteQuestion,
        voteQuestion,
    } from "../lib/questions";
    import type { Question } from "../models/question";

    const dispatch = createEventDispatcher();

    let showModal;
    export let question: Question;

    function edit() {
        dispatch("edit", {
            id: question.Id,
        });
    }
</script>

<div class="container" class:complete={question.Answered}>
    <div>{question.Text}</div>
    <div>{question.Id}</div>
    <div>Votes: {question.Votes}</div>
    <button on:click={() => voteQuestion(question.Id)}>UpVote</button>
    <button on:click={() => edit()}> bearbeiten </button>
    <button on:click={() => answerQuestion(question.Id)}>Answered</button>
    <button on:click={() => deleteQuestion(question.Id)}>Delete</button>
</div>

<style>
    .container {
        display: flex;
        gap: 16px;
        padding: 16px;
        background-color: #242424;
        color: #fcfcfc;
        box-shadow: var(--shadow-medium);
        border-radius: 0.5rem;
        justify-content: flex-end;
        place-items: center;
    }

    .complete {
        background-color: hsla(131, 79%, 66%, 0.5);
    }

    @media (prefers-color-scheme: light) {
        .container {
            color: #213547;
            background-color: #fcfcfc;
        }
        button {
            background-color: #eee;
        }
    }
</style>
