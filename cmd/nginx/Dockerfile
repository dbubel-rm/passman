FROM nginx:latest
RUN rm /etc/nginx/conf.d/default.conf
RUN rm /etc/nginx/nginx.conf
COPY default.conf /etc/nginx/conf.d/default.conf
COPY nginx.conf /etc/nginx/nginx.conf

COPY index.html /opt/htdocs/engineerbeard/index.html
COPY index2.html /opt/htdocs/onepercentergame/index2.html

