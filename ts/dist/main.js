"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const express_1 = __importDefault(require("express"));
const app = (0, express_1.default)();
app.post("/new", (req, res) => {
    const rawUrl = req.query.url;
    if (rawUrl && rawUrl.trim().length != 0) {
        try {
            const url = new URL(rawUrl);
            res.send(url);
        }
        catch (e) {
            res.status(409);
        }
    }
});
app.listen(4000);
//# sourceMappingURL=main.js.map