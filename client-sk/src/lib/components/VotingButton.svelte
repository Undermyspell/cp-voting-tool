<script lang="ts">
	import type { Question } from '$lib/models/question';
	import { undoVoteQuestion, voteQuestion } from '$lib/questions';
	import { ArrowUpS, ArrowDownS } from '@steeze-ui/remix-icons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { css } from 'styled-system/css';
	
	export let question: Question;
</script>

<div class={css({display: 'flex', flexDirection: 'column', gap: '1'})}>
	<button
		class={css({
			gap: 1,
			display: 'flex',
			flexDirection: 'column',
			alignItems: 'center',
			borderRadius: 'sm',
			fontWeight: 'bold',
			backgroundColor: 'gray.100',
			padding: '1',
      width: '72px',
      cursor: 'pointer',
      '&:hover': {
        backgroundColor: 'gray.500',
        color: 'white'
      }
		})}
		on:click={question.Voted
			? () => undoVoteQuestion(question.Id)
			: () => voteQuestion(question.Id)}
	>
		{#if question.Voted}
			<Icon theme="solid" src={ArrowDownS} size="20" />
		{/if}
		{#if !question.Voted}
			<Icon theme="solid" src={ArrowUpS} size="20" />
		{/if}
		<p class={css({ fontSize: 'xl' })}>{question.Votes}</p>
    </button>
	<p>{question.Anonymous ? 'Anonym' : question.Creator}</p>
</div>
