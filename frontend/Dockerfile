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


# --- Production Runner Stage (Nginx) ---
FROM nginx:alpine AS prod

# Copy static files from build stage
COPY --from=build /app/dist/frontend/browser /usr/share/nginx/html

# Custom nginx config to handle SPA routing
RUN printf 'server {\n    listen 80;\n    location / {\n        root /usr/share/nginx/html;\n        index index.html index.htm;\n        try_files $uri $uri/ /index.html;\n    }\n}\n' > /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
