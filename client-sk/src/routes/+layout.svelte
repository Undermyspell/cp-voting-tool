<script>
	import { authenticate } from '$lib/auth/auth';
	import HeaderBar from '$lib/components/HeaderBar.svelte';
	import '../index.css';
	import { container } from '../../styled-system/patterns';
	import { css } from 'styled-system/css';
	import { page } from '$app/stores';
	import { activeSessison, getSession } from '$lib/session';

	const linkClass = css({
		padding: '2',
		fontWeight: 'semibold',
		'&.active': {
			color: 'blue.600'
		}
	});
</script>

<HeaderBar></HeaderBar>

{#await authenticate()}
	<p>logging in</p>
{:then _}
	{#if $activeSessison === false}
		<p class={css({ textAlign: 'center', fontSize: '2xl' })}>keine aktive Q & A Session</p>
	{:else}
		<main class={container()}>
			<nav class={css({ display: 'flex', gap: '4', marginTop: '2', marginBottom: '4' })}>
				<a class={linkClass} class:active={$page.url.pathname === '/'} href="/">Unbeantwortet</a>
				<a class={linkClass} class:active={$page.url.pathname === '/complete'} href="/complete"
					>Beantwortet</a
				>
			</nav>
			<slot />
		</main>
	{/if}
{/await}
