const http=require("https")
const express = require('express')
const app = express()
const jwt = require('jsonwebtoken');
const fs = require('fs')
const bodyParser = require('body-parser');
const randomstring = require("randomstring");
const fetch = require('node-fetch');
const NodeCache = require( "node-cache" );
const currencyCache = new NodeCache();
const _ = require('lodash')

app.use(bodyParser.json());
app.use(bodyParser.urlencoded({
    extended: true
  }));

var bodies = [] //local non persistent user registry
var currencyConverterApiKey="b251c2e302c482b85955"
var apiEfishery="https://stein.efishery.com/v1/storages/5e1edf521073e315924ceab4/list"
var apiCurrency="https://free.currconv.com/api/v7/convert?q=IDR_USD&compact=ultra&apiKey=b251c2e302c482b85955"
var currency={}



  app.get('/httpget', async function (req, res) {

    try{
        const request = async () => {
            
            const stainRes = await fetch(apiEfishery)
            const jsonStain = await stainRes.json()
            var curRes
            var usdRate
            if (currencyCache.get("IDR_USD")==undefined){
            curRes = await fetch(apiCurrency)
            jsonCur= await curRes.json()
            usdRate=jsonCur["IDR_USD"]
            console.log("cache is not set for currency, retreiving from API")
            currencyCache.set("IDR_USD",usdRate,10)
            }else{
                console.log("usd rate retreived from cache")
                usdRate=currencyCache.get("IDR_USD")
            }
            
             for (var i = 0, len = Object.keys(jsonStain).length; i < len; i++) {        
                if (jsonStain[i]["price"]!=null && /^\d+$/.test(jsonStain[i]["price"])){
                    jsonStain[i]["usd_price"]= jsonStain[i]["price"]*usdRate
                   // console.log(jsonStain[i]["usd_price"])
                }
            }
            res.send(jsonStain)
        }    
      await request();
    }catch (e){
    res.send(e)
    }
    

  })

  app.get('/aggregate', function (req, res) {
    
    privateClaim=isAuthorized(req,res)
    console.log(privateClaim['role'])
    if (verifyPassword(bodies,privateClaim) && privateClaim['role']=="Administrator"){
        
        try{
            const request = async () => {
                const stainRes = await fetch(apiEfishery)
                const jsonStain = await stainRes.json()
                for (var i = 0, len = Object.keys(jsonStain).length; i < len; i++) {  
                 if (jsonStain[i]["area_provinsi"]==null){
                     delete jsonStain[i]
                 }else{
                    if (jsonStain[i]['timestamp']!=null){
                    let date
                        if(jsonStain[i]['timestamp'].length==10){
                            date=new Date(jsonStain[i]['timestamp']*1000).getDate()
                            
                        }else if(jsonStain[i]['timestamp'].length==13){
                            
                            date=new Date(jsonStain[i]['timestamp'].substring(0, 10)*1000).getDate()
                        }
                        if(date>=1 && date < 8){
                            jsonStain[i]['week']="week 1"
                        }else if(date>=8 && date < 15){
                            jsonStain[i]['week']="week 2"
                        }else if(date>=15 && date < 22){
                            jsonStain[i]['week']="week 3"
                        }else if(date>=22){
                            jsonStain[i]['week']="week 4"
                        }
                    }
                 }
                }
                newJson=_.groupBy(jsonStain, 'area_provinsi');
                _.forEach(newJson, function(value,key) {
                    newJson[key]= _.groupBy(newJson[key],'week')
                  }) 
                  _.forEach(newJson, function(value,key) {
                    _.forEach(value, function(obj,key2) {           
                           value[key2]= _.maxBy(obj, 'price')
                       })
                   })            
                res.send(newJson)
            }    
           request();
        }catch (e){
        res.send(e)
        }


    }else{
        var response={};
        response["error"]="Not Authorized"
        res.send(response)
    }
  })



