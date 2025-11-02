# API Functions Usage Guide

This document shows how to use the reusable API functions from `$lib/api.js`.

## Available Functions

All functions return a Promise that resolves with the API data:

- `checkSSL(domain)` - Check SSL certificate information
- `checkHTTP3(domain)` - Check HTTP/3 support
- `checkDNS(domain)` - Get DNS information
- `checkIP(ip)` - Check IP address information
- `checkWebSettings(domain)` - Check web server settings
- `checkComprehensive(domain)` - Get all information at once

## Usage in Svelte Components

### Example 1: HTTP/3 Check

```svelte
<script>
  import { checkHTTP3 } from '$lib/api.js';
  
  let data = $state(null);
  let loading = $state(false);
  
  async function checkDomain(domain) {
    loading = true;
    data = await checkHTTP3(domain);
    loading = false;
  }
</script>
```

### Example 2: SSL Check

```svelte
<script>
  import { checkSSL } from '$lib/api.js';
  
  let sslData = $state(null);
  let loading = $state(false);
  
  async function getSSL(domain) {
    loading = true;
    sslData = await checkSSL(domain);
    loading = false;
  }
</script>
```

### Example 3: Form Integration

```svelte
<script>
  import { checkComprehensive } from '$lib/api.js';
  
  let results = $state(null);
  let loading = $state(false);
  
  async function handleSubmit(event) {
    event.preventDefault();
    loading = true;
    
    const formData = new FormData(event.target);
    const domain = String(formData.get('domain'));
    
    results = await checkComprehensive(domain);
    loading = false;
  }
</script>

<form onsubmit={handleSubmit}>
  <input name="domain" placeholder="e.g., google.com" />
  <button type="submit" disabled={loading}>
    {loading ? 'Loading...' : 'Check'}
  </button>
</form>
```

## Error Handling

All functions automatically handle errors and return an object with an `error` property:

```javascript
const result = await checkHTTP3('example.com');

if (result.error) {
  console.error('Error:', result.error);
} else {
  // Use result.data or result directly
  console.log('Success:', result);
}
```

## Custom API Calls

You can also use the generic `fetchApiData` function for custom queries:

```javascript
import { fetchApiData } from '$lib/api.js';

// Call any endpoint with custom parameters
const data = await fetchApiData('/cek/dns', { domain: 'example.com' });
const ipData = await fetchApiData('/cek/ip', { ip: '8.8.8.8' });
```

## Available Endpoints

The following endpoints are available:
- `/cek` - Comprehensive check
- `/http3` - HTTP/3 support check
- `/cek/ssl` - SSL certificate check
- `/cek/dns` - DNS information
- `/cek/ip` - IP information
- `/cek/web-settings` - Web server settings

