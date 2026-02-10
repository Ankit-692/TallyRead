📚 TallyRead Backend

TallyRead is a robust Go-based REST API designed to help readers track their literary journeys. Users can search for books via the Google Books API, manage a personal library, and update reading statuses (Reading, Planning, Completed, or Dropped).
🚀 Key Features

    Book Discovery: Integrated with Google Books API.

    Smart Caching: Redis implementation to cache book search results for 48 hours, optimizing performance and saving API tokens.

    Authentication: Secure JWT-based auth stored in HTTP-only cookies (7-day expiry).

    Library Management: Full CRUD operations for personal book collections.

    Mailing System: Password recovery flow using gopkg.in/mail.v2 and custom templates.

    Security: Middleware-protected routes and rate-limiting.

🛠️ Tech Stack

    Language: Go (v1.24.0)

    Framework: Gin Gonic

    Database: PostgreSQL (via Supabase)

    Cache: Redis

    Authentication: JWT (JSON Web Tokens)

    External APIs: Google Books API

📂 Project Structure
Plaintext

├── config/       # Configuration loaders (SMTP, etc.)
├── db/           # Database connection and initialization
├── middlewares/  # JWT Authentication & Authorization logic
├── models/       # Database schemas (Books, Users, ResetRequests)
├── routes/       # API route definitions and controllers
├── templates/    # HTML mail templates for password resets
├── utils/        # Helpers: Mailer, JWT, Hashing, Rate Limiter
└── main.go       # Application entry point

⚙️ Setup & Installation
1. Prerequisites

    Go: 1.24.0 or higher

    Redis: Local or managed instance

    PostgreSQL: Supabase instance or local DB

2. Environment Variables

Create a .env file in the root directory and fill in your credentials:
Code snippet

# Server & Auth
secretkey=your_jwt_secret
FRONTEND_URL=http://localhost:4200

# Database (Supabase/Postgres)
DB_HOST=your_host
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=your_db_name

# Redis
REDIS_URL=your_redis_url

# Google Books API
GOOGLE_BOOKS_API_KEY=your_api_key

# SMTP Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_EMAIL=your_email@gmail.com
SMTP_PASSWORD=your_app_password

3. Run the App
Bash

go mod download
go run main.go

🛣️ API Endpoints
Public Routes
Method	Endpoint	Description
POST	/register	Register a new user
POST	/login	Login and receive JWT cookie
POST	/forgot-Password	Request a password reset link
POST	/reset-password	Update password using token
Protected Routes (Requires JWT Cookie)
Method	Endpoint	Description
GET	/api/searchBooks	Search Google Books (Cached via Redis)
POST	/api/addBook	Add a book to your library
GET	/api/getAllBooks	Fetch user's entire library
POST	/api/book/:id	Update book status or details
GET	/api/me	Check auth status (used by Angular AuthGuard)
POST	/api/logout	Clear the auth cookie
🛡️ Frontend Integration (Angular)

This backend is designed to work with an Angular frontend using:

    WithCredentials: Set to true in your HTTP Interceptor to allow cookie exchange.

    AuthGuard: Protects UI routes by calling the /api/me endpoint to verify the session.
