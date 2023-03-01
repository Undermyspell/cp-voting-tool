<script lang="ts">
  import svelteLogo from "./assets/svelte.svg";
  import { authenticate } from "./lib/auth/auth";
  import Questions from "./components/Questions.svelte";
  import Session from "./components/Session.svelte";
</script>

{#await authenticate()}
  <div>logging in</div>
{:then _}
  <main>
    <Session />
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
