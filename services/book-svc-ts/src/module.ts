import { Module } from '@nestjs/common';
import { BookController } from './routes';
import { connect } from 'mongoose';

@Module({
  controllers: [BookController],
  providers: [{
    provide: 'MONGO_INIT',
    useFactory: async () => {
      const uri = process.env.MONGO_URI || 'mongodb://localhost:27017/books';
      await connect(uri);
      return true;
    }
  }]
})
export class AppModule {}
