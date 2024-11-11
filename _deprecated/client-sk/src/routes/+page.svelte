<script>
	import Questions from '$lib/components/Questions.svelte';
	import AddQuestion from '$lib/components/AddQuestion.svelte';
</script>

<div class="flex mt-4 space-y-4 flex-col max-h-full overflow-hidden mb-8">
	<AddQuestion />
	<Questions></Questions>
</div>
