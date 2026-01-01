# Pastebin Lite

Pastebin Lite is a simple Pastebin-like application that allows users to create and store text snippets. The project consists of a Go backend and a Next.js frontend.

---

## Tech Stack

### Backend

* Go
* net/http
* PostgreSQL

### Frontend

* Next.js
* TypeScript

---

## GitHub Repository

[https://github.com/sana-mukhtar/pastebin_lite](https://github.com/sana-mukhtar/pastebin_lite)

---

## Prerequisites

Ensure the following are installed on your system:

* Go 1.21 or later
* PostgreSQL
* Node.js 18 or later
* npm

---

## Running the App Locally

Follow the steps below to run the backend and frontend locally.

---

## Backend Setup (Go + PostgreSQL)

### 1. Clone the Repository

```bash
git clone https://github.com/sana-mukhtar/pastebin_lite.git
cd pastebin_lite
```

---

### 2. Create PostgreSQL Database

Start PostgreSQL and create a database:

```sql
CREATE DATABASE pastebin;
```

---

### 3. Configure Environment Variables

Create a `.env` file in the project root (do not commit this file):

```env
PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/pastebin?sslmode=disable
```

Update username and password if required.

---

### 4. Run the Backend Server

```bash
go mod tidy
go run .
```

Backend will start at:

```
http://localhost:8080
```

You can verify the backend is running by opening:

```
http://localhost:8080/health
```

---

## API Endpoint

### Create Paste

```
POST /api/paste/create
```

Request Body:

```json
{
  "content": "Hello World"
}
```

---

## Frontend Setup (Next.js)

### 1. Navigate to UI Folder

```bash
cd ui
```

---

### 2. Install Dependencies

```bash
npm install
```

---

### 3. Configure Frontend Environment Variable

Create a `.env.local` file inside the `ui` directory:

```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

---

### 4. Run Frontend

```bash
npm run dev
```

Frontend will start at:

```
http://localhost:3000
```

---

## Using the Application

1. Ensure PostgreSQL is running
2. Start the backend server
3. Start the frontend server
4. Open `http://localhost:3000` in your browser
5. Enter text and click on **Create Paste**
6. After creating the paste, the application generates a **unique URL**
7. Open the generated URL to view the created paste
8. Refresh max views times, after max views is 0, the text becomes unavailable.

---

## Notes

* `.env` and `.env.local` files should not be committed to version control
* Backend must be running before starting the frontend
* Database credentials should be managed via environment variables only

---

## Author

Sana Mukhtar
