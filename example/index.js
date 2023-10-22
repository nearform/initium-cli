const express = require('express');
const bodyParser = require('body-parser');
const axios = require('axios');

const app = express();

app.use(bodyParser.json());

app.get('/', (req, res) => {
    res.status(200).send({"message": "Welcome to Initium CLI"});
});

app.post('/', async(req, res) => {
    const { mode, to_host, number } = req.body;

    console.log('mode: %s - to_host: %s - number: %s', mode, to_host, number);

    if (mode === 'sender') {
        const generatedNumber = Math.floor(Math.random() * 1000);
        console.log('Generated number is %s', generatedNumber);

        try {
            const resp = await axios.post(to_host, {
                "mode": 'receiver',
                "to_host": req.headers.host,
                "number": generatedNumber
            });

            res.send(resp.data);
        } catch (error) {
            console.log(error)
            res.status(500).send({error: 'Failed to send data'});
        }
    } else if (mode === 'receiver') {
        const generatedNumber = Math.floor(Math.random() * 1000);
        console.log('Generated number is %s', generatedNumber);

        const sum = generatedNumber + number;

        res.send({"sum": sum})
    } else {
        res.status(200).send({"message": "Welcome to Initium"});
    }
})

const PORT=8080;
app.listen(PORT, () => {
	console.log(`Server is running on ${PORT}`);
});
