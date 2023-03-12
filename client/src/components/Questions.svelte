<script lang="ts">
    import AddQuestion from "./AddQuestion.svelte";
    import Question from "./Question.svelte";
    import Modal from "./Modal.svelte";
    import {
        getQuestions,
        questions,
        getQuestion,
        updateQuestion,
    } from "../lib/questions";
    import { activeSessison } from "../lib/session";

    let showModal = false;
    let activeQuestion = { Text: "", Anonymous: true, Id: "" };

    function editMessage(event) {
        activeQuestion = getQuestion(event.detail.id);
        showModal = true;
    }

    async function saveEdit() {
        showModal = false;
        await updateQuestion(activeQuestion);
    }

    $: promise = getQuestions();
</script>

{#if !$activeSessison}
    <div>keine aktive Q & A Session</div>
{:else}
    <AddQuestion />
    <div class="container">
        {#if $questions.length === 0}
            <div>keine Fragen vorhanden</div>
        {:else}
            <div class="flex flex-col gap-4">
                {#each $questions as question}
                    <Question on:edit={editMessage} {question} />
                {/each}
            </div>
            <Modal bind:showModal>
                <button
                    type="button"
                    class="btn btn-sm variant-filled-success"
                    slot="action"
                    on:click={() => saveEdit()}>Speichern</button
                >
                <h3 slot="header" class="pb-4 text-token">Frage bearbeiten</h3>

                <div class="space-y-4 pb-4">
                    <textarea
                        class="textarea text-token resize-none"
                        rows="4"
                        cols="80"
                        bind:value={activeQuestion.Text}
                    />
                    <label class="text-token">
                        <input
                            type="checkbox"
                            class="checkbox"
                            bind:checked={activeQuestion.Anonymous}
                        />
                        Frage anonym stellen
                    </label>
                </div>
            </Modal>
        {/if}
    </div>
{/if}
