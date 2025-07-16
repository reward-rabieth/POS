# 🍽️ Point of Sale(POS) SYSTEM

A modern, lightweight Point of Sale (POS) system built using:

    Go as the backend

    PocketBase as the local database and authentication provider

    HTMX for reactive UI

    Gomponents for declarative HTML rendering

    TailwindCSS for styling

    🔐 Auth is powered by PocketBase with JWT tokens, stored in secure HTTP-only cookies.
---

## ✨ Features

-  **User Authentication** (Login/Logout with JWT token)
-  **Order Management** (Create orders, view order history)
-  **Inventory Management** (Add/update menu items and stock levels)

-  **PocketBase Integration** (as backend with API access rules)
-  **HTMX + Tailwind** for snappy, minimal frontend experience
-  Fully rendered using `Gomponents` in Go (no JS frontend frameworks)

---

## 🛠️ Tech Stack

| Layer       | Technology                       |
|-------------|----------------------------------|
| Frontend    | [HTMX](https://htmx.org/), [Tailwind CSS](https://tailwindcss.com/) |
| Backend     | [Go](https://golang.org/), [Gomponents](https://github.com/maragudk/gomponents) |
| Database/API| [PocketBase](https://pocketbase.io/) (SQLite + REST API) |
| Auth        | PocketBase email + password auth with JWT cookie session |

---

## 🚀 Getting Started

### Prerequisites

- Go 1.20+
- [PocketBase](https://pocketbase.io/) (`./pocketbase serve`)
- Make sure your `items`, `orders`, and `users` collections are created in PocketBase



```bash
git clone https://github.com/your-username/possystem.git
cd possystem
```

###   Download PocketBase Binary

Download from https://pocketbase.io/docs/ and place the binary in your project root.

```bash
chmod +x pocketbase

#Start PocketBase

./pocketbase serve

```

once you run ./pocketbase serve, PocketBase automatically creates the following folders in your project root:

```bash
possystem/
├── pb_data/
├── pb_migrations/
├── pb_public/

```



Start the Go Backend
```bash
go run cmd/app/main.go
```


API ENDPOINTS

| Endpoint      | Method | Description                 |
| ------------- | ------ | --------------------------- |
| `/login`      | POST   | Log in and set token cookie |
| `/logout`     | GET    | Clears session cookie       |
| `/items`      | GET    | List all items              |
| `/items/new`  | POST   | Create a new item           |
| `/orders`     | GET    | List all orders             |
| `/orders/new` | POST   | Create new order            |
