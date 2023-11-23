<script lang="ts">
	import { answeredQuestions, getQuestions } from '$lib/questions';
	import { css } from 'styled-system/css';
	import { flex } from 'styled-system/patterns';

	$: promise = getQuestions();
	$: questions = answeredQuestions;
</script>

<div class={css({ padding: '4' })}>
	{#if $questions.length === 0}
		<p class={css({ textAlign: 'center', fontSize: '2xl' })}>
			keine beantwortete Fragen vorhanden.
		</p>
	{:else}
		<div class={flex({ gap: '4', width: '100%' })}>
			{#each $questions as question, i}
				<div
					class={css({
						padding: '4',
						display: 'flex',
						flexDirection: 'column',
						width: 'full',
						gap: '2'
					})}
				>
					<div>
						{question.Text}
					</div>
					<hr class={css({ color: 'gray.300' })} />
					<div class={flex({ width: '100%', justifyContent: 'end', gap: '4' })}>
						<p>Votes: {question.Votes}</p>
						<p>{question.Anonymous ? 'Anonym' : question.Creator}</p>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>
