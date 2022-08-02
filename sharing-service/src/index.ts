import express from 'express'
import dotenv from 'dotenv'
import * as redis from 'redis'


dotenv.config()

const PORT = process.env.PORT || 3000

const app = express()


const rcon = redis.createClient({
  'url': 'redis://redis:6379'
})

rcon.connect().then( () => {
  console.log("connected to redis")
  rcon.set("test", "test")

})



app.get('/:link', (req, res) => {
  res.send({"link": req.params.link})
})


app.listen(PORT, () => {
  console.log(`Starting listener on port ${PORT}`)
})
