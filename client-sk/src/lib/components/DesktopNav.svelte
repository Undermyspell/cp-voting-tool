<script lang="ts">
	import { isAdmin } from '$lib/auth/auth';
	import { export2csv } from '$lib/export';
	import { activeSessison, startSession, stopSession } from '$lib/session';
	import { Button } from '$lib/components/ui/button';
</script>

<div class="hidden sm:flex">
	{#if !!$isAdmin}
		{#if $activeSessison === false}
			<Button variant="ghost" on:click={startSession}>Fragerunde starten</Button>
		{:else}
			<Button variant="ghost" on:click={stopSession}>Fragerunde beenden</Button>
			<Button variant="ghost" on:click={export2csv}>Fragen exportieren</Button>
		{/if}
	{/if}
</div>
