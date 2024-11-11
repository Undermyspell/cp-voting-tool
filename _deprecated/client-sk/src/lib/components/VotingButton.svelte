<script lang="ts">
	import type { Question } from '$lib/models/question';
	import { undoVoteQuestion, voteQuestion } from '$lib/questions';
	import { Heart } from '@steeze-ui/remix-icons';
	import { Icon } from '@steeze-ui/svelte-icon';

	export let question: Question;
</script>

<div
	class="cursor-pointer hover:scale-105 flex sm:shadow-md rounded-top-sm sm:rounded-sm flex-col items-center px-4 py-2 bg-white dark:bg-gray-900"
>
	<button
		on:click={question.Voted
			? () => undoVoteQuestion(question.Id)
			: () => voteQuestion(question.Id)}
	>
		{#if question.Voted}
			<Icon theme="solid" color="red" src={Heart} size="32" />
		{/if}
		{#if !question.Voted}
			<Icon theme="outline" src={Heart} size="32" />
		{/if}
	</button>
	<p>{question.Votes}</p>
</div>
