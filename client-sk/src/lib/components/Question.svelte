<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { isAdmin, isSessionAdmin } from '$lib/auth/auth';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { answerQuestion } from '$lib/questions';
	import type { Question } from '$lib/models/question';
	import { CheckDouble, DeleteBin, Edit } from '@steeze-ui/remix-icons';
	import { css } from 'styled-system/css';
	import { flex, hstack } from 'styled-system/patterns';
	import { textButton } from 'styled-system/recipes';
	import VotingButton from './VotingButton.svelte';

	const dispatch = createEventDispatcher<{
		edit: {id: string},
		delete: {id: string}
	}>();

	let popupModal = false;
	export let question: Question;

	function onEdit() {
		dispatch('edit', {
			id: question.Id
		});
	}

	function onDelete() {
		dispatch('delete', {
			id: question.Id
		})
	}
</script>

<div
	class={flex({
		direction: 'column',
		gap: 4,
		padding: 4
	})}
>
	<div
		class={flex({
			gap: 4,
			direction: { base: 'column', sm: 'row' },
			alignItems: { base: 'center', sm: 'start' }
		})}
	>
		<VotingButton {question} />

		<div
			class={flex({
				direction: { base: 'column', sm: 'row' },
				width: '100%',
				justifyContent: 'space-between',
				alignContent: 'center',
				gap: { base: 4, sm: '0' }
			})}
		>
			<div class={css({ width: '100%', whiteSpace: 'pre-wrap' })}>{question.Text}</div>
			<div
				class={flex({
					direction: { base: 'row', sm: 'column' },
					justifyContent: { base: 'space-between', sm: 'start' },
					gap: '2'
				})}
			>
				{#if question.Owned}
					<button class={textButton()} on:click={() => onEdit()}>
						<div class={hstack({ gap: '2' })}>
							<Icon src={Edit} size="20" />
							<div>Editieren</div>
						</div>
					</button>
				{/if}
				{#if $isAdmin || $isSessionAdmin}
					<button
						class={textButton({ color: 'green' })}
						on:click={() => answerQuestion(question.Id)}
					>
						<div class={hstack({ gap: 2 })}>
							<Icon class="mr-2" src={CheckDouble} size="20" /> Beantworten
						</div>
					</button>
				{/if}
				{#if question.Owned}
					<button class={textButton({ color: 'red' })} on:click={() => onDelete()}
						><div class={hstack({ gap: 2 })}>
							<Icon src={DeleteBin} size="20" /> Löschen
						</div></button
					>
				{/if}
			</div>
		</div>
	</div>
</div>

<style>
</style>
