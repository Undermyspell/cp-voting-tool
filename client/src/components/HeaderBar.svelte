<script lang="ts">
    import { activeSessison, startSession, stopSession, userOnline } from "../lib/session";
    import pollerrLogo from "../assets/logo.svg";
    import { Icon } from "@steeze-ui/svelte-icon";
    import { Github } from "@steeze-ui/remix-icons";

    import {
        DarkMode,
        Navbar,
        NavBrand,
        NavHamburger,
        NavLi,
        NavUl,
    } from "flowbite-svelte";
    import { isAdmin } from "../lib/auth/auth";
    import { export2csv } from "../lib/export";
</script>

<Navbar let:hidden let:toggle>
    <NavBrand href="/">
        <img src={pollerrLogo} class="h-8" alt="" />
        <span
            class="ml-4 self-center whitespace-nowrap text-xl font-semibold dark:text-white"
        >
            Pollerr
        </span>
    </NavBrand>
    <NavHamburger on:click={toggle} />
    <NavUl {hidden}>
        {#if $isAdmin}
            {#if !$activeSessison}
            <NavLi class="cursor-pointer" on:click={startSession}
                >Fragerunde starten</NavLi
            >
            {:else}
            <NavLi class="cursor-pointer" on:click={stopSession}
                >Fragerunde beenden</NavLi
            >
            <NavLi class="cursor-pointer" on:click={export2csv}>Fragen exportieren</NavLi>
            {/if}
        {/if}
    </NavUl>

    <div class="flex flex-row items-center gap-4">
        <DarkMode />
        <div class="dark:text-white">online: {$userOnline}</div>
        <a
            href="https://github.com/Undermyspell/cp-voting-tool"
            target="_blank"
        >
            <Icon src={Github} theme="solid" size="20" class="cursor-pointer" />
        </a>
    </div>
</Navbar>
