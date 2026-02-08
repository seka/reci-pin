# Development stage
FROM node:20-alpine

WORKDIR /app

# Enable corepack for Yarn 4 support
RUN corepack enable

# Copy package files and yarn config
COPY package.json yarn.lock .yarnrc.yml ./
COPY .yarn ./.yarn

# Install dependencies
RUN yarn install

# Copy source code
COPY . .

# Expose port
EXPOSE 4200

# Start development server
CMD ["yarn", "ng", "serve", "--host", "0.0.0.0", "--port", "4200"]
