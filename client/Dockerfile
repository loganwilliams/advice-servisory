FROM node:9.4.0-alpine

WORKDIR /client

COPY package*.json ./

RUN npm install
ENV NODE_PATH=/client/node_modules

COPY . .

EXPOSE 3000

CMD ["npm", "start"]
