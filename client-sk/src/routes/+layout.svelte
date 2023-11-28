<script>
	import '../app.pcss';
	import { authenticate } from '$lib/auth/auth';
	import HeaderBar from '$lib/components/HeaderBar.svelte';
	import '../app.pcss';
	import { page } from '$app/stores';
	import { activeSessison } from '$lib/session';
	import { ModeWatcher } from 'mode-watcher';
	import { answeredCount, getQuestions, unAnsweredCount } from '$lib/questions';
	import { cn } from '$lib/utils';

</script>

<HeaderBar></HeaderBar>
<ModeWatcher />

{#await authenticate()}
	<main class="bg-gray-100 dark:bg-gray-800 flex min-h-screen">
		<div class="container">
			<p>logging in</p>
		</div>
	</main>
{:then _}
	<main class="relative min-h-screen bg-gray-100 dark:bg-gray-800">
		<div class="flex flex-col container">
			{#if $activeSessison === false}
				<div class="text-center my-8 text-2xl">
					<p>keine aktive Q &amp; A Session</p>
				</div>
			{:else}
				<nav class="flex space-x-8 py-4 text-md sm:text-xl">
					<a
						class={cn(
							'transition-colors hover:text-foreground/80',
							$page.url.pathname === '/' ? 'text-foreground' : 'text-foreground/60'
						)}
						href="/">Unbeantwortet ({$unAnsweredCount})</a
					>
					<a
						class={cn(
							'transition-colors hover:text-foreground/80',
							$page.url.pathname === '/complete' ? 'text-foreground' : 'text-foreground/60'
						)}
						href="/complete">Beantwortet ({$answeredCount})</a
					>
				</nav>
				<div class="flex-1">
					<slot />
				</div>
			{/if}
		</div>
	</main>
{/await}
