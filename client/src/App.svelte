<script lang="ts">
  import svelteLogo from "./assets/svelte.svg";
  import { authenticate, isAdmin } from "./lib/auth/auth";
  import Questions from "./components/Questions.svelte";

  // Finally, your application's global stylesheet (sometimes labeled 'app.css')
  import "./app.postcss";
  import HeaderBar from "./components/HeaderBar.svelte";
  import Session from "./components/Session.svelte";
  import { P } from "flowbite-svelte";
</script>

<div class="bg-white dark:bg-gray-800 h-full">
  <HeaderBar />
  {#await authenticate()}
    <P size="xl" class="text-center">logging in</P>
  {:then _}
    <main class="container mx-auto p-8 space-y-8 overlow-y-auto h-full">
      {#if isAdmin}
        <Session />
      {/if}
      <Questions />
    </main>
  {/await}
</div>
