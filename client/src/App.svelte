<script lang="ts">
  import svelteLogo from "./assets/svelte.svg";
  import { authenticate, isAdmin } from "./lib/auth/auth";
  import Questions from "./components/Questions.svelte";
  import Session from "./components/Session.svelte";
</script>

{#await authenticate()}
  <div>logging in</div>
{:then _}
  <main>
    {#if isAdmin}
      <Session />
    {/if}
    <div class="questions">
      <Questions />
    </div>
  </main>
{/await}

<style>
  .questions {
    margin-top: 4rem;
  }
</style>
