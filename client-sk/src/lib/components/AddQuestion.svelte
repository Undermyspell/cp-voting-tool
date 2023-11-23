<script lang="ts">
    import { Icon } from "@steeze-ui/svelte-icon";
    import { MailAdd, RadioButton } from "@steeze-ui/remix-icons";
	import { postQuestion } from "$lib/questions";
	import { Constants } from "$lib/constants";
	import { css } from "styled-system/css";
	import { textarea, button } from "styled-system/recipes";
	import { hstack } from "styled-system/patterns";

    let value = "";
    let anonymous = true;
    const addNewQuestion = async () => {
        await postQuestion(value, anonymous);
        value = "";
    };

    $: showMaxLengthHint = value.length === Constants.QuestionMaxLength;
</script>

<div class={css({ display: "flex", flexDirection: "column", gap: 4 })}>
    <textarea
        id="message"
        class={textarea()}
        rows="3"
        maxlength={Constants.QuestionMaxLength}
        bind:value
        placeholder="Stelle eine Frage..."
    />

    <div class={css({ display: "flex", gap: 4, justifyContent: "end", alignItems: "center" })}>
        {#if showMaxLengthHint}
            <p color="red" class={css({ fontSize: "sm", color: "red" })}
                ><span
                    >{`Die Frage muss kürzer als ${Constants.QuestionMaxLength} Zeichen sein`}</span
                ></p
            >
        {/if}
        <label>
            <input type=checkbox bind:checked={anonymous}/>
            Frage anonym stellen
        </label>

        <button on:click={addNewQuestion} class={button()}>
            <div class={hstack({ gap: 2 })}>
                <Icon src={MailAdd} size="20" />
                <div>Posten</div>
            </div>
        </button>
    </div>
</div>
