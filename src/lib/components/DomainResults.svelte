<script>
	import ResultCard from './ResultCard.svelte';
	import * as Accordion from "$lib/components/ui/accordion/index.js";

	/** @type {{
		ssl?: any;
		http3?: any;
		dns?: any;
		web_settings?: any;
		email_config?: any;
		error?: any;
	}} */
	export let number = /** @type {any} */ (null);
</script>

{#if number}
	<div class="w-full max-w-6xl mx-auto bg-white rounded-lg shadow-lg border border-gray-200">
		<div class="p-6 border-b border-gray-200">
			<h2 class="text-xl font-semibold text-gray-800">Domain Analysis Results</h2>
		</div>
		
		<div class="p-4 sm:p-6">
			{#if number.error}
				<ResultCard type="error" data={number} />
			{:else}
				<Accordion.Root type="multiple" orientation="vertical" class="w-full gap-4" value="item-1">
					<Accordion.Item value="item-1" class="">
						<Accordion.Trigger class="">SSL Certificate</Accordion.Trigger>
						<Accordion.Content class="flex flex-col gap-4 text-balance">
							<ResultCard type="ssl" data={number.ssl} />
						</Accordion.Content>
					</Accordion.Item>
					<Accordion.Item value="item-2" class="">
						<Accordion.Trigger class="">HTTP/3 Support</Accordion.Trigger>
						<Accordion.Content class="flex flex-col gap-4 text-balance">
							<ResultCard type="http3" data={number.http3} />
						</Accordion.Content>
					</Accordion.Item>
					<Accordion.Item value="item-3" class="">
						<Accordion.Trigger class="">DNS Information</Accordion.Trigger>
						<Accordion.Content class="flex flex-col gap-4 text-balance">
							<ResultCard type="dns" data={number.dns} />
						</Accordion.Content>
					</Accordion.Item>
					<Accordion.Item value="item-4" class="">
						<Accordion.Trigger class="">Web Server</Accordion.Trigger>
						<Accordion.Content class="flex flex-col gap-4 text-balance">
							<ResultCard type="web_settings" data={number.web_settings} />
						</Accordion.Content>
					</Accordion.Item>
					<Accordion.Item value="item-5" class="">
						<Accordion.Trigger class="">Email Config</Accordion.Trigger>
						<Accordion.Content class="flex flex-col gap-4 text-balance">
							<ResultCard type="email_config" data={number.email_config} />
						</Accordion.Content>
					</Accordion.Item>
				</Accordion.Root>
			{/if}
		</div>
	</div>
{/if}
