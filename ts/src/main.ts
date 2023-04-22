import express from "express";

const app = express();

app.post("/new", (req, res) => {
    const rawUrl = req.query.url as string
    if (rawUrl && rawUrl.trim().length != 0) {
        try {
            const url = new URL(rawUrl)
            res.send(url)
        } catch(e) {
            res.status(409)
        }
    }
});
app.listen(4000)
