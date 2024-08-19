const express = require('express');
const bodyParser = require('body-parser');
const path = require('path');
const fs = require('fs');
const csv = require('csv-parser');
const app = express();
const PORT = process.env.PORT || 3000;
app.use(express.static(path.join(__dirname, 'public')));
app.use(bodyParser.json());
app.get('/', (req, res) => {
    res.sendFile(path.join(__dirname, 'public', 'index.html'));
});
app.get('/excercises',(req,res)=>{
    const results = [];
    const csvFilePath = path.join(__dirname, 'data','excercise.csv');
    console.log(csvFilePath)
    fs.createReadStream(csvFilePath)
      .pipe(csv())
      .on('data', (data) => results.push(data))
      .on('end', () => 
        // Results will be an array of objects representing each row of the CSV
        res.json(results));
      });

      app.post('/gridLog', (req, res) => {
        const gridLog = req.body.gridLog;
    
        // Define the file path where you want to save the data
        const filePath = path.join(__dirname, 'gridLog.json');
    
        // Save the data to the file
        fs.writeFile(filePath, JSON.stringify(gridLog, null, 2), (err) => {
            if (err) {
                console.error('Error writing file:', err);
                return res.status(500).json({ error: 'Failed to save data' });
            }
            res.status(200).json({ message: 'Data saved successfully' });
        });
    });
    

app.listen(PORT,'0.0.0.0', () => {
    console.log(`Server is running on http://localhost:${PORT}`);
});
