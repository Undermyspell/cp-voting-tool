<script lang="ts">
    import { postQuestion } from "../lib/questions";
    import { Icon } from "@steeze-ui/svelte-icon";
    import { MailAdd } from "@steeze-ui/remix-icons";
    let value = "";
    let anonymous = true;
    const addNewQuestion = async () => {
        await postQuestion(value, anonymous);
        value = "";
    };
</script>

<div class="flex flex-col gap-4 mb-8">
    <label class="label">
        <strong>Frage eingeben</strong>
        <textarea
            class="textarea resize-none"
            rows="4"
            bind:value
            name="Text1"
        />
        <label class="text-token">
            <input type="checkbox" class="checkbox" bind:checked={anonymous} />
            Frage anonym stellen
        </label>
    </label>
    <button
        class="btn variant-filled self-center"
        type="button"
        on:click={addNewQuestion}
        ><span><Icon src={MailAdd} size="20" /></span><span>Posten</span
        ></button
    >
</div>
