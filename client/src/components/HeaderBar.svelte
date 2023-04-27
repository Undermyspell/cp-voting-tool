<script lang="ts">
    import { startSession, stopSession, userOnline } from "../lib/session";
    import pollerrLogo from "../assets/logo.svg";

    import {
        DarkMode,
        Navbar,
        NavBrand,
        NavHamburger,
        NavLi,
        NavUl,
    } from "flowbite-svelte";
    import { isAdmin } from "../lib/auth/auth";
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
        {#if isAdmin()}
            <NavLi class="cursor-pointer" on:click={startSession}
                >Fragerunde starten</NavLi
            >
            <NavLi class="cursor-pointer" on:click={stopSession}
                >Fragerunde beenden</NavLi
            >
        {/if}
    </NavUl>

    <div class="flex flex-row items-center gap-4">
        <DarkMode />
        <div class=" dark:text-white">online: {$userOnline}</div>
    </div>
</Navbar>
