const express = require('express');
const app = express();

// This file is to test exclude directories feature
// @api_generate_doc
/**
 * remove app 
 */ 
app.delete('/app', (req, res) => {
    // Simulating a user list
    const users = [
        { id: 1, name: 'John Doe' },
        { id: 2, name: 'Jane Smith' }
    ];
    res.json(users);
});

