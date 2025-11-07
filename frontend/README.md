# Frontend

Simple Vue SPA for managing invoices. Supports backends SWFP.

## What It Does

A single-page application that provides:

- **Dashboard** with financial overview and statistics
- **Invoice management** - create, edit, view, delete
- **Line items** - add/edit products with quantity, price, VAT
- **PDF generation** - download invoices as PDF
- **Smart filtering** - search, sort, paginate invoices
- **Real-time calculations** - automatic totals with VAT

## Technologies

- **Vue 3** - framework
- **Vite** - build tool
- **Vuetify 3** - component library
- **Vue Router** - client-side routing
- **Axios** - HTTP client with interceptors
- **Material Design Icons** - icon set

## Architecture

```shell
src/
├── assets/           # Static files (images, SVGs)
├── components/       # Reusable UI components
│   ├── InvoiceCard.vue           # Invoice preview card
│   ├── InvoiceForm.vue           # Main form wrapper
│   ├── InvoiceHeaderFields.vue   # Invoice metadata fields
│   ├── InvoiceItemsManager.vue   # Line items CRUD
│   ├── InvoiceItemForm.vue       # Item add/edit form
│   ├── InvoiceItemRow.vue        # Single item display
│   ├── InvoiceSummary.vue        # Totals calculator
│   └── InvoiceStatusBadge.vue    # Status chip
├── views/            # Page-level components
│   ├── DashboardView.vue         # Stats & recent invoices
│   ├── InvoiceListView.vue       # Grid with filters
│   ├── InvoiceDetailView.vue     # Read-only invoice
│   └── InvoiceFormView.vue       # Create/edit wrapper
├── router/           # Route definitions
├── services/         # API layer
│   └── api.js        # Axios instance + invoice service
├── plugins/          # Vuetify config
├── App.vue           # Root component
└── main.js           # App entry point
```

## Key Features

### API Integration

- Centralized Axios instance with interceptors
- Automatic error handling and user-friendly messages
- Request/response logging for debugging
- Service layer abstracts HTTP calls

### Form Handling

- Two-way binding with `v-model`
- Validation with Vuetify rules
- Separate forms for invoice headers and line items
- Real-time calculations for VAT and totals

### State Management

- No Vuex/Pinia - uses composables and props
- Parent components manage data fetching
- Child components emit events for mutations

## Running Locally

### Development Server

```bash
cd frontend
npm install
npm run dev
```

Access at `http://localhost:5173`

### Production Build

```bash
npm run build
npm run preview  # Test production build
```

Output in `dist/` directory.

### Environment Variables

Create `.env` file:

```env
VITE_API_URL=http://localhost:8080/api
VITE_APP_VERSION=1.0.0
```

- `VITE_API_URL` - Backend API base URL
- `VITE_APP_VERSION` - Displayed in UI (injected at build time)

## Docker

### Build Image

```bash
docker build \
  --build-arg NODE_VERSION=24.9.0 \
  --build-arg NGINX_VERSION=3.18 \
  --build-arg APP_VERSION=1.0.0 \
  -t invoice-frontend .
```

### Run Container

```bash
docker run -p 80:80 invoice-frontend
```

Multi-stage build: Node.js for compilation, Nginx for serving.

## API Service

All backend communication goes through `src/services/api.js`.

Error handling is automatic - user-friendly messages from interceptors.

## Nginx Config

Production uses Nginx for:

- Gzip compression
- Cache headers for static assets
- Security headers (XSS, frame options)
- SPA history mode (redirect to index.html)
- Health check endpoint at `/health`
