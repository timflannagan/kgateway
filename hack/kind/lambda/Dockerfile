FROM public.ecr.aws/lambda/nodejs:20

COPY hello-function.js package.json ./

RUN npm install

CMD [ "app.handler" ]
