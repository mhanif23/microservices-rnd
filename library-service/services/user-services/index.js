const express = require('express');
const mongoose = require('mongoose');
const userRoutes = require('./routes/users');

const app = express();
const PORT = 3000;

app.use(express.json());
app.use('/users', userRoutes);

mongoose.connect(process.env.MONGO_URI || 'mongodb://localhost:27017/users', {
  useNewUrlParser: true,
  useUnifiedTopology: true
}).then(() => {
  console.log('Connected to MongoDB');
  app.listen(PORT, () => console.log(`User Service running on port ${PORT}`));
}).catch(err => console.error(err));
