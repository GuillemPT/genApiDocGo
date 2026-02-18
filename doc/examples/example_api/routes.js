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
 * Fun facts and quirks
 */
app.get('/about', (req, res) => {
    res.send('🎉 I turn coffee into code and bugs into features! Still learning to center divs though... 🤷');
});

// @api_generate_doc
/** 
 * Contact route
 * super description
 */
app.get('/contact', (req, res) => {
    res.send('This is the Contact Page!');
});

// @api_generate_doc
/**
 * Professional background and experience
 */
app.get('/cv', (req, res) => {
    const cv = {
        name: 'API Developer',
        experience: '5+ years in backend development',
        skills: ['Node.js', 'Express', 'API Design', 'Documentation'],
        education: 'Computer Science',
        languages: ['JavaScript', 'TypeScript', 'Python']
    };
    res.json(cv);
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
