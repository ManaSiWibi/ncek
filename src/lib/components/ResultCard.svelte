<script>
	import OGImageResults from './OGImageResults.svelte';

	/** @type {{
		type?: 'dns' | 'ssl' | 'http3' | 'ip' | 'web_settings' | 'email_config' | 'og_image' | 'raw' | 'error';
		data?: any;
		isOpen?: boolean;
		toggleAccordion?: ((id: string) => void) | null;
		accordionId?: string;
	}} */
	let {
		type = 'dns',
		data = null,
		isOpen = false,
		toggleAccordion = null,
		accordionId = ''
	} = $props();
	
	function getDisplayData() {
		if (!data) return [];
		
		if (type === 'dns') {
			return [
				{ label: 'Domain', value: data.domain },
				{ label: 'IPv4', value: data.ipv4, isArray: true },
				{ label: 'IPv6', value: data.ipv6, isArray: true },
				{ label: 'MX', value: data.mx, isArray: true },
				{ label: 'TXT', value: data.txt, isArray: true },
				{ label: 'CNAME', value: data.cname, isArray: true },
				{ label: 'NS', value: data.ns, isArray: true },
			];
		}
		
		if (type === 'ssl') {
			return [
				{ label: 'Domain', value: data.domain },
				{ label: 'Subject', value: data.subject },
				{ label: 'Issuer', value: data.issuer },
				{ label: 'Expires', value: data.not_after ? new Date(data.not_after).toLocaleDateString() : null },
				{ label: 'Key Size', value: data.key_size ? `${data.key_size} bits` : null },
				{ label: 'Signature', value: data.signature_algorithm },
				{ label: 'Public Key', value: data.public_key_algorithm },
			];
		}
		
		if (type === 'http3') {
			return [
				{ label: 'Domain', value: data.domain },
				{ label: 'Protocol', value: data.protocol },
				{ label: 'Status', value: data.status },
				{ label: 'Details', value: data.details },
			];
		}
		
		if (type === 'web_settings') {
			return [
				{ label: 'Domain', value: data.domain },
				{ label: 'Status', value: data.status_code },
				{ label: 'Response', value: data.response_time_ms ? `${data.response_time_ms}ms` : null },
				{ label: 'Server', value: data.server },
				{ label: 'Content Type', value: data.content_type },
				{ label: 'Size', value: data.content_length ? `${data.content_length} bytes` : null },
				{ label: 'Modified', value: data.last_modified },
				{ label: 'ETag', value: data.etag },
				{ label: 'Redirect', value: data.redirect_url },
				{ label: 'HSTS Enabled', value: data.hsts?.enabled ? '✓' : '✗' },
				{ label: 'HSTS Max-Age', value: data.hsts?.max_age ? `${data.hsts.max_age} seconds` : null },
				{ label: 'HSTS Include SubDomains', value: data.hsts?.include_subdomains ? 'Yes' : 'No' },
				{ label: 'HSTS Preload', value: data.hsts?.preload ? 'Yes' : 'No' },
				{ label: 'HSTS Details', value: data.hsts?.details },
				{ label: 'Headers', value: data.headers, isObject: true },
			];
		}
		
		if (type === 'ip') {
			return [
				{ label: 'Input', value: data.input },
				{ label: 'Type', value: data.is_domain ? 'Domain' : 'IP Address' },
				{ label: 'IP Address', value: data.ip },
				{ label: 'Resolved IPs', value: data.resolved_ips, isArray: true },
				{ label: 'Country', value: data.country },
				{ label: 'Region', value: data.region },
				{ label: 'City', value: data.city },
				{ label: 'ISP', value: data.isp },
				{ label: 'Organization', value: data.organization },
				{ label: 'Timezone', value: data.timezone },
			];
		}
		
		if (type === 'email_config') {
			return [
				{ label: 'Domain', value: data.domain },
				{ label: 'SPF Configured', value: data.spf?.configured ? '✓' : '✗' },
				{ label: 'SPF Record', value: data.spf?.record },
				{ label: 'DKIM Configured', value: data.dkim?.configured ? '✓' : '✗' },
				{ label: 'DKIM Selectors', value: data.dkim?.selectors, isArray: true },
				{ label: 'DMARC Configured', value: data.dmarc?.configured ? '✓' : '✗' },
				{ label: 'DMARC Policy', value: data.dmarc?.policy },
				{ label: 'DMARC Record', value: data.dmarc?.record },
				{ label: 'BIMI Configured', value: data.bimi?.configured ? '✓' : '✗' },
				{ label: 'BIMI Logo', value: data.bimi?.logo_url },
			];
		}
		
		if (type === 'og_image') {
			return null; // Handled separately
		}
		
		if (type === 'raw') {
			return null; // Handled separately
		}
		
		if (type === 'error') {
			return null; // Handled separately
		}
		
		return [];
	}
	
	const displayData = $derived(getDisplayData());
	
	function getIcon() {
		switch(type) {
			case 'dns': return '🌐';
			case 'ssl': return '🔒';
			case 'http3': return '⚡';
			case 'ip': return '🌍';
			case 'web_settings': return '⚙️';
			case 'email_config': return '📧';
			case 'raw': return '📄';
			case 'error': return '✗';
			default: return '📋';
		}
	}
	
	function getTitle() {
		switch(type) {
			case 'dns': return 'DNS Information';
			case 'ssl': return 'SSL Certificate';
			case 'ip': return 'IP Information';
			case 'http3': return 'HTTP/3 Support';
			case 'web_settings': return 'Web Server';
			case 'email_config': return 'Email Config';
			case 'raw': return 'Raw Data';
			case 'error': return 'Error';
			default: return 'Information';
		}
	}
	
	function getSubtitle() {
		if (data?.domain) return data.domain;
		if (type === 'web_settings') return 'Server Settings';
		if (type === 'raw') return 'Complete API response';
		return null;
	}
	
	// const needsAccordion = $derived(['web_settings', 'email_config', 'raw'].includes(type));
	
	/**
	 * @param {any} item
	 */
	function shouldRender(item) {
		if (item?.value === undefined || item?.value === null) return false;
		if (Array.isArray(item.value) && item.value.length === 0) return false;
		return true;
	}
