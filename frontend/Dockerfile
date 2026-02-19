# Stage 1: Build (Keep your existing build steps)
FROM node:lts-alpine AS build-stage
ARG VITE_API_BASE_URL
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# Stage 2: Serve with Nginx
FROM nginx:stable-alpine as production-stage

# --- New: Copy our custom nginx config ---
COPY nginx.conf /etc/nginx/conf.d/default.conf
# -----------------------------------------

COPY --from=build-stage /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]