<script lang="ts">
	import { answeredQuestions, getQuestions } from '$lib/questions';
	import { Separator } from '$lib/components/ui/separator';

	$: getQuestions();
	$: questions = answeredQuestions;
</script>

<div class="p-4">
	{#if $questions.length === 0}
		<p class="text-center text-2xl">keine beantwortete Fragen vorhanden.</p>
	{:else}
		<div class="flex flex-col w-full spsace-y-4">
			{#each $questions as question, i}
				<div class="p-4 flex flex-col w-full space-y-4">
					<div>
						{question.Text}
					</div>
					<Separator class="dark:bg-gray-500 my-4" />
					<div class="flex justify-end space-x-4">
						<p>Votes: {question.Votes}</p>
						<p>{question.Anonymous ? 'Anonym' : question.Creator}</p>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>
