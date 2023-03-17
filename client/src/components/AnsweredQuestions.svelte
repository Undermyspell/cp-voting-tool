<script lang="ts">
    import { answeredQuestions } from "../lib/questions";
    import { activeSessison } from "../lib/session";
    import { Button, Checkbox, Hr, Modal, P, Textarea } from "flowbite-svelte";
    import { Icon } from "@steeze-ui/svelte-icon";
    import { CheckDouble } from "@steeze-ui/remix-icons";
</script>

{#if !$activeSessison}
    <P size="2xl" class="text-center">keine aktive Q & A Session</P>
{:else}
    <div class="py-4">
        {#if $answeredQuestions.length === 0}
            <div>keine Fragen beantwortet bisher</div>
        {:else}
            <div class="flex flex-col gap-4 ">
                {#each $answeredQuestions as question, i}
                    <div class="border border-gray-600 p-4 rounded">
                        <div class="flex w-full items-center gap-4 ">
                            <Icon
                                src={CheckDouble}
                                size="48"
                                class="text-green-500"
                            />
                            <P>{question.Text}</P>
                        </div>
                        <Hr class="my-2" height="h-px" />
                        <div class="flex w-full justify-end gap-4">
                            <P>Votes: {question.Votes}</P>
                            <P
                                >{question.Anonymous
                                    ? "Anonym"
                                    : question.Creator}</P
                            >
                        </div>
                    </div>
                {/each}
            </div>
        {/if}
    </div>
{/if}
