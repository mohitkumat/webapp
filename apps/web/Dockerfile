# --- build stage ---
FROM node:24.21.0-alpine AS builder
WORKDIR /src

COPY package.json package-lock.json* ./
RUN npm install

COPY . .

# Vite bakes VITE_* vars in at build time, not at container start.
ARG VITE_API_URL=http://127.0.0.1:8080
ENV VITE_API_URL=$VITE_API_URL
RUN npm run build

# --- runtime stage ---
FROM nginx:1.27-alpine
COPY --from=builder /src/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
