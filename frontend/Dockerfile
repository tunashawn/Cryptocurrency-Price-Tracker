FROM node:alpine

WORKDIR /app

# Copy package files and install dependencies
COPY package*.json ./
RUN npm ci

# Copy source code
COPY . .

# Build the application
RUN npm run build

CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0", "--port", "3000"]
