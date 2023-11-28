<script>
	import { userOnline } from '$lib/session';
	import pollerrLogo from '$lib/logo.svg';
	import { Github } from '@steeze-ui/remix-icons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import DarkModeToggle from './DarkModeToggle.svelte';
	import MobileNav from './MobileNav.svelte';
	import DesktopNav from './DesktopNav.svelte';
</script>

<header
	class="supports-[backdrop-filter]:bg-background/60 sticky top-0 z-50 border-b bg-background/95 backdrop-blur"
>
	<div class="px-4 flex h-14 items-center">
		<div class="mr-6 flex items-center space-x-2">
			<img class="h-8 w-8" src={pollerrLogo} alt="" />
			<span class="hidden font-bold sm:inline-block text-base"> Pollerr </span>
		</div>

		<div class="flex w-full items-center space-x-4">
			<MobileNav></MobileNav>
			<DesktopNav></DesktopNav>
			<div class="flex justify-end items-center flex-1 space-x-4">
				<DarkModeToggle></DarkModeToggle>
				<div class="dark:text-white">online: {$userOnline}</div>
				<a href="https://github.com/Undermyspell/cp-voting-tool" target="_blank">
					<Icon src={Github} theme="solid" size="20" class="cursor-pointer" />
				</a>
			</div>
		</div>
	</div>
</header>
