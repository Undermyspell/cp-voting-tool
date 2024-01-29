<script lang="ts">
	import { Icon } from '@steeze-ui/svelte-icon';
	import { MailAdd } from '@steeze-ui/remix-icons';
	import { postQuestion } from '$lib/questions';
	import { Constants } from '$lib/constants';
	import { Button } from '$lib/components/ui/button';
	import { Textarea } from '$lib/components/ui/textarea';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Label } from '$lib/components/ui/label';
	import { mediaQuery } from 'svelte-legos';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as Drawer from '$lib/components/ui/drawer';
	import Rules from './Rules.svelte';
	import { Divide } from 'lucide-svelte';
	import Separator from './ui/separator/separator.svelte';

	let value = '';
	let anonymous = true;
	const addNewQuestion = async () => {
		await postQuestion(value, anonymous);
		value = '';
		open = false;
	};
	let open = false;
	const isDesktop = mediaQuery('(min-width: 768px)');
	$: showMaxLengthHint = value.length === Constants.QuestionMaxLength;
</script>

{#if $isDesktop}
	<Dialog.Root bind:open>
		<Dialog.Trigger asChild let:builder>
			<div class="flex w-full items-center justify-center">
				<Button size="lg" builders={[builder]}>
					<Icon class="mr-2" src={MailAdd} size="16" />
					<div>Frage stellen</div>
				</Button>
			</div>
		</Dialog.Trigger>
		<Dialog.Content class="sm:max-w-[1024px]">
			<Dialog.Header>
				<Dialog.Title>Frage stellen</Dialog.Title>
				<Dialog.Description>Stelle eine neue Frage zur Diskussion</Dialog.Description>
			</Dialog.Header>
			<div class="flex gap-16">
				<Textarea
					id="message"
					rows={3}
					maxlength={Constants.QuestionMaxLength}
					bind:value
					placeholder="Stelle eine Frage..."
				/>
				<Rules></Rules>
			</div>
			{#if showMaxLengthHint}
					<div>{`Die Frage muss kürzer als ${Constants.QuestionMaxLength} Zeichen sein`}</div>
			{/if}
			<Separator class="my-4" />
			<div class="flex justify-end gap-4 items-center">
				<div class="flex items-center space-x-2">
					<Checkbox id="anonymous" bind:checked={anonymous} />
					<Label for="anonymous">Frage anonym stellen?</Label>
				</div>

				<Button on:click={addNewQuestion}   disabled={value === ''}>
					<Icon class="mr-2" src={MailAdd} size="16" />
					<div>Posten</div>
				</Button>
			</div>
		</Dialog.Content>
	</Dialog.Root>
{:else}
	<Drawer.Root bind:open>
		<Drawer.Trigger asChild let:builder>
			<div class="flex w-full items-center justify-center">
				<Button size="lg" builders={[builder]}>
					<Icon class="mr-2" src={MailAdd} size="16" />
					<div>Frage stellen</div>
				</Button>
			</div>
		</Drawer.Trigger>
		<Drawer.Content>
			<div>
				<Drawer.Header class="text-left">
					<Drawer.Title>Frage stellen</Drawer.Title>
					<Drawer.Description>Stelle eine neue Frage zur Diskussion</Drawer.Description>
				</Drawer.Header>
				

				<div class="flex flex-col space-y-4 items-center">
					<Rules></Rules>
				<Textarea
					id="message"
					rows={3}
					maxlength={Constants.QuestionMaxLength}
					bind:value
					placeholder="Stelle eine Frage..."
				/>
					{#if showMaxLengthHint}
						<p color="red">
							<span>{`Die Frage muss kürzer als ${Constants.QuestionMaxLength} Zeichen sein`}</span>
						</p>
					{/if}
					<div class="flex items-center space-x-2">
						<Checkbox id="anonymous" bind:checked={anonymous} />
						<Label for="anonymous">Frage anonym stellen?</Label>
					</div>
				</div>
				<Separator class="my-4"></Separator>
				<Drawer.Footer class="mt-4 flex flex-col space-y-4">
					<Drawer.Close asChild let:builder>
						<Button disabled={value === ''} builders={[builder]} on:click={addNewQuestion}>
							<Icon class="mr-2" src={MailAdd} size="16" />
							<div>Posten</div>
						</Button>
						<Button variant="outline" builders={[builder]}>Abbrechen</Button>
					</Drawer.Close>
				</Drawer.Footer>
			</div>
		</Drawer.Content>
	</Drawer.Root>
{/if}
