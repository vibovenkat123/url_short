"use strict";
exports.__esModule = true;
var express_1 = require("express");
var app = (0, express_1["default"])();
app.post("/new", function (req, res) {
    var rawUrl = req.query.url;
    if (rawUrl && rawUrl.trim().length != 0) {
        var url = new URL(rawUrl);
        res.send(url);
    }
});
