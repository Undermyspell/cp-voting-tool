<script lang="ts">
  import { authenticate, isAdmin } from "./lib/auth/auth";
  import Questions from "./components/Questions.svelte";

  import "./app.postcss";
  import HeaderBar from "./components/HeaderBar.svelte";
  import { P, TabItem, Tabs } from "flowbite-svelte";
  import { Icon } from "@steeze-ui/svelte-icon";
  import { Chat2, ChatCheck } from "@steeze-ui/remix-icons";
  import AnsweredQuestions from "./components/AnsweredQuestions.svelte";
</script>

<div class="h-full">
  <HeaderBar />
  {#await authenticate()}
    <P size="xl" class="text-center">logging in</P>
  {:then _}
    <main class="container mx-auto md:p-8 space-y-8 overlow-y-auto h-full">
      <Tabs
        style="underline"
        contentClass="pt-2 bg-white rounded-lg dark:bg-gray-800"
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
</div>
