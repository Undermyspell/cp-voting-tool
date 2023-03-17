<script lang="ts">
    import { userOnline } from "../lib/session";
    import { DarkMode } from "flowbite-svelte";
    import { isAdmin } from "../lib/auth/auth";
    import Session from "./Session.svelte";
</script>

<nav class="bg-white border-gray-200 px-2 sm:px-4 py-2.5 dark:bg-gray-900">
    <div class="flex items-center justify-between mx-auto">
        <div
            class="self-center text-xl font-semibold whitespace-nowrap dark:text-white"
        >
            Pollerr
        </div>
        <div class="flex flex-row items-center gap-4">
            <DarkMode />
            <div class=" dark:text-white">online: {$userOnline}</div>
        </div>
    </div>
</nav>
