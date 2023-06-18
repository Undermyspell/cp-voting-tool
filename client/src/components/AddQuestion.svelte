<script lang="ts">
    import { postQuestion } from "../lib/questions";
    import { Icon } from "@steeze-ui/svelte-icon";
    import { MailAdd } from "@steeze-ui/remix-icons";
    import { Textarea, Checkbox, Label, Button, Helper } from "flowbite-svelte";
    import { Constants } from "../lib/constants";

    let value = "";
    let anonymous = true;
    const addNewQuestion = async () => {
        await postQuestion(value, anonymous);
        value = "";
    };

    $: showMaxLengthHinth = value.length === Constants.QuestionMaxLength;
</script>

<div class="flex flex-col gap-4 mb-8">
    <Label for="message" class="mb-2">Frage eingeben</Label>
    <Textarea
        id="message"
        class="resize-none"
        rows="3"
        maxlength={Constants.QuestionMaxLength}
        bind:value
        placeholder="Stelle eine Frage..."
    />

    <div class="flex justify-around sm:justify-end gap-4 items-center">
        {#if showMaxLengthHinth}
            <Helper color="red" class="mr-auto text-sm font-medium"
                ><span
                    >{`Die Frage muss kürzer als ${Constants.QuestionMaxLength} Zeichen sein`}</span
                ></Helper
            >
        {/if}
        <Checkbox bind:checked={anonymous}>Frage anonym stellen</Checkbox>

        <Button color="blue" size="xs" on:click={addNewQuestion}>
            <span><Icon src={MailAdd} size="20" /></span><span class="ml-2"
                >Posten</span
            >
        </Button>
    </div>
</div>
