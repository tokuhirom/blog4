<script lang="ts">
	import { debounce } from '../utils';
	let { onSearch } = $props();
	let keyword = '';

	const debouncedSearch = debounce((newKeyword: string) => {
		onSearch(newKeyword);
	}, 1000);

	function onInput(event: Event) {
		const target = event.target as HTMLInputElement;
		if (target) {
			keyword = target.value;
		}
		debouncedSearch(keyword);
	}
</script>

<input type="text" placeholder="Search entries..." class="search-input" oninput={onInput} />

<style>
	.search-input {
		margin-bottom: 1rem;
		width: 100%;
		border-radius: 0.375rem;
		border: 1px solid #d1d5db;
		padding: 0.5rem;
	}
</style>
