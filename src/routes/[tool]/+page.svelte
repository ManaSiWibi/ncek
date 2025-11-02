<script lang="ts">
	import ResultCard from '$lib/components/ResultCard.svelte';
	import { fetchApiData } from '$lib/api.js';
	import { page } from '$app/stores';
	
	let { data } = $props();
	
	// Make config reactive to data changes
	const config = $derived(data?.config);
	
	// Make sure we react to route changes - derive the tool from page params
	const currentTool = $derived($page.params.tool);
	
	/** @type {any} */
	let result = $state<any>(null);
	let loading = $state(false);
	
	// Reset state when tool changes
	$effect(() => {
		// Track the current tool to trigger effect on route change
		currentTool;
		result = null;
		loading = false;
	});
	
	const pageTitle = $derived(config ? `${config.title} | NCEK` : 'Tool | NCEK');
	const pageDescription = $derived(config?.description || 'Network and website checking tool');
	
	// Type-safe result card type
	const resultCardType = $derived(config.resultCard as 'ip' | 'raw' | 'dns' | 'ssl' | 'http3' | 'web_settings' | 'email_config' | 'og_image' | undefined);
	
	async function checkTool(event: SubmitEvent) {
		event.preventDefault();
		if (!config) return;
		
		loading = true;
		result = null;
		const form = event.target as HTMLFormElement;
		const formData = new FormData(form);
		const value = String(formData.get(config.formField) || '');
		
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

<svelte:head>
	<title>{pageTitle}</title>
	<meta name="description" content={pageDescription} />
	<link rel="canonical" href={`https://tools.kenadera.org${$page.url.pathname}`} />
</svelte:head>

{#key currentTool}
	{#if config}
		<div class="min-h-screen flex py-6 sm:py-8 md:py-10 px-4">
			<div class="w-full max-w-7xl mx-auto">
				<!-- Header -->
				<header class="text-center mb-6 sm:mb-8">
					<h1 class="text-2xl sm:text-3xl md:text-4xl font-bold text-gray-900 mb-2 sm:mb-3 md:mb-4">
						{config.title}
					</h1>
					<p class="text-sm sm:text-base md:text-lg text-gray-600">
						{config.description}
					</p>
				</header>

				<!-- Form -->
				<div class="w-full max-w-2xl mx-auto">
					<form onsubmit={checkTool} class="space-y-6" aria-label={`${config.title} form`}>
						<div>
							<label for="value" class="block text-sm font-medium text-gray-700 mb-2">
								Input
							</label>
							<input
								id="value"
								name={config.formField}
								type="text"
								placeholder={config.placeholder}
								autocomplete={config.formField === 'ip' ? 'off' : 'url'}
								required
								aria-required="true"
								aria-label={config.placeholder}
								class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition-all duration-200"
								disabled={loading}
								aria-describedby="input-help"
							/>
							<p id="input-help" class="sr-only">Enter {config.placeholder}</p>
						</div>
						
						<button
							type="submit"
							disabled={loading}
							aria-label={loading ? 'Checking, please wait' : `Run ${config.title} check`}
							aria-busy={loading}
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
							<ResultCard type={resultCardType} data={result} />
						{/if}
					</div>
				{/if}
			</div>
		</div>
	{:else}
		<div class="min-h-screen flex items-center justify-center">
			<p class="text-gray-600">Tool not found</p>
		</div>
	{/if}
{/key}
