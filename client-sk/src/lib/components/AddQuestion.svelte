<script lang="ts">
	import { Icon } from '@steeze-ui/svelte-icon';
	import { MailAdd } from '@steeze-ui/remix-icons';
	import { postQuestion } from '$lib/questions';
	import { Constants } from '$lib/constants';
	import { Button } from '$lib/components/ui/button';
	import { Textarea } from '$lib/components/ui/textarea';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Label } from '$lib/components/ui/label';

	let value = '';
	let anonymous = true;
	const addNewQuestion = async () => {
		await postQuestion(value, anonymous);
		value = '';
	};

	$: showMaxLengthHint = value.length === Constants.QuestionMaxLength;
</script>

<div class="flex flex-col p-4 rounded-sm shadow-md gap-4 bg-white dark:bg-gray-900">
	<Textarea
		id="message"
		rows={3}
		maxlength={Constants.QuestionMaxLength}
		bind:value
		placeholder="Stelle eine Frage..."
	/>

	<div class="flex justify-end gap-4 items-center">
		{#if showMaxLengthHint}
			<p color="red">
				<span>{`Die Frage muss kürzer als ${Constants.QuestionMaxLength} Zeichen sein`}</span>
			</p>
		{/if}
		<div class="flex items-center space-x-2">
			<Checkbox id="anonymous" bind:checked={anonymous} />
			<Label for="anonymous">Frage anonym stellen?</Label>
		</div>

		<Button variant="outline" on:click={addNewQuestion}>
			<Icon class="mr-2" src={MailAdd} size="16" />
			<div>Posten</div>
		</Button>
	</div>
</div>
