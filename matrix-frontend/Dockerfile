FROM nginx:1.27-alpine
ENV GO_API_URL=http://host.docker.internal:8080
ENV NGINX_ENVSUBST_FILTER=GO_API_URL
COPY nginx.conf /etc/nginx/templates/default.conf.template
COPY index.html styles.css app.js config.js /usr/share/nginx/html/
EXPOSE 80
