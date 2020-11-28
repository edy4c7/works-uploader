import * as functions from 'firebase-functions';
import express from 'express'
//@ts-ignore  
import basicAuth from 'basic-auth-connect'
import { compareSync } from 'bcrypt'
import { join } from 'path'

const app = express()

app.all('/*', basicAuth((user: string, password: string) => {
  return user === functions.config().basic.user
  && compareSync(password, functions.config().basic.password)
}))

app.use(express.static(join(__dirname, '../web/dist')))

exports.app = functions.https.onRequest(app)
