<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import { isAdmin, isSessionAdmin } from "../lib/auth/auth";
    import { Icon } from "@steeze-ui/svelte-icon";
    import {
        answerQuestion,
        deleteQuestion,
        voteQuestion,
    } from "../lib/questions";
    import type { Question } from "../models/question";
    import { ArrowUp, ArrowUpCircle } from "@steeze-ui/remix-icons";

    const dispatch = createEventDispatcher();

    let showModal;
    export let question: Question;

    function edit() {
        dispatch("edit", {
            id: question.Id,
        });
    }
</script>

<div class="card">
    <section class="flex gap-4 p-2">
        {#if !question.Answered}
            <div class="flex flex-col items-center justify-center">
                <div
                    class:opacity-20={question.Voted}
                    class={question.Voted
                        ? "cursor-not-allowed"
                        : "cursor-pointer"}
                    class:cursor-not-allowed={question.Voted}
                    on:click={question.Voted
                        ? null
                        : () => voteQuestion(question.Id)}
                >
                    <Icon src={ArrowUpCircle} size="48" />
                </div>
                <h3>{question.Votes}</h3>
            </div>
        {/if}
        <div class="questiontext">{question.Text}</div>
    </section>
    <hr class="opacity-50" />
    <footer class="card-footer flex justify-between gap-4 p-2">
        <div>{question.Anonymous ? "Anonym" : question.Creator}</div>
        <div>
            {#if question.Owned}
                <button
                    type="button"
                    class="btn btn-sm variant-filled"
                    on:click={() => edit()}
                >
                    Bearbeiten
                </button>
            {/if}
            {#if isAdmin || isSessionAdmin}
                <button
                    type="button"
                    class="btn btn-sm variant-filled-success"
                    on:click={() => answerQuestion(question.Id)}
                    >Beantworten</button
                >
            {/if}
            {#if question.Owned}
                <button
                    type="button"
                    class="btn btn-sm variant-filled-error"
                    on:click={() => deleteQuestion(question.Id)}>Löschen</button
                >
            {/if}
        </div>
    </footer>
</div>

<style>
</style>
