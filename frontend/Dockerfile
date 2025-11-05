ARG NODE_VERSION=24.9.0
ARG NGINX_VERSION=3.18

ARG APP_VERSION
ARG GIT_COMMIT
ARG BUILD_DATE

FROM node:${NODE_VERSION} AS build

LABEL version=${APP_VERSION} \
      commit=${GIT_COMMIT} \
      build_date=${BUILD_DATE}

WORKDIR /app

COPY package.json ./
RUN npm install

COPY . .
RUN npm run build

FROM nginx:${NGINX_VERSION} AS runtime

RUN rm /etc/nginx/conf.d/default.conf

COPY nginx.conf /etc/nginx/conf.d/default.conf

COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
