FROM golang:1.25.0 AS go-builder

WORKDIR /server

COPY ./api/go.mod ./api/go.sum ./
RUN go mod download

COPY ./api ./

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o api .

FROM node:22-slim

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable

WORKDIR /nel

COPY ./package.json ./pnpm-lock.yaml ./pnpm-workspace.yaml ./.npmrc ./nuxt.config.ts ./tsconfig.json ./
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile

COPY ./app ./app
COPY ./public ./public
RUN pnpm run build

RUN mkdir /database
VOLUME ["/database"]

COPY --from=go-builder /server/api /nel/api

EXPOSE 3000

CMD ["sh", "-c", "/nel/api & node /nel/.output/server/index.mjs; wait"]
