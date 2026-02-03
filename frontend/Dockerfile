# Development stage
FROM node:20-alpine

WORKDIR /app

# Copy package files
COPY package.json yarn.lock ./

# Install dependencies
RUN yarn install

# Copy source code
COPY . .

# Expose port
EXPOSE 4200

# Start development server
CMD ["yarn", "start", "--", "--host", "0.0.0.0", "--port", "4200"]
