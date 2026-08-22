# Internal Whiteboard Tool — Architecture & Build Plan

*A Miro-style collaborative whiteboard for internal company use, built on Go (Fiber) + Next.js + shadcn/ui.*

---

## Table of Contents

1. [Overview](#overview)
2. [Scope for v1](#scope-for-v1)
3. [Design Philosophy](#design-philosophy)
4. [Tech Stack](#tech-stack)
5. [Architecture Overview](#architecture-overview)
6. [Real-Time Collaboration Design](#real-time-collaboration-design)
7. [Data Model](#data-model)
8. [Authentication and Authorization](#authentication-and-authorization)
9. [REST API Reference](#rest-api-reference)
10. [WebSocket Protocol](#websocket-protocol)
11. [Backend Structure (Go + Fiber)](#backend-structure-go--fiber)
12. [Frontend Structure (Next.js + shadcn/ui)](#frontend-structure-nextjs--shadcnui)
13. [Security Considerations](#security-considerations)
14. [Testing Strategy](#testing-strategy)
15. [Deployment](#deployment)
16. [Development Roadmap](#development-roadmap)
17. [Appendix: Dependencies](#appendix-dependencies)

---

## Overview

We're building a self-hosted, internal-only whiteboard: an infinite-ish canvas where employees create sticky notes, shapes, text, and connectors, and see each other's changes live. No public access, no self-signup — this is a tool for people already inside the company, provisioned by an admin.

Two requirements shape every decision in this document:

1. **The code has to stay simple.** Whoever maintains this in a year, possibly not you, should be able to open any file and understand it without a whiteboard-collaboration PhD.
2. **The product has to feel fast with multiple people on the same board at once.** This is the metric that actually matters — not feature count, not visual polish.

These two goals sound like they pull in opposite directions (usually, "handle lots of concurrent editors well" means "adopt a complicated CRDT library"). Section 3 explains why they don't have to here, and Section 6 shows the design that gets us both.

---

## Scope for v1

**In scope:**
- Boards containing sticky notes, rectangles, ellipses, lines/arrows, free text, and images
- Pan/zoom canvas
- Real-time multi-user editing, live cursors, and presence ("who's on this board right now")
- Admin-managed user accounts (create, disable, reset password)
- Board sharing with view/edit permissions

**Out of scope for v1** (call these out now so nobody's surprised later):
- Character-by-character collaborative text editing (CRDT-style). We handle simultaneous text edits with a simple lock instead — see Section 6.
- Comment threads, @mentions, notifications
- Cross-session version history / infinite undo (a simple per-session undo stack is a cheap stretch goal, but a full history browser is not v1)
- Templates gallery, stickers marketplace
- Native mobile apps (the web app should be usable on a tablet, but that's it)
- Voice/video chat
- Public/guest links — everyone who touches a board is a provisioned employee
- Offline editing with later sync

---

## Design Philosophy

There is exactly **one hot path** in this whole application: relaying a change one person made to everyone else looking at the same board. Everything else — logging in, listing boards, renaming things, managing users — happens at human typing speed and doesn't need to be clever.

So the plan deliberately splits effort unevenly:

- **The hot path (WebSocket message fan-out) is kept lean, native, and carefully commented.** No ORM, no frameworks-on-frameworks — just goroutines, channels, and a map. This is the ~150 lines of code in the entire project that actually deserve close attention.
- **Everything else leans on high-level, boring, well-documented tools** — GORM for CRUD, shadcn/ui for every screen, Zustand for client state. These aren't the fastest possible choices, but they don't need to be; they need to be easy to read six months from now.

> **Rule of thumb used throughout this plan:** if a piece of code runs once per user action (a click, a form submit), optimize for readability. If it runs once per keystroke-of-mouse-movement across every connected client (a drag, a cursor move), optimize for speed. Almost nothing in this app falls in the second bucket.

---

## Tech Stack

| Layer | Choice | Why |
|---|---|---|
| Backend framework | **Go + Fiber v3** | Built on fasthttp, one of the fastest HTTP stacks available in Go; Express-like API keeps handler code short. (Fiber v2 is an equally solid, more battle-tested fallback if you'd rather not be on the newer major version — the Hub/WebSocket design in this doc works the same either way.) |
| Real-time transport | **Native WebSockets** (`fiber/v3/middleware/websocket`) | Lower latency than polling or SSE-for-input; no need for Socket.IO's extra protocol layer since we control both ends of the connection. |
| Database | **PostgreSQL** | Reliable, easy to self-host internally, and JSONB gives us a flexible spot for per-element styling without a migration for every new shape property. |
| ORM | **GORM** | Cuts a lot of boilerplate out of the non-critical CRUD paths (users, boards, sharing). Not used on the hot path — see Design Philosophy above. |
| Frontend framework | **Next.js 16 (App Router)** | File-based routing keeps "where does this page live" obvious; Turbopack (default since v16) makes local dev fast. |
| UI components | **shadcn/ui + Tailwind** | You own the generated component code (it's copied into your repo, not hidden in node_modules), so tweaking anything later is just editing a file. This is the main lever for "don't build everything from scratch." |
| Client state | **Zustand** | A fraction of Redux's boilerplate; a plain store with a few functions, which is exactly the "simple to read" bar we want. |
| Forms | **react-hook-form + zod** | shadcn's own recommended pairing — every form (login, create user, share board) follows the same small pattern. |
| Auth | **JWT + bcrypt** | Stateless (no session table to manage), simple to reason about, and easy to swap for company SSO later if you ever want to — see Section 8. |

---

## Architecture Overview

```mermaid
flowchart LR
    subgraph Client["Browser (Next.js + shadcn/ui)"]
        UI["Board UI"]
        WSClient["WebSocket client"]
    end

    subgraph Server["Go + Fiber backend"]
        REST["REST API"]
        WSH["WebSocket handler"]
        Hub["In-memory Hub, one per active board"]
    end

    DB[("PostgreSQL")]

    UI -->|"REST: login, boards, users"| REST
    WSClient -->|"WS: element changes"| WSH
    WSH --> Hub
    Hub -->|"broadcast to other clients"| WSH
    WSH -->|"push updates"| WSClient
    REST --> DB
    Hub -.->|"async, batched persistence"| DB
```

Two separate channels do two separate jobs, and neither one is allowed to block the other:

- **REST** handles anything that isn't time-critical: authentication, listing/creating/renaming boards, admin user management, sharing permissions. Normal request/response, normal GORM queries.
- **WebSocket** handles only board-editing events once a board is open. The server keeps one lightweight in-memory "Hub" per board that's currently being viewed by someone; it exists purely to relay messages between the browsers connected to that board, and it disappears once everyone leaves (see Section 6).

The database is the source of truth for anything durable, but it's never on the critical path of "user moves a shape and sees it move." That path stays entirely in memory until it's flushed asynchronously.

---

## Real-Time Collaboration Design

This is the section that actually earns the "efficient" requirement, so it's worth walking through the reasoning, not just the result.

### The model: independent elements, not a shared document

A whiteboard isn't really "one document" the way a Google Doc is. It's a bag of mostly-independent objects — sticky notes, rectangles, arrows — each with its own position and content. That difference matters a lot for how hard collaboration is:

- In a text document, two people typing in the *same paragraph* at the *same time* is common, and merging their keystrokes correctly requires real machinery (Operational Transform or a CRDT like Yjs).
- On a whiteboard, two people editing the *same shape* at the *same instant* is rare. The much more common case — two people moving two *different* shapes at the same time — has **zero conflict potential**, because the shapes don't share any state.

So the design treats every element as an independent row with its own ID, and every user action (create, move, resize, restyle, delete) becomes a small event about *one* element. The server's job is mostly to relay that event to everyone else looking at the same board — it doesn't need to understand or merge anything.

For the rare case where two people *do* touch the same element:
- **Moving/resizing/restyling:** last-write-wins, ordered by the server's receipt time. Each element carries a `version` counter that increments on every change, purely so a client can detect "I'm about to apply a stale update" and skip it if a newer one already arrived out of order.
- **Editing text inside the same sticky note at the same time:** this is the one case LWW would feel actively broken (keystrokes clobbering each other), so instead of reaching for a CRDT, we use a **soft lock**. The moment someone focuses a text field, the client sends `element:lock`; everyone else's UI shows "Priya is editing this" and treats that element as read-only until `element:unlock` fires (on blur, or automatically if that client disconnects mid-edit). Two extra message types, no merge logic.

> **Why not a CRDT library like Yjs?** Yjs (and CRDTs generally) are genuinely the right tool when you need real-time character-level merging — and that's exactly what Figma/Miro use under the hood for text. But adopting one means the whole team now needs to understand CRDT documents, providers, and awareness protocols to touch anything collaboration-related. For a shape-based internal tool where same-element collisions are rare and text collisions are handled by a lock, that complexity buys us very little. If this ever grows into a product where rich collaborative text editing is core, Yjs is the natural next step — but it isn't needed to hit "multiple people, same board, feels instant" today.

### Efficiency techniques

These are the concrete choices that make the above fast, in rough order of impact:

- **Push, don't poll.** WebSocket means an update reaches other clients in one network hop, as soon as it happens — no polling interval to wait out.
- **One Hub per board, created lazily, torn down when empty.** The server only holds state (and runs a goroutine) for boards someone is actively viewing. A `map[boardID]*Hub` behind a mutex is the whole "room" system — no external message broker needed at this scale.
- **A buffered outbound channel per client.** Each connected client has its own `chan []byte` for outgoing messages. If one person's connection is slow, the Hub drops *that client* rather than blocking the broadcast loop for everyone else.
- **Deltas, not full state.** Every message describes one small change (`element:update {id, x, y}`), never the whole board. A new client gets exactly one full snapshot on connect; after that, only diffs.
- **Throttle high-frequency events.** Dragging a shape or moving a cursor can generate far more events than a human eye can distinguish. The client throttles these to roughly 20–30 messages/second before sending — smooth to look at, a fraction of the network traffic of sending on every animation frame.
- **Persistence is asynchronous and batched.** The Hub marks a board "dirty" on any change; a background goroutine flushes changed elements to Postgres every few seconds (or on a debounce after edits stop), instead of writing to the database on every keystroke-equivalent event. The database is never in the path between "I moved this" and "you see it moved."

### Message flow for a typical edit

```mermaid
sequenceDiagram
    participant A as User A's browser
    participant S as Fiber server (Hub)
    participant B as User B's browser

    A->>S: element:update {id, x, y}
    S->>S: apply patch to in-memory board state
    S-->>B: broadcast element:update
    B->>B: move the shape locally
    Note over S: every few seconds
    S->>S: flush dirty elements to Postgres
```

---

## Data Model

Five tables carry the whole app. Nothing here is exotic on purpose.

- **users** — one row per person. Only admins create these; there is no signup form anywhere.
- **boards** — a board's name and owner.
- **board_members** — who can see/edit which board, and at what permission level.
- **board_elements** — every shape/note/line on every board. This is the table the real-time layer reads and writes most, and it stays intentionally flat (no nested structures) so it's easy to reason about.

```sql
CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username      TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    display_name  TEXT NOT NULL,
    role          TEXT NOT NULL DEFAULT 'member' CHECK (role IN ('admin', 'member')),
    is_active     BOOLEAN NOT NULL DEFAULT true,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE boards (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       TEXT NOT NULL,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE board_members (
    board_id UUID NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
    user_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role     TEXT NOT NULL CHECK (role IN ('owner', 'editor', 'viewer')),
    PRIMARY KEY (board_id, user_id)
);

CREATE TABLE board_elements (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    board_id   UUID NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
    type       TEXT NOT NULL CHECK (type IN ('sticky_note','rectangle','ellipse','text','line','arrow','image')),
    x          DOUBLE PRECISION NOT NULL,
    y          DOUBLE PRECISION NOT NULL,
    width      DOUBLE PRECISION NOT NULL,
    height     DOUBLE PRECISION NOT NULL,
    rotation   DOUBLE PRECISION NOT NULL DEFAULT 0,
    z_index    INTEGER NOT NULL DEFAULT 0,
    content    TEXT,                  -- text/sticky-note body
    style      JSONB NOT NULL DEFAULT '{}',  -- color, font size, stroke width, etc.
    version    INTEGER NOT NULL DEFAULT 1,   -- bumped on every change, used for stale-update detection
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_board_elements_board_id ON board_elements(board_id);
CREATE INDEX idx_board_members_user_id ON board_members(user_id);
```

`style` being JSONB (rather than a dozen nullable columns) is the one deliberate "flexible" choice here — it means adding a new visual property to a shape later (say, a shadow, or a dashed-border option) never needs a migration.

---

## Authentication and Authorization

No signup page exists anywhere in this app. The only way an account is created is:

1. An admin logs in and opens **Admin → Users → New User**.
2. The admin sets a username, a temporary password, and a role (`admin` or `member`).
3. The admin tells the new user their credentials directly (Slack DM, in person, however your team already shares this kind of thing) — there's no invite-email system to build or maintain.
4. The user logs in and, optionally, changes their password from their profile menu.

That's the entire flow. It's boring on purpose.

**Mechanics:**
- Passwords are hashed with **bcrypt** before they ever touch the database — the server never stores or logs a plaintext password.
- On successful login, the server issues a **JWT** containing the user's ID and role, with a sensible expiry (e.g. 7 days). No refresh-token rotation, no server-side session table — for an internal tool, "log in again after a week" is a perfectly fine trade for not having to build and maintain session infrastructure.
- Every protected REST route runs through a `RequireAuth` middleware that validates the token and attaches `userID` / `role` to the request context. Admin-only routes additionally run `RequireAdmin`.
- **Board-level permission is checked on every single WebSocket message, not just at connection time.** If someone's access to a board is revoked mid-session, the next message they send is rejected rather than trusted because they were allowed in five minutes ago.
- One practical detail worth flagging: browsers can't attach custom headers to a WebSocket handshake, so while REST calls send the JWT as an `Authorization: Bearer …` header, the WebSocket connection sends it as a `?token=` query parameter instead. The same `RequireAuth` middleware checks both, so the rest of the code doesn't need to know the difference.

If the company later wants SSO (say, integrating with an existing identity provider) instead of admin-created passwords, only the login endpoint changes — everything downstream of "we have a valid JWT with a user ID and role" stays exactly the same.

---

## REST API Reference

| Method | Path | Description | Access |
|---|---|---|---|
| POST | `/api/auth/login` | Log in with username + password, returns a JWT | Public |
| POST | `/api/auth/logout` | Client-side token discard (stateless, nothing to invalidate server-side) | Authenticated |
| GET | `/api/users/me` | Current user's profile | Authenticated |
| PATCH | `/api/users/me/password` | Change own password | Authenticated |
| GET | `/api/admin/users` | List all users | Admin |
| POST | `/api/admin/users` | Create a new user | Admin |
| PATCH | `/api/admin/users/:id` | Update a user (disable, reset password, change role) | Admin |
| DELETE | `/api/admin/users/:id` | Remove a user | Admin |
| GET | `/api/boards` | List boards the current user can access | Authenticated |
| POST | `/api/boards` | Create a new board | Authenticated |
| GET | `/api/boards/:id` | Board metadata + full element snapshot | Board member |
| PATCH | `/api/boards/:id` | Rename / update board metadata | Owner or editor |
| DELETE | `/api/boards/:id` | Delete a board | Owner or admin |
| POST | `/api/boards/:id/members` | Share the board with another user | Owner |
| DELETE | `/api/boards/:id/members/:userId` | Revoke a user's access | Owner |
| GET | `/ws/boards/:id` | Upgrades to the real-time WebSocket channel | Board member |

---

## WebSocket Protocol

One JSON envelope shape for every message, client→server and server→client alike:

```json
{
  "type": "element:update",
  "boardId": "b7e1...",
  "senderId": "u-42...",
  "payload": { "id": "el-9...", "changes": { "x": 320, "y": 180 } },
  "ts": 1755878400000
}
```

**Message types:**

| Type | Direction | Payload | Notes |
|---|---|---|---|
| `element:create` | both | full element object | New shape added |
| `element:update` | both | `{ id, changes }` | Partial patch — only changed fields |
| `element:delete` | both | `{ id }` | |
| `element:lock` | both | `{ id }` | Sender started editing this element's text |
| `element:unlock` | both | `{ id }` | Sender finished (blur, or auto-sent on disconnect) |
| `cursor:move` | both | `{ userId, x, y, color }` | Throttled client-side; never persisted |
| `presence:join` / `presence:leave` | server→client | `{ userId, displayName, color }` | Broadcast when the Hub's client set changes |
| `sync:snapshot` | server→client | `{ elements: [...], members: [...] }` | Sent once, right after connecting |

The server mostly just validates and relays — it doesn't need type-specific logic beyond "apply this to in-memory state, then broadcast to everyone else in the Hub."

---

## Backend Structure (Go + Fiber)

```
backend/
├── cmd/
│   └── server/
│       └── main.go            # entry point: config, DB connection, Fiber app, routes
├── internal/
│   ├── auth/
│   │   ├── handler.go          # POST /api/auth/login
│   │   ├── middleware.go       # RequireAuth, RequireAdmin
│   │   └── jwt.go              # token issue/parse helpers
│   ├── users/
│   │   ├── handler.go          # admin user management handlers
│   │   ├── model.go             # GORM User struct
│   │   └── repository.go        # DB queries
│   ├── boards/
│   │   ├── handler.go          # board + membership CRUD handlers
│   │   ├── model.go
│   │   └── repository.go
│   ├── elements/
│   │   ├── model.go             # GORM BoardElement struct
│   │   └── repository.go        # batched upsert for the async persistence flush
│   ├── realtime/
│   │   ├── hub.go               # Hub + Client — the ~150 lines that matter most
│   │   ├── manager.go           # map[boardID]*Hub, created/torn down lazily
│   │   └── ws_handler.go        # wires an upgraded connection into its board's Hub
│   └── database/
│       └── db.go                # GORM + Postgres connection, migrations
├── pkg/config/
│   └── config.go                # env var loading
├── go.mod
└── go.sum
```

### The Hub (the important part)

```go
package realtime

import "sync"

// Client is one connected browser tab. Send is a buffered outbound queue —
// if this client is slow to read, the Hub drops it instead of stalling
// everyone else on the board.
type Client struct {
	Conn   *websocket.Conn
	Send   chan []byte
	UserID string
}

// Hub relays messages between every client currently viewing one board.
// One Hub exists per board that's actively being viewed, and only for as
// long as someone is viewing it.
type Hub struct {
	BoardID    string
	clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
	mu         sync.RWMutex
}

func NewHub(boardID string) *Hub {
	return &Hub{
		BoardID:    boardID,
		clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte, 256),
	}
}

// Run owns this Hub's state and must be started in its own goroutine.
// Because only this goroutine ever touches `clients`, we don't need to
// lock around reads of the map from within this loop — only Register/
// Unregister/Broadcast, which are channels, cross goroutine boundaries.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			h.mu.Unlock()

		case message := <-h.Broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Send <- message:
					// queued for delivery
				default:
					// this client's queue is full — it's too slow to keep
					// up, so drop it rather than block everyone else
					close(client.Send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}
```

```go
package realtime

import "sync"

// Manager lazily creates one Hub per board and removes it once the last
// client leaves, so idle boards don't cost any memory or goroutines.
type Manager struct {
	hubs map[string]*Hub
	mu   sync.Mutex
}

func NewManager() *Manager {
	return &Manager{hubs: make(map[string]*Hub)}
}

func (m *Manager) GetOrCreateHub(boardID string) *Hub {
	m.mu.Lock()
	defer m.mu.Unlock()

	if hub, ok := m.hubs[boardID]; ok {
		return hub
	}

	hub := NewHub(boardID)
	m.hubs[boardID] = hub
	go hub.Run()
	return hub
}
```

```go
package realtime

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/websocket"

	"internal-whiteboard/internal/auth"
)

// RegisterRoutes wires /ws/boards/:id into the app. RequireAuth runs first,
// so only a logged-in, valid session ever reaches the websocket upgrade.
func RegisterRoutes(app *fiber.App, manager *Manager, jwtSecret string) {
	app.Get("/ws/boards/:id", auth.RequireAuth(jwtSecret), websocket.New(func(c *websocket.Conn) {
		boardID := c.Params("id")
		userID, _ := c.Locals("userID").(string)

		hub := manager.GetOrCreateHub(boardID)
		client := &Client{Conn: c, Send: make(chan []byte, 32), UserID: userID}

		hub.Register <- client
		defer func() { hub.Unregister <- client }()

		go writePump(client)
		readPump(hub, client) // blocks until this connection closes
	}))
}

func readPump(hub *Hub, client *Client) {
	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			break // client disconnected
		}
		// Production version: parse the JSON envelope, confirm this user
		// still has edit access to the board, apply the change to
		// in-memory state and mark the board dirty for the next async
		// persistence flush. Simplified here to keep the pattern visible.
		hub.Broadcast <- msg
	}
}

func writePump(client *Client) {
	for msg := range client.Send {
		if err := client.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
}
```

### Auth middleware

```go
package auth

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

// RequireAuth accepts a JWT either as an "Authorization: Bearer …" header
// (used by normal REST calls) or as a ?token= query parameter (used by the
// WebSocket handshake, since browsers can't set custom headers there).
// On success it stores userID/role on the context for later handlers.
func RequireAuth(secret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		tokenStr := c.Query("token")
		if tokenStr == "" {
			tokenStr = strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		}
		if tokenStr == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "missing token")
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Locals("userID", claims["sub"])
		c.Locals("role", claims["role"])
		return c.Next()
	}
}

// RequireAdmin should be chained after RequireAuth on admin-only routes.
func RequireAdmin(c fiber.Ctx) error {
	if c.Locals("role") != "admin" {
		return fiber.NewError(fiber.StatusForbidden, "admin only")
	}
	return c.Next()
}
```

### A GORM model + handler, for scale reference

```go
package users

import "time"

// User is the only way an account exists — there's no self-signup path
// anywhere in the app, so every row here was created by an admin.
type User struct {
	ID           string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Username     string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	DisplayName  string
	Role         string `gorm:"default:member"` // "admin" or "member"
	IsActive     bool   `gorm:"default:true"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
```

```go
// CreateUser lets an admin provision a new account.
func CreateUser(db *gorm.DB) fiber.Handler {
	type request struct {
		Username    string `json:"username" validate:"required,min=3"`
		Password    string `json:"password" validate:"required,min=8"`
		DisplayName string `json:"displayName"`
		Role        string `json:"role" validate:"oneof=admin member"`
	}

	return func(c fiber.Ctx) error {
		var req request
		if err := c.Bind().Body(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "could not hash password")
		}

		user := User{
			Username:     req.Username,
			PasswordHash: string(hash),
			DisplayName:  req.DisplayName,
			Role:         req.Role,
		}
		if err := db.Create(&user).Error; err != nil {
			return fiber.NewError(fiber.StatusConflict, "username already exists")
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"id":       user.ID,
			"username": user.Username,
		})
	}
}
```

---

## Frontend Structure (Next.js + shadcn/ui)

```
frontend/
├── app/
│   ├── login/page.tsx
│   ├── dashboard/page.tsx        # board list
│   ├── board/[id]/page.tsx        # the canvas
│   ├── admin/users/page.tsx       # admin user management
│   ├── layout.tsx
│   └── globals.css
├── components/
│   ├── ui/                        # shadcn-generated components (button.tsx, dialog.tsx, ...)
│   ├── canvas/
│   │   ├── Canvas.tsx              # pan/zoom container
│   │   ├── ElementRenderer.tsx     # switches on element.type
│   │   ├── StickyNote.tsx
│   │   ├── Toolbar.tsx
│   │   └── PresenceCursors.tsx
│   ├── admin/UserTable.tsx
│   └── shared/AppSidebar.tsx
├── hooks/
│   ├── useBoardSocket.ts
│   └── useAuth.ts
├── lib/
│   ├── api.ts                     # fetch wrapper for REST calls
│   └── store.ts                    # Zustand store(s)
├── types/index.ts                  # shared TS types
├── next.config.js
└── package.json
```

### Making shadcn do the actual work

The instruction to "let shadcn do most of the work" maps concretely onto specific screens:

| UI area | shadcn components |
|---|---|
| Login page | `Form`, `Input`, `Button`, `Card`, `Label` |
| Admin user list | `Table`, `Badge` (role/status), `DropdownMenu` (row actions), `Dialog` (create/edit) |
| Create/edit user | `Dialog`, `Form`, `Input`, `Select` |
| Dashboard (board list) | `Card` (per-board tile), `Button`, `Input` (search), `Skeleton` (loading) |
| App shell | `Sidebar` (or `NavigationMenu`), `Avatar`, `DropdownMenu` (user menu) |
| Canvas toolbar | `ToggleGroup` (tool selection), `Tooltip`, `Separator`, `Popover` (color/style picker) |
| Sticky-note text editing | `Textarea` directly, as the shape's editor |
| Presence indicator | `Avatar` + `Tooltip` |
| Share board dialog | `Dialog`, `Input`, `Select`, `Button` |
| Delete confirmations | `AlertDialog` |
| Notifications | `Sonner` (toast) |

A useful pattern that comes up constantly: shadcn components aren't only for chrome around the canvas — they can *be* canvas elements. A sticky note is just a styled `<div>` with a shadcn `Textarea` inside it:

```tsx
// components/canvas/StickyNote.tsx
// A single sticky note. Position/size come from the store; this component
// only handles rendering plus the text-editing interaction.
import { Textarea } from "@/components/ui/textarea";
import { useBoardStore, type BoardElement } from "@/lib/store";

export function StickyNote({ element }: { element: BoardElement }) {
  const setElement = useBoardStore((s) => s.setElement);

  return (
    <div
      className="absolute rounded-md bg-yellow-200 p-2 shadow-md"
      style={{ left: element.x, top: element.y, width: element.width, height: element.height }}
    >
      <Textarea
        className="h-full w-full resize-none border-none bg-transparent"
        value={element.content ?? ""}
        onChange={(e) => setElement({ ...element, content: e.target.value })}
      />
    </div>
  );
}
```

### The canvas itself

Rather than reaching for `<canvas>` pixel drawing or raw SVG (both of which mean writing your own hit-testing and event handling), elements are plain absolutely-positioned `<div>`s inside one container, and pan/zoom is a single CSS transform. This keeps the canvas approachable to anyone who already knows regular web development — no new rendering model to learn — and is fast enough for the element counts an internal team's boards will realistically have.

```tsx
// components/canvas/Canvas.tsx (simplified)
// Pan/zoom is one CSS transform on an inner layer — no canvas or SVG math,
// which is what keeps this file easy to read.
import { useState } from "react";

export function Canvas({ children }: { children: React.ReactNode }) {
  const [transform, setTransform] = useState({ x: 0, y: 0, scale: 1 });
  // ...wheel handler adjusts `scale`; pointer-drag on empty space adjusts x/y...

  return (
    <div className="relative h-full w-full overflow-hidden bg-muted">
      <div
        style={{
          transform: `translate(${transform.x}px, ${transform.y}px) scale(${transform.scale})`,
          transformOrigin: "0 0",
        }}
      >
        {children}
      </div>
    </div>
  );
}
```

### Client state and the socket hook

```typescript
// lib/store.ts
// All board state the UI needs, in one flat map keyed by element id.
import { create } from "zustand";

export type BoardElement = {
  id: string;
  type: "sticky_note" | "rectangle" | "ellipse" | "text" | "line" | "arrow";
  x: number;
  y: number;
  width: number;
  height: number;
  content?: string;
  style?: Record<string, string>;
};

type BoardState = {
  elements: Record<string, BoardElement>;
  setElement: (el: BoardElement) => void;
  removeElement: (id: string) => void;
  loadSnapshot: (elements: BoardElement[]) => void;
};

export const useBoardStore = create<BoardState>((set) => ({
  elements: {},
  setElement: (el) => set((s) => ({ elements: { ...s.elements, [el.id]: el } })),
  removeElement: (id) =>
    set((s) => {
      const next = { ...s.elements };
      delete next[id];
      return { elements: next };
    }),
  loadSnapshot: (elements) =>
    set(() => ({ elements: Object.fromEntries(elements.map((e) => [e.id, e])) })),
}));
```

```typescript
// hooks/useBoardSocket.ts
// One WebSocket connection per open board. Keeps the wire format and
// reconnect logic in one place so components don't touch WebSocket directly.
import { useEffect, useRef, useState, useCallback } from "react";
import { useBoardStore } from "@/lib/store";

export function useBoardSocket(boardId: string, token: string) {
  const wsRef = useRef<WebSocket | null>(null);
  const [connected, setConnected] = useState(false);
  const { setElement, removeElement, loadSnapshot } = useBoardStore();

  useEffect(() => {
    const ws = new WebSocket(
      `${process.env.NEXT_PUBLIC_WS_URL}/ws/boards/${boardId}?token=${token}`
    );
    wsRef.current = ws;

    ws.onopen = () => setConnected(true);
    ws.onclose = () => setConnected(false);
    ws.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      switch (msg.type) {
        case "sync:snapshot":
          loadSnapshot(msg.payload.elements);
          break;
        case "element:create":
        case "element:update":
          setElement(msg.payload);
          break;
        case "element:delete":
          removeElement(msg.payload.id);
          break;
        // cursor:move, presence:*, element:lock/unlock handled similarly
      }
    };

    return () => ws.close();
  }, [boardId, token]);

  const send = useCallback((type: string, payload: unknown) => {
    wsRef.current?.send(JSON.stringify({ type, payload }));
  }, []);

  return { connected, send };
}
```

---

## Security Considerations

- Passwords hashed with **bcrypt** (default cost); never logged, never stored in plaintext, never returned in any API response.
- JWTs use a reasonable expiry (e.g. 7 days) rather than complex refresh-token rotation — an appropriate trade-off for an internal tool, not a public product.
- **HTTPS/WSS enforced in production** via a reverse proxy (see Deployment) — never plain HTTP/WS outside local development.
- **CORS locked to the known frontend origin** — no wildcard.
- **Board permission is re-checked on every WebSocket message**, not only at connection time, so a revoked user is cut off immediately rather than at their next reconnect.
- **Rate-limit the login endpoint** (Fiber ships a built-in limiter middleware) to blunt brute-force attempts.
- **Validate input sizes** on element creation/update (max string lengths, reasonable bounds on width/height/position) so a malformed or malicious payload can't bloat storage or break rendering for everyone else on the board.

---

## Testing Strategy

Keep this proportional to the app's actual risk surface — the Hub's broadcast logic is the one piece of code where a bug would be genuinely hard to notice by eye, so it gets a real test:

```go
// hub_test.go
// The one guarantee the whole real-time feature depends on: a message
// from one client reaches every other client on the same board.
func TestHub_BroadcastsToOtherClients(t *testing.T) {
	hub := NewHub("test-board")
	go hub.Run()

	a := &Client{Send: make(chan []byte, 4)}
	b := &Client{Send: make(chan []byte, 4)}
	hub.Register <- a
	hub.Register <- b

	hub.Broadcast <- []byte(`{"type":"element:update"}`)

	select {
	case msg := <-b.Send:
		if !strings.Contains(string(msg), "element:update") {
			t.Fatalf("unexpected message: %s", msg)
		}
	case <-time.After(time.Second):
		t.Fatal("client B never received the broadcast")
	}
}
```

Beyond that:
- Standard Go unit tests for REST handlers (auth, board CRUD, user management), using Fiber's test utilities.
- One integration test that opens two WebSocket connections against a running test server and asserts a message from one arrives at the other — this is the real end-to-end version of the unit test above.
- Frontend: keep it light on purpose. Manual QA plus the simplest possible sanity check — **open the same board in two browser tabs and confirm changes sync both ways** — catches the overwhelming majority of real-time bugs faster than writing a heavy test suite around it.

---

## Deployment

Internal tool, modest scale — a single small VM running three containers is enough; no need for Kubernetes here.

```yaml
# docker-compose.yml
services:
  postgres:
    image: postgres:16
    environment:
      POSTGRES_DB: whiteboard
      POSTGRES_USER: whiteboard
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - pgdata:/var/lib/postgresql/data

  backend:
    build: ./backend
    environment:
      DATABASE_URL: postgres://whiteboard:${DB_PASSWORD}@postgres:5432/whiteboard
      JWT_SECRET: ${JWT_SECRET}
    depends_on:
      - postgres

  frontend:
    build: ./frontend
    environment:
      NEXT_PUBLIC_API_URL: https://whiteboard.internal.company.com/api
      NEXT_PUBLIC_WS_URL: wss://whiteboard.internal.company.com/ws
    depends_on:
      - backend

volumes:
  pgdata:
```

A simple reverse proxy handles TLS and routes the two paths to the backend:

```
whiteboard.internal.company.com {
    handle /api/* {
        reverse_proxy backend:3000
    }
    handle /ws/* {
        reverse_proxy backend:3000
    }
    handle {
        reverse_proxy frontend:3000
    }
}
```

Other notes:
- Secrets (`DB_PASSWORD`, `JWT_SECRET`) come from environment variables, never hardcoded.
- A daily `pg_dump` cron job is a cheap, sufficient backup strategy at this scale.

---

## Development Roadmap

Rough phases, not a commitment — adjust to your team's size and pace.

| Phase | Focus | Rough duration* |
|---|---|---|
| 0 | Repo scaffolding: Go module + Fiber skeleton, Next.js + shadcn init, Docker Compose for local Postgres | 2–3 days |
| 1 | Auth + admin user management (backend + admin UI) | 3–5 days |
| 2 | Board CRUD, sharing/permissions, dashboard UI | 3–4 days |
| 3 | Canvas fundamentals, single-user: create/move/resize/delete shapes, saved via REST | 5–7 days |
| 4 | Real-time collaboration: Hub, WebSocket hook, broadcasting, presence/cursors, throttling | 5–7 days |
| 5 | Hardening: reconnect handling, per-message permission checks, rate limiting, loading/error states | 3–5 days |
| 6 | Testing, deployment, internal rollout to a small pilot group | 3–5 days |

*For a team of 1–2 developers; treat these as a starting estimate, not a deadline.

Phase 0 is the natural place to start — it's almost entirely scaffolding, so it's a good way to confirm the tech stack choices above feel right before Phase 1 builds on top of them.

---

## Appendix: Dependencies

**Backend (`go.mod`):**
```
github.com/gofiber/fiber/v3
github.com/gofiber/fiber/v3/middleware/websocket   // or github.com/gofiber/contrib/v3/websocket — check current docs
gorm.io/gorm
gorm.io/driver/postgres
github.com/golang-jwt/jwt/v5
golang.org/x/crypto/bcrypt
github.com/go-playground/validator/v10
github.com/joho/godotenv
```

**Frontend (`package.json`):**
```
next
react / react-dom
zustand
react-hook-form
zod
tailwindcss
(shadcn/ui components are generated into your repo via `npx shadcn@latest add <component>`,
 rather than installed as one opaque dependency)
```

Scaffolding commands to get started once you're ready to build:
```bash
# Frontend
npx create-next-app@latest frontend --typescript --tailwind --app
cd frontend && npx shadcn@latest init

# Backend
mkdir backend && cd backend && go mod init internal-whiteboard
go get github.com/gofiber/fiber/v3 gorm.io/gorm gorm.io/driver/postgres github.com/golang-jwt/jwt/v5
```
