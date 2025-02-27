FROM node:18-alpine

WORKDIR /usr/src/mdm-frontend

COPY package.json yarn.lock ./

RUN yarn 

COPY . .

RUN yarn build

EXPOSE 3000

CMD ["yarn","preview", "--no-open"]