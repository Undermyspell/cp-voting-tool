<script>
  import Questions from "$lib/components/Questions.svelte";
  import AddQuestion from "$lib/components/AddQuestion.svelte";
import { css } from "styled-system/css";
	import { vstack } from "styled-system/patterns";
</script>
<div class={css({display: 'flex', flexDirection: 'column', gap: '8'})}>
  <AddQuestion  />
  <Questions></Questions>
</div>