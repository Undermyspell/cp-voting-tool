<script lang="ts">
    import AddQuestion from "./AddQuestion.svelte";
    import Question from "./Question.svelte";
    import {
        getQuestions,
        questions,
        getQuestion,
        updateQuestion,
    } from "../lib/questions";
    import { activeSessison } from "../lib/session";
    import { Button, Checkbox, Hr, Modal, P, Textarea } from "flowbite-svelte";

    let showModal = false;
    let activeQuestion = { Text: "", Anonymous: true, Id: "" };

    function editMessage(event) {
        activeQuestion = getQuestion(event.detail.id);
        showModal = true;
    }

    async function saveEdit() {
        await updateQuestion(activeQuestion);
    }

    $: promise = getQuestions();
</script>

{#if !$activeSessison}
    <P size="2xl" class="text-center">keine aktive Q & A Session</P>
{:else}
    <div class="flex flex-col">
        <AddQuestion />
        <div class="gap-4">
            {#if $questions.length === 0}
                <P size="2xl" class="text-center">keine Fragen vorhanden</P>
            {:else}
                <div class="flex-col flex gap-4">
                    {#each $questions as question, i}
                        <Question on:edit={editMessage} {question} />
                    {/each}
                </div>
                <Modal bind:open={showModal} title="Frage bearbeiten" autoclose>
                    <div class="space-y-4 pb-4">
                        <Textarea
                            class="resize-none"
                            rows="4"
                            cols="80"
                            bind:value={activeQuestion.Text}
                        />
                        <Checkbox bind:checked={activeQuestion.Anonymous}
                            >Frage anonym stellen</Checkbox
                        >
                    </div>
                    <svelte:fragment slot="footer">
                        <div class="flex gap-4 w-full justify-end">
                            <Button color="alternative">Abbrechen</Button>
                            <Button on:click={() => saveEdit()}
                                >Speichern</Button
                            >
                        </div>
                    </svelte:fragment>
                </Modal>
            {/if}
        </div>
    </div>
{/if}
