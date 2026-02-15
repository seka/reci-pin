# Base Stage: Install dependencies
FROM node:20-alpine AS base

WORKDIR /app

# Enable corepack for Yarn 4 support
RUN corepack enable

# Copy package files and yarn config
COPY package.json yarn.lock .yarnrc.yml ./
COPY .yarn ./.yarn

# Install dependencies
RUN yarn install

# --- Development Stage ---
FROM base AS dev

COPY . .

EXPOSE 4200

# Start development server
CMD ["yarn", "ng", "serve", "--host", "0.0.0.0", "--port", "4200"]


# --- Production Build Stage ---
FROM base AS build

COPY . .

# Build the application (Production)
# This generates dist/frontend/browser and dist/frontend/server
RUN yarn ng build --configuration production


# --- Production Runner Stage (SSR) ---
FROM node:20-alpine AS prod

WORKDIR /app

# Copy output from build stage
COPY --from=build /app/dist/frontend/server ./server
COPY --from=build /app/dist/frontend/browser ./browser

# The SSR server usually runs on port 4000 by default
EXPOSE 4000

# Start the SSR server
CMD ["node", "server/server.mjs"]
