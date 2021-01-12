# build environment
FROM node:15.5.1-alpine3.10 as build
WORKDIR /app
ENV PATH /app/node_modules/.bin:$PATH
COPY website/package*.json /app/
RUN npm install --silent --only=production
RUN npm install @vue/cli --silent
COPY website /app
RUN npm run build

# production environment
FROM nginx:latest
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]