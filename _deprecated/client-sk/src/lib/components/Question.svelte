<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { isAdmin, isSessionAdmin } from '$lib/auth/auth';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { answerQuestion } from '$lib/questions';
	import type { Question } from '$lib/models/question';
	import { CheckDouble, DeleteBin } from '@steeze-ui/remix-icons';
	import VotingButton from './VotingButton.svelte';
	import { Button } from '$lib/components/ui/button';

	const dispatch = createEventDispatcher<{
		edit: { id: string };
		delete: { id: string };
	}>();

	export let question: Question;

	function onEdit() {
		dispatch('edit', {
			id: question.Id
		});
	}

	function onDelete() {
		dispatch('delete', {
			id: question.Id
		});
	}
</script>

<div class="flex flex-col sm:flex-row w-full">
	<div class="flex-0 sm:mr-2">
		<VotingButton {question} />
	</div>
	<div
		class="flex flex-col sm:flex-row justify-between space-x-4 bg-white dark:bg-gray-900 p-4 shadow-md rounded-bottom-sm sm:rounded-sm w-full"
	>
		<div>
			<div class="whitespace-pre-wrap">{question.Text}</div>
			<p class="text-sm italic mt-4">{question.Anonymous ? 'Anonym' : question.Creator}</p>
		</div>
		<div class="items-center flex space-x-8 sm:space-x-4 mt-4 sm:mt-0">
			{#if question.Owned}
				<Button variant="outline" on:click={() => onEdit()}>
					<div>Editieren</div>
				</Button>
			{/if}
			{#if $isAdmin || $isSessionAdmin}
				<Button
					variant="outline"
					class="bg-green-500/20 text-green-600 border-green-500/50 hover:text-white hover:bg-green-600"
					size="icon"
					on:click={() => answerQuestion(question.Id)}
				>
					<Icon src={CheckDouble} size="20" />
				</Button>
			{/if}
			{#if question.Owned}
				<Button
					variant="outline"
					class="bg-red-500/20 text-red-600 border-red-500/50 hover:text-white hover:bg-red-600"
					size="icon"
					on:click={() => onDelete()}
				>
					<Icon src={DeleteBin} size="20" />
				</Button>
			{/if}
		</div>
	</div>
</div>

<style>
</style>
