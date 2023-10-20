FROM nginx:stable-alpine AS production
COPY dist /usr/share/nginx/html
CMD ["nginx", "-g", "daemon off;"]