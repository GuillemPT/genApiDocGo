const express = require('express');
const app = express();
const PORT = process.env.PORT || 3000;

/**
 * Middleware to parse JSON
 */ 
app.use(express.json());


/**
 * Home route
 */
app.get('/', (req, res) => {
    res.send('Welcome to the Home Page!');
});

// @api_generate_doc
/**
 * About route
 */
app.get('/about', (req, res) => {
    res.send('This is the About Page!');
});

// @api_generate_doc
/** 
 * Contact route
 * super description
 */
app.get('/contact', (req, res) => {
    res.send('This is the Contact Page!');
});

/**
 * Users route
 */ 
app.get('/users', (req, res) => {
    // Simulating a user list
    const users = [
        { id: 1, name: 'John Doe' },
        { id: 2, name: 'Jane Smith' }
    ];
    res.json(users);
});

/** 
 * Start the server
 */
app.listen(PORT, () => {
    console.log(`Server is running on http://localhost:${PORT}`);
});
