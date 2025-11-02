<script>
	/** @type {{
		data?: any;
	}} */
	let { data = null } = $props();

	// Helper function to get display values with fallbacks
	function getTitle() {
		return data?.og_title || data?.twitter_title || data?.meta_title || 'No title';
	}

	function getDescription() {
		return data?.og_description || data?.twitter_description || data?.meta_description || '';
	}

	function getImage() {
		return data?.image_url || data?.twitter_image || data?.image_secure || '';
	}

	function getSiteName() {
		return data?.og_site_name || data?.domain || '';
	}

	function getURL() {
		return data?.og_url || data?.url || '';
	}

	function getTwitterSite() {
		return data?.twitter_site || '';
	}

	/**
	 * @param {string} str
	 * @param {number} maxLen
	 */
	function truncate(str, maxLen = 200) {
		if (!str) return '';
		return str.length > maxLen ? str.substring(0, maxLen) + '...' : str;
	}

	/**
	 * @param {Event} event
	 */
	function handleImageError(event) {
		const target = event.target;
		if (target instanceof HTMLImageElement) {
			target.style.display = 'none';
		}
	}
</script>

{#if data?.error}
	<div class="bg-red-50 border border-red-200 rounded-lg p-4 sm:p-6">
		<div class="flex items-center space-x-2 text-red-600">
			<span class="text-xl sm:text-2xl">✗</span>
			<h3 class="text-lg sm:text-xl font-semibold">Error</h3>
		</div>
		<p class="mt-2 text-sm sm:text-base text-red-700 break-words">{data.error}</p>
	</div>
{:else if data}
	<div class="space-y-6">
		<!-- Platform Previews -->
		<section aria-labelledby="previews-heading">
			<h2 id="previews-heading" class="sr-only">Social Media Platform Previews</h2>
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4" role="list">
				<!-- Twitter/X Preview -->
				<article class="bg-white border border-gray-300 rounded-lg overflow-hidden shadow-sm" role="listitem" aria-label="Twitter/X preview">
					<div class="bg-gray-50 px-3 py-2 border-b border-gray-200">
						<div class="flex items-center space-x-2">
							<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
								<path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z"/>
							</svg>
							<span class="text-sm font-semibold text-gray-700">Twitter/X</span>
						</div>
					</div>
				<div class="p-4">
					{#if getImage()}
						<div class="mb-3 rounded-lg overflow-hidden bg-gray-100">
							<img 
								src={getImage()} 
								alt={data?.twitter_image_alt || getTitle()}
								class="w-full h-48 object-cover"
								onerror={handleImageError}
							/>
						</div>
					{/if}
					<div class="space-y-2">
						{#if getSiteName()}
							<div class="text-xs text-gray-500">{getSiteName()}</div>
						{/if}
						<h3 class="text-sm font-semibold text-gray-900 line-clamp-2">{getTitle()}</h3>
						{#if getDescription()}
							<p class="text-sm text-gray-600 line-clamp-3">{truncate(getDescription(), 150)}</p>
						{/if}
					</div>
				</div>
			</article>

			<!-- LinkedIn Preview -->
			<article class="bg-white border border-gray-300 rounded-lg overflow-hidden shadow-sm" role="listitem" aria-label="LinkedIn preview">
				<div class="bg-blue-50 px-3 py-2 border-b border-gray-200">
					<div class="flex items-center space-x-2">
						<svg class="w-5 h-5 text-blue-600" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
							<path d="M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.046c.477-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433c-1.144 0-2.063-.926-2.063-2.065 0-1.138.92-2.063 2.063-2.063 1.14 0 2.064.925 2.064 2.063 0 1.139-.925 2.065-2.064 2.065zm1.782 13.019H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451C23.2 24 24 23.227 24 22.271V1.729C24 .774 23.2 0 22.222 0h.003z"/>
						</svg>
						<span class="text-sm font-semibold text-gray-700">LinkedIn</span>
					</div>
				</div>
				<div class="p-4">
					{#if getImage()}
						<div class="mb-3 rounded-lg overflow-hidden bg-gray-100">
							<img 
								src={getImage()} 
								alt={getTitle()}
								class="w-full h-48 object-cover"
								onerror={handleImageError}
							/>
						</div>
					{/if}
					<div class="space-y-2">
						<h3 class="text-sm font-semibold text-gray-900 line-clamp-2">{getTitle()}</h3>
						{#if getDescription()}
							<p class="text-sm text-gray-600 line-clamp-3">{truncate(getDescription(), 150)}</p>
						{/if}
						{#if getSiteName()}
							<div class="text-xs text-gray-500">{getSiteName()}</div>
						{/if}
					</div>
				</div>
			</article>

			<!-- Slack Preview -->
			<article class="bg-white border border-gray-300 rounded-lg overflow-hidden shadow-sm" role="listitem" aria-label="Slack preview">
				<div class="bg-purple-50 px-3 py-2 border-b border-gray-200">
					<div class="flex items-center space-x-2">
						<svg class="w-5 h-5 text-purple-600" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
							<path d="M5.042 15.165a2.528 2.528 0 0 1-2.52 2.523A2.528 2.528 0 0 1 0 15.165a2.527 2.527 0 0 1 2.522-2.52h2.52v2.52zM6.313 15.165a2.527 2.527 0 0 1 2.521-2.52 2.527 2.527 0 0 1 2.521 2.52v6.313A2.528 2.528 0 0 1 8.834 24a2.528 2.528 0 0 1-2.521-2.522v-6.313zM8.834 5.042a2.528 2.528 0 0 1-2.521-2.52A2.528 2.528 0 0 1 8.834 0a2.528 2.528 0 0 1 2.521 2.522v2.52H8.834zM8.834 6.313a2.528 2.528 0 0 1 2.521 2.521 2.528 2.528 0 0 1-2.521 2.521H2.522A2.528 2.528 0 0 1 0 8.834a2.528 2.528 0 0 1 2.522-2.521h6.312zM18.956 5.042a2.528 2.528 0 0 1-2.521-2.52A2.528 2.528 0 0 1 18.956 0a2.528 2.528 0 0 1 2.522 2.522v2.52h-2.522zM18.956 6.313a2.528 2.528 0 0 1 2.522 2.521 2.528 2.528 0 0 1-2.522 2.521h-6.313A2.528 2.528 0 0 1 10.043 8.834a2.528 2.528 0 0 1 2.522-2.521h6.391zM13.043 18.956a2.528 2.528 0 0 1 2.521 2.522A2.528 2.528 0 0 1 13.043 24a2.528 2.528 0 0 1-2.521-2.522v-2.521h2.521zM13.043 17.688a2.528 2.528 0 0 1-2.521-2.522 2.528 2.528 0 0 1 2.521-2.521h6.313A2.528 2.528 0 0 1 21.878 15.166a2.528 2.528 0 0 1-2.522 2.522h-6.313z"/>
						</svg>
						<span class="text-sm font-semibold text-gray-700">Slack</span>
					</div>
				</div>
				<div class="p-4">
					{#if getImage()}
						<div class="mb-3 rounded-lg overflow-hidden bg-gray-100">
							<img 
								src={getImage()} 
								alt={getTitle()}
								class="w-full h-48 object-cover"
								onerror={handleImageError}
							/>
						</div>
					{/if}
					<div class="space-y-2">
						<h3 class="text-sm font-semibold text-gray-900 line-clamp-2">{getTitle()}</h3>
						{#if getDescription()}
							<p class="text-sm text-gray-600 line-clamp-2">{truncate(getDescription(), 120)}</p>
						{/if}
						{#if getSiteName()}
							<div class="text-xs text-gray-500">{getSiteName()}</div>
						{/if}
					</div>
				</div>
			</article>
		</div>
		</section>

		<!-- Image Details -->
		{#if getImage()}
			<section aria-labelledby="image-details-heading" class="bg-white border border-gray-200 rounded-lg p-4 sm:p-6">
				<h2 id="image-details-heading" class="text-lg font-semibold text-gray-900 mb-4">Image Details</h2>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
					<div>
						<span class="font-medium text-gray-700">Image URL:</span>
						<div class="mt-1 text-gray-600 break-all">{getImage()}</div>
					</div>
					{#if data.image_width && data.image_height}
						<div>
							<span class="font-medium text-gray-700">Dimensions:</span>
							<div class="mt-1 text-gray-600">{data.image_width} × {data.image_height}</div>
						</div>
					{/if}
					{#if data.content_type}
						<div>
							<span class="font-medium text-gray-700">Content Type:</span>
							<div class="mt-1 text-gray-600">{data.content_type}</div>
						</div>
					{/if}
					{#if data.size}
						<div>
							<span class="font-medium text-gray-700">Size:</span>
							<div class="mt-1 text-gray-600">{(data.size / 1024).toFixed(2)} KB</div>
						</div>
					{/if}
					<div>
						<span class="font-medium text-gray-700">Accessible:</span>
						<div class="mt-1">
							<span class={`inline-flex items-center px-2 py-1 rounded text-xs font-medium ${data.accessible ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}`}>
								{data.accessible ? '✓ Yes' : '✗ No'}
							</span>
						</div>
					</div>
				</div>
			</section>
		{/if}

		<!-- All Detected Tags -->
		<section aria-labelledby="tags-heading" class="bg-white border border-gray-200 rounded-lg p-4 sm:p-6">
			<h2 id="tags-heading" class="text-lg font-semibold text-gray-900 mb-4">All Detected Tags</h2>
			
			{#if data.all_meta_tags && Object.keys(data.all_meta_tags).length > 0}
				<div class="mb-6">
					<h4 class="text-sm font-semibold text-gray-700 mb-2">Open Graph Tags</h4>
					<div class="bg-gray-50 rounded-lg p-3 space-y-2 text-xs">
						{#each Object.keys(data.all_meta_tags) as key}
							<div class="flex flex-col sm:flex-row sm:items-start">
								<span class="font-medium text-gray-900 min-w-[150px]">{key}:</span>
								<span class="text-gray-600 break-all sm:ml-2">{data.all_meta_tags[key]}</span>
							</div>
						{/each}
					</div>
				</div>
			{/if}

			{#if data.all_twitter_tags && Object.keys(data.all_twitter_tags).length > 0}
				<div class="mb-6">
					<h4 class="text-sm font-semibold text-gray-700 mb-2">Twitter Card Tags</h4>
					<div class="bg-gray-50 rounded-lg p-3 space-y-2 text-xs">
						{#each Object.keys(data.all_twitter_tags) as key}
							<div class="flex flex-col sm:flex-row sm:items-start">
								<span class="font-medium text-gray-900 min-w-[150px]">{key}:</span>
								<span class="text-gray-600 break-all sm:ml-2">{data.all_twitter_tags[key]}</span>
							</div>
						{/each}
					</div>
				</div>
			{/if}

			{#if (!data.all_meta_tags || Object.keys(data.all_meta_tags).length === 0) && (!data.all_twitter_tags || Object.keys(data.all_twitter_tags).length === 0)}
				<div class="text-sm text-gray-500 italic">No meta tags detected</div>
			{/if}
		</section>
	</div>
{/if}

