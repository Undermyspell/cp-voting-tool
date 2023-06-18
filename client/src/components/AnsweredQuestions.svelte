<script lang="ts">
    import { answeredQuestions } from "../lib/questions";
    import { activeSessison } from "../lib/session";
    import { Hr, P } from "flowbite-svelte";
    import { Icon } from "@steeze-ui/svelte-icon";
    import { CheckDouble } from "@steeze-ui/remix-icons";
</script>

{#if !$activeSessison}
    <P size="2xl" class="text-center">keine aktive Q & A Session</P>
{:else}
    <div class="py-4">
        {#if $answeredQuestions.length === 0}
            <P size="2xl" class="text-center"
                >keine beantwortete Fragen vorhanden.</P
            >
        {:else}
            <div class="flex flex-col gap-4">
                {#each $answeredQuestions as question, i}
                    <div
                        class="border border-gray-400 shadow-lg shadow-gray-400/40 dark:border-gray-900 dark:shadow-gray-900/40 p-4 rounded"
                    >
                        <div class="flex w-full items-center gap-4">
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
