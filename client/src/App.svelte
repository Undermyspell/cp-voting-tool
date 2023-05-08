<script lang="ts">
  import { authenticate } from "./lib/auth/auth";
  import Questions from "./components/Questions.svelte";

  import "./app.postcss";
  import HeaderBar from "./components/HeaderBar.svelte";
  import FooterBar from "./components/FooterBar.svelte";
  import { P, TabItem, Tabs } from "flowbite-svelte";
  import { Icon } from "@steeze-ui/svelte-icon";
  import { Chat2, ChatCheck } from "@steeze-ui/remix-icons";
  import AnsweredQuestions from "./components/AnsweredQuestions.svelte";
</script>

<div class="q-container h-full">
  <header>
    <HeaderBar />
  </header>

  {#await authenticate()}
    <P size="xl" class="text-center">logging in</P>
  {:then _}
    <main class="mx-auto md:p-8 h-full w-full max-w-[80rem] overflow-hidden">
      <Tabs
        style="underline"
        contentClass="pt-2 bg-white rounded-lg dark:bg-gray-800 h-full tab"
      >
        <TabItem open>
          <div slot="title" class="flex items-center gap-2">
            <Icon src={Chat2} size="20" />
            <div>Aktive Fragen</div>
          </div>
          <Questions />
        </TabItem>
        <TabItem>
          <div slot="title" class="flex items-center gap-2">
            <Icon src={ChatCheck} size="20" />
            <div>Beantwortete Fragen</div>
          </div>
          <AnsweredQuestions />
        </TabItem>
      </Tabs>
    </main>
  {/await}

  <footer class="flex justify-center items-center">
    <FooterBar />
  </footer>
</div>

<style>
  .q-container {
    display: grid;
    grid-template-areas:
      "header"
      "content"
      "footer";
    grid-template-rows: auto 1fr 4rem;
  }

  header {
    grid-area: header;
  }

  main {
    grid-area: content;
  }

  footer {
    grid-area: footer;
  }
</style>
