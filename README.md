# Jolt

Jolt is a minimalistic marketplace web application designed for lightning-fast trading and robust identity verification. Built with Go and templ, it will offer a seamless and secure platform for buyers and sellers to conduct transactions with confidence.

## Features
- **Go + templ + HTMX Stack**: Leverages the power of Go with the simplicity of templ and the reactivity of HTMX for efficient server-side rendering.
- **Hot Module Replacement (HMR) via templ proxy**: Supports rapid development with instant updates.
- **Tailwind CSS Integration**: Utilizes Tailwind for responsive and customizable designs.

## Future Features
- **Lightning-Fast Trading**: Optimized for speed to facilitate quick transactions.
- **Aggressive Identity Verification**: Ensures the authenticity of both buyers and sellers.

## Tech Stack

- **Backend**: Go
- **Frontend**: templ (HTML templating engine for Go) + HTMX
- **Styling**: Tailwind CSS
- **Development**: Air (for Go live-reloading)

## Getting Started

### Prerequisites

- Go (version 1.23 or later)
- Node.js and npm/pnpm/yarn
- Make

### Installation

1. Clone the repository:
   ```
   git clone https://github.com/DillonEnge/jolt.git
   cd jolt
   ```

2. Install Go dependencies:
   ```
   go mod tidy
   ```

3. Install Node.js dependencies (choose one):
   ```
   npm install
   # or
   pnpm install
   # or
   yarn install
   ```

### Environment Variables

Jolt requires the following environment variables to be set:

- `DATABASE_PASSWORD`: The password for your database user
- `DATABASE_URL`: The connection string for your database.
- `PORT`: The port number on which the application will run.
- `CASDOOR_APPLICATION_NAME`: The name of your Casdoor application.
- `CASDOOR_CLIENT_ID`: The client ID for your Casdoor application.
- `CASDOOR_CLIENT_SECRET`: The client secret for your Casdoor application.
- `CASDOOR_ENDPOINT`: The endpoint URL for your Casdoor server.
- `CASDOOR_ORGANIZATION_NAME`: The name of your Casdoor organization.
- `CASDOOR_REDIRECT_URI`: The redirect URI for Casdoor authentication.

We recommend using `direnv` to manage your environment variables. Follow these steps:

1. Install `direnv` if you haven't already. (Visit [direnv.net](https://direnv.net) for installation instructions)

2. Add the following to ~/.config/direnv/direnv.toml

   ```
   [global]
   load_dotenv = true
   ```

2. Create a `.env` file in the root of your project:

   ```
   DATABASE_PASSWORD=testpass
   DATABASE_URL=postgresql://username:password@localhost:5432/marketwise
   PORT=8080
   CASDOOR_APPLICATION_NAME=your_application_name
   CASDOOR_CLIENT_ID=your_client_id
   CASDOOR_CLIENT_SECRET=your_client_secret
   CASDOOR_ENDPOINT=https://your-casdoor-endpoint.com
   CASDOOR_ORGANIZATION_NAME=your_organization_name
   CASDOOR_REDIRECT_URI=http://localhost:8080/callback
   ```

   Replace the placeholder values with your actual configuration details.

3. Allow direnv to load the `.env` file:

   ```
   direnv allow .
   ```

4. Add `.env` to your `.gitignore` file to avoid committing sensitive information:

   ```
   echo ".env" >> .gitignore
   ```

Now, whenever you enter the project directory, these environment variables will be automatically loaded.

### Database Setup

1. Start the PostgreSQL instance using Docker Compose:
   ```
   docker compose start db -d
   ```
   This command will start a PostgreSQL container in detached mode.

2. Run database migrations:
   ```
   make migrate
   ```
   This command will set up your database schema and apply any necessary migrations.

### Running the Application

To start the development server with live-reloading for Go, templ, and Tailwind CSS:

```
make dev
```

This command will:
- Start the Go server using Air for live-reloading
- Watch for changes in templ files
- Proxy to the browser for HMR via templ
- Watch for changes in Tailwind CSS files and rebuild as necessary

Open your browser and navigate to `http://localhost:8080` (or the port you specified in the `PORT` environment variable) to view the application.

## Development

The `make dev` command sets up a complete development environment with hot-reloading. You can now edit your Go files, templ templates, and Tailwind CSS, and see the changes reflected immediately in your browser.

## License

Jolt is open-sourced software licensed under a [personal use license](LICENSE.md).

---

Jolt - Trade with confidence, verified at the speed of light.
