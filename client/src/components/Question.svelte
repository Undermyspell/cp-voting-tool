<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import { isAdmin, isSessionAdmin } from "../lib/auth/auth";
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
    <button on:click={() => voteQuestion(question.Id)} class="votingbubble">
        <div class="arrow-up" />
        <div class="votes">{question.Votes}</div>
    </button>
    <div class="content">
        <div class="questiontext">{question.Text}</div>
        <div class="action-panel">
            {#if isAdmin || isSessionAdmin || question.Owned}
                <button on:click={() => edit()}> Bearbeiten </button>
            {/if}
            {#if isAdmin || isSessionAdmin}
                <button on:click={() => answerQuestion(question.Id)}
                    >Beantwortet</button
                >
            {/if}
            {#if isAdmin || isSessionAdmin || question.Owned}
                <button on:click={() => deleteQuestion(question.Id)}
                    >Löschen</button
                >
            {/if}
        </div>
    </div>
</div>

<style>
    .container {
        display: flex;
        gap: 16px;
        padding: 8px;
        background-color: #343434;
        color: #fcfcfc;
        box-shadow: var(--shadow-medium);
        border-radius: 0.5rem;
        place-items: center;
    }

    .content {
        width: 100%;
    }

    .arrow-up {
        width: 0;
        height: 0;
        border-left: 8px solid transparent;
        border-right: 8px solid transparent;

        border-bottom: 8px solid white;
    }

    .votingbubble {
        display: flex;
        cursor: pointer;
        border-radius: 4px;
        width: 3.5rem;
        align-items: center;
        border: 1px solid gray;
        flex-direction: column;
        padding: 8px;
    }

    .action-panel {
        display: flex;
        justify-content: end;
        gap: 8px;
    }

    .votingbubble:hover {
        background-color: #444;
    }

    .questiontext {
        display: block;
        text-align: left;
        width: 100%;
        margin-bottom: 8px;
    }

    .votes {
        font-size: 1.5rem;
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
