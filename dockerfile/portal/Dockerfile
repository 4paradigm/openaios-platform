FROM node AS build
COPY ./src/portal /
RUN npm install
RUN npm run build

FROM nginx:latest
ENV TZ=Asia/Shanghai
RUN mkdir -p /www
COPY --from=build /dist /www/
COPY ./dockerfile/portal/nginx.conf /etc/nginx/nginx.conf
EXPOSE 80

