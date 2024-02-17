FROM node:20 as build-env

WORKDIR /app/
COPY . .
RUN npm i
RUN npm run build

FROM nginx:1-alpine
COPY --from=build-env /app/nginx.conf /etc/nginx/
COPY --from=build-env /app/dist /usr/share/nginx/html