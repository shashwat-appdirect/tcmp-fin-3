# AppDirect India AI Workshop - Event Registration System

A modern React SPA with Golang backend for event registration, featuring real-time attendee tracking, session management, and an admin panel.

## Features

- **Event Registration**: Live attendee count with registration form
- **Sessions & Speakers**: Display sessions with associated speakers (many-to-many relationship)
- **Location**: Event details with embedded Google Maps
- **Admin Panel**: Password-protected admin interface with:
  - Attendee management (view, delete)
  - Speaker management (add, update, delete)
  - Session management (add, update, delete)
  - Statistics visualization (pie chart by designation)

## Tech Stack

### Frontend
- React 18 with TypeScript
- Vite for build tooling
- Tailwind CSS for styling
- Recharts for data visualization
- Modern animations and responsive design

### Backend
- Golang 1.21
- Google Firestore for database
- Gorilla Mux for routing
- CORS enabled for frontend integration

## Project Structure

```
tcmp-fin-3/
├── frontend/          # React SPA
│   ├── src/
│   │   ├── components/
│   │   ├── App.tsx
│   │   ├── api.ts
│   │   └── main.tsx
│   └── package.json
├── backend/           # Golang REST API
│   ├── handlers/
│   ├── models/
│   ├── firestore/
│   ├── main.go
│   └── Dockerfile
└── README.md
```

## Setup Instructions

### Prerequisites
- Node.js 18+ and npm/yarn
- Go 1.21+
- Google Cloud service account JSON file
- Access to Google Firestore

### Backend Setup

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Copy `.env.example` to `.env`:
   ```bash
   cp .env.example .env
   ```

4. Update `.env` with your configuration:
   ```env
   GOOGLE_APPLICATION_CREDENTIALS=/path/to/your/service-account.json
   ADMIN_PASSWORD=your-secure-password
   PORT=8080
   ```

5. Run the server:
   ```bash
   go run main.go
   ```

The backend will start on `http://localhost:8080`

### Frontend Setup

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Copy `.env.example` to `.env`:
   ```bash
   cp .env.example .env
   ```

4. Update `.env` with your backend URL (if different from default):
   ```env
   VITE_API_URL=http://localhost:8080
   ```

5. Start the development server:
   ```bash
   npm run dev
   ```

The frontend will start on `http://localhost:3000`

## API Endpoints

### Public Endpoints
- `GET /api/attendees` - Get all attendees with count
- `POST /api/attendees` - Register new attendee
- `GET /api/sessions` - Get all sessions with speakers
- `GET /api/speakers` - Get all speakers

### Admin Endpoints (Require X-Admin-Password header)
- `POST /api/admin/login` - Admin authentication
- `GET /api/admin/attendees` - Get full attendee list
- `DELETE /api/admin/attendees/:id` - Delete attendee
- `POST /api/admin/sessions` - Add/update session
- `DELETE /api/admin/sessions/:id` - Delete session
- `POST /api/admin/speakers` - Add/update speaker
- `DELETE /api/admin/speakers/:id` - Delete speaker
- `GET /api/admin/stats` - Get designation statistics

## Firestore Structure

The application uses the following Firestore structure:
- Main collection: `events`
- Subcollections:
  - `events/{client_id}/attendees`
  - `events/{client_id}/sessions`
  - `events/{client_id}/speakers`

The `client_id` is automatically extracted from the service account JSON file.

## Deployment

### Backend (Google Cloud Run)

1. Build the Docker image:
   ```bash
   docker build -t gcr.io/YOUR_PROJECT_ID/event-registration-backend ./backend
   ```

2. Push to Google Container Registry:
   ```bash
   docker push gcr.io/YOUR_PROJECT_ID/event-registration-backend
   ```

3. Deploy to Cloud Run:
   ```bash
   gcloud run deploy event-registration-backend \
     --image gcr.io/YOUR_PROJECT_ID/event-registration-backend \
     --platform managed \
     --region us-central1 \
     --set-env-vars ADMIN_PASSWORD=your-password \
     --set-env-vars GOOGLE_APPLICATION_CREDENTIALS=/path/to/credentials
   ```

### Frontend

Build for production:
```bash
cd frontend
npm run build
```

The `dist` folder contains the production build that can be deployed to any static hosting service (e.g., Firebase Hosting, Netlify, Vercel).

## Security Notes

- Never commit `.env` files or service account JSON files to version control
- Use strong passwords for admin access
- Ensure CORS is properly configured for production
- Use environment variables for all sensitive configuration

## Designations

The registration form includes the following default tech designations:
- Software Engineer
- Senior Software Engineer
- Tech Lead
- Engineering Manager
- Product Manager
- Product Designer
- UX Designer
- Data Scientist
- DevOps Engineer
- QA Engineer
- Full Stack Developer
- Frontend Developer
- Backend Developer
- Mobile Developer
- Architect
- CTO
- Other

## License

This project is for the AppDirect India AI Workshop event.

