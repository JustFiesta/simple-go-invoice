ARG NODE_VERSION=24.9.0

FROM node:${NODE_VERSION}
WORKDIR /app
COPY . .
RUN npm install
CMD ["npm", "run", "dev", "--", "--host"]
