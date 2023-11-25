<script lang="ts">
	import Question from './Question.svelte';
	import {
		getQuestions,
		questions,
		getQuestion,
		updateQuestion,
		deleteQuestion
	} from '$lib/questions';
	import { Constants } from '$lib/constants';
	import { css } from 'styled-system/css';
	import { flex, hstack, vstack } from 'styled-system/patterns';
	import Modal from './Modal.svelte';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { ErrorWarning, Close, DeleteBin, Check } from '@steeze-ui/remix-icons';
	import { textButton, textarea } from 'styled-system/recipes';

	let showEditModal = false;
	let showDeleteModal = false;
	let activeQuestion: Pick<Question, 'Text' | 'Anonymous' | 'Id'> = {
		Text: '',
		Anonymous: true,
		Id: ''
	};

	function handleDeleteQuestion(event: CustomEvent<{ id: string }>) {
		activeQuestion = getQuestion(event.detail.id) ?? { Text: '', Anonymous: true, Id: '' };
		showDeleteModal = true;
	}

	function handleEditMessage(event: CustomEvent<{ id: string }>) {
		activeQuestion = getQuestion(event.detail.id) ?? { Text: '', Anonymous: true, Id: '' };
		showEditModal = true;
	}

	async function saveEdit() {
		if (activeQuestion) {
			await updateQuestion(activeQuestion);
		}
		showEditModal = false;
	}

	$: promise = getQuestions();
	$: showMaxLengthHint = activeQuestion.Text.length === Constants.QuestionMaxLength;
</script>

<div class={flex({ gap: 2, direction: 'column' })}>
	{#if $questions.length === 0}
		<p class="text-center">keine Fragen vorhanden</p>
	{:else}
		{#each $questions as question, i}
			<Question on:delete={handleDeleteQuestion} on:edit={handleEditMessage} {question} />
			{#if i < $questions.length - 1}
				<hr class={css({ borderColor: 'gray.300' })} />
			{/if}
		{/each}
		<Modal bind:show={showDeleteModal}>
			<div slot="header" class={flex({ justifyContent: 'center' })}>
				<h2 class={css({ fontSize: '2xl' })}>Frage löschen</h2>
			</div>
			<div class={vstack({ gap: '4', padding: '8' })}>
				<Icon src={ErrorWarning} size="64" />
				<div>Willst du die Frage wirklich löschen?</div>
			</div>
			<div slot="actions">
				<div class={flex({ gap: '4', justifyContent: 'space-between', paddingTop: '2' })}>
					<button class={textButton()} on:click={() => (showDeleteModal = false)}>
						<div class={hstack({ gap: 2 })}>
							<Icon src={Close} size="20" /> Abbrechen
						</div>
					</button>
					<button
						class={textButton({ color: 'red' })}
						on:click={() => deleteQuestion(activeQuestion.Id)}
					>
						<div class={hstack({ gap: 2 })}>
							<Icon src={DeleteBin} size="20" /> Löschen
						</div>
					</button>
				</div>
			</div>
		</Modal>
		<Modal bind:show={showEditModal}>
			<div slot="header" class={flex({ justifyContent: 'center' })}>
				<h2 class={css({ fontSize: '2xl' })}>Frage editieren</h2>
			</div>
			<div class={vstack({ gap: '4', padding: '8' })}>
				<textarea
					class={textarea({ resize: 'none' })}
					rows="4"
					cols="80"
					maxlength={Constants.QuestionMaxLength}
					bind:value={activeQuestion.Text}
				/>
				{#if showMaxLengthHint}
					<p color="red" class={css({ fontSize: 'sm', color: 'red' })}>
						<span>{`Die Frage muss kürzer als ${Constants.QuestionMaxLength} Zeichen sein`}</span>
					</p>
				{/if}
				<label>
					<input type="checkbox" bind:checked={activeQuestion.Anonymous} />
					Frage anonym stellen
				</label>
			</div>
			<div slot="actions">
				<div class={flex({ gap: '4', justifyContent: 'space-between', paddingTop: '2' })}>
					<button class={textButton()} on:click={() => (showEditModal = false)}>
						<div class={hstack({ gap: 2 })}>
							<Icon src={Close} size="20" /> Abbrechen
						</div>
					</button>
					<button class={textButton({ color: 'green' })} on:click={() => saveEdit()}>
						<div class={hstack({ gap: 2 })}>
							<Icon src={Check} size="20" /> Speichern
						</div>
					</button>
				</div>
			</div>
		</Modal>
	{/if}
</div>
