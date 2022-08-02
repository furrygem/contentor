FROM node

WORKDIR /app
COPY . /app

ENV NODE_ENV "development"

RUN yarn install
RUN yarn build

CMD ["yarn", "start:dev"]
EXPOSE 3000
