<script lang="ts">
	import Question from './Question.svelte';
	import {
		getQuestions,
		questions,
		getQuestion,
		updateQuestion,
		deleteQuestion,
		isAutosortActive,
		updateAutosort
	} from '$lib/questions';
	import { Constants } from '$lib/constants';
	import Modal from './Modal.svelte';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { ErrorWarning, Close, DeleteBin, Check } from '@steeze-ui/remix-icons';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Label } from '$lib/components/ui/label';
	import { Button } from '$lib/components/ui/button';
	import { Textarea } from '$lib/components/ui/textarea';

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

	async function deleteQuestionAndClose() {
		if (activeQuestion) {
			await deleteQuestion(activeQuestion.Id);
		}
		showDeleteModal = false;
	}

	$: getQuestions();
	$: autoSortActive = $isAutosortActive;
	$: showMaxLengthHint = activeQuestion.Text.length === Constants.QuestionMaxLength;
</script>

{#if $questions.length === 0}
	<p class="text-center text-2xl">Keine Fragen vorhanden</p>
{:else}
	<div class="flex items-center space-x-2">
		<Checkbox
			id="autosort"
			on:click={() => updateAutosort(!autoSortActive)}
			bind:checked={autoSortActive}
		/>
		<Label for="autosort">Fragen automatisch sortieren?</Label>
	</div>
	<div class="mt-4 space-y-4 max-h-full overflow-y-auto">
		{#each $questions as question (question.Id)}
			<Question on:delete={handleDeleteQuestion} on:edit={handleEditMessage} {question} />
		{/each}
	</div>
	<Modal bind:show={showDeleteModal}>
		<div slot="header">
			<h2 class="text-2xl">Frage löschen</h2>
		</div>
		<div class="flex my-8 items-center flex-col space-y-4">
			<Icon src={ErrorWarning} size="64" />
			<div>Willst du die Frage wirklich löschen?</div>
		</div>
		<div slot="actions">
			<div class="flex justify-end space-x-4 mt-8">
				<Button variant="outline" on:click={() => (showDeleteModal = false)}>
					<Icon class="mr-2" src={Close} size="20" /> Abbrechen
				</Button>
				<Button
					class="bg-red-500/20 text-red-600 border-red-500/50 hover:text-white hover:bg-red-600"
					on:click={() => deleteQuestionAndClose()}
				>
					<Icon class="mr-2" src={DeleteBin} size="20" /> Löschen
				</Button>
			</div>
		</div>
	</Modal>
	<Modal bind:show={showEditModal}>
		<div slot="header">
			<h2 class="text-2xl">Frage editieren</h2>
		</div>
		<div class="flex my-8 flex-col space-y-4">
			<Textarea
				maxlength={Constants.QuestionMaxLength}
				rows={8}
				cols={80}
				bind:value={activeQuestion.Text}
			/>
			{#if showMaxLengthHint}
				<p color="red">
					<span>{`Die Frage muss kürzer als ${Constants.QuestionMaxLength} Zeichen sein`}</span>
				</p>
			{/if}
			<div class="flex items-center space-x-2">
				<Checkbox id="anonymous" bind:checked={activeQuestion.Anonymous} />
				<Label for="anonymous">Frage anonym stellen?</Label>
			</div>
		</div>
		<div slot="actions">
			<div class="flex justify-end space-x-4 mt-8">
				<Button variant="outline" on:click={() => (showEditModal = false)}>
					<Icon class="mr-2" src={Close} size="20" /> Abbrechen
				</Button>
				<Button on:click={() => saveEdit()}>
					<Icon class="mr-2" src={Check} size="20" /> Speichern
				</Button>
			</div>
		</div>
	</Modal>
{/if}
