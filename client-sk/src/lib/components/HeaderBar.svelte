<script>
	import { isAdmin } from '$lib/auth/auth';
	import { export2csv } from '$lib/export';
	import { activeSessison, startSession, stopSession, userOnline } from '$lib/session';
	import pollerrLogo from '$lib/logo.svg';
	import { FontSize, Github } from '@steeze-ui/remix-icons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { css } from 'styled-system/css';
	import { textButton } from 'styled-system/recipes';
</script>
<nav class={css({ display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '8px' })}>
	<div class={css({ display: 'flex', alignItems: 'center' })}>
		<img class={css({ height: '2rem' })} src={pollerrLogo} alt="" />
		<span class={css({ marginLeft: '16px', whiteSpace: 'nowrap', fontSize: 'xl', fontWeight: 'bold' })}>
			Pollerr
		</span>
	</div>

	<div class={css({ display: 'flex', alignItems: 'center', gap: '16px' })}>
    		{#if $isAdmin}
			{#if !$activeSessison}
				<button class={textButton()} on:click={startSession}>Fragerunde starten</button>
			{:else}
				<button class={textButton()} on:click={stopSession}>Fragerunde beenden</button>
				<button class={textButton()} on:click={export2csv}>Fragen exportieren</button>
			{/if}
		{/if}
		<div class="dark:text-white">online: {$userOnline}</div>
		<a href="https://github.com/Undermyspell/cp-voting-tool" target="_blank">
			<Icon src={Github} theme="solid" size="20" class="cursor-pointer" />
		</a>
	</div>
</nav>
