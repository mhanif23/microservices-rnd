import { Controller, Get, Post, Param, Query, Body } from '@nestjs/common';
import { Schema, model } from 'mongoose';

const BookSchema = new Schema({
  _id: { type: String, required: true }, // UUID string
  title: { type: String, required: true },
  authors: [{ type: String }],
  isbn: { type: String, unique: true, sparse: true },
  categories: [{ type: String }],
  tags: [{ type: String }],
  language: String,
  edition: Number,
  media: { coverUrl: String }
}, { timestamps: true });

BookSchema.index({ isbn: 1 }, { unique: true, sparse: true });
BookSchema.index({ title: 'text', authors: 1, categories: 1, tags: 1 });

const Book = model('books', BookSchema);

@Controller()
export class BookController {
  @Get('/health')
  health() { return { status: 'ok' }; }

  @Get('/books')
  async list(@Query('query') q?: string) {
    const filter: any = q ? { $text: { $search: q } } : {};
    return await Book.find(filter).limit(100).lean();
  }

  @Get('/books/:id')
  async get(@Param('id') id: string) {
    const doc = await Book.findById(id).lean();
    if (!doc) return { code: 'NOT_FOUND', message: 'not found' };
    return doc;
  }

  @Get('/books/isbn/:isbn')
  async byIsbn(@Param('isbn') isbn: string) {
    const doc = await Book.findOne({ isbn }).lean();
    if (!doc) return { code: 'NOT_FOUND', message: 'not found' };
    return doc;
  }

  @Post('/books')
  async create(@Body() body: any) {
    // Expect body contains _id (UUID string)
    if (!body._id || !body.title) return { code: 'BAD_REQUEST', message: 'missing _id or title' };
    const created = await Book.create(body);
    return created;
  }
}
