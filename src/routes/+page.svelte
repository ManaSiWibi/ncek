
<script>
	import DomainResults from '$lib/components/DomainResults.svelte';
	import { checkComprehensive } from '$lib/api.js';
	
	/** @type {any} */
	let number = $state();
	let loading = $state(false);

	/**
	 * @param {any} event
	 */
	async function roll(event) {
		event.preventDefault();
		loading = true;
		number = null;
		
		const data = new FormData(event.target);
		const domain = String(data.get('description') || 'google.com');
		
		try {
			number = await checkComprehensive(domain);
		} catch (error) {
			console.error('Error fetching data:', error);
			number = { error: 'Failed to fetch data' };
		} finally {
			loading = false;
		}
	}

	const tools = [
		{
			title: 'SSL Certificate',
			href: '/ssl',
			description: 'Check SSL certificate validity and expiry dates',
			icon: '🔒'
		},
		{
			title: 'DNS Lookup',
			href: '/dns',
			description: 'Lookup A, AAAA, CNAME, MX, TXT, and NS records',
			icon: '🌐'
		},
		{
			title: 'HTTP/3 Test',
			href: '/http3',
			description: 'Check if a website supports HTTP/3 and QUIC',
			icon: '⚡'
		},
		{
			title: 'IP Lookup',
			href: '/ip',
			description: 'Get information about IP addresses',
			icon: '📍'
		},
		{
			title: 'Email Config',
			href: '/email-config',
			description: 'Check SPF, DKIM, DMARC, and BIMI records',
			icon: '📧'
		},
		{
			title: 'Web Settings',
			href: '/web-settings',
			description: 'Check web server headers and settings',
			icon: '⚙️'
		},
		{
			title: 'Sitemap Check',
			href: '/sitemap',
			description: 'Check and analyze sitemap.xml files',
			icon: '🗺️'
		}
	];
</script>

<div class="w-full">
	<!-- Hero Section -->
	<section class="py-12 sm:py-16 md:py-20 lg:py-24">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="text-center max-w-4xl mx-auto">
				<h1 class="text-4xl sm:text-5xl md:text-6xl font-bold text-gray-900 mb-6">
					Network & Website
					<span class="text-blue-600">Check</span>
				</h1>
				<p class="text-lg sm:text-xl md:text-2xl text-gray-600 mb-8 leading-relaxed">
					Comprehensive tools to analyze SSL certificates, DNS records, HTTP/3 support, email configuration, and more. 
					All in one place, completely free.
				</p>
				
				<!-- Main Check Form -->
				<div class="max-w-2xl mx-auto mt-10">
					<form onsubmit={roll} class="space-y-4">
						<div class="flex flex-col sm:flex-row gap-3">
							<input
								id="domain"
								name="description"
								type="text"
								placeholder="Enter domain or IP address (e.g., google.com)"
								autocomplete="off"
								class="flex-1 px-5 py-4 text-base border-2 border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all duration-200 shadow-sm"
								disabled={loading}
							/>
							<button
								type="submit"
								disabled={loading}
								class="bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white font-semibold py-4 px-8 rounded-lg transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 shadow-md hover:shadow-lg whitespace-nowrap"
							>
								{loading ? 'Checking...' : 'Run Comprehensive Check'}
							</button>
						</div>
					</form>
				</div>
			</div>
		</div>
	</section>

	<!-- Results Section -->
	{#if number}
		<section class="py-8 pb-16">
			<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
				<DomainResults {number} />
			</div>
		</section>
	{/if}

	<!-- Features Section -->
	<section class="py-16 sm:py-20 md:py-24 bg-gray-50">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="text-center mb-12 md:mb-16">
				<h2 class="text-3xl sm:text-4xl md:text-5xl font-bold text-gray-900 mb-4">
					Powerful Tools for Network Analysis
				</h2>
				<p class="text-lg sm:text-xl text-gray-600 max-w-3xl mx-auto">
					Get detailed insights into your website's security, performance, and configuration
				</p>
			</div>

			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 md:gap-8">
				{#each tools as tool}
					<a
						href={tool.href}
						data-sveltekit-reload
						class="group bg-white rounded-xl p-6 shadow-md hover:shadow-xl transition-all duration-200 border border-gray-200 hover:border-blue-300"
					>
						<div class="flex items-start gap-4">
							<div class="text-4xl flex-shrink-0">{tool.icon}</div>
							<div class="flex-1">
								<h3 class="text-xl font-semibold text-gray-900 mb-2 group-hover:text-blue-600 transition-colors">
									{tool.title}
								</h3>
								<p class="text-gray-600 text-sm leading-relaxed">
									{tool.description}
								</p>
							</div>
						</div>
					</a>
				{/each}
			</div>
		</div>
	</section>

	<!-- Additional Features -->
	<section class="py-16 sm:py-20 md:py-24">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="grid grid-cols-1 md:grid-cols-3 gap-8 md:gap-12">
				<div class="text-center">
					<div class="inline-flex items-center justify-center w-16 h-16 bg-blue-100 rounded-full mb-4">
						<svg class="w-8 h-8 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
						</svg>
					</div>
					<h3 class="text-xl font-semibold text-gray-900 mb-2">Secure & Private</h3>
					<p class="text-gray-600">
						All checks are performed securely without storing any personal data
					</p>
				</div>

				<div class="text-center">
					<div class="inline-flex items-center justify-center w-16 h-16 bg-blue-100 rounded-full mb-4">
						<svg class="w-8 h-8 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
						</svg>
					</div>
					<h3 class="text-xl font-semibold text-gray-900 mb-2">Fast & Reliable</h3>
					<p class="text-gray-600">
						Get results in seconds with our optimized network checking tools
					</p>
				</div>

				<div class="text-center">
					<div class="inline-flex items-center justify-center w-16 h-16 bg-blue-100 rounded-full mb-4">
						<svg class="w-8 h-8 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
						</svg>
					</div>
					<h3 class="text-xl font-semibold text-gray-900 mb-2">Completely Free</h3>
					<p class="text-gray-600">
						No registration required. Use all features without any limitations
					</p>
				</div>
			</div>
		</div>
	</section>
</div>
