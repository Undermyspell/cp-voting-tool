<script lang="ts">
	import * as Sheet from '$lib/components/ui/sheet';
	import { Button } from '$lib/components/ui/button';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { Menu } from '@steeze-ui/remix-icons';
	import { isAdmin } from '$lib/auth/auth';
	import { export2csv } from '$lib/export';
	import { activeSessison, startSession, stopSession } from '$lib/session';

	let open = false;
</script>

<Sheet.Root bind:open>
	<Sheet.Trigger asChild let:builder>
		<Button
			builders={[builder]}
			variant="ghost"
			class="mr-2 px-0 text-base hover:bg-transparent focus-visible:bg-transparent focus-visible:ring-0 focus-visible:ring-offset-0 md:hidden"
		>
			<Icon src={Menu} class="h-5 w-5" />
			<span class="sr-only">Toggle Menu</span>
		</Button>
	</Sheet.Trigger>
	<Sheet.Content side="left" class="pr-0">
		<div class="my-4 h-[calc(100vh-8rem)] pb-10 pl-6 overflow-auto">
			<div class="flex flex-col space-y-4">
				{#if !!$isAdmin}
					{#if $activeSessison === false}
						<Button variant="ghost" on:click={startSession}>Fragerunde starten</Button>
					{:else}
						<Button variant="ghost" on:click={stopSession}>Fragerunde beenden</Button>
						<Button variant="ghost" on:click={export2csv}>Fragen exportieren</Button>
					{/if}
				{/if}
			</div>
		</div>
	</Sheet.Content>
</Sheet.Root>
