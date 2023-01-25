import express, { NextFunction, Request, Response } from "express"
import bodyParser from "body-parser"
import cors from "cors"
express.response
const app = express()

app.use(cors())
app.use(bodyParser.json())
app.use(bodyParser.urlencoded({ extended: false }))

app.get("/status", (request, response) => response.json({ clients: clients.length }))

let clients: { id: number; response: Response }[] = []
const facts: any[] = []

function eventsHandler(request: Request, response: Response, next: NextFunction) {
	console.log("Connection opened")
	const headers = {
		"Content-Type": "text/event-stream",
		Connection: "keep-alive",
		"Cache-Control": "no-cache"
	}
	response.writeHead(200, headers)

	const data = `data: ${JSON.stringify(facts)}\n\n`

	response.write(data)

	const clientId = Date.now()

	const newClient = {
		id: clientId,
		response
	}

	clients.push(newClient)

	request.on("close", () => {
		console.log(`${clientId} Connection closed`)
		clients = clients.filter((client) => client.id !== clientId)
	})
}

app.get("/events", eventsHandler)

function sendEventsToAll(newFact: any) {
	clients.forEach((client) => client.response.write(`data: ${JSON.stringify(newFact)}\n\n`))
}

app.post("/fact", (request, respsonse, next) => {
	const newFact = request.body
	facts.push(newFact)
	respsonse.json(newFact)
	return sendEventsToAll(newFact)
})

const PORT = 3333

app.listen(PORT, () => {
	console.log(`SSE server listening at http://localhost:${PORT}`)
})
