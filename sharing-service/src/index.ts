import express from 'express'
import dotenv from 'dotenv'
import * as redis from 'redis'


dotenv.config()

const PORT = process.env.PORT || 3000

const app = express()


type URLInfo = {
  owner: string;
  internalID: string;
}


const rcon = redis.createClient({
  'url': 'redis://redis:6379'
})


app.get('/:link', async (req, res) => {
  let exists = await rcon.exists(req.params.link)
  if (!exists) {
    res.status(404)
    res.send("not found")
    return
  }
  let val = await rcon.hGetAll(req.params.link) as URLInfo
  console.log(val)
  res.setHeader("X-Accel-Redirect", `/api/objects/${val.internalID}`)
  res.send()
  return
})

app.get("/test/internal/:filename", (req, res) => {
  res.setHeader("X-Accel-Redirect", `/api/objects/${req.params.filename}`)
  res.send()
})

rcon.on('error', (err) => console.log('Redis Client Error', err))


rcon.connect().then( () => {
  console.log("Connected to redis")
  app.listen(PORT, () => {
    console.log(`Starting listener on port ${PORT}`)
  })
})
