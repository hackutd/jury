FROM node as client-builder
WORKDIR /client

ENV REACT_APP_JURY_NAME=$REACT_APP_JURY_NAME
ENV REACT_APP_JURY_URL=$REACT_APP_JURY_URL

RUN yarn install --frozen-lockfile

CMD [ "yarn", "run", "start" ]
