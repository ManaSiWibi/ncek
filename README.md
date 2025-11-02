# NetCheck

A comprehensive network checking tool built with SvelteKit (frontend) and Go (backend API).

## Features

- **SSL Certificate Analysis**: Check certificate validity, issuer, expiration, and key details
- **HTTP/3 Support Detection**: Test if a domain supports HTTP/3 protocol
- **DNS Information**: Get A, AAAA, CNAME, MX, TXT, and NS records
- **IP Information**: Basic IP address validation and connection testing
- **Web Server Settings**: Analyze HTTP headers, server information, and response details
- **Email Authentication**: Check SPF, DKIM, DMARC, and BIMI records
- **Comprehensive Check**: Get all information in a single request

## Running with Docker (Recommended)

### Quick Start

```bash
# Build and run all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

The application will be available at:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080

### Building Individual Services

#### Backend Only
```bash
cd api
docker build -t netcheck-backend .
docker run -p 8080:8080 netcheck-backend
```

#### Frontend Only
```bash
docker build -t netcheck-frontend .
docker run -p 3000:3000 -e BACKEND_URL=http://your-backend-url:8080 netcheck-frontend
```

## Local Development

### Prerequisites

- Node.js 20+ (for frontend)
- Go 1.24+ (for backend)

### Backend Setup

```bash
cd api

# Install dependencies
go mod tidy

# Run the server
go run main.go
```

The API will be available at `http://localhost:8080`

### Frontend Setup

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# or start the server and open the app in a new browser tab
npm run dev -- --open
```

The frontend will be available at `http://localhost:5173` (or the port shown in terminal)

### Build for Production

```bash
# Build frontend
npm run build

# Preview production build
npm run preview

# Build backend (from api directory)
cd api
go build -o netcheck main.go
```

## Project Structure

```
ncek/
├── api/              # Go backend API
│   ├── main.go       # Main backend server
│   └── Dockerfile    # Backend Docker configuration
├── src/              # SvelteKit frontend
│   ├── lib/          # Shared components and utilities
│   └── routes/       # Application routes
├── Dockerfile        # Frontend Docker configuration
└── docker-compose.yml # Docker Compose configuration
```

## Environment Variables

Create a `.env` file in the project root to configure all environment variables:

```bash
# Generate a secure API secret key
openssl rand -hex 32

# Create .env file from template
cp .env.example .env

# Edit .env and update API_SECRET_KEY with your generated key
```

### Required Variables

- `API_SECRET_KEY`: **REQUIRED** - Secret key for API authentication between frontend and backend. Generate one using: `openssl rand -hex 32`

### Optional Variables (with defaults)

**Backend:**
- `GIN_MODE`: Gin mode (`debug` or `release`, default: `release`)

**Frontend:**
- `NODE_ENV`: Node environment (default: `production`)
- `HOST`: Host to bind to (default: `0.0.0.0`)
- `PORT`: Internal container port (default: `3000`, fixed)
- `FRONTEND_PORT`: External port mapping (default: `3001`)
- `BACKEND_URL`: Backend API URL (default: `http://backend:8080`)
- `ALLOWED_ORIGIN`: Allowed CORS origin (default: `http://localhost:3001`)

All variables are read from the `.env` file automatically by Docker Compose.

**Note:** The `.env` file is already in `.gitignore` to prevent committing secrets. Copy `.env.example` to `.env` and update the values.

## API Documentation

See [api/README.md](api/README.md) for detailed API documentation.

## License

This project is open source and available under the MIT License.

## Dynamic tool pages (config-driven)

To avoid cluttering `src/routes` with many near-identical pages, the app supports dynamic tool pages driven by a central registry.

- Registry: `src/lib/tools.js`
- Dynamic route: `src/routes/[tool]/+page.svelte` and `+page.js`

Adding a new tool page:
1. Add an entry in `src/lib/tools.js`:
```js
export const tools = {
  // existing...
  newtool: {
    title: 'My New Tool',
    description: 'Short description here',
    placeholder: 'e.g., example.com',
    // Name of the form field your page expects
    formField: 'description',
    // Build API params that will be sent to /routes/api/+server.js
    apiParams: (value) => ({ type: 'my-endpoint', domain: value }),
    // ResultCard type ('ssl' | 'http3' | 'dns' | 'ip' | 'web_settings' | 'email_config' | 'raw')
    resultCard: 'raw',
    // Optional transform on the API response before rendering
    transform: (data) => data
  }
}
```
2. Visit `/newtool` in the browser. No new files are needed under `src/routes`.

Notes:
- The API call is fulfilled by the existing proxy at `src/routes/api/+server.js`.
- Use `resultCard: 'raw'` to show JSON, or a supported `ResultCard` type for structured display.
