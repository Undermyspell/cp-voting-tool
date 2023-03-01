<script lang="ts">
    import { postQuestion } from "../lib/questions";
    let value = "";

    const addNewQuestion = async () => {
        await postQuestion(value);
        value = "";
    };
</script>

<div class="container">
    <textarea bind:value name="Text1" cols="80" rows="5" />
    <button on:click={addNewQuestion}>Frage stellen</button>
</div>

<style>
    .container {
        display: flex;
        gap: 16px;
        align-items: center;
    }
</style>
