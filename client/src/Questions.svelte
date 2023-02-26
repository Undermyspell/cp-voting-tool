<script lang="ts">
    import { each } from "svelte/internal";
    import { get } from "svelte/store";
    import { accessToken, idToken } from "./lib/auth/auth";
    import { getQuestions, questions } from "./lib/questions";

    const eventSource = new EventSource("http://localhost:3333/api/v1/events", {
        headers: {
            Authorization: `Bearer ${get(idToken)}`,
        },
    } as any);

    $: promise = getQuestions(eventSource);
</script>

{#each $questions as question}
    <div>{JSON.stringify(question)}</div>
{/each}
