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
                <h2 slot="header">Frage bearbeiten</h2>
                <button slot="action" on:click={() => saveEdit()}
                    >Speichern</button
                >
                <div class="container">
                    <textarea
                        bind:value={activeQuestion.Text}
                        cols="80"
                        rows="5"
                    />
                    <label>
                        <input
                            type="checkbox"
                            bind:checked={activeQuestion.Anonymous}
                        />
                        Frage anonym stellen
                    </label>
                </div>
            </Modal>
        {/if}
    </div>
{/if}
