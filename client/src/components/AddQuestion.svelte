<script lang="ts">
    import { postQuestion } from "../lib/questions";
    import { Icon } from "@steeze-ui/svelte-icon";
    import { MailAdd } from "@steeze-ui/remix-icons";
    import { Textarea, Checkbox, Label, Button } from "flowbite-svelte";

    let value = "";
    let anonymous = true;
    const addNewQuestion = async () => {
        await postQuestion(value, anonymous);
        value = "";
    };
</script>

<div class="flex flex-col gap-4 mb-8">
    <Label for="message" class="mb-2">Frage eingeben</Label>
    <Textarea
        id="message"
        class="resize-none"
        rows="3"
        maxlength="500"
        bind:value
        placeholder="Stelle eine Frage..."
    />

    <div class="flex justify-around sm:justify-end gap-4 items-center">
        <Checkbox bind:checked={anonymous}>Frage anonym stellen</Checkbox>

        <Button color="blue" size="xs" on:click={addNewQuestion}>
            <span><Icon src={MailAdd} size="20" /></span><span class="ml-2"
                >Posten</span
            >
        </Button>
    </div>
</div>
