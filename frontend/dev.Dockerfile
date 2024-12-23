FROM node:latest

RUN mkdir /app
WORKDIR /app

CMD ["npm", "run", "dev", "--", "--host"]
