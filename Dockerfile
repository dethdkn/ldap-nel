# ----------- Stage 1: Build Go backend -----------
FROM golang:1.24.4 AS go-builder

WORKDIR /server

# Copy Go module files and download dependencies
COPY ./backend/go.mod ./backend/go.sum ./
RUN go mod download

# Copy Go source code
COPY ./backend .

# Build Go binary with CGO enabled targeting linux/amd64
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o server .

# ----------- Stage 2: Build Nuxt frontend and prepare final image -----------
FROM node:22-slim

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable

WORKDIR /app

COPY ./frontend/package.json ./frontend/pnpm-lock.yaml ./frontend/pnpm-workspace.yaml ./frontend/.npmrc ./
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile

COPY ./frontend .
RUN pnpm run build

# Create database directory
RUN mkdir /database
VOLUME ["/database"]

# Copy Go binary from go-builder stage
COPY --from=go-builder /server/server /server/server

# Expose Nuxt port
EXPOSE 3000

# Start Go server in background and Nuxt SSR server in foreground
CMD ["sh", "-c", "/server/server & node /app/.output/server/index.mjs; wait"]