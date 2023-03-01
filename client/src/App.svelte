<script lang="ts">
  import svelteLogo from "./assets/svelte.svg";
  import { authenticate } from "./lib/auth/auth";
  import Questions from "./Questions.svelte";
</script>

{#await authenticate()}
  <div>logging in</div>
{:then _}
  <main>
    <Questions />
  </main>
{/await}

<style>
  .logo {
    height: 6em;
    padding: 1.5em;
    will-change: filter;
    transition: filter 300ms;
  }
  .logo:hover {
    filter: drop-shadow(0 0 2em #646cffaa);
  }
  .logo.svelte:hover {
    filter: drop-shadow(0 0 2em #ff3e00aa);
  }
  .read-the-docs {
    color: #888;
  }
</style>
