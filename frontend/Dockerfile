FROM node:latest

RUN mkdir /app
WORKDIR /app

COPY . /app

RUN npm install --silent

RUN npm run build && npm install -g serve

CMD serve -s build
