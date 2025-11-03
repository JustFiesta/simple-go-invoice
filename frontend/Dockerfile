ARG NODE_VERSION=24.9.0
ARG NGINX_VERSION=3.18

FROM node:${NODE_VERSION} AS build

WORKDIR /app

COPY package.json ./
RUN npm install

COPY . .
RUN npm run build

FROM nginx:${NGINX_VERSION} AS runtime

COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
