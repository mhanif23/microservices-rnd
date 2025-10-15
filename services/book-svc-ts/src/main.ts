import 'reflect-metadata';
import { NestFactory } from '@nestjs/core';
import { AppModule } from './module';
import * as express from 'express';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  app.use(express.json());
  const port = process.env.PORT || 3001;
  await app.listen(port as number);
  console.log(`book-svc listening on ${port}`);
}
bootstrap();