// main()
function isAuthorized(req, res, next) {
    if (typeof req.headers.authorization !== "undefined") {
        // retrieve the authorization header and parse out the
        // JWT using the split function
        let token = req.headers.authorization.split(" ")[1];
        console.log(token)
        //let privateKey = fs.readFileSync('./private.pem', 'utf8');
        // Here we validate that the JSON Web Token is valid and has been 
        // created using the same private pass phrase
    
        jwt.verify(token, "secret", { algorithm: "HS256" }, (err, user) => {  
            // if there has been an error...
            // console.log(user)
            if (err) {  
                // shut them out!
                res.status(500).json({ error: "Not Authorized" });
                throw new Error("Not Authorized");
            }
            // if the JWT is valid, allow them to hit
            // the intended endpoint
        });

        return jwt.verify(token, "secret", { algorithm: "HS256" })

    } else {
        // No authorization header exists on the incoming
        // request, return not authorized and throw a new error 
        res.status(500).json({ error: "Not Authorized" });
        throw new Error("Not Authorized");
    }
}

function checkBodies(bodies,reqbody){
    for (var i = 0, len = bodies.length; i < len; i++) {
         if(bodies[i]["phone"]==reqbody["phone"] && bodies[i]["name"]==reqbody["name"] && bodies[i]["role"]==reqbody["role"] ){
            return true
        }       
    }
    return false
}

function verifyPassword(bodies,body){
    for (var i = 0, len = bodies.length; i < len; i++) {
        if(bodies[i]["phone"]==body["phone"] && bodies[i]["password"]==body["password"] ){
           return true
       }       
   }
   return false
}

app.get('/', function (req, res) {
    
    res.send('Hello World')
  })



app.get('/privateClaims', function (req, res) {
    privateClaim=isAuthorized(req,res)
    if (verifyPassword(bodies,privateClaim)){
    delete privateClaim["password"]
    delete privateClaim["iat"]
    res.send(privateClaim)
    }else{
        var response={};
        response["error"]="No Data Found"
        res.send(response)
    }
  })


app.get('/fetch', async function (req, res) {
    privateClaim=isAuthorized(req,res)
    if (verifyPassword(bodies,privateClaim)){
        
        try{
            const request = async () => {
                
                const stainRes = await fetch(apiEfishery)
                const jsonStain = await stainRes.json()
                var curRes
                var usdRate
                if (currencyCache.get("IDR_USD")==undefined){
                curRes = await fetch(apiCurrency)
                jsonCur= await curRes.json()
                usdRate=jsonCur["IDR_USD"]
                console.log("cache is not set for currency, retreiving from API")
                currencyCache.set("IDR_USD",usdRate,10)
                }else{
                    console.log("usd rate retreived from cache")
                    usdRate=currencyCache.get("IDR_USD")
                }
                
                 for (var i = 0, len = Object.keys(jsonStain).length; i < len; i++) {        
                    if (jsonStain[i]["price"]!=null && /^\d+$/.test(jsonStain[i]["price"])){
                        jsonStain[i]["usd_price"]= ""+jsonStain[i]["price"]*usdRate
                       // console.log(jsonStain[i]["usd_price"])
                    }
                }
                res.send(jsonStain)
            }    
          await request();
        }catch (e){
        res.send(e)
        }

        
    }else{
        var response={};
        response["error"]="Invalid JWT"
        res.send(response)
    }
})


app.post('/jwt', (req, res) => {
    //let privateKey = fs.readFileSync('./private.pem', 'utf8'); 
    // console.log(req.body["name"])

    console.log(checkBodies(bodies,req.body))
    if (checkBodies(bodies,req.body)==false){
    req.body["password"]=randomstring.generate({
        length: 4,
      });
    bodies.push(req.body)
    }
    console.log(bodies)
    let token = jwt.sign(req.body, "secret", { algorithm: 'HS256'});
    req.body["jwt"]=token
    res.send(req.body);
})

var config=JSON.parse(fs.readFileSync('./startup_port.json'))
var port=config['port']
app.listen(config['port'], () => console.log(`Fetching app listening at http://localhost:${port}`))

 