</script>

{#if type === 'og_image'}
	<OGImageResults data={data} />
{:else if type === 'error' && data?.error}
	<div class="col-span-full p-4 sm:p-6 bg-red-50 border border-red-200 rounded-lg">
		<div class="flex items-center space-x-2 text-red-600">
			<span class="text-xl sm:text-2xl">✗</span>
			<h3 class="text-lg sm:text-xl font-semibold">Error</h3>
		</div>
		<p class="mt-2 text-sm sm:text-base text-red-700 break-words">{data.error}</p>
	</div>
{:else if data}
	<div class="bg-gray-50 border border-gray-200 rounded-lg overflow-hidden w-full">
		<div class="px-3 py-2.5 sm:px-4 sm:py-3 md:px-6 md:py-4 bg-white">
			<div class="text-xs sm:text-sm space-y-1.5 sm:space-y-2">
				{#each displayData ?? [] as item}
					{#if shouldRender(item)}
						<div class="break-words">
							<span class="font-medium text-gray-900">{item.label}:</span>
							{#if 'isArray' in item && item.isArray}
								<div class="ml-2 mt-1 sm:ml-3 sm:mt-1.5 space-y-0.5 sm:space-y-1">
									{#each item.value as value}
										<div class="text-gray-600 text-xs sm:text-sm break-all">{value}</div>
									{/each}
								</div>
							{:else}
								<span class="ml-1 sm:ml-2 text-gray-600 break-all">{item.value}</span>
							{/if}
						</div>
					{/if}
				{/each}
				{#if data.error}
					<div class="text-red-600 text-xs sm:text-sm"><span class="font-medium">Error:</span> {data.error}</div>
				{/if}
			</div>
		</div>
	</div>
{/if}

