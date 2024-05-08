FROM node:alpine
EXPOSE 3000
WORKDIR /app/json-server
COPY db.json .
RUN npm i -g json-server
CMD [ "json-server", "-w", "db.json" ]