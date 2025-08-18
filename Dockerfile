FROM golang:1.25.0 AS go-builder

WORKDIR /server

COPY ./backend/go.mod ./backend/go.sum ./
RUN go mod download

COPY ./backend ./

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o api .

FROM node:22-slim

WORKDIR /nel
COPY ./package.json ./

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable

WORKDIR /nel/app

COPY ./frontend/package.json ./frontend/pnpm-lock.yaml ./frontend/pnpm-workspace.yaml ./frontend/.npmrc ./
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile

COPY ./frontend ./
RUN pnpm run build

RUN mkdir /database
VOLUME ["/database"]

COPY --from=go-builder /server/api /nel/server/api

EXPOSE 3000

CMD ["sh", "-c", "/nel/server/api & node /nel/app/.output/server/index.mjs; wait"]
