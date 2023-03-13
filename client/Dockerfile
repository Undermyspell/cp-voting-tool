FROM node:19-alpine As build

RUN npm i -g pnpm

WORKDIR /app

ENV VITE_API_BASE_URL=https://cp-voting-tool-api.azurewebsites.net

COPY pnpm-lock.yaml /app/
RUN pnpm fetch
COPY . .
RUN pnpm install -r --offline
RUN pnpm run build

FROM nginx:stable-alpine

COPY ./nginx/default.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist /usr/share/nginx/html

WORKDIR /app

RUN chown -R nginx:nginx /app && chmod -R 755 /app && \
    chown -R nginx:nginx /var/cache/nginx && \
    chown -R nginx:nginx /var/log/nginx && \
    chown -R nginx:nginx /etc/nginx/conf.d && \
    chown -R nginx:nginx /usr/share/nginx/html

RUN touch /var/run/nginx.pid && \
    chown -R nginx:nginx /var/run/nginx.pid

USER nginx

CMD ["nginx", "-g", "daemon off;"]
