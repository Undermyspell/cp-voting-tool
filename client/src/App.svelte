<script lang="ts">
  import svelteLogo from "./assets/svelte.svg";
  import { authenticate, isAdmin } from "./lib/auth/auth";
  import Questions from "./components/Questions.svelte";
  import Session from "./components/Session.svelte";
  import { AppShell } from "@skeletonlabs/skeleton";

  // Your selected Skeleton theme:
  import "@skeletonlabs/skeleton/themes/theme-skeleton.css";

  // This contains the bulk of Skeletons required styles:
  import "@skeletonlabs/skeleton/styles/all.css";

  // Finally, your application's global stylesheet (sometimes labeled 'app.css')
  import "./app.postcss";
  import HeaderBar from "./components/HeaderBar.svelte";
</script>

<AppShell>
  <svelte:fragment slot="header">
    <HeaderBar />
  </svelte:fragment>
  <slot>
    {#await authenticate()}
      <div>logging in</div>
    {:then _}
      <main class="container mx-auto p-8 space-y-8">
        {#if isAdmin}
          <Session />
        {/if}
        <div class="questions">
          <Questions />
        </div>
      </main>
    {/await}
  </slot>
  <svelte:fragment slot="footer">Footer</svelte:fragment>
</AppShell>

<style>
  .questions {
    margin-top: 4rem;
  }
</style>
