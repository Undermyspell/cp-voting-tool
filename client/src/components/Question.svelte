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
    import { Button, Modal, P } from "flowbite-svelte";
    import {
        AlarmWarning,
        ArrowUp,
        ArrowUpCircle,
        ArrowUpS,
        DeleteBin,
        ErrorWarning,
    } from "@steeze-ui/remix-icons";

    const dispatch = createEventDispatcher();

    let popupModal = false;
    export let question: Question;

    function edit() {
        dispatch("edit", {
            id: question.Id,
        });
    }
</script>

<div class="flex w-full">
    <section class="p-2">
        {#if !question.Answered}
            <div class="flex flex-col items-center">
                <Button
                    pill={true}
                    outline={true}
                    size="xl"
                    class="!p-2"
                    on:click={question.Voted
                        ? null
                        : () => voteQuestion(question.Id)}
                >
                    <Icon src={ArrowUpS} size="32" />
                </Button>
                <P size="xl">{question.Votes}</P>
            </div>
        {/if}
    </section>
    <div class="w-full flex flex-col justify-between gap-4 p-2">
        <P size="base">{question.Text}</P>
        <div
            class="flex-col sm:flex-row gap-2 sm:gap-0 flex items-center justify-between"
        >
            <P size="sm">{question.Anonymous ? "Anonym" : question.Creator}</P>
            <div
                class="flex-col sm:flex-row flex gap-4 sm:gap-2 w-full sm:w-auto"
            >
                {#if question.Owned}
                    <Button outline size="sm" on:click={() => edit()}
                        >Bearbeiten</Button
                    >
                {/if}
                {#if isAdmin || isSessionAdmin}
                    <Button
                        outline
                        size="sm"
                        color="green"
                        on:click={() => answerQuestion(question.Id)}
                        >Beantworten</Button
                    >
                {/if}
                {#if question.Owned}
                    <Button
                        outline
                        size="sm"
                        color="red"
                        on:click={() => (popupModal = true)}>Löschen</Button
                    >
                {/if}
            </div>
        </div>
    </div>
    <Modal bind:open={popupModal} size="xs" autoclose>
        <div class="text-center">
            <div class="flex flex-col gap-4 p-8 items-center">
                <Icon src={ErrorWarning} size="64" />
                <P size="base">Willst du die Frage wirklich löschen?</P>
            </div>
            <div class="flex w-full justify-end md:gap-4">
                <Button color="alternative">Abbrechen</Button>
                <Button
                    color="red"
                    class="mr-2"
                    on:click={() => deleteQuestion(question.Id)}>Löschen</Button
                >
            </div>
        </div>
    </Modal>
</div>

<style>
</style>
