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
    import { Button, Hr, Modal, P } from "flowbite-svelte";
    import {
        CheckDouble,
        DeleteBin,
        Edit,
        ErrorWarning,
        ThumbUp,
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

<div
    class="flex flex-col gap-4 p-4 rounded border border-gray-400 shadow-lg shadow-gray-400/40 dark:border-gray-900 dark:shadow-gray-900/40"
>
    <div class="flex gap-4">
        <div class="flex flex-col items-center">
            <Button
                outline={true}
                size="sm"
                class="!p-2"
                disabled={question.Voted || null}
                on:click={question.Voted
                    ? null
                    : () => voteQuestion(question.Id)}
            >
                <Icon src={ThumbUp} size="20" />
            </Button>
            <P size="lg">{question.Votes}</P>
        </div>
        <div class="flex w-full">
            <P class="grow whitespace-pre-line" size="base">{question.Text}</P>
            {#if question.Owned}
                <Button
                    pill={true}
                    outline
                    class="!p-2 self-start"
                    on:click={() => edit()}
                    ><Icon src={Edit} size="20" /></Button
                >
            {/if}
        </div>
    </div>
    <Hr height="h-px" />

    <div class="flex-col sm:flex-row gap-2 flex items-center justify-between">
        <P size="sm" class="sm:w-full"
            >{question.Anonymous ? "Anonym" : question.Creator}</P
        >

        <div
            class="flex gap-4 justify-between sm:justify-end w-full sm:gap-2 sm:w-full justify-end"
        >
            {#if $isAdmin || $isSessionAdmin}
                <Button
                    size="xs"
                    color="green"
                    on:click={() => answerQuestion(question.Id)}
                    ><Icon class="mr-2" src={CheckDouble} size="20" /> Beantworten</Button
                >
            {/if}
            {#if question.Owned}
                <Button
                    size="xs"
                    color="red"
                    on:click={() => (popupModal = true)}
                    ><Icon
                        class="mr-2"
                        src={DeleteBin}
                        size="20"
                    />Löschen</Button
                >
            {/if}
        </div>
    </div>
    <Modal bind:open={popupModal} size="xs" autoclose>
        <div class="text-center">
            <div class="flex flex-col gap-4 p-8 items-center">
                <Icon src={ErrorWarning} size="64" />
                <P size="base">Willst du die Frage wirklich löschen?</P>
            </div>
            <div class="flex w-full justify-end gap-4">
                <Button color="alternative">Abbrechen</Button>
                <Button
                    color="red"
                    class="mr-2"
                    on:click={() => deleteQuestion(question.Id)}
                    ><Icon class="mr-2" src={DeleteBin} size="20" /> Löschen</Button
                >
            </div>
        </div>
    </Modal>
</div>

<style>
</style>
