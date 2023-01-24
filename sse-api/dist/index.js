"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const express_1 = __importDefault(require("express"));
const body_parser_1 = __importDefault(require("body-parser"));
const cors_1 = __importDefault(require("cors"));
express_1.default.response;
const app = (0, express_1.default)();
app.use((0, cors_1.default)());
app.use(body_parser_1.default.json());
app.use(body_parser_1.default.urlencoded({ extended: false }));
app.get("/status", (request, response) => response.json({ clients: clients.length }));
let clients = [];
const facts = [];
function eventsHandler(request, response, next) {
    const headers = {
        "Content-Type": "text/event-stream",
        Connection: "keep-alive",
        "Cache-Control": "no-cache"
    };
    response.writeHead(200, headers);
    const data = `data: ${JSON.stringify(facts)}\n\n`;
    response.write(data);
    const clientId = Date.now();
    const newClient = {
        id: clientId,
        response
    };
    clients.push(newClient);
    request.on("close", () => {
        console.log(`${clientId} Connection closed`);
        clients = clients.filter((client) => client.id !== clientId);
    });
}
app.get("/events", eventsHandler);
function sendEventsToAll(newFact) {
    clients.forEach((client) => client.response.write(`data: ${JSON.stringify(newFact)}\n\n`));
}
app.post("/fact", (request, respsonse, next) => {
    const newFact = request.body;
    facts.push(newFact);
    respsonse.json(newFact);
    return sendEventsToAll(newFact);
});
const PORT = 3333;
app.listen(PORT, () => {
    console.log(`SSE server listening at http://localhost:${PORT}`);
});
