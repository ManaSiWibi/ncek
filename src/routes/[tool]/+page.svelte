<script>
	import ResultCard from '$lib/components/ResultCard.svelte';
	import { fetchApiData } from '$lib/api.js';
	
	let { data } = $props();
	const config = data?.config;
	
	/** @type {any} */
	let result = $state(null);
	let loading = $state(false);
	
	if (!config) {
		throw new Error('Tool not found');
	}
	
	/**
	 * @param {any} event
	 */
	async function checkTool(event) {
		event.preventDefault();
		loading = true;
		result = null;
		const form = new FormData(event.target);
		const value = String(form.get(config.formField) || '');
		
		try {
			// All checks now use API
			const res = await fetchApiData('/api', config?.apiParams(value) || {});
			result = config?.transform ? config.transform(res) : res;
		} catch (err) {
			const errorMessage = err instanceof Error ? err.message : String(err);
			result = { error: 'Request failed', details: errorMessage };
		} finally {
			loading = false;
		}
	}
</script>

<div class="min-h-screen flex py-6 sm:py-8 md:py-10 px-4">
	<div class="w-full max-w-7xl mx-auto">
		<!-- Header -->
		<div class="text-center mb-6 sm:mb-8">
			<h1 class="text-2xl sm:text-3xl md:text-4xl font-bold text-gray-900 mb-2 sm:mb-3 md:mb-4">
				{config.title}
			</h1>
			<p class="text-sm sm:text-base md:text-lg text-gray-600">
				{config.description}
			</p>
		</div>

		<!-- Form -->
		<div class="w-full max-w-2xl mx-auto">
			<form onsubmit={checkTool} class="space-y-6">
				<div>
					<label for="value" class="block text-sm font-medium text-gray-700 mb-2">
						Input
					</label>
					<input
						id="value"
						name={config.formField}
						type="text"
						placeholder={config.placeholder}
						autocomplete="off"
						class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition-all duration-200"
						disabled={loading}
					/>
				</div>
				
				<button
					type="submit"
					disabled={loading}
					class="w-full bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white font-semibold py-3 px-4 rounded-lg transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
				>
					{loading ? 'Checking...' : 'Submit'}
				</button>
			</form>
		</div>

		<!-- Results -->
		{#if result}
			<div class="mt-8">
				{#if config.resultCard === 'raw'}
					<div class="bg-white rounded-lg shadow-md p-6 overflow-auto">
						<pre class="text-sm">{JSON.stringify(result, null, 2)}</pre>
					</div>
				{:else}
					<ResultCard type={config.resultCard} data={result} />
				{/if}
			</div>
		{/if}
	</div>
</div>